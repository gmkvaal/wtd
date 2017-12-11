package controller

import (
	"time"
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/gmkvaal/wtd/app/shared/encoding"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	)

/*
var conf = &oauth2.Config{
	RedirectURL:    "http://127.0.0.1:9090/GoogleCallback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}
*/

type OAuthConfig struct {
	Google Google
}

type Google struct {
	RedirectURL string
	ClientID string
	ClientSecret string
	Scopes string
}

var conf = &oauth2.Config{
}

func ConfigureOAuth(gConfig Google) {
	conf = &oauth2.Config{
		RedirectURL:    gConfig.RedirectURL,
		ClientID:     gConfig.ClientID,
		ClientSecret: gConfig.ClientSecret,
		Scopes:       []string{gConfig.Scopes},
		Endpoint:     google.Endpoint,
	}
}

// handleGoogleLogin generates a random oAuth code and redirects to Google login.
// The oAuth code is saved as a cookie which is read in handleGoogleCallback
func HandleOAuth2Login(w http.ResponseWriter, r *http.Request) {

	fmt.Println(conf)

	oAuthState := encoding.RandToken()
	cookie := http.Cookie{Name: "oAuthState", Value: oAuthState, Expires: time.Now().Add(1 * time.Minute)}
	http.SetCookie(w, &cookie)

	url := conf.AuthCodeURL(oAuthState) //, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback first verifies that the user is authenticating
// with the random code generated in handleGoogleLogin
func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	stateFromCookie, err := r.Cookie("oAuthState")
	if err != nil {
		log.Fatal(err)
	}

	if state != string(stateFromCookie.Value) {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", stateFromCookie, state)
		http.Redirect(w, r, "/Login", http.StatusTemporaryRedirect)
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

	_, err  = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	//err = a.saveUserCredentials(w, body, token.AccessToken)
	//if err != nil {
	//	log.Fatal(err)
	//}

	fmt.Println("save user cred ok")

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}