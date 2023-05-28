package types

import (
	"time"
)

type Pirg struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Owner     User      `json:"owner"`
	Admins    []User    `json:"admins"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PirgCreate struct {
	Name   string `json:"name"`
	Owner  User   `json:"owner"`
	Admins []User `json:"admins"`
}

type PirgStub struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
