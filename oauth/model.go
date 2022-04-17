package oauth

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ClientConfig struct {
	RedirectUri       string
	ClientSecret      string
	ClientId          string
	Scopes            []string
	AuthorizeEndpoint string
	TokenEndPoint     string
}

type OAuthSession struct {
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

var config *ClientConfig
var session *OAuthSession
