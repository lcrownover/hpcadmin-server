package auth

// none of this is used, but it's good to have for reference

// import (
// 	"context"
// 	"crypto/tls"
// 	"fmt"
// 	"log/slog"
// 	"net"
// 	"net/http"
// 	"net/url"
// 	"os"
//
// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/microsoft"
// )
//
// type AuthHandler struct {
// 	Ctx          context.Context
// 	Logger       *slog.Logger
// 	ListenAddr   string
// 	Oauth2Config *oauth2.Config
// 	HttpClient   *http.Client
// 	HttpServer   *http.Server
// 	HttpMux      *http.ServeMux
// 	Token        string
// 	AuthDoneCh   chan struct{}
// }
//
// func NewAuthHandler(logger *slog.Logger) *AuthHandler {
// 	ctx := context.Background()
// 	listener, err := net.Listen("tcp", ":0")
// 	if err != nil {
// 		panic(err)
// 	}
// 	authPort := listener.Addr().(*net.TCPAddr).Port
//
// 	redirectURL := fmt.Sprintf("http://localhost:%d/oauth/callback", authPort)
// 	fmt.Println(redirectURL)
// 	tenantID, found := os.LookupEnv("TENANT_ID")
// 	if !found {
// 		panic("TENANT_ID not found")
// 	}
// 	clientID, found := os.LookupEnv("CLIENT_ID")
// 	if !found {
// 		panic("CLIENT_ID not found")
// 	}
// 	conf := &oauth2.Config{
// 		ClientID:    clientID,
// 		Endpoint:    microsoft.AzureADEndpoint(tenantID),
// 		RedirectURL: redirectURL,
// 		Scopes:      []string{"openid", "profile", "offline_access"},
// 	}
//
// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
// 	sslclient := &http.Client{Transport: tr}
// 	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslclient)
// 	mux := http.NewServeMux()
// 	server := &http.Server{Addr: fmt.Sprintf("localhost:%d", authPort), Handler: mux}
// 	return &AuthHandler{
// 		Ctx:          ctx,
// 		Logger:       logger,
// 		Oauth2Config: conf,
// 		HttpClient:   nil,
// 		HttpMux:      mux,
// 		HttpServer:   server,
// 		Token:        "",
// 		AuthDoneCh:   make(chan struct{}, 1),
// 	}
// }
//
// func (h *AuthHandler) GetAuthenticationURL() string {
// 	return h.Oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
// }
//
// func (h *AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	h.Logger.Info("callbackHandler called")
// 	queryParts, err := url.ParseQuery(r.URL.RawQuery)
// 	if err != nil {
// 		h.Logger.Error(err.Error())
// 	}
// 	code := queryParts["code"][0]
//
// 	// Exchange will do the handshake to retrieve the initial access token.
// 	tok, err := h.Oauth2Config.Exchange(h.Ctx, code)
// 	if err != nil {
// 		h.Logger.Error(err.Error())
// 	}
// 	h.Token = tok.AccessToken
//
// 	// The HTTP Client returned by conf.Client will refresh the token as necessary.
// 	client := h.Oauth2Config.Client(h.Ctx, tok)
// 	h.HttpClient = client
//
// 	// // use the client to connect to HPCAdmin, I assume it loads the bearer token
// 	// // into the Authorization header automatically
// 	// serverURL := "https://hpcadmin.talapas.uoregon.edu"
// 	// resp, err := client.Get("http://" + serverURL + "/oauth/check")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// } else {
// 	// 	log.Println(color.CyanString("Authentication successful"))
// 	// }
// 	// defer resp.Body.Close()
//
// 	// show succes page
// 	successHTML := `
// <h1>Authentication Success</h1>
// <p>You are authenticated and can now return to the CLI.</p>
// `
// 	fmt.Fprint(w, successHTML)
// 	h.AuthDoneCh <- struct{}{}
// }


// this is the associated main

// func main() {
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
// 	h := auth.NewAuthHandler(logger)
//
// 	fmt.Println(color.CyanString("You will now be taken to your browser for authentication"))
//
// 	time.Sleep(1 * time.Second)
//
// 	url := h.GetAuthenticationURL()
// 	fmt.Printf("Authentication URL: %s\n", url)
// 	browser.OpenURL(url)
//
// 	time.Sleep(1 * time.Second)
//
// 	h.HttpMux.HandleFunc("/oauth/callback", h.CallbackHandler)
// 	err := h.HttpServer.ListenAndServe()
// 	if err != nil {
// 		h.Logger.Error(err.Error())
// 	}
//
// 	fmt.Printf("server info: %+v\n", h.HttpServer)
//
// 	// go func() {
// 	// 	http.HandleFunc("/oauth/callback", h.callbackHandler)
// 	// 	h.server.ListenAndServe()
// 	// }()
//
// 	// wait for token to be set
// 	select {
// 	case <-h.AuthDoneCh:
// 		h.HttpServer.Shutdown(h.Ctx)
// 		break
// 	case <-time.After(1 * time.Minute):
// 		h.HttpServer.Shutdown(h.Ctx)
// 		fmt.Println(color.RedString("Authentication timed out"))
// 		os.Exit(1)
// 	}
//
// 	fmt.Println("Authorization Token: ")
// 	fmt.Printf("%s\n", h.Token)
// }

