package api

import (
	"context"
	"net/http"

	"github.com/lcrownover/hpcadmin-server/internal/keys"
)

type APIRole int

const (
	Admin APIRole = iota
	User
	Unauthorized
)

func (r APIRole) String() string {
	return [...]string{"admin", "user", "unauthorized"}[r]
}

func APIRoleFromString(role string) APIRole {
	switch role {
	case "admin":
		return Admin
	case "user":
		return User
	default:
		return Unauthorized
	}
}

// RoleVerifier middleware ensures that the Role is set and valid in context.
func RoleVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleStr, ok := r.Context().Value(keys.RoleKey).(string)
		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusUnauthorized)
			return
		}
		if roleStr == "" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusUnauthorized)
			return
		}
		role := APIRoleFromString(roleStr)
		if role == Unauthorized {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, keys.RoleKey, role.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
