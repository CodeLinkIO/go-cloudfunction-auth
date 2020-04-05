package cloudfunction

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/CodeLinkIO/go-cloudfunction-auth/jws"
	"golang.org/x/oauth2"
)

type jwtAccessTokenSource struct {
	email, audience string
	pk              *rsa.PrivateKey
	pkID            string
}

func (ts *jwtAccessTokenSource) Token() (*oauth2.Token, error) {
	iat := time.Now()
	exp := iat.Add(time.Hour)
	cs := &jws.ClaimSet{
		Iss: ts.email,
		Aud: GOOGLE_TOKEN_URL,
		Iat: iat.Unix(),
		Exp: exp.Unix(),

		PrivateClaims: map[string]interface{}{
			"target_audience": ts.audience,
		},
	}
	hdr := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
		KeyID:     string(ts.pkID),
	}
	msg, err := jws.Encode(hdr, cs, ts.pk)
	if err != nil {
		return nil, fmt.Errorf("google: could not encode JWT: %v", err)
	}
	return &oauth2.Token{AccessToken: msg, TokenType: "Bearer", Expiry: exp}, nil
}

type googleTokenSource struct {
	GoogleToken *oauth2.Token
}

func (ts *googleTokenSource) Token() (*oauth2.Token, error) {
	return ts.GoogleToken, nil
}
