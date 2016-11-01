package routes

import (
  "gopkg.in/gin-gonic/gin.v1"
  "net/http"
)

// Try local cache first
// If cached version is not found or fails, try upstream
// Sync distTags with npm registry remote
func DistTags(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"latest": "1.0.0"})
}
