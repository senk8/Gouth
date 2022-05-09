package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/senk8/oauth-entities/pkg/client"
)

const (
	authzEndpoint = "https://twitter.com/i/oauth2/authorize"
	tokenEndpoint = "https://api.twitter.com/2/oauth2/token"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config := &client.Config{
		ClientID:      os.Getenv("CLIENT_ID"),
		ClientSecret:  os.Getenv("CLIENT_SECRET"),
		RedirectURI:   os.Getenv("REDIRECT_URI"),
		Scopes:        []string{"tweet.read", "users.read", "list.read", "list.write", "offline.access"},
		AuthzEndpoint: authzEndpoint,
		TokenEndPoint: tokenEndpoint,
	}
	ctx := context.Background()
	oauth := client.New(config)
	tokenResponse, err := oauth.ExecFlow(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tokenResponse.AccessToken)
}
