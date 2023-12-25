package auth

import (
	"context"
	"net/http"

	"github.com/lcrownover/hpcadmin-server/internal/data"
	"github.com/lcrownover/hpcadmin-server/internal/keys"
)

// APIKeyLoader middleware checks the provided api key against the database
// and sets the role if found
func (m *Middleware) APIKeyLoader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			// api key wasnt passed, so we'll just continue
			// not setting role
			next.ServeHTTP(w, r)
			return
		}
		// api key was passed
		ctx := r.Context()

		// store the api key in the context
		ctx = context.WithValue(ctx, keys.APIKey, apiKey)

		// lets check the cache
		cachedRole := ac.LookupCachedAPIKey(apiKey)

		// if the role is not unknown,
		// that means it's a valid role
		if cachedRole != "unknown" {
			// api key and valid role was found in cache,
			// so we'll set the role, cache it, and continue
			ctx = context.WithValue(ctx, keys.RoleKey, cachedRole)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// not found in cache, so we'll check the database
		apiKeyEntry, err := data.GetAPIKeyEntry(m.db, apiKey)
		if err != nil {
			// error getting api key entry from database
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if apiKeyEntry == nil {
			// api key wasnt found in the database
			// cache the unknown key and continue
			ac.CacheAPIKey(apiKey, "unknown")
			next.ServeHTTP(w, r)
			return
		}

		// api key found in database, cache it and continue
		ac.CacheAPIKey(apiKey, apiKeyEntry.Role)
		ctx = context.WithValue(ctx, keys.RoleKey, apiKeyEntry.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
