package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
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
