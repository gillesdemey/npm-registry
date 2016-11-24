package server

import (
	"net/http"

	"github.com/gillesdemey/npm-registry/auth"
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
	auth := auth.NewHtpasswdProvider("store/htpasswd")

	// Attach storage and renderer on every request
	n.UseFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		var ctx = req.Context()
		ctx = context.WithValue(ctx, "storage", storage)
		ctx = context.WithValue(ctx, "renderer", render)
		ctx = context.WithValue(ctx, "auth", auth)
		next(w, req.WithContext(ctx))
	})

	// favicon requests
	router.Get("/favicon.ico", routes.Noop)

	router.Get("/-/ping", routes.Ping)

	// TODO: logout
	router.Put("/-/user/{user}", routes.Login)

	// Print the username config to standard output.
	router.Get("/-/whoami", routes.CreateMiddleware(
		routes.ValidateToken,
		routes.Whoami,
	))

	// tarballs
	tarballHandlerFunc := routes.CreateMiddleware(
		routes.ValidateToken,
		routes.GetTarball,
	)
	router.Get("/{scope}/{pkg}/-/{filename}", tarballHandlerFunc)
	router.Get("/{pkg}/-/{filename}", tarballHandlerFunc)

	// packages
	pkgMetaHandlerFunc := routes.CreateMiddleware(
		routes.ValidateToken,
		routes.GetPackageMetadata,
	)
	router.Get("/{scope}/{pkg}", pkgMetaHandlerFunc) // scoped package
	router.Get("/{pkg}", pkgMetaHandlerFunc)

	// publish
	publishPackageHandlerFunc := routes.CreateMiddleware(
		routes.ValidateToken,
		routes.PublishPackage,
	)
	router.Put("/{scope}/{pkg}", publishPackageHandlerFunc)
	router.Put("/{pkg}", publishPackageHandlerFunc)

	// root
	router.Get("/", routes.Root)

	n.UseHandler(router)
	return n
}
