package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	keys "github.com/lcrownover/hpcadmin-server/internal"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type AuthHandler struct {
	dbConn       *sql.DB
	oauth2Config *oauth2.Config
	tokenCh      chan string
	tokenCtx     *context.Context
	tokenTimeout time.Duration
}

func newAuthHandler(ctx context.Context) *AuthHandler {
	dbConn := ctx.Value(keys.DBConnKey).(*sql.DB)
	tenantID, found := os.LookupEnv("TENANT_ID")
	if !found {
		panic("TENANT_ID not found")
	}
	clientID, found := os.LookupEnv("CLIENT_ID")
	if !found {
		panic("CLIENT_ID not found")
	}
	clientSecret, found := os.LookupEnv("CLIENT_SECRET")
	if !found {
		panic("CLIENT_SECRET not found")
	}

	var redirectURL = fmt.Sprintf("http://%s/auth/callback", ctx.Value(keys.ListenAddrKey).(string))
	var oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     microsoft.AzureADEndpoint(tenantID),
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "offline_access"},
	}
	return &AuthHandler{
		dbConn:       dbConn,
		oauth2Config: oauth2Config,
		tokenCh:      make(chan string, 1),
		tokenTimeout: 5 * time.Minute,
	}
}

func AuthRouter(ctx context.Context) http.Handler {
	r := chi.NewRouter()
	h := newAuthHandler(ctx)
	r.Get("/", h.Authenticate)
	r.Get("/callback", h.Callback)
	return r
}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	url := h.oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	// browser.OpenURL(url)
	// redirect to url
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken := token.AccessToken
	respHTML := `
	<html>
		<head>
			<title>Auth Callback</title>
		</head>
		<body>
			<h1>Auth Callback</h1>
			<p>You have been logged in successfully. You can now close this browser tab.</p>
		</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("token", accessToken)
	fmt.Fprint(w, respHTML)
}

// func (h *AuthHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
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
