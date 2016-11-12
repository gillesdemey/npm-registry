package routes

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"regexp"

	"github.com/gillesdemey/npm-registry/auth"
	"github.com/gillesdemey/npm-registry/model"
)

// Create or verify a user named <username>
func Login(w http.ResponseWriter, req *http.Request) {
	var err error
	var login model.Login

	render := RendererFromContext(req.Context())
	storage := StorageFromContext(req.Context())

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

	err = storage.StoreUserToken(token, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, http.StatusCreated, map[string]string{"token": token})
}

// Return the username associated with the NPM token
func Whoami(w http.ResponseWriter, req *http.Request) {
	render := RendererFromContext(req.Context())
	storage := StorageFromContext(req.Context())

	re := regexp.MustCompile("(?i)Bearer ")
	authHeader := req.Header.Get("Authorization")
	token := re.ReplaceAllString(authHeader, "")

	log.Info("Whoami request")

	username, err := storage.RetrieveUsernameFromToken(token)

	if err != nil {
		log.WithFields(log.Fields{
			"token": token,
		}).Info("Whoami request failed: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"username": username})
}
