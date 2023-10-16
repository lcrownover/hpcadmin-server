package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

	// TODO(lcrown): This should be read from env, or config file
	dbRequest := data.DBRequest{
		Host:       "localhost",
		Port:       5432,
		User:       "postgres",
		Password:   "postgres",
		DBName:     "hpcadmin_test",
		DisableSSL: true,
	}

	dbConn, err := data.NewDBConn(dbRequest)
	if err != nil {
		log.Fatal(err)
	}

	// TODO(lcrown): this should be in config file
	host, found := os.LookupEnv("HOST")
	if !found {
		host = "localhost"
	}
	port, found := os.LookupEnv("PORT")
	if !found {
		port = "3333"
	}
	listenAddr := fmt.Sprintf("%s:%s", host, port)

	ctx := context.Background()
	ctx = context.WithValue(ctx, keys.DBConnKey, dbConn)
	ctx = context.WithValue(ctx, keys.ListenAddrKey, listenAddr)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/oauth", api.OauthRouter(ctx))

	r.Mount("/admin", api.AdminRouter())

	r.Mount("/api/v1", func(ctx context.Context) http.Handler {
		r := chi.NewRouter()
		r.Mount("/users", api.UsersRouter(ctx))
		r.Mount("/pirgs", api.PirgsRouter(ctx))
		return r
	}(ctx))

	if *docs != "" {
		api.GenerateDocs(r, *docs)
		return
	}

	docgen.PrintRoutes(r)

	fmt.Println("Listening on " + listenAddr)
	http.ListenAndServe(listenAddr, r)
}
