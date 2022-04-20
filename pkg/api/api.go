package api

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/senk8/go-twitter-oauth-client/pkg/oauth"
	"log"
	"net/http"
	"os"
)


func Run() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	var config *oauth.ClientConfig
	var session *oauth.PKCESession

	config = &oauth.ClientConfig{
		RedirectUri:       os.Getenv("REDIRECT_URI"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientId:          os.Getenv("CLIENT_ID"),
		Scopes:            []string{"tweet.read", "users.read", "list.read", "list.write", "offline.access"},
		AuthorizeEndpoint: "https://twitter.com/i/oauth2/authorize",
		TokenEndPoint:     "https://api.twitter.com/2/oauth2/token",
	}

	srv := &http.Server{
		Addr: "127.0.0.1:3000",
	}
	http.HandleFunc("/login", loginHandler(config, session))
	http.HandleFunc("/callback", callbackHandler(config, session))
	log.Fatal(srv.ListenAndServe())
}

func loginHandler(config *oauth.ClientConfig, session *oauth.PKCESession) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request)	{
		session = oauth.CreatePKCESession()
		authURL := session.BuildAuthURL(config)
		w.Header().Set("Location", authURL)
		w.WriteHeader(http.StatusFound)
		return
	}
}

func callbackHandler(config *oauth.ClientConfig, session *oauth.PKCESession) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request)	{
		state := r.URL.Query().Get("state")
		if session.State != state {
			log.Println("invalid state")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		req := oauth.CreateTokenRequest(code, config, session.CodeVerifier)

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
}
