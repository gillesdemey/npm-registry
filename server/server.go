package server

import (
	"github.com/gillesdemey/npm-registry/routes"
	"github.com/gillesdemey/npm-registry/storage"
	"gopkg.in/gin-gonic/gin.v1"
)

func New(router *gin.Engine, storage storage.StorageEngine) *gin.Engine {
	router.GET("/", routes.Root)
	router.GET("/-/ping", routes.Ping)

	// TODO: logout
	router.PUT("/-/user/:user", func(c *gin.Context) {
		c.Set("storage", storage)
		routes.Login(c)
	})

	// Print the username config to standard output.
	router.GET("/-/whoami", func(c *gin.Context) {
		c.Set("storage", storage)
		routes.Whoami(c)
	})

	// dist-tags
	router.GET("/-/package/:name/dist-tags", routes.DistTags)

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

	return router
}
