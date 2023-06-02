package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lcrownover/hpcadmin-server/internal/db"
	"github.com/lcrownover/hpcadmin-server/internal/types"
)

type UserHandler struct {
	dbConn *sql.DB
}

// UserCtx middleware is used to load a User object from /users/{username} requests
// and then attach it to the request context. In case of failure the request is aborted
// and a 404 error response is sent to the client.
func (h *UserHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *types.UserResponse
		var err error

		username := chi.URLParam(r, "username")
		if username == "" {
			render.Render(w, r, ErrNotFound)
			return
		}
		user, err = db.GetUserByUsername(h.dbConn, username)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAllUsers returns all existing users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []*types.UserResponse

	users, err := db.GetAllUsers(h.dbConn)
	if err != nil {
		render.Render(w, r, ErrInternalServer)
		return
	}

	if err := render.Render(w, r, NewAPIResponse(users)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// GetUserById returns a single user by id, but is not currently used
// func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
// 	user := r.Context().Value(types.UserKey).(*types.User)

// 	if err := render.Render(w, r, NewAPIResponse(user)); err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}
// }

// GetUserByUsername returns a single user by username, which is basically the primary key
func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(types.UserKey).(*types.UserResponse)

	if err := render.Render(w, r, NewAPIResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &types.UserRequest{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newUser, err := db.NewUser(h.dbConn, user)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewAPIResponse(newUser))
}
