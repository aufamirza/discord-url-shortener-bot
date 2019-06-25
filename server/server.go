package server

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
)

func Start() {
	const port = "8080"

	router := httprouter.New()
	router.GET("/health", Health)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(fmt.Sprintf("error: %v", err))
	}

	log.Println("started web server")

	err = http.Serve(listener, router)
	if err != nil {
		log.Fatal(fmt.Sprintf("error: %v", err))
	}
}

func Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
