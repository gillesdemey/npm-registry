package main

import (
	"github.com/gillesdemey/npm-registry/server"
	"github.com/gillesdemey/npm-registry/storage"
	"github.com/gorilla/pat"
	"net/http"
)

func main() {
	router := pat.New()
	storage := storage.NewFSStorage()

	server := server.New(router, storage)
	http.ListenAndServe(":8080", server)
}
