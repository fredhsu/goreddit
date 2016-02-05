package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/oauth2"
)

// Token acquired from reddit
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func oauthToken(id, secret, redirect string) {
	conf := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Scopes:       []string{"read"},
		RedirectURL:  redirect,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.reddit.com/api/v1/authorize",
			TokenURL: "https://www.reddit.com",
		},
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state")
	fmt.Printf("Visit the URL for the auth dialog:\n %v", url)
	// Need to make HTTP POST to get the token https://www.reddit.com/api/v1/access_token
	accessTokenURL := "https://www.reddit.com/api/v1/access_token"
	client := &http.Client{}
	//resp, err := client.Post(accessTokenUrl, "application/json", )
	code := "ZCcT2wlrSdUxFMellmJixAU5GcU"
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, redirect)
	body := []byte(data)
	req, err := http.NewRequest("POST", accessTokenURL, bytes.NewBuffer(body))
	req.Header.Set("User-Agent", "Golang")
	if err != nil {
		log.Fatal(err)
	}
	// ...
	req.SetBasicAuth(id, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
}

func scriptToken(id, secret, username, password string) []byte {
	client := &http.Client{}
	data := fmt.Sprintf("grant_type=password&username=%s&password=%s", username, password)
	body := []byte(data)
	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", bytes.NewBuffer(body))
	req.Header.Set("User-Agent", "Golang")
	if err != nil {
		log.Fatal(err)
	}
	// ...
	req.SetBasicAuth(id, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	return responseBody

}

func apiCall(url string, token Token) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	authstring := fmt.Sprintf("%s %s", token.TokenType, token.AccessToken)
	req.Header.Set("Authorization", authstring)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseBody
}

func main() {
	// TODO: Grab id, secret, username, password from ENV
	id := "Zw_HDcw9YrkU_Q"
	secret := "sBwII_4zfoIjzIS4C2TlxZpN9f4"
	//redirect := "http://www.example.com/unused/redirect/uri"
	username := "fredlhsu"
	password := "thegames"
	b := scriptToken(id, secret, username, password)
	token := Token{}
	err := json.Unmarshal(b, &token)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("token: ", token)
	apiResponse := apiCall("https://oauth.reddit.com/api/v1/me", token)
	fmt.Println(string(apiResponse))
	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.
	/*
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Fatal(err)
		}

		//client := conf.Client(oauth2.NoContext, tok)
		//client.Get("")
	*/
}
