package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azi "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

const JWT_ADMIN_ROLE = "Role.Admin"
const JWT_USER_ROLE = "Role.User"

type AzureAuthHandlerOptions struct {
	TenantID  string
	ClientID  string
	ConfigDir string
}

type authHandler struct {
	TenantID  string
	ClientID  string
	TokenPath string
	Token     *jwt.Token
}

func NewAzureAuthHandler(opts AzureAuthHandlerOptions) *authHandler {
	h := &authHandler{
		TenantID:  opts.TenantID,
		ClientID:  opts.ClientID,
		TokenPath: opts.ConfigDir + "/token",
	}
	return h
}

func (h *authHandler) LoadToken() (*jwt.Token, error) {
	// Load a local token
	slog.Debug("loading local token")
	token, found := h.loadLocalToken(h.TokenPath)

	// If the local token is expired or invalid, get a new one
	if !found || !h.tokenValid(token) {
		slog.Debug("token is expired or invalid, getting a new one")
		// Get an Azure access token
		access_token, err := h.getAzureAccessToken()
		if err != nil {
			fmt.Printf("Error getting access token: %v\n", err)
			os.Exit(1)
		}
		token, err = h.getJWTFromToken(access_token.Token)
		if err != nil {
			fmt.Printf("Error getting JWT from token: %v\n", err)
			os.Exit(1)
		}
		h.Token = token

		// Save the token to a file
		slog.Debug("saving token")
		err = h.saveLocalToken(h.TokenPath, token)
		if err != nil {
			fmt.Printf("Error saving token: %v\n", err)
			os.Exit(1)
		}
	}

	// Show some of the properties of the user
	// fmt.Println("User: ", jwtToken.Claims.(jwt.MapClaims)["name"])
	// fmt.Println("User: ", jwtToken.Claims.(jwt.MapClaims)["preferred_username"])
	// fmt.Println("User: ", jwtToken.Claims.(jwt.MapClaims)["roles"])

	// Check if user is admin
	// isAdmin := h.checkAdmin(token)
	// if isAdmin {
	// 	fmt.Println("User is an admin")
	// }

	return token, nil
}

func (h *authHandler) checkAdmin(token *jwt.Token) bool {
	for _, role := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
		r := role.(string)
		if r == JWT_ADMIN_ROLE {
			return true
		}
	}
	return false
}

func (h *authHandler) tokenValid(token *jwt.Token) bool {
	return token.Claims.(jwt.MapClaims).VerifyExpiresAt(jwt.TimeFunc().Unix(), true) && token.Valid
}

func (h *authHandler) getJWTFromToken(token string) (*jwt.Token, error) {
	slog.Debug("parsing token to JWT")
	slog.Debug("fetching keyset from azure for validation")
	keySet, err := jwk.Fetch(context.Background(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")
	if err != nil {
		return nil, err
	}
	slog.Debug("validating jwt token")
	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}

		return publickey, nil
	})

	if err != nil {
		return nil, err
	}

	return tokenData, nil
}

func (h *authHandler) getAzureAccessToken() (*azcore.AccessToken, error) {
	slog.Debug("getting azure credential from interactive browser flow")
	cred, err := azi.NewInteractiveBrowserCredential(&azi.InteractiveBrowserCredentialOptions{
		ClientID: h.ClientID,
		TenantID: h.TenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %v", err)
	}

	scope := fmt.Sprintf("%s/.default", h.ClientID)
	slog.Debug("using credential to fetch token")
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	return &token, nil
}

func (h *authHandler) loadLocalToken(filepath string) (*jwt.Token, bool) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, false
	}
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, false
	}
	token, err := h.getJWTFromToken(string(b))
	if err != nil {
		return nil, false
	}
	return token, true
}

func (h *authHandler) saveLocalToken(filepath string, token *jwt.Token) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(token.Raw)
	if err != nil {
		return err
	}
	os.Chmod(filepath, 0600)
	return nil
}
