package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

const JWT_ADMIN_ROLE = "Role.Admin"
const JWT_USER_ROLE = "Role.User"

func CheckAdmin(token *jwt.Token) bool {
	for _, role := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
		r := role.(string)
		if r == JWT_ADMIN_ROLE {
			return true
		}
	}
	return false
}

func JWTTokenIsValid(token *jwt.Token) bool {
	slog.Debug("checking if token is valid", "method", "JWTTokenIsValid")
	tokenIsValid := token.Claims.(jwt.MapClaims).VerifyExpiresAt(jwt.TimeFunc().Unix(), true) && token.Valid
	return tokenIsValid
}

func GetJWTFromTokenString(token string) (*jwt.Token, error) {
	slog.Debug("parsing token to JWT", "method", "GetJWTFromToken")

	slog.Debug("fetching keyset from azure for validation", "method", "GetJWTFromToken")
	keySet, err := jwk.Fetch(context.Background(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")
	if err != nil {
		slog.Debug(fmt.Sprintf("error fetching keyset from azure: %v", err), "method", "GetJWTFromToken")
		return nil, err
	}

	slog.Debug("validating jwt token", "method", "GetJWTFromToken")
	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			slog.Debug(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]), "method", "GetJWTFromToken")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		slog.Debug("getting kid from header", "method", "GetJWTFromToken")
		kid, ok := token.Header["kid"].(string)
		if !ok {
			slog.Debug("kid header not found", "method", "GetJWTFromToken")
			return nil, fmt.Errorf("kid header not found")
		}

		slog.Debug("looking up key", "method", "GetJWTFromToken")
		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			slog.Debug(fmt.Sprintf("key %v not found", kid), "method", "GetJWTFromToken")
			return nil, fmt.Errorf("key %v not found", kid)
		}

		slog.Debug("parsing public key", "method", "GetJWTFromToken")
		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {
			slog.Debug(fmt.Sprintf("could not parse pubkey: %v", err), "method", "GetJWTFromToken")
			return nil, fmt.Errorf("could not parse pubkey")
		}

		return publickey, nil
	})

	if err != nil {
		slog.Debug(fmt.Sprintf("error parsing token: %v", err), "method", "GetJWTFromToken")
		return nil, err
	}

	return tokenData, nil
}
