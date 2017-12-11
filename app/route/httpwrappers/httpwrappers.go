package httpwrappers

import (
	"net/http"
	"github.com/gmkvaal/wtd/app/model"
	"github.com/gmkvaal/wtd/app/controller"
)



func Handler (fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w,r)
	}
}

func AuthHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenFromUser, err := r.Cookie("gAppToken")
		if err != nil {
			controller.RedirectToLogin(w, r)
			return
		}

		emailFromUser, err := r.Cookie("gAppEmail")
		if err != nil {
			controller.RedirectToLogin(w, r)
			return
		}

		accepted, err := model.CheckIfUserIsValidated(emailFromUser.Value, tokenFromUser.Value, 11)
		if err != nil {
			controller.RedirectToLogin(w, r)
			return
		}

		if accepted {
			fn(w, r)
			return
		} else {
			controller.RedirectToLogin(w, r)
		}
	}
}