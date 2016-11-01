package storage

import (
	"github.com/gillesdemey/npm-registry/model"
)

type StorageEngine interface {
	initialize() error
	StoreTarball() error
	RetrieveTarball() ([]byte, error)
	RetrieveUser() (model.User, error)
	StoreUserToken(username string, token string) error
}
