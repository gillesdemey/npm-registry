package main

import (
	"github.com/gillesdemey/npm-registry/server"
	"github.com/gillesdemey/npm-registry/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	router := gin.Default()
	storage := storage.NewFSStorage()

	server := server.New(router, storage)
	server.Run(":8080")
}
