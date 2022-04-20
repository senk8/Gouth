package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (session *OAuthSession) buildAuthURL(config *ClientConfig) string {
	scopesString := strings.Join(config.Scopes, " ")
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

func getRandomString(l int) []byte {
	b := make([]byte, l)
	rand.Read(b)
	return b
}

func createOAuthSession() *OAuthSession {
	// state
	b := getRandomString(80)
	state := base64.RawURLEncoding.EncodeToString(b)

	// code verifier
	b = getRandomString(80)
	codeVerifier := base64.RawURLEncoding.EncodeToString(b)

	// code challenge
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	hashed := h.Sum(nil)
	codeChallenge := base64.RawURLEncoding.EncodeToString(hashed[:])
	codeChallengeMethod := "S256"

	return &OAuthSession{
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session = createOAuthSession()
	authURL := session.buildAuthURL(config)
	w.Header().Set("Location", authURL)
	w.WriteHeader(http.StatusFound)
	return
}
