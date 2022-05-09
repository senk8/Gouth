package client

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/senk8/oauth-entities/pkg/util"
)

const (
	codeVerifierLength  = 32
	stateLength         = 32
	codeChallengeMethod = "S256"
)

type Session struct {
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func newSession() *Session {
	r := util.GetRandomBytes(stateLength)
	state := base64.RawURLEncoding.EncodeToString(r)

	r = util.GetRandomBytes(codeVerifierLength)
	codeVerifier := base64.RawURLEncoding.EncodeToString(r)

	h := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(h[:])

	fmt.Println(len(codeChallenge))

	return &Session{
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}
