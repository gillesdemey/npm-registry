package routes

import (
	"github.com/gillesdemey/npm-registry/storage"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
	"net/http"
)

func Root(w http.ResponseWriter, req *http.Request) {
	render := RendererFromContext(req.Context())
	render.Text(w, http.StatusOK, "Running npm registry")
}

// Ping the configured or given npm registry and verify authentication.
func Ping(w http.ResponseWriter, req *http.Request) {
	render := RendererFromContext(req.Context())
	render.JSON(w, http.StatusOK, map[string]string{})
}

func StorageFromContext(c context.Context) storage.Engine {
	return c.Value("storage").(storage.Engine)
}

func RendererFromContext(c context.Context) *render.Render {
	return c.Value("renderer").(*render.Render)
}
