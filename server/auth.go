package server

import (
  "github.com/satori/go.uuid"
)

func (r *NPMRegistry) Login (username string, password string) (string, error) {
  token := uuid.NewV4().String()
  err := r.Storage.StoreUserToken(token, username)

  if err != nil {
    return "", err
  }

  return token, nil
}

func (r *NPMRegistry) Logout (username string, password string) error {
  return nil
}
