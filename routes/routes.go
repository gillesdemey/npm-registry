package routes

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func Root(c *gin.Context) {
	c.String(http.StatusOK, "Running npm registry")
}

// Ping the configured or given npm registry and verify authentication.
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
