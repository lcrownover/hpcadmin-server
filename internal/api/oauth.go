package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lcrownover/hpcadmin-lib/pkg/oauth"
	"github.com/lcrownover/hpcadmin-server/internal/auth"
	"github.com/lcrownover/hpcadmin-server/internal/config"
	"github.com/lcrownover/hpcadmin-server/internal/keys"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var ac *auth.AuthCache

func init() {
	ac = auth.NewAuthCache()
}

type OauthHandler struct {
	dbConn       *sql.DB
	oauth2Config *oauth2.Config
	tokenCh      chan string
	tokenTimeout time.Duration
}

func newOauthHandler(ctx context.Context) *OauthHandler {
	dbConn := ctx.Value(keys.DBConnKey).(*sql.DB)
	tenantID := ctx.Value(keys.ConfigKey).(*config.ServerConfig).Oauth.TenantID
	clientID := ctx.Value(keys.ConfigKey).(*config.ServerConfig).Oauth.ClientID
	clientSecret := ctx.Value(keys.ConfigKey).(*config.ServerConfig).Oauth.ClientSecret

	var redirectURL = fmt.Sprintf("http://%s/oauth/callback", ctx.Value(keys.ListenAddrKey).(string))
	var oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     microsoft.AzureADEndpoint(tenantID),
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "offline_access"},
	}
	return &OauthHandler{
		dbConn:       dbConn,
		oauth2Config: oauth2Config,
		tokenCh:      make(chan string, 1),
		tokenTimeout: 5 * time.Minute,
	}
}

func OauthRouter(ctx context.Context) http.Handler {
	r := chi.NewRouter()
	h := newOauthHandler(ctx)
	r.Get("/", h.Authenticate)
	r.Get("/url", h.GetAuthURL)
	r.Get("/callback", h.Callback)
	return r
}

func (h *OauthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	url := h.oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *OauthHandler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	url := h.oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	w.Write([]byte(url))
}

func (h *OauthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken := token.AccessToken

	respHTML := `
<h1>Authentication Success</h1>
<p>You are authenticated and can now return to the CLI.</p>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("token", accessToken)
	fmt.Fprint(w, respHTML)
}

// OauthLoader middleware ensures that a JWT token was passed and it's a valid token.
func OauthLoader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerString := r.Header.Get("Authorization")
		if bearerString == "" {
			// bearer token wasn't passed
			// so we wont load a role or anything
			next.ServeHTTP(w, r)
		} else {
			if len(bearerString) < len("Bearer ") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			tokenString := bearerString[len("Bearer "):]
			jwtToken, isValid, err := ac.TokenIsValid(tokenString)
			if err != nil || !isValid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), keys.JWTTokenKey, jwtToken)
			role := oauth.GetJWTRoleFromToken(jwtToken)
			ctx = context.WithValue(ctx, keys.RoleKey, role)
			if role != "admin" {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// TODO(lcrown): implement api key auth
// func handleAPIKey(apiKey string, w http.ResponseWriter, r *http.Request) (string, error) {
// 	role, err := ac.GetRoleFromAPIKey(apiKey)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 		return "", err
// 	}
// 	ctx := context.WithValue(r.Context(), keys.RoleKey, role)
// 	return role, nil
// }

// func (h *OauthHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
// 	// just for testing, get the stuff from env vars
// 	tenantID, found := os.LookupEnv("TENANT_ID")
// 	if !found {
// 		panic("TENANT_ID not found")
// 	}
// 	clientID, found := os.LookupEnv("CLIENT_ID")
// 	if !found {
// 		panic("CLIENT_ID not found")
// 	}
// 	clientSecret, found := os.LookupEnv("CLIENT_SECRET")
// 	if !found {
// 		panic("CLIENT_SECRET not found")
// 	}
// 	// confidential clients have a credential, such as a secret or a certificate
// 	cred, err := confidential.NewCredFromSecret(clientSecret)
// 	if err != nil {
// 		panic("Failed to create cred from client secret")
// 	}
// 	tenantLoginURL := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)
// 	confidentialClient, err := confidential.New(tenantLoginURL, clientID, cred)
// 	if err != nil {
// 		panic("Failed to create confidential client")
// 	}

// 	scopes := []string{"user.read"}
// 	result, err := confidentialClient.AcquireTokenSilent(context.TODO(), scopes)
// 	if err != nil {
// 		// cache miss, authenticate with another AcquireToken... method
// 		result, err = confidentialClient.AcquireTokenByCredential(context.TODO(), scopes)
// 		if err != nil {
// 			panic("Failed to acquire token")
// 		}
// 	}
// 	accessToken := result.AccessToken
// 	w.Write([]byte(accessToken))
// }
