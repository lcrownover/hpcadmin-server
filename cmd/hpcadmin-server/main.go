package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"

	"github.com/lcrownover/hpcadmin-server/internal/api"
	"github.com/lcrownover/hpcadmin-server/internal/data"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var docs = flag.String("docs", "", "Generate router documentation")

func main() {
	var err error

	flag.Parse()

	connStr := "postgresql://postgres:postgres@localhost/hpcadmin?sslmode=disable"
	dbConn, err := data.GetDBConnection(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	u := api.NewUserHandler(dbConn)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", u.GetAllUsers)
		r.Post("/", u.CreateUser)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(u.UserCtx)
			r.Get("/", u.GetUser)
			r.Put("/", u.UpdateUser)
			r.Delete("/", u.DeleteUser)
		})
	})

	r.Mount("/admin", api.AdminRouter())

    if *docs != "" {
        api.GenerateDocs(r, *docs)
        return
    }

	docgen.PrintRoutes(r)

	fmt.Println("Listening on :3333")
	http.ListenAndServe(":3333", r)
}
