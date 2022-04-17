package oauth

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if session.State != state {
		log.Println("invalid state")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")

	values := url.Values{}
	values.Set("code", code)
	values.Add("grant_type", "authorization_code")
	values.Add("redirect_uri", config.RedirectUri)
	values.Add("code_verifier", session.CodeVerifier)

	req, err := http.NewRequest(http.MethodPost, config.TokenEndPoint, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.ClientId, config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Print(string(body))

	return
}
