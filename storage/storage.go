package storage

import (
	"io"

	"github.com/gillesdemey/npm-registry/model"
)

type Engine interface {
	initialize() error
	StoreTarball() error
	RetrieveTarball() ([]byte, error)
	StoreMetadata(pkg string) error
	RetrieveMetadata(pkg string, writer io.Writer) error
	RetrieveUser(token string) (model.User, error)
	StoreUserToken(username string, token string) error
	RetrieveUsernameFromToken(token string) (string, error)
}
