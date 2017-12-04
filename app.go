package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"time"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}


func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	//r := mux.NewRouter()
	a.Router.HandleFunc("/GoogleLogin", handleGoogleLogin)
	a.Router.HandleFunc("/GoogleCallback", a.handleGoogleCallback)
	a.Router.HandleFunc("/index", a.verifiedHandler(index))
}


// saveUserCredentials saves the user email, token, and duration of validity both to
// the DB and as a cookie. This is matched for user validation.
func (a *App) saveUserCredentials(w http.ResponseWriter, body []byte, token string) {

	validUntil := time.Now().Add(3600*time.Second)  // User token is valid for one hour

	ui := ExtractUserData(body)

	rewrite, err := checkIfEmailAlreadyInDB(a.DB, ui.Email, token, validUntil.Unix())
	if err != nil {
		log.Fatal(err)
	}

	if rewrite {
		updateTokenInDB(a.DB, ui.Email, token, validUntil.Unix())
	} else {
		saveUserToDB(a.DB, ui.Email, token, validUntil.Unix())
	}

	cookie := http.Cookie{Name: "gAppToken", Value: token, Expires: validUntil}
	http.SetCookie(w, &cookie)

	cookie = http.Cookie{Name: "gAppEmail", Value: ui.Email, Expires: validUntil}
	http.SetCookie(w, &cookie)

}

func (a *App) verifiedHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenFromUser, err := r.Cookie("gAppToken")
		if err != nil {
			return
		}

		emailFromUser, err := r.Cookie("gAppEmail")
		if err != nil {
			return
		}

		err = checkIfUserIsValidated(a.DB, emailFromUser.Value, tokenFromUser.Value, 11)
		if err != nil {
			return
		}

		fn(w, r)
	}
}


func index(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("gAppToken")
	fmt.Fprint(w, cookie)

	cookie, _ = r.Cookie("gAppEmail")
	fmt.Fprint(w, cookie)
}