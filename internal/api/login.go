package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoginRouter(ctx context.Context) http.Handler {
	r := chi.NewRouter()
	// h := newLoginHandler(ctx)
	// r.Get("/", h.Authenticate)
	return r
}


