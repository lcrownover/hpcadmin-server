package api

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lcrownover/hpcadmin-server/internal/data"
	"github.com/lcrownover/hpcadmin-server/internal/keys"
)

type UserResponse struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstname"`
	LastName   string    `json:"lastname"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func (u *UserResponse) Bind(r *http.Request) error {
	return nil
}

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newUserResponse(u *data.User) *UserResponse {
	return &UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

// newUserResponseList converts a list of UserResponse objects into a list of render.Renderer objects
func newUserResponseList(users []*data.User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		list = append(list, newUserResponse(user))
	}
	return list
}

type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func (u *UserRequest) Bind(r *http.Request) error {
	if u.Username == "" || u.Email == "" || u.FirstName == "" || u.LastName == "" {
		return fmt.Errorf("missing required User fields: %+v", u)
	}
	// add in more checks like alphanumeric, length, etc.
	return nil
}

func newUserRequest(u *data.User) *UserRequest {
	return &UserRequest{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

type UserHandler struct {
	dbConn *sql.DB
}

func UsersRouter(ctx context.Context) http.Handler {
	r := chi.NewRouter()
	h := newUserHandler(ctx)
	r.Get("/", h.GetAllUsers)
	r.Post("/", h.CreateUser)
	r.Route("/{userID}", func(r chi.Router) {
		r.Use(h.UserCtx)
		r.Get("/", h.GetUser)
		r.Put("/", h.UpdateUser)
		r.Delete("/", h.DeleteUser)
	})
	return r
}

func newUserHandler(ctx context.Context) *UserHandler {
	dbConn := ctx.Value(keys.DBConnKey).(*sql.DB)
	return &UserHandler{dbConn: dbConn}
}

// GetAllUsers returns all existing users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	searchUsername := r.URL.Query().Get("username")
	// username query parameter exists, so we are looking for a specific user
	// TODO(lcrown): why are both arms of this if statement running???
	if searchUsername != "" {
		slog.Debug("getting user by username", "package", "api", "method", "GetAllUsers")
		user, err := data.GetUserByUsername(h.dbConn, searchUsername)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		resp := newUserResponse(user)
		if err := render.Render(w, r, resp); err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
	} else {
		// username query parameter doesn't exist, so we are looking for all users
		slog.Debug("getting all users", "package", "api", "method", "GetAllUsers")
		var users []*data.User

		users, err := data.GetAllUsers(h.dbConn)
		if err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		resp := newUserResponseList(users)
		if err := render.RenderList(w, r, resp); err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("creating new user", "package", "api", "method", "CreateUser")
	userReq := &UserRequest{}
	if err := render.Bind(r, userReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	dataUser := data.UserRequest(*userReq)

	newUser, err := data.CreateUser(h.dbConn, &dataUser)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	resp := newUserResponse(newUser)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, resp)
}

// UserCtx middleware is used to load a User object from /users/{username} requests
// and then attach it to the request context. In case of failure the request is aborted
// and a 404 error response is sent to the client.
func (h *UserHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *data.User
		var err error

		userIdParam := chi.URLParam(r, "userID")
		slog.Debug("loading specific user ctx", "id", userIdParam, "package", "api", "method", "UserCtx")
		userId, err := strconv.Atoi(userIdParam)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		user, err = data.GetUserById(h.dbConn, userId)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keys.UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser returns the user in the request context
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("getting user", "package", "api", "method", "GetUser")
	user := r.Context().Value(keys.UserKey).(*data.User)
	resp := newUserResponse(user)
	if err := render.Render(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("updating user", "package", "api", "method", "UpdateUser")
	// existing user comes from the request context because
	// `userID` is part of the URL.
	user := r.Context().Value(keys.UserKey).(*data.User)
	// create a new UserRequest object
	// so that it contains all the fields of the existing user
	// then bind the request body to it so that the new values
	// from the request body are updated in the UserRequest object
	userReq := newUserRequest(user)
	if err := render.Bind(r, userReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	dataUserRequest := data.UserRequest(*userReq)
	err := data.UpdateUser(h.dbConn, user.Id, &dataUserRequest)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	updatedUser, err := data.GetUserById(h.dbConn, user.Id)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	resp := newUserResponse(updatedUser)
	render.Status(r, http.StatusOK)
	render.Render(w, r, resp)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("deleting user", "package", "api", "method", "DeleteUser")
	user := r.Context().Value(keys.UserKey).(*data.User)
	err := data.DeleteUser(h.dbConn, user.Id)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	render.Status(r, http.StatusNoContent)
}
