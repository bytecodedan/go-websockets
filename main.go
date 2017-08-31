package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log *logger

type logger struct {
	*logrus.Logger
}

func main() {
	address := flag.String("address", "localhost:8080", "host and port")
	flag.Parse()

	initWebSocket()

	mux := mux.NewRouter()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.Handle("/ws", HandleWSConnections).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, mux)
	logrus.Infof("HTTP server started on %v", *address)
	err := http.ListenAndServe(*address, loggedRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
