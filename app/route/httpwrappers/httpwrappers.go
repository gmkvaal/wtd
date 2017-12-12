package httpwrappers

import (
	"net/http"

	"github.com/gmkvaal/wtd/app/controller"
	"github.com/gmkvaal/wtd/app/shared/sessions"
)

func Handler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func AuthHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		accepted := sessions.GetSessionValue(r, "authenticated")

		if accepted.(bool) {
			fn(w, r)
			return
		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
	}
}
