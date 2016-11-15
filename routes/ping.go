package routes

import "net/http"

func Ping(w http.ResponseWriter, req *http.Request) {
	render := RendererFromContext(req.Context())
	render.JSON(w, http.StatusOK, map[string]string{})
}
