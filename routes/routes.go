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

func StorageFromContext(c context.Context) storage.Engine {
	return c.Value("storage").(storage.Engine)
}

func RendererFromContext(c context.Context) *render.Render {
	return c.Value("renderer").(*render.Render)
}
