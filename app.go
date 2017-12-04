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
func (a *App) saveUserCredentials(w http.ResponseWriter, body []byte, token string) error {

	validUntil := time.Now().Add(3600*time.Second)  // User token is valid for one hour

	ui := ExtractUserData(body)

	fmt.Println("checkifemail")
	rewrite, err := checkIfEmailAlreadyInDB(a.DB, ui.Email, token, validUntil.Unix())
	if err != nil {
		return err
	}

	if rewrite {
		fmt.Println("rewrite")
		err = updateTokenInDB(a.DB, ui.Email, token, validUntil.Unix())
	} else {
		fmt.Println("save to DB")
		err = saveUserToDB(a.DB, ui.Email, token, validUntil.Unix())
	}

	if err != nil {
		return err
	}

	cookie := http.Cookie{Name: "gAppToken", Value: token, Expires: validUntil}
	http.SetCookie(w, &cookie)

	cookie = http.Cookie{Name: "gAppEmail", Value: ui.Email, Expires: validUntil}
	http.SetCookie(w, &cookie)

	return nil
}

func (a *App) lolHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectToLogin(w, r)
	}
}

func (a *App) verifiedHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenFromUser, err := r.Cookie("gAppToken")
		if err != nil {
			fmt.Println("e1")
			redirectToLogin(w, r)
			return
		}

		emailFromUser, err := r.Cookie("gAppEmail")
		if err != nil {
			fmt.Println("e2")
			redirectToLogin(w, r)
			return
		}

		accepted, err := checkIfUserIsValidated(a.DB, emailFromUser.Value, tokenFromUser.Value, 11)
		if err != nil {
			fmt.Println("e3")
			redirectToLogin(w, r)
			return
		}

		if accepted {
			fmt.Println("accepted")
			fn(w, r)
			return
		} else {
			fmt.Println("e4")
			redirectToLogin(w, r)
		}
	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/GoogleLogin", http.StatusSeeOther)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w,"Welcome to the index")
}