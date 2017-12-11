package server

import (
	"net/http"
	"log"
	"fmt"
)

type Server struct {
	Hostname string
	HTTPPort int
}


func Run(httpHandlers http.Handler, s Server) {
	startHTTP(httpHandlers, s)
}

func startHTTP(handlers http.Handler, s Server) {
	log.Fatal(http.ListenAndServe(HttpAddress(s), handlers))
}

func HttpAddress(s Server) string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.HTTPPort)
}