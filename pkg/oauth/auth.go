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

const (
	CODE_VERIFIER_LENGTH = 32
	STATE_LENGTH         = 32
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
	q.Set("client_id", config.ClientID)
	q.Set("redirect_uri", config.RedirectURI)
	q.Set("scope", scopesString)
	q.Set("state", session.State)
	q.Set("code_challenge", session.CodeChallenge)
	q.Set("code_challenge_method", session.CodeChallengeMethod)

	u.RawQuery = q.Encode()
	escapedUrl := regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(u.String(), "$1%20")

	return escapedUrl
}

func CreatePKCESession() *PKCESession {
	r := util.GetRandomString(STATE_LENGTH)
	state := base64.RawURLEncoding.EncodeToString(r)

	r = util.GetRandomString(CODE_VERIFIER_LENGTH)
	codeVerifier := base64.RawURLEncoding.EncodeToString(r)

	hashed := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hashed[:])
	codeChallengeMethod := "S256"

	return &PKCESession{
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}
