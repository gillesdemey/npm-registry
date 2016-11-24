package auth

import (
	"errors"
	log "github.com/Sirupsen/logrus"

	"github.com/satori/go.uuid"
)

func (p *HtpasswdProvider) Login(username string, password string) (string, error) {
	logger := log.WithFields(log.Fields{
		"username": username,
	})

	err := p.File.CompareUsernameAndPassword(username, password)
	if err != nil {
		logger.Info("Login attempt failed")
		return "", errors.New("invalid credentials")
	}

	logger.Info("Login attempt successful")
	token := uuid.NewV4().String()

	return token, nil
}
