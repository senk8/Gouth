package oauth

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type ClientConfig struct {
	RedirectUri       string
	ClientSecret      string
	ClientId          string
	Scopes            []string
	AuthorizeEndpoint string
	TokenEndPoint     string
}

type OAuthSession struct {
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

var config *ClientConfig
var session *OAuthSession

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config = &ClientConfig{
		RedirectUri:       os.Getenv("REDIRECT_URI"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientId:          os.Getenv("CLIENT_ID"),
		Scopes:            []string{"tweet.read", "users.read", "list.read", "list.write", "offline.access"},
		AuthorizeEndpoint: "https://twitter.com/i/oauth2/authorize",
		TokenEndPoint:     "https://api.twitter.com/2/oauth2/token",
	}
}

func Run() {
	srv := &http.Server{
		Addr: "127.0.0.1:3000",
	}
	http.HandleFunc("/login", login)
	http.HandleFunc("/callback", callback)
	log.Fatal(srv.ListenAndServe())
}
