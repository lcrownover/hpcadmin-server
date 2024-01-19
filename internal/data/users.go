package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	Id        int
	Username  string
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRequest struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
}

func GetAllUsers(db *sql.DB) ([]*User, error) {
	slog.Debug("getting all users from database", "package", "data", "method", "GetAllUsers")
	var users []*User
	rows, err := db.Query("SELECT id, username, email, firstname, lastname, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetUserById(db *sql.DB, id int) (*User, error) {
	slog.Debug("querying database for user", "package", "data", "method", "GetUserById")
	var user User
	err := db.QueryRow("SELECT id, username, email, firstname, lastname, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.Id, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func CreateUser(db *sql.DB, user *UserRequest) (*User, error) {
	slog.Debug("creating new user in database", "package", "data", "method", "CreateUser")
	var newUser User
	err := db.QueryRow("INSERT INTO users (username, email, firstname, lastname) VALUES ($1, $2, $3, $4) RETURNING id, username, email, firstname, lastname, created_at, updated_at", user.Username, user.Email, user.FirstName, user.LastName).Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.FirstName, &newUser.LastName, &newUser.CreatedAt, &newUser.UpdatedAt)
	return &newUser, err
}

func UpdateUser(db *sql.DB, userId int, user *UserRequest) error {
	slog.Debug("updating user in database", "package", "data", "method", "UpdateUser")
	res, err := db.Exec("UPDATE users SET username = $1, email = $2, firstname = $3, lastname = $4 WHERE id = $5 RETURNING id, username, email, firstname, lastname, created_at, updated_at", user.Username, user.Email, user.FirstName, user.LastName, userId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("expected to update 1 row, updated %d rows", count)
	}

	return err
}

func DeleteUser(db *sql.DB, id int) error {
	slog.Debug("deleting user from database", "package", "data", "method", "DeleteUser")
	res, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return err
	}
	return nil
}
