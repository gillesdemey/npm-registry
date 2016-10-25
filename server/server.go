package server

import (
	"github.com/gillesdemey/npm-registry/storage"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type NPMRegistry struct {
	Router  *gin.Engine
	Storage storage.StorageEngine
}

func New(router *gin.Engine, storage storage.StorageEngine) *NPMRegistry {
	registry := &NPMRegistry{
		Router:  router,
		Storage: storage,
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Running npm registry")
	})

	// Ping the configured or given npm registry and verify authentication.
	router.GET("/-/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	// login
	// TODO: logout
	router.PUT("/-/user/:user", func(c *gin.Context) {
		var json Login

		if err := c.BindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		}

		token, err := registry.Login(json.Username, json.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		} else {
			c.JSON(http.StatusCreated, gin.H{"token": token})
		}
	})

	// Print the username config to standard output.
	router.GET("/-/whoami", func(c *gin.Context) {
		c.String(http.StatusOK, "gilles")
	})

	// dist-tags
	router.GET("/-/package/:name/dist-tags", func(c *gin.Context) {

	})

	router.PUT("/-/package/:name/dist-tags/:tag", func(c *gin.Context) {

	})

	router.DELETE("/-/package/:name/dist-tags/:tag", func(c *gin.Context) {

	})

	// packages
	// router.GET("/:pkg", func(c *gin.Context) {})

	// tarballs
	// router.GET("/:pkg/-/:filename", func (c *gin.Context) {})

	// publish
	// router.PUT("/:pkg", func (c *gin.Context) {})

	return registry
}

func (r *NPMRegistry) Run(port string) {
	r.Router.Run(port)
}
