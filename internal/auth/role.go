package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/lcrownover/hpcadmin-server/internal/keys"
)

//
// Valid API Roles:
//
// - Admin 		-- can do anything
// - User 		-- can do anything that the associated user_id can do
// - Unknown 	-- can do nothing
//

// RoleVerifier middleware ensures that the Role is set and valid in context.
func (m *Middleware) RoleVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(keys.RoleKey).(string)
		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusUnauthorized)
			return
		}
		if role == "" {
			slog.Debug("role is empty", "package", "auth", "method", "RoleVerifier")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusUnauthorized)
			return
		}
		if role == "unknown" {
			slog.Debug("role is unknown", "package", "auth", "method", "RoleVerifier")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusUnauthorized)
			return
		}
		slog.Debug("role is valid", "package", "auth", "method", "RoleVerifier")
		ctx := r.Context()
		ctx = context.WithValue(ctx, keys.RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
