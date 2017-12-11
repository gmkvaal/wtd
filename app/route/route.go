package route

import (
	"github.com/gorilla/mux"

	"github.com/gmkvaal/wtd/app/route/httpwrappers"
	"github.com/gmkvaal/wtd/app/controller"
)

func Routes () *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", httpwrappers.Handler(controller.HandleOAuth2Login))
	router.HandleFunc("/oauthcallback", httpwrappers.Handler(controller.HandleOAuth2Callback))
	router.HandleFunc("/index", httpwrappers.AuthHandler(controller.Index))

	return router
}

