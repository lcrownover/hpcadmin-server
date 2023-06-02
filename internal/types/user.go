package types

import (
	"fmt"
	"net/http"
	"time"
)

const UserKey key = "UserKey"

type UserResponse struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserResponse) Bind(r *http.Request) error {
	return nil
}

type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (u *UserRequest) Bind(r *http.Request) error {
	if u.Username == "" || u.Email == "" || u.Firstname == "" || u.Lastname == "" {
		return fmt.Errorf("missing required User fields: %+v", u)
	}
	// add in more checks like alphanumeric, length, etc.
	return nil
}

type UserStub struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
