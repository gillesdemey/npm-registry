package storage

import "io"

type TarballStoreRetriever interface {
	StoreTarball(filename string, reader io.Reader) error
	RetrieveTarball(pkg string, filename string, writer io.Writer) error
}

type MetaDataStoreRetriever interface {
	StoreMetadata(pkg string, data io.Reader) error
	RetrieveMetadata(pkg string, writer io.Writer) error
}

type UserStoreRetriever interface {
	StoreUserToken(token, username string) error
}

type TokenRetriever interface {
	RetrieveUsernameFromToken(token string) (string, error)
}

type Engine interface {
	TarballStoreRetriever
	MetaDataStoreRetriever
	UserStoreRetriever
	TokenRetriever
}
