package auth

import (
	"database/sql"
	"net/http"

	"github.com/lcrownover/hpcadmin-server/internal/keys"
)

type Middleware struct {
	db *sql.DB
}

func NewMiddleware(db *sql.DB) *Middleware {
	return &Middleware{db: db}
}

// AdminOnly middleware restricts access to just administrators.
func (m *Middleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(keys.RoleKey).(string)
		if !ok || !(role == "admin") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
