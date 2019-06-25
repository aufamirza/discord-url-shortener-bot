package server

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"os"
)

func Start(stop chan os.Signal) {
	const port = "8080"

	router := httprouter.New()

	//route for health check
	router.GET("/", Health)

	//route for URL redirects
	router.GET("/:id", ReturnShortURL)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(fmt.Sprintf("error: %v", err))
	}

	//start in goroutine so we can wait for OS signal
	go func() {
		err := http.Serve(listener, router)
		if err != nil {
			log.Fatal(fmt.Sprintf("error: %v", err))
		}
	}()

	log.Println(fmt.Sprintf("started web server on port %v", port))

	//block until OS signal
	<-stop
}

//handle healthcheck
func Health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

//handle server URL redirect
func ReturnShortURL(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	const location = "https://github.com/fraserdarwent"

	w.Header().Add("Location", location)
	w.WriteHeader(http.StatusMovedPermanently)
}
