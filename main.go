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

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }

func main() {
	address := flag.String("address", "localhost:8080", "host and port")
	flag.Parse()

	log = &logger{logrus.New()}
	InitWS()

	mux := mux.NewRouter()
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.Handle("/ws", HandleWSConnections).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, mux)
	log.Infof("HTTP server started on %v", *address)
	err := http.ListenAndServe(*address, loggedRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
