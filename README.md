# go-cloudfunction-auth

## Introduction

`go-cloudfunction-auth` contains a client implementation for Google Cloud Function following Google's Oauth2 spec at https://developers.google.com/identity/protocols/oauth2/service-account

Google provides their own official [oauth2](https://godoc.org/golang.org/x/oauth2) and higher level [Cloud API](https://github.com/googleapis/google-cloud-go), which include authentication part for most of all services. Unfortunately, these libraries does not support well authenticating with cloud function. There are also not so many useful documents or articles relate to this topic. Investigating on this take me much longer than my expectation. Therefore, I decide to publish a library here in case other people are also stuck at this step.

## Reference

The implementation inside this library follows Google's [official developers guideline](https://developers.google.com/identity/protocols/oauth2/service-account):

![Google OAuth2 Flow](https://developers.google.com/accounts/images/serviceaccount.png)

- Create and sign JWT using a custom wrapper of [google oauth2](https://github.com/golang/oauth2/tree/master/google)
- Send the signed JWT to Google's server to get an authenticated token
- Use the received token to create a custom http Client. This client attachs the token in all requests by default.

## Sample

See the full sample at https://github.com/CodeLinkIO/go-cloudfunction-auth/blob/master/samples/main.go

```go
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
```

## License

MIT