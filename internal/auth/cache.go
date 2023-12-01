package auth

import (
	"log/slog"

	"github.com/golang-jwt/jwt"
)

// AuthCache is the cache for the auth service
// It caches the user's token and the user's data

type AuthCache struct {
	userCache map[string]TokenCache
}

type TokenCache struct {
	tokenString string
	validUntil  int64
	jwtToken    *jwt.Token
}

func NewAuthCache() *AuthCache {
	return &AuthCache{
		userCache: make(map[string]TokenCache),
	}
}

func (a *AuthCache) TokenIsValid(token string) (*jwt.Token, bool, error) {
	// if the token is in our cache, it's valid and it hasn't expired, return it
	slog.Debug("checking if token is in cache", "method", "TokenIsValid")
	if cache, ok := a.userCache[token]; ok {
		slog.Debug("token is in cache", "method", "TokenIsValid")
		if cache.jwtToken.Valid && cache.validUntil > jwt.TimeFunc().Unix() {
			slog.Debug("token is valid and not expired", "method", "TokenIsValid")
			return a.userCache[token].jwtToken, true, nil
		}
	}

	// otherwise, check if the token is valid and return it
	slog.Debug("token is not in cache, parsing token", "method", "TokenIsValid")
	jwtToken, err := GetJWTFromTokenString(token)
	if err != nil {
		return nil, false, err
	}
	slog.Debug("token parsed, checking if token is valid", "method", "TokenIsValid")
	isValid := JWTTokenIsValid(jwtToken)
	if !isValid {
		slog.Debug("token is not valid, failing authentication", "method", "TokenIsValid")
		return nil, false, nil
	}
	slog.Debug("token is valid, adding to cache", "method", "TokenIsValid")
	a.userCache[token] = TokenCache{
		tokenString: token,
		validUntil:  int64(jwtToken.Claims.(jwt.MapClaims)["exp"].(float64)),
		jwtToken:    jwtToken,
	}
	slog.Debug("token added to cache, returning success", "method", "TokenIsValid")
	return jwtToken, true, nil
}
