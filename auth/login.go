package auth

import (
	"fmt"
	log "github.com/Sirupsen/logrus"

	"github.com/satori/go.uuid"
)

func Login(username string, password string) (string, error) {
	logger := log.WithFields(log.Fields{
		"username": username,
	})

	htpasswdFile := NewHtpasswdFile("store/htpasswd")
	err := htpasswdFile.CompareUsernameAndPassword(username, password)
	if err != nil {
		logger.Info("Login attempt failed")
		return "", fmt.Errorf("invalid credentials")
	}

	logger.Info("Login attempt successful")
	token := uuid.NewV4().String()

	return token, nil
}
