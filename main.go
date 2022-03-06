package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	authorizeEndpoint = "https://twitter.com/i/oauth2/authorize"
	tokenEndPoint     = "https://api.twitter.com/2/oauth2/token"
	listEndPoint      = "https://api.twitter.com/2/lists"
)

type Config struct {
	RedirectUri  string
	ClientSecret string
	ClientId     string
}

type OAuthSessionData struct {
	ClientId            string
	RedirectUri         string
	Scopes              []string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type MyHandler struct {
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
	fmt.Fprintf(w, r.URL.String())
}

func (data *OAuthSessionData) PrepareAuthEndPoint() string {
	scopesString := strings.Join(data.Scopes, " ")
	u, err := url.Parse(authorizeEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", data.ClientId)
	q.Set("redirect_uri", data.RedirectUri)
	q.Set("scope", scopesString)
	q.Set("state", data.State)
	q.Set("code_challenge", data.CodeChallenge)
	q.Set("code_challenge_method", data.CodeChallengeMethod)

	u.RawQuery = q.Encode()
	//req, err := http.Get(u.String())

	str := regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(u.String(), "$1%20")

	return str
}

func oauth2UserTokenWithCredential(client *http.Client, data *OAuthSessionData) {
	authEndPoint := data.PrepareAuthEndPoint()

	req, err := http.NewRequest(http.MethodGet, authEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}
	_ = byteArray

	fmt.Println(authEndPoint)
	//fmt.Printf("%#v", string(byteArray))1
}

func InitTLS() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"Senk"},
		OrganizationalUnit: []string{"Gouth"},
		CommonName:         "Gouth Mock",
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	certOut.Close()

	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk),
	})
	keyOut.Close()
}

func FrontChannel(client *http.Client) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	redirectUri := os.Getenv("REDIRECT_URI")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientId := os.Getenv("CLIENT_ID")

	handler := MyHandler{}
	config := Config{
		RedirectUri:  redirectUri,
		ClientSecret: clientSecret,
		ClientId:     clientId,
	}

	c := 300
	b := make([]byte, c)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	scopes := []string{"tweet.read", "users.read", "list.read", "list.write"}
	codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
	codeChallengeMethod := "s256"

	auth := OAuthSessionData{
		ClientId:            config.ClientId,
		RedirectUri:         config.RedirectUri,
		Scopes:              scopes,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	server := http.Server{
		Addr:    "127.0.0.1:3000",
		Handler: &handler,
	}

	oauth2UserTokenWithCredential(client, &auth)

	err = server.ListenAndServe()
	//err := server.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

func BackChannel(client *http.Client) {

}

func main() {
	client := new(http.Client)

	FrontChannel(client)
	BackChannel(client)

	/*
	   jsonString := "{\"name\":\"test v2 create list\"}"
	   body := bytes.NewBuffer([]byte(jsonString))
	   endPoint := "https://api.twitter.com/2/lists"

	   client := new(http.Client)

	   req, err := http.NewRequest(http.MethodPost, endPoint, body)
	   if err != nil {
	       log.Fatal(err)
	   }

	   req.Header.Set("Content-type", "application/json")
	   req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	   dump, _ := httputil.DumpRequest(req, true)
	   fmt.Println(string(dump))

	   resp, err := client.Do(req)
	   if err != nil {
	       log.Fatal(err)
	   }
	   defer resp.Body.Close()

	   byteArray, err := ioutil.ReadAll(resp.Body)
	   if err != nil {
	       panic("Error")
	   }
	   fmt.Printf("%#v", string(byteArray))
	*/

	// req.Header.Set("Accept", "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8")
}
