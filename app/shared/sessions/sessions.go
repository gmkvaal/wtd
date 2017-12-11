package main

import (
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

var Store *sessions.CookieStore

type Session struct {
	Name string
}

func ConfigureSession(s Session) {
	Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	sessions.NewSession(Store, s.Name)

}
