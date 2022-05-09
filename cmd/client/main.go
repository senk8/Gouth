package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/senk8/go-twitter-oauth-client/pkg/oauth"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config := &oauth.ClientConfig{
		RedirectURI:       os.Getenv("REDIRECT_URI"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientID:          os.Getenv("CLIENT_ID"),
		Scopes:            []string{"tweet.read", "users.read", "list.read", "list.write", "offline.access"},
		AuthorizeEndpoint: "https://twitter.com/i/oauth2/authorize",
		TokenEndPoint:     "https://api.twitter.com/2/oauth2/token",
	}
	client := oauth.CreateOAuthClient(config)

	token, err := client.GenerateAccessToken(context.Background())
	if err != nil {
		log.Fatalf("failed to generate access token: %v", err)
	}

	log.Println(token)
}
