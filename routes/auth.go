package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gillesdemey/npm-registry/auth"
	"github.com/gillesdemey/npm-registry/model"
	"github.com/gillesdemey/npm-registry/storage"
)

// Create or verify a user named <username>
func Login(w http.ResponseWriter, req *http.Request) {
	var err error
	var login model.Login

	render := RendererFromContext(req.Context())
	auth := req.Context().Value("auth").(auth.AuthProvider)

	decoder := json.NewDecoder(req.Body)
	if err = decoder.Decode(&login); err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}
	defer req.Body.Close()

	username := login.Username
	password := login.Password

	token, err := auth.Login(username, password)
	if err != nil {
		render.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	storage := req.Context().Value("storage").(storage.UserStoreRetriever)
	err = storage.StoreUserToken(token, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, http.StatusCreated, map[string]string{"token": token})
}

// Return the username associated with the NPM token
func Whoami(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	render := RendererFromContext(req.Context())

	username := req.Context().Value("user").(string)

	render.JSON(w, http.StatusOK, map[string]string{"username": username})
}
