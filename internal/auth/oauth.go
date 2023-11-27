package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/browser"

	"github.com/lcrownover/hpcadmin-server/internal/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

const JWT_ADMIN_ROLE = "Role.Admin"
const JWT_USER_ROLE = "Role.User"

type AzureAuthHandlerOptions struct {
	TenantID            string
	ClientID            string
	ConfigDir           string
	SkipTLSVerification bool
}

type AuthHandler struct {
	Ctx          context.Context
	ListenAddr   string
	Oauth2Config *oauth2.Config
	HttpClient   *http.Client
	HttpServer   *http.Server
	HttpMux      *http.ServeMux
	Token        string
	AuthDoneCh   chan struct{}
}

func getRandomPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port
}

func newAzureOauth2Config(AuthPort int, TenantID string, ClientID string) *oauth2.Config {
	redirectURL := fmt.Sprintf("http://localhost:%d/oauth/callback", AuthPort)
	slog.Debug("redirectURL", "value", redirectURL, "method", "newAzureOauth2Config")
	scopes := []string{fmt.Sprintf("%s/.default", ClientID)}
	slog.Debug("scopes", "value", scopes, "method", "newAzureOauth2Config")
	return &oauth2.Config{
		ClientID:    ClientID,
		Endpoint:    microsoft.AzureADEndpoint(TenantID),
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}
}

func NewAuthHandler(opts AzureAuthHandlerOptions) *AuthHandler {
	ctx := context.Background()
	authPort := getRandomPort()
	slog.Debug("authPort", "value", authPort, "method", "NewAuthHandler")

	// oauth2 config includes things like the client id,
	// the auth endpoint, redirectURL, and scopes
	oauthConfig := newAzureOauth2Config(authPort, opts.TenantID, opts.ClientID)

	// register a custom http client that maybe skips SSL verification
	// and store it in ctx
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.SkipTLSVerification},
	}
	sslclient := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslclient)

	// create a new http server and mux
	mux := http.NewServeMux()
	server := &http.Server{Addr: fmt.Sprintf(":%d", authPort), Handler: mux}
	return &AuthHandler{
		Ctx:          ctx,
		Oauth2Config: oauthConfig,
		HttpClient:   sslclient,
		HttpMux:      mux,
		HttpServer:   server,
		Token:        "",
		AuthDoneCh:   make(chan struct{}, 1),
	}
}

func (h *AuthHandler) GetAuthenticationURL() string {
	return h.Oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (h *AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("callbackHandler called", "method", "CallbackHandler")
	slog.Debug("parsing query string", "method", "CallbackHandler")
	queryParts, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		slog.Error(err.Error())
	}
	code := queryParts["code"][0]

	// Exchange will do the handshake to retrieve the initial access token.
	slog.Debug("exchanging code for token", "method", "CallbackHandler")
	tok, err := h.Oauth2Config.Exchange(h.Ctx, code)
	if err != nil {
		slog.Debug(err.Error())
		fmt.Fprintf(w, "Authentication Code exchange failed")
		os.Exit(1)
	}
	h.Token = tok.AccessToken

	// The HTTP Client returned by conf.Client will refresh the token as necessary.
	client := h.Oauth2Config.Client(h.Ctx, tok)
	h.HttpClient = client

	// show succes page
	slog.Debug("showing success page", "method", "CallbackHandler")
	successHTML := `
<h1>Authentication Success</h1>
<p>You are authenticated and can now return to the CLI.</p>
`
	fmt.Fprint(w, successHTML)
	slog.Debug("sending auth done signal", "method", "CallbackHandler")
	h.AuthDoneCh <- struct{}{}
	slog.Debug("callbackHandler finished", "method", "CallbackHandler")
}

func (h *AuthHandler) Authenticate() error {
	var err error
	util.InfoPrint("You will now be taken to your browser for authentication")

	time.Sleep(1 * time.Second)

	url := h.GetAuthenticationURL()
	err = browser.OpenURL(url)
	if err != nil {
		return fmt.Errorf("error opening browser: %v", err)
	}

	time.Sleep(1 * time.Second)

	go func() {
		h.HttpMux.HandleFunc("/oauth/callback", h.CallbackHandler)
		slog.Debug("Starting server", "method", "Authenticate")
		err := h.HttpServer.ListenAndServe()
		if err != nil {
			// This is normal behavior when the server shuts down
			slog.Error("Server no longer listening", "method", "Authenticate")
		}
	}()

	for n := 0; n < 1; {
		select {
		case <-h.AuthDoneCh:
			slog.Debug("Authentication successful, shutting down server", "method", "Authenticate")
			h.HttpServer.Shutdown(h.Ctx)
			slog.Debug("Server shut down", "method", "Authenticate")
			n++
		case <-time.After(1 * time.Minute):
			slog.Debug("Authentication failed, shutting down server", "method", "Authenticate")
			h.HttpServer.Shutdown(h.Ctx)
			slog.Debug("Server shut down", "method", "Authenticate")
			return fmt.Errorf("authentication timed out")
		}
	}

	return nil
}
