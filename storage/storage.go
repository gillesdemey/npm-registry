package storage

import (
	"github.com/gillesdemey/npm-registry/model"
)

type StorageEngine interface {
	initialize() error
	StoreTarball() error
	RetrieveTarball() ([]byte, error)
	RetrieveUser(token string) (model.User, error)
	StoreUserToken(username string, token string) error
	RetrieveUsernameFromToken(token string) (string, error)
}
