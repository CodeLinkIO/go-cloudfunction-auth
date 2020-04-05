package cloudfunction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CodeLinkIO/go-cloudfunction-auth/internal"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const GOOGLE_TOKEN_URL = "https://oauth2.googleapis.com/token"

func JWTAccessTokenSourceFromJSON(jsonKey []byte, audience string) (oauth2.TokenSource, error) {
	cfg, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, fmt.Errorf("google: could not parse JSON key: %v", err)
	}
	pk, err := internal.ParseKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("google: could not parse key: %v", err)
	}
	ts := &jwtAccessTokenSource{
		email:    cfg.Email,
		audience: audience,
		pk:       pk,
		pkID:     cfg.PrivateKeyID,
	}
	tok, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return oauth2.ReuseTokenSource(tok, ts), nil
}

type TokenResponse struct {
	IdToken string `json:"id_token"`
}

func Authenticate(tokenSource oauth2.TokenSource) (token oauth2.Token, err error) {
	jwt, err := tokenSource.Token()
	if err != nil {
		return
	}

	client := &http.Client{Timeout: time.Second * 10}
	payload := strings.NewReader("grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer&assertion=" + jwt.AccessToken)
	req, _ := http.NewRequest("POST", GOOGLE_TOKEN_URL, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	tokenRes := &TokenResponse{}
	err = json.Unmarshal(body, tokenRes)
	if err != nil {
		fmt.Println(err.Error())
	}
	token = oauth2.Token{
		AccessToken: tokenRes.IdToken,
	}
	return
}

func NewClient(jwtSource oauth2.TokenSource) *http.Client {
	token, err := Authenticate(jwtSource)
	if err != nil {
		fmt.Printf("cannot authenticate with google: %v", err)
		os.Exit(1)
	}

	return &http.Client{
		Transport: &oauth2.Transport{
			Base: http.DefaultClient.Transport,
			Source: &googleTokenSource{
				GoogleToken: &token,
			},
		},
	}
}
