package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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

type ClientConfig struct {
	RedirectUri       string
	ClientSecret      string
	ClientId          string
	AuthorizeEndpoint string
	TokenEndPoint     string
}

type OAuthSession struct {
	Scopes              []string
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

var config *ClientConfig

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config = &ClientConfig{
		RedirectUri:       os.Getenv("REDIRECT_URI"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientId:          os.Getenv("CLIENT_ID"),
		AuthorizeEndpoint: "https://twitter.com/i/oauth2/authorize",
		TokenEndPoint:     "https://api.twitter.com/2/oauth2/token",
	}
}

func (session *OAuthSession) buildAuthURI(config *ClientConfig) string {
	scopesString := strings.Join(session.Scopes, " ")
	u, err := url.Parse(config.AuthorizeEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", config.ClientId)
	q.Set("redirect_uri", config.RedirectUri)
	q.Set("scope", scopesString)
	q.Set("state", session.State)
	q.Set("code_challenge", session.CodeChallenge)
	q.Set("code_challenge_method", session.CodeChallengeMethod)

	u.RawQuery = q.Encode()
	escapedUrl := regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(u.String(), "$1%20")

	return escapedUrl
}

func createOAuthSession(scopes []string) *OAuthSession {
	// state
	c := 128
	b := make([]byte, c)
	rand.Read(b)
	state := base64.RawURLEncoding.EncodeToString(b)

	// code verifier
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		log.Fatal(err)
	}
	bytes := n.Bytes()
	codeVerifier := base64.RawURLEncoding.EncodeToString(bytes)

	// code challenge
	codeVerifierHash := sha256.Sum256(bytes)
	codeChallenge := base64.RawURLEncoding.EncodeToString(codeVerifierHash[:])
	codeChallengeMethod := "S256"

	return &OAuthSession{
		Scopes:              scopes,
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	scopes := []string{"tweet.read", "users.read", "list.read", "list.write", "offline.access"}
	auth := createOAuthSession(scopes)
	authURL := auth.buildAuthURI(config)
	w.Header().Set("Location", authURL)
	w.WriteHeader(http.StatusFound)
	return
}

func Run() {
	srv := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/login", login)
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}