package types

import (
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

type UserCreate struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UserStub struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
