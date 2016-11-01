package auth

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
)

func Login(username string, password string) (string, error) {
	log.Printf("Login attempt for %s@%s", username, password)
	if username != "foo" || password != "bar" {
		return "", fmt.Errorf("invalid credentials")
	}

	log.Printf("Login attempt valid for %s@%s", username, password)
	token := uuid.NewV4().String()

	return token, nil
}
