package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lcrownover/hpcadmin-server/internal/db"
	"github.com/lcrownover/hpcadmin-server/internal/types"
)

type UserHandler struct {
	dbConn *sql.DB
}

func (h *UserHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *types.User
		var err error

		if userID := chi.URLParam(r, "userID"); userID != "" {
			id, err := strconv.Atoi(userID)
			if err != nil {
				render.Render(w, r, ErrNotFound)
				return
			}
			user, err = db.GetUserById(h.dbConn, id)
			if err != nil {
				render.Render(w, r, ErrNotFound)
			}
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

func (a *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

}

func (a *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the user
	// context because this handler is a child of the UserCtx
	// middleware. The worst case, the recoverer middleware will save us.
	user := r.Context().Value("user").(*types.User)

	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func (a *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &types.UserCreate{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newUser, err := db.NewUser(a.dbConn, user)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(newUser))
}

type UserResponse struct {
	User *types.User `json:"user,omitempty"`
}

func NewUserResponse(user *types.User) *UserResponse {
	resp := &UserResponse{User: user}
	return resp
}

func (ur *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UserRequest struct {
	User *types.UserCreate `json:"user,omitempty"`
}

func (ur *UserRequest) Bind(r *http.Request) error {
	if ur.User == nil {
		return fmt.Errorf("missing required User fields: %+v", ur)
	}
	return nil
}
