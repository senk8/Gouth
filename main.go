package main

import "github.com/senk8/Gouth/oauth"

/*
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

		res, err := client.Do(req)
		if err != nil {

		}

	return str
}
*/

func main() {
	oauth.Run()
}
