package server

import (
	"discord-url-shortener-bot/persistence"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"os"
)

var URLStore persistence.URLStore
var port string

func Start(stop chan os.Signal, newPort string, newURLStore persistence.URLStore) {
	//configure port
	port = newPort

	//configure URLStore
	URLStore = newURLStore

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
	id := p.ByName("id")

	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	URL := URLStore.Get(id)
	if URL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Location", URL)
	w.WriteHeader(http.StatusMovedPermanently)
}
