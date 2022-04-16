package oauth

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func buildTokenURL(code string) string {
	u, err := url.Parse(config.TokenEndPoint)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("code", code)
	q.Set("grant_type", "authorization_code")
	q.Set("redirect_uri", config.RedirectUri)
	q.Set("client_id", config.ClientId)
	q.Set("code_verifier", session.CodeVerifier)
	u.RawQuery = q.Encode()

	return regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(u.String(), "$1%20")
}

func callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	url := buildTokenURL(code)

	req, err := http.NewRequest("POST", url, nil)
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

	return
}
