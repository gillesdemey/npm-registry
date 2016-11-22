package routes

import (
	"github.com/gillesdemey/npm-registry/storage"
	"github.com/urfave/negroni"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
	"net/http"
)

func Root(w http.ResponseWriter, req *http.Request) {
	render := RendererFromContext(req.Context())
	render.Text(w, http.StatusOK, "Running npm registry")
}

func Noop(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func StorageFromContext(c context.Context) storage.Engine {
	return c.Value("storage").(storage.Engine)
}

func RendererFromContext(c context.Context) *render.Render {
	return c.Value("renderer").(*render.Render)
}

// CreateMiddleware takes a variable number of negroni.HandlerFuncs and returns a single
// http.HandlerFunc to pass to gorilla.pat
func CreateMiddleware(handlers ...negroni.HandlerFunc) http.HandlerFunc {
	n := negroni.New()
	for _, handlerFunc := range handlers {
		n.Use(handlerFunc)
	}
	return n.ServeHTTP
}
