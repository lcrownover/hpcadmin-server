package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// A completely separate router for administrator routes
func AdminRouter(ctx context.Context) chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: index"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: list accounts.."))
	})
	r.Get("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "admin: view user id %v", chi.URLParam(r, "userId"))
	})
	return r
}
