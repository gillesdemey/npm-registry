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

	logger.Info("Login attempt")
	if username != "foo" || password != "bar" {
		logger.Info("Login attempt failed")
		return "", fmt.Errorf("invalid credentials")
	}

	logger.Info("Login successful")
	token := uuid.NewV4().String()

	return token, nil
}
