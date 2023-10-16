package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/browser"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type authHandler struct {
	ctx        context.Context
	logger     *slog.Handler
	listenAddr string
	conf       *oauth2.Config
	client     *http.Client
	server     *http.Server
	token      string
	authDoneCh chan struct{}
}

func newAuthHandler(listenAddr string, logger *slog.Handler) *authHandler {
	ctx := context.Background()
	redirectURL := fmt.Sprintf("http://%s/oauth/callback", listenAddr)
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
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     microsoft.AzureADEndpoint(tenantID),
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "offline_access"},
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslclient := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslclient)
	server := &http.Server{Addr: listenAddr}
	return &authHandler{
		ctx:        ctx,
		logger:     logger,
		listenAddr: listenAddr,
		conf:       conf,
		client:     nil,
		server:     server,
		token:      "",
		authDoneCh: make(chan struct{}, 1),
	}
}

func (h *authHandler) callbackHandler(w http.ResponseWriter, r *http.Request) {
	queryParts, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.logger.Error(err.Error())
	}
	code := queryParts["code"][0]

	// Exchange will do the handshake to retrieve the initial access token.
	tok, err := h.conf.Exchange(h.ctx, code)
	if err != nil {
		h.logger.Error(err.Error())
	}
	h.token = tok.AccessToken

	// The HTTP Client returned by conf.Client will refresh the token as necessary.
	client := h.conf.Client(h.ctx, tok)
	h.client = client

	// // use the client to connect to HPCAdmin, I assume it loads the bearer token
	// // into the Authorization header automatically
	// serverURL := "https://hpcadmin.talapas.uoregon.edu"
	// resp, err := client.Get("http://" + serverURL + "/oauth/check")
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println(color.CyanString("Authentication successful"))
	// }
	// defer resp.Body.Close()

	// show succes page
	successHTML := `
<h1>Authentication Success</h1>
<p>You are authenticated and can now return to the CLI.</p>
`
	fmt.Fprint(w, successHTML)
	h.authDoneCh <- struct{}{}
}

func main() {
	listenAddr := "localhost:36664"
	logger := slog.NewTextHandler(os.Stderr, nil)

	h := newAuthHandler(listenAddr, logger)
	url := h.conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Println(color.CyanString("You will now be taken to your browser for authentication"))
	time.Sleep(1 * time.Second)
	browser.OpenURL(url)
	time.Sleep(1 * time.Second)
	log.Printf("Authentication URL: %s\n", url)

	go func() {
		http.HandleFunc("/oauth/callback", h.callbackHandler)
		h.server.ListenAndServe()
	}()

	// wait for token to be set
	select {
	case <-h.authDoneCh:
		break
	case <-time.After(5 * time.Minute):
		log.Fatal("Authentication timed out")
	}

	fmt.Println("Authorization Token: ")
	fmt.Printf("%s\n", h.token)
}
