package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CodeLinkIO/go-cloudfunction-auth/cloudfunction"

	"golang.org/x/oauth2/google"
)

func main() {
	baseUrl := "your-cloudfunction-baseurl"
	ctx := context.Background()
	targetAudience := baseUrl
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		fmt.Printf("cannot get credentials: %v", err)
		os.Exit(1)
	}

	jwtSource, err := cloudfunction.JWTAccessTokenSourceFromJSON(credentials.JSON, targetAudience)
	if err != nil {
		fmt.Printf("cannot create jwt source: %v", err)
		os.Exit(1)
	}

	client := cloudfunction.NewClient(jwtSource)
	res, err := client.Get(baseUrl + "/cloudfunction-sub-page")
	if err != nil {
		fmt.Printf("cannot fetch result: %v", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("cannot read response: %v", err)
		os.Exit(1)
	}
	println(string(body))
}
