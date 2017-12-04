package main

import (
	"fmt"
	"log"
	"net/http"
	"golang.org/x/oauth2"
	"os"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"crypto/rand"
	"encoding/base64"
	"time"
	"io/ioutil"
)



var conf = &oauth2.Config{
		RedirectURL:    "http://127.0.0.1:9090/GoogleCallback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
		}


func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// handleGoogleLogin generates a random oAuth code and redirects to Google login.
// The oAuth code is saved as a cookie which is read in handleGoogleCallback
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {

	oAuthState := randToken()
	cookie := http.Cookie{Name: "oAuthState", Value: oAuthState, Expires: time.Now().Add(1 * time.Minute)}
	http.SetCookie(w, &cookie)

	url := conf.AuthCodeURL(oAuthState) //, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback first verifies that the user is authenticating with the random code
// generated in handleGoogleLogin
func (a *App) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	stateFromCookie, err := r.Cookie("oAuthState")
	if err != nil {
		log.Fatal(err)
	}

	if state != string(stateFromCookie.Value) {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", stateFromCookie, state)
		http.Redirect(w, r, "/GoogleLogin", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err  := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = a.saveUserCredentials(w, body, token.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}



