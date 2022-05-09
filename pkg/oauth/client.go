package oauth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type OAuthClient struct {
	config  *ClientConfig
	session *PKCESession
	ch      chan string
}

type ClientConfig struct {
	ClientID          string
	RedirectURI       string
	ClientSecret      string
	Scopes            []string
	AuthorizeEndpoint string
	TokenEndPoint     string
}

const (
	LISTEN_ADDR = "127.0.0.1:3000"
)

func CreateOAuthClient(config *ClientConfig) *OAuthClient {
	ch := make(chan string)
	return &OAuthClient{
		config:  config,
		session: nil,
		ch:      ch,
	}
}

func (c *OAuthClient) authHandler(w http.ResponseWriter, r *http.Request) {
	c.session = CreatePKCESession()
	authURL := c.session.BuildAuthURL(c.config)
	w.Header().Set("Location", authURL)
	w.WriteHeader(http.StatusFound)
	return
}

func (c *OAuthClient) callbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if c.session.State != state {
		log.Println("invalid state")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	req := CreateTokenRequest(code, c.config, c.session.CodeVerifier)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		log.Fatal(err)
	}
	c.ch <- tokenResponse.AccessToken
	return
}

func (c *OAuthClient) GenerateAccessToken(ctx context.Context) (string, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", c.authHandler)
	mux.HandleFunc("/callback", c.callbackHandler)
	srv := &http.Server{
		Addr:    LISTEN_ADDR,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	accessToken := <-c.ch
	if err := srv.Shutdown(ctx); err != nil {
		return "", err
	}

	return accessToken, nil

}
