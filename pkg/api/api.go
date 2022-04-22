package api

import (
	"encoding/json"
	"github.com/senk8/go-twitter-oauth-client/pkg/oauth"
	"log"
	"net/http"
)

type OAuth2ClientServer struct {
	config *oauth.ClientConfig
	session *oauth.PKCESession
}

func (srv *OAuth2ClientServer) authHandler(w http.ResponseWriter, r *http.Request){
	srv.session = oauth.CreatePKCESession()
	authURL := srv.session.BuildAuthURL(srv.config)
	w.Header().Set("Location", authURL)
	w.WriteHeader(http.StatusFound)
	return
}

func (srv *OAuth2ClientServer) callbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if srv.session.State != state {
		log.Println("invalid state")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	req := oauth.CreateTokenRequest(code, srv.config, srv.session.CodeVerifier)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var tokenResponse oauth.TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		log.Fatal(err)
	}
	return
}

func Run() {
	srv := OAuth2ClientServer{
		config: oauth.CreateClientConfig(),
		session: nil,
	}
	http.HandleFunc("/auth", srv.authHandler)
	http.HandleFunc("/callback", srv.callbackHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}