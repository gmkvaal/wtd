package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"net/http"
)

var (
	Store *sessions.CookieStore
	Name string
)

type Session struct {
	Name string
	Options sessions.Options
}

func ConfigureSession(s Session) {
	Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	Store.Options = &s.Options
	Name = s.Name
}

func SaveSession(w http.ResponseWriter, r *http.Request, value interface{}, valueName string) {
	sess, _ := Store.Get(r, Name)
	sess.Values[valueName] = value
	sess.Save(r, w)
}

func GetSessionValue(r *http.Request, valueName string) interface{}  {
	sess, _ := Store.Get(r, Name)
	return sess.Values[valueName]
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	sess, _ := Store.Get(r, Name)

	if sess.Values["authenticated"] != nil {
		delete(sess.Values, "authenticated")
		sess.Save(r, w)
	}

}