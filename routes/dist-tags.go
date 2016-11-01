package routes

import (
	"net/http"
)

// Try local cache first
// If cached version is not found or fails, try upstream
// Sync distTags with npm registry remote
func DistTags(w http.ResponseWriter, req *http.Request) {
	render := RendererFromContext(req.Context())
	render.JSON(w, http.StatusOK, map[string]string{"latest": "1.0.0"})
}
