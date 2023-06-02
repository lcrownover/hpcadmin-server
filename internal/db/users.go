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

func GetAllUsers(db *sql.DB) ([]*types.UserResponse, error) {
	var users []*types.UserResponse
	rows, err := db.Query("SELECT id, username, email, firstname, lastname, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user types.UserResponse
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetUserById(db *sql.DB, id int) (*types.UserResponse, error) {
	var user types.UserResponse
	err := db.QueryRow("SELECT id, username, email, firstname, lastname, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.Id, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func GetUserByUsername(db *sql.DB, username string) (*types.UserResponse, error) {
	var user types.UserResponse
	err := db.QueryRow("SELECT id, username, email, firstname, lastname, created_at, updated_at FROM users WHERE username = $1", username).Scan(&user.Id, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func NewUser(db *sql.DB, user *types.UserRequest) (*types.UserResponse, error) {
	var newUser types.UserResponse
	err := db.QueryRow("INSERT INTO users (username, email, firstname, lastname) VALUES ($1, $2, $3, $4) RETURNING id, username, email, firstname, lastname, created_at, updated_at", user.Username, user.Email, user.Firstname, user.Lastname).Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Firstname, &newUser.Lastname, &newUser.CreatedAt, &newUser.UpdatedAt)
	return &newUser, err
}
