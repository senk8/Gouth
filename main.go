package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	authorizeEndpoint = "https://twitter.com/i/oauth2/authorize"
	tokenEndPoint     = "https://api.twitter.com/2/oauth2/token"
)

type ClientConfig struct {
	RedirectUri       string
	ClientSecret      string
	ClientId          string
	Scopes            []string
	AuthorizeEndpoint string
	TokenEndPoint     string
}

type OAuthSessionData struct {
	ClientId            string
	RedirectUri         string
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func createClientConfig(scopes []string, authorizeEndPoint string, tokenEndPoint string) *ClientConfig {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config := &ClientConfig{
		RedirectUri:       os.Getenv("REDIRECT_URI"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientId:          os.Getenv("CLIENT_ID"),
		Scopes:            scopes,
		AuthorizeEndpoint: tokenEndPoint,
		TokenEndPoint:     authorizeEndPoint,
	}

	return config
}

func createOAuthSession(config *ClientConfig) *OAuthSessionData {
	// state
	c := 300
	b := make([]byte, c)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	// code verifier
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		log.Fatal(err)
	}
	bytes := n.Bytes()
	codeVerifier := base64.RawURLEncoding.EncodeToString(bytes)

	// code challenge
	codeVerifierHash := sha256.Sum256(bytes)
	codeChallenge := base64.StdEncoding.EncodeToString(codeVerifierHash[:])
	codeChallengeMethod := "s256"

	return &OAuthSessionData{
		ClientId:            config.ClientId,
		RedirectUri:         config.RedirectUri,
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}

func (data *OAuthSessionData) buildAuthURI() string {
	scopesString := strings.Join(data.Scopes, " ")
	u, err := url.Parse(authorizeEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", data.ClientId)
	q.Set("redirect_uri", data.RedirectUri)
	q.Set("scope", scopesString)
	q.Set("state", data.State)
	q.Set("code_challenge", data.CodeChallenge)
	q.Set("code_challenge_method", data.CodeChallengeMethod)

	u.RawQuery = q.Encode()
	escapedUrl := regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(u.String(), "$1%20")

	return escapedUrl
}

func FrontChannel(config *ClientConfig) {
	auth := createOAuthSession(config)
	authEndPoint := auth.buildAuthURI()
	fmt.Println(authEndPoint)
}

func BackChannel(client *http.Client, authCode string, data *OAuthSessionData, config *ClientConfig) string {
	u, err := url.Parse(tokenEndPoint)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("code", authCode)
	q.Set("grant_type", config.ClientId)
	q.Set("redirect_uri", config.RedirectUri)
	q.Set("client_id", config.ClientId)
	q.Set("code_verifier", data.CodeVerifier)
	u.RawQuery = q.Encode()

	str := regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(u.String(), "$1%20")

	req, err := http.NewRequest("POST", str, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.ClientId, config.ClientSecret)
	req.Header.Set("Content-Type: ", "application/x-www-form-urlencoded")

	/*
		res, err := client.Do(req)
		if err != nil {

		}
	*/

	return str
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hoge")
}

func main() {
	scopes := []string{"tweet.read", "users.read", "list.read", "list.write", "offline.access"}
	config := createClientConfig(scopes)
	FrontChannel(config)
	srv := &http.Server{
		Addr: "127.0.0.1:3000",
	}
	http.HandleFunc("/login", login)
	err := srv.ListenAndServe()
	if err != nil {

	}
}
