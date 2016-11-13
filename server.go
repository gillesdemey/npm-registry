package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gillesdemey/npm-registry/server"
	engines "github.com/gillesdemey/npm-registry/storage-engines"
	"github.com/gorilla/pat"
)

func main() {
	router := pat.New()
	storage := engines.NewFSStorage("store/")

	server := server.New(router, storage)

	log.Info("Listening on http://0.0.0.0:8080/")
	http.ListenAndServe(":8080", server)
}
