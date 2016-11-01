package routes

import (
  "gopkg.in/gin-gonic/gin.v1"
  "net/http"
)

func Root(c *gin.Context) {
  c.String(http.StatusOK, "Running npm registry")
}

func Ping(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{})
}
