package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/lcrownover/hpcadmin-server/internal/db"
	"github.com/lcrownover/hpcadmin-server/internal/types"
	"net/http"
)

func Run(dbConn *sql.DB) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", GetAllUsers)
		r.Post("/", CreateUser)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(UserCtx)
			// r.Get("/", db.GetUserById(dbConn, userID))
		})
	})

	fmt.Println("Listening on :3333")
	http.ListenAndServe(":3333", r)

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) { }

func GetUserById(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the user
	// context because this handler is a child of the UserCtx
	// middleware. The worst case, the recoverer middleware will save us.
	user := r.Context().Value("user").(*types.User)

	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *types.User
		var err error

		if userID := chi.URLParam(r, "userID"); userID != "" {
			user, err = db.GetUserById(dbConn, userID)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &types.UserCreate{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user := data.User
	dbNewArticle(user)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(user))
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
