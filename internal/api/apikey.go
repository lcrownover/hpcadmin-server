package api

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/lcrownover/hpcadmin-server/internal/data"
	"github.com/lcrownover/hpcadmin-server/internal/keys"
)

// APIKeyLoader middleware checks the provided api key against the database
// and sets the role if found
func APIKeyLoader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(lcrown): implement auth cache for api key
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			// api key wasnt passed, so we'll just continue
			// not setting role
			next.ServeHTTP(w, r)
		} else {
			ctx := context.WithValue(r.Context(), keys.APIKey, apiKey)

			dbConn := ctx.Value(keys.DBConnKey).(*sql.DB)
			apiKeyToken, err := data.GetAPIKeyToken(dbConn, apiKey)
			if err != nil {
				slog.Debug("api key not loaded", "method", "APIKeyLoader")
			}
			if apiKeyToken != nil {
				ctx = context.WithValue(ctx, keys.RoleKey, role)
				if role != "admin" {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}

// TODO(lcrown): implement API Key auth
