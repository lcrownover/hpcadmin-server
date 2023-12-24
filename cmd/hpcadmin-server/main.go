package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"

	"github.com/lcrownover/hpcadmin-server/internal/api"
	"github.com/lcrownover/hpcadmin-server/internal/auth"
	"github.com/lcrownover/hpcadmin-server/internal/data"
	"github.com/lcrownover/hpcadmin-server/internal/keys"
	"github.com/lcrownover/hpcadmin-server/internal/util"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var docs = flag.String("docs", "", "Generate router documentation")
var debug = flag.Bool("debug", false, "Enable debug mode")

func main() {
	var err error

	flag.Parse()

	util.ConfigureLogging(*debug)

	slog.Debug("Starting hpcadmin-server", "method", "main")

	// TODO(lcrown): This should be read from env, or config file
	dbRequest := data.DBRequest{
		Host:       "localhost",
		Port:       5432,
		User:       "hpcadmin",
		Password:   "superfancytestpasswordthatnobodyknows&",
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

	authCache := auth.NewAuthCache()

	ctx := context.Background()
	ctx = context.WithValue(ctx, keys.DBConnKey, dbConn)
	ctx = context.WithValue(ctx, keys.ListenAddrKey, listenAddr)
	ctx = context.WithValue(ctx, keys.AuthCacheKey, authCache)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// public routes for logging in and simple homepage
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		r.Mount("/login", api.LoginRouter(ctx)) // TODO(lcrown)
		r.Mount("/oauth", api.OauthRouter(ctx))
	})

	// private routes for authenticated users
	r.Group(func(r chi.Router) {
		r.Use(api.AuthVerifier)
		r.Mount("/api/v1", func(ctx context.Context) http.Handler {
			r := chi.NewRouter()
			r.Mount("/users", api.UsersRouter(ctx))
			r.Mount("/pirgs", api.PirgsRouter(ctx))
			return r
		}(ctx))
	})

	// admin routes for authenticated admins
	r.Group(func(r chi.Router) {
		r.Use(api.AuthVerifier)
		r.Use(api.AdminOnly)
		r.Mount("/admin", api.AdminRouter())
	})

	if *docs != "" {
		api.GenerateDocs(r, *docs)
		return
	}

	docgen.PrintRoutes(r)

	fmt.Println("Listening on " + listenAddr)
	http.ListenAndServe(listenAddr, r)
}
