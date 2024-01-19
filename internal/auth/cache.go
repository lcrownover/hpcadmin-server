package auth

import (
	"log/slog"

	"github.com/golang-jwt/jwt"
	"github.com/lcrownover/hpcadmin-lib/pkg/oauth"
)

// AuthCache is the cache for the auth service
// It caches the user's token and the user's data

type AuthCache struct {
	JWTTokenCache map[string]TokenCache
	APITokenCache map[string]APIKeyCache
}

type TokenCache struct {
	TokenString string
	ValidUntil  int64
	Role        string
	JWTToken    *jwt.Token
}

type APIKeyCache struct {
	Key  string
	Role string
}

func NewAuthCache() *AuthCache {
	return &AuthCache{
		JWTTokenCache: make(map[string]TokenCache),
		APITokenCache: make(map[string]APIKeyCache),
	}
}

// TokenIsValid checks if the token is valid and returns it if it is
func (a *AuthCache) TokenIsValid(token string) (*jwt.Token, bool, error) {
	// if the token is in our cache, it's valid and it hasn't expired, return it
	jwtToken, ok, err := a.LookupCachedToken(token)
	if err != nil {
		return nil, false, err
	}
	if ok {
		return jwtToken, true, nil
	}

	// otherwise, check if the token is valid and return it
	slog.Debug("token is not in cache, parsing token", "package", "auth", "method", "TokenIsValid")
	jwtToken, err = oauth.GetJWTFromTokenString(token)
	if err != nil {
		return nil, false, err
	}

	slog.Debug("token parsed, checking if token is valid", "package", "auth", "method", "TokenIsValid")
	isValid := oauth.JWTTokenIsValid(jwtToken)
	if !isValid {
		slog.Debug("token is not valid, failing authentication", "package", "auth", "method", "TokenIsValid")
		return nil, false, nil
	}
	slog.Debug("token is valid", "package", "auth", "method", "TokenIsValid")

	a.CacheJWTToken(token, jwtToken)

	slog.Debug("token added to cache, returning success", "package", "auth", "method", "TokenIsValid")
	return jwtToken, true, nil
}

// LookupCachedToken checks if the token is in the cache and returns it if it is
func (a *AuthCache) LookupCachedToken(token string) (*jwt.Token, bool, error) {
	slog.Debug("checking if token is in cache", "package", "auth", "method", "TokenIsValid")
	if cache, ok := a.JWTTokenCache[token]; ok {
		slog.Debug("token is in cache", "package", "auth", "method", "TokenIsValid")
		if cache.JWTToken.Valid && cache.ValidUntil > jwt.TimeFunc().Unix() {
			slog.Debug("token is valid and not expired", "package", "auth", "method", "TokenIsValid")
			return a.JWTTokenCache[token].JWTToken, true, nil
		}
	}
	return nil, false, nil
}

// CacheJWTToken adds the token to the cache
func (a *AuthCache) CacheJWTToken(token string, jwtToken *jwt.Token) {
	slog.Debug("adding to cache", "package", "auth", "method", "TokenIsValid")
	a.JWTTokenCache[token] = TokenCache{
		TokenString: token,
		ValidUntil:  int64(jwtToken.Claims.(jwt.MapClaims)["exp"].(float64)),
		JWTToken:    jwtToken,
	}
}

// LookupCachedAPIKey checks if the api key is in the cache 
// returns the role if found, "unknown" if not found
func (a *AuthCache) LookupCachedAPIKey(key string) string {
	slog.Debug("checking if api key is in cache", "package", "auth", "method", "LookupCachedAPIKey")
	if cache, ok := a.APITokenCache[key]; ok {
		slog.Debug("api key is in cache", "package", "auth", "method", "LookupCachedAPIKey")
		return cache.Role
	}
	slog.Debug("api key not found in cache", "package", "auth", "method", "LookupCachedAPIKey")
	return "unknown"
}

// CacheAPIKey adds the api key to the cache
func (a *AuthCache) CacheAPIKey(key string, role string) {
	slog.Debug("adding api key to cache", "role", role, "package", "auth", "method", "CacheAPIKey")
	a.APITokenCache[key] = APIKeyCache{
		Key:  key,
		Role: role,
	}
	slog.Debug("cached api key", "package", "auth", "method", "CacheAPIKey")
}
