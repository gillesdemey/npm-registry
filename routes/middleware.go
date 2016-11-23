package routes

import (
  "golang.org/x/net/context"
  "net/http"
  "regexp"
	"github.com/gillesdemey/npm-registry/storage"

  log "github.com/Sirupsen/logrus"
)

func ValidateToken(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	storage := req.Context().Value("storage").(storage.TokenRetriever)

  re := regexp.MustCompile("(?i)Bearer ")
	authHeader := req.Header.Get("Authorization")
	token := re.ReplaceAllString(authHeader, "")

	username, err := storage.RetrieveUsernameFromToken(token)
  if err != nil {
		log.WithFields(log.Fields{
			"token": token,
		}).Info("token validation failed: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

  ctx := req.Context()
  ctx = context.WithValue(ctx, "user", username)
  ctx = context.WithValue(ctx, "token", token)

  next(w, req.WithContext(ctx))
}
