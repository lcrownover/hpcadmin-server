package types

import (
	"net/http"
	"time"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Bind(r *http.Request) error {
	return nil
}

type UserCreate struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (u *UserCreate) Bind(r *http.Request) error {
	return nil
}

type UserStub struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
