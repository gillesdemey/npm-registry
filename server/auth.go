package server

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
)

type Login struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *NPMRegistry) Login(username string, password string) (string, error) {
	log.Printf("Logging in with %s %s", username, password)
	if username != "foo" || password != "bar" {
		return "", fmt.Errorf("invalid credentials")
	}

	token := uuid.NewV4().String()
	err := r.Storage.StoreUserToken(token, username)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *NPMRegistry) Logout(username string, password string) error {
	return nil
}
