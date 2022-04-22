package oauth

import (
	"github.com/joho/godotenv"
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

func CreateClientConfig() *ClientConfig{
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &ClientConfig{
		RedirectUri:       os.Getenv("REDIRECT_URI"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientId:          os.Getenv("CLIENT_ID"),
		Scopes:            []string{"tweet.read", "users.read", "list.read", "list.write", "offline.access"},
		AuthorizeEndpoint: "https://twitter.com/i/oauth2/authorize",
		TokenEndPoint:     "https://api.twitter.com/2/oauth2/token",
	}
}
