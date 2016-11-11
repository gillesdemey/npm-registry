package storage

import (
	"io"
)

type TarballStoreRetriever interface {
	StoreTarball() error
	RetrieveTarball(pkg string, writer io.Writer) error
}

type MetaDataStoreRetriever interface {
	StoreMetadata(pkg string, data io.Reader) error
	RetrieveMetadata(pkg string, writer io.Writer) error
}

type UserStoreRetriever interface {
	StoreUser(pkg string) error
	StoreUserToken(token, username string) error
	RetrieveUser(token string, writer io.Writer) error
}

type Engine interface {
	TarballStoreRetriever
	MetaDataStoreRetriever
	UserStoreRetriever
	RetrieveUsernameFromToken(token string) (string, error) // Do we realy need this?
}
