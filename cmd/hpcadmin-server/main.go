package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"

	keys "github.com/lcrownover/hpcadmin-server/internal"
	"github.com/lcrownover/hpcadmin-server/internal/api"
	"github.com/lcrownover/hpcadmin-server/internal/data"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var docs = flag.String("docs", "", "Generate router documentation")

func main() {
	var err error

	flag.Parse()

	dbRequest := data.DBRequest{
		Driver:     "postgres",
		Host:       "localhost",
		Port:       5432,
		User:       "postgres",
		Password:   "postgres",
		DBName:     "hpcadmin",
		DisableSSL: true,
	}

	dbConn, err := data.NewDBConn(dbRequest)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, keys.DBConnKey, dbConn)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Mount("/admin", api.AdminRouter())

	r.Mount("/users", api.UsersRouter(ctx))
	r.Mount("/pirgs", api.PirgsRouter(ctx))

	if *docs != "" {
		api.GenerateDocs(r, *docs)
		return
	}

	docgen.PrintRoutes(r)

	fmt.Println("Listening on :3333")
	http.ListenAndServe(":3333", r)
}
