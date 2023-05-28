package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/lcrownover/hpcadmin-server/internal/types"
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &UserCreate{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	article := data.Article
	dbNewArticle(article)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewArticleResponse(article))
}
}
