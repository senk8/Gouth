package oauth

type ClientConfig struct {
	RedirectUri       string
	ClientSecret      string
	ClientId          string
	Scopes            []string
	AuthorizeEndpoint string
	TokenEndPoint     string
}
