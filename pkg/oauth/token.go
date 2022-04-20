package oauth

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func CreateTokenRequest(code string, config *ClientConfig, codeVerifier string) *http.Request {
	values := url.Values{}
	values.Set("code", code)
	values.Add("grant_type", "authorization_code")
	values.Add("redirect_uri", config.RedirectUri)
	values.Add("code_verifier", codeVerifier)

	req, err := http.NewRequest(http.MethodPost, config.TokenEndPoint, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.ClientId, config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
