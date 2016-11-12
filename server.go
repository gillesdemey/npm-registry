package main

import (
	"net/http"

	"github.com/gillesdemey/npm-registry/server"
	"github.com/gillesdemey/npm-registry/storage-engines"
	"github.com/gorilla/pat"
)

func main() {
	router := pat.New()
	storage := storageengines.NewFSStorage("registry_data/")

	server := server.New(router, storage)
	http.ListenAndServe(":8080", server)
}
