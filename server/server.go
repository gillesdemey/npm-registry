package server

import (
	"net/http"

	"github.com/gillesdemey/npm-registry/routes"
	"github.com/gillesdemey/npm-registry/storage"
	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
)

func New(router *pat.Router, storage storage.Engine) *negroni.Negroni {
	n := negroni.Classic()
	render := render.New()

	// Attach storage and renderer on every request
	n.UseFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		var ctx = req.Context()
		ctx = context.WithValue(ctx, "storage", storage)
		ctx = context.WithValue(ctx, "renderer", render)
		next(w, req.WithContext(ctx))
	})

	// favicon requests
	router.Get("/favicon.ico", routes.Noop)

	router.Get("/-/ping", routes.Ping)

	// TODO: logout
	router.Put("/-/user/{user}", routes.Login)

	// Print the username config to standard output.
	router.Get("/-/whoami", negroni.New(
		negroni.HandlerFunc(routes.ValidateToken),
		negroni.HandlerFunc(routes.Whoami),
	).ServeHTTP)

	// tarballs
	router.Get("/{scope}/{pkg}/-/{filename}", routes.GetTarball)
	router.Get("/{pkg}/-/{filename}", routes.GetTarball)

	// packages
	router.Get("/{scope}/{pkg}", routes.GetPackageMetadata) // scoped package
	router.Get("/{pkg}", routes.GetPackageMetadata)

	// publish
	router.Put("/{scope}/{pkg}", routes.PublishPackage)
	router.Put("/{pkg}", routes.PublishPackage)

	// root
	router.Get("/", routes.Root)

	n.UseHandler(router)
	return n
}
