package server

import (
	"github.com/gillesdemey/npm-registry/routes"
	"github.com/gillesdemey/npm-registry/storage"
	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
	"net/http"
)

func New(router *pat.Router, storage storage.StorageEngine) *negroni.Negroni {
	n := negroni.Classic()
	render := render.New()

	// Attach storage and renderer on every request
	n.Use(negroni.HandlerFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		var ctx = req.Context()
		ctx = context.WithValue(ctx, "storage", storage)
		ctx = context.WithValue(ctx, "renderer", render)
		next(w, req.WithContext(ctx))
	}))

	router.Get("/-/ping", routes.Ping)

	// TODO: logout
	router.Put("/-/user/{user}", routes.Login)

	// Print the username config to standard output.
	router.Get("/-/whoami", routes.Whoami)

	// dist-tags
	router.Get("/-/package/{name}/dist-tags", routes.DistTags)

	router.Put("/-/package/{name}/dist-tags/:tag", func(w http.ResponseWriter, r *http.Request) {})

	router.Delete("/-/package/{name}/dist-tags/:tag", func(w http.ResponseWriter, r *http.Request) {})

	// packages
	router.Get("/{pkg}", routes.GetPackageMetadata)

	// tarballs
	router.Get("/{pkg}/-/:filename", func(w http.ResponseWriter, r *http.Request) {})

	// publish
	router.Put("/{pkg}", func(w http.ResponseWriter, r *http.Request) {})

	// root
	router.Get("/", routes.Root)

	n.UseHandler(router)
	return n
}
