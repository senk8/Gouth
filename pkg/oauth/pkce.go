package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/senk8/go-twitter-oauth-client/pkg/util"
	"log"
	"net/url"
	"regexp"
	"strings"
)

type PKCESession struct {
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func (session *PKCESession) BuildAuthURL(config *ClientConfig) string {
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

func CreatePKCESession() *PKCESession {
	// state
	b := util.GetRandomString(80)
	state := base64.RawURLEncoding.EncodeToString(b)

	// code verifier
	b = util.GetRandomString(80)
	codeVerifier := base64.RawURLEncoding.EncodeToString(b)

	// code challenge
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	hashed := h.Sum(nil)
	codeChallenge := base64.RawURLEncoding.EncodeToString(hashed[:])
	codeChallengeMethod := "S256"

	return &PKCESession{
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}
