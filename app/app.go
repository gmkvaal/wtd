package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"

	"log"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}


func (a *App) Initialize(user, password, dbname string) {




	a.Router = mux.NewRouter()

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

/*

// saveUserCredentials saves the user email, token, and duration of validity both to
// the DB and as a cookie. This is matched for user validation.
func (a *App) saveUserCredentials(w http.ResponseWriter, body []byte, token string) error {

	validUntil := time.Now().Add(3600*time.Second)  // User token is valid for one hour

	ui := ExtractUserData(body)

	rewrite, err := model.CheckIfEmailAlreadyInDB(ui.Email, token, validUntil.Unix())
	if err != nil {
		return err
	}

	if rewrite {
		fmt.Println("rewrite")
		err = model.UpdateTokenInDB(ui.Email, token, validUntil.Unix())
	} else {
		fmt.Println("save to DB")
		err = model.SaveUserToDB(ui.Email, token, validUntil.Unix())
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





*/