package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "github.com/go-chi/docgen"
	"github.com/go-chi/render"

	db "github.com/lcrownover/hpcadmin-server/internal/db"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	var err error

	flag.Parse()

	connStr := "postgresql://postgres:postgres@localhost/hpcadmin"
	dbConn, err := db.GetDBConnection(connStr)
	if err != nil {
		fmt.Errorf("Failed to connect to database: {}", err.Error())
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", api.GetAllUsers)
		r.Post("/", database.CreateUser)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(UserCtx)
			r.Get("/", db.GetUserById(dbConn, userID))
		})
	})

	fmt.Println("Listening on :3333")
	http.ListenAndServe(":3333", r)

}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value("article").(*Article)

	if err := render.Render(w, r, NewArticleResponse(article)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *User
		var err error

		if userID := chi.URLParam(r, "userID"); userID != "" {
			user, err = db.GetUserById(dbConn, userID)
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
