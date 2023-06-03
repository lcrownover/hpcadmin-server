package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lcrownover/hpcadmin-server/internal/data"
)

type key string

const UserKey key = "UserKey"

type UserResponse struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserResponse) Bind(r *http.Request) error {
	return nil
}

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
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

type UserStub struct {
	Id       int
	Username string
}

type UserHandler struct {
	dbConn *sql.DB
}

func NewUserHandler(dbConn *sql.DB) *UserHandler {
	return &UserHandler{dbConn: dbConn}
}

// GetAllUsers returns all existing users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []*data.User

	users, err := data.GetAllUsers(h.dbConn)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.RenderList(w, r, newUserResponseList(users)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// GetUserById returns a single user by id, but is not currently used
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey).(*data.User)

	if err := render.Render(w, r, newUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &UserRequest{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	dataUser := data.UserRequest(*user)

	newUser, err := data.CreateUser(h.dbConn, &dataUser)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, newUserResponse(newUser))
}

// UserCtx middleware is used to load a User object from /users/{username} requests
// and then attach it to the request context. In case of failure the request is aborted
// and a 404 error response is sent to the client.
func (h *UserHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *data.User
		var err error

		userIdParam := chi.URLParam(r, "userID")
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

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser returns the user in the request context
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey).(*data.User)
	if err := render.Render(w, r, newUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey).(*data.User)
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

	render.Status(r, http.StatusOK)
	render.Render(w, r, newUserResponse(updatedUser))
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey).(*data.User)
	err := data.DeleteUser(h.dbConn, user.Id)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	render.Status(r, http.StatusNoContent)
}

// Helpers
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

func newUserRequest(u *data.User) *UserRequest {
	return &UserRequest{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
