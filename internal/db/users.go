package db

import (
	"database/sql"

	"github.com/lcrownover/hpcadmin-server/internal/types"
	_ "github.com/lib/pq"
)

func GetDBConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	return db, err
}

func GetAllUsers() (*[]types.User, error) {
	return nil, nil
}

func GetUserById(db *sql.DB, id int) (*types.User, error) {
	var user types.User
	err := db.QueryRow("SELECT id, username, email, firstname, lastname, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.Id, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func NewUser(db *sql.DB, user *types.UserCreate) (*types.User, error) {
	var newUser types.User
	err := db.QueryRow("INSERT INTO users (username, email, firstname, lastname) VALUES ($1, $2, $3, $4) RETURNING id, username, email, firstname, lastname, created_at, updated_at", user.Username, user.Email, user.Firstname, user.Lastname).Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Firstname, &newUser.Lastname, &newUser.CreatedAt, &newUser.UpdatedAt)
	return &newUser, err
}
