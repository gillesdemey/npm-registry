package storage

import (
	"time"
)

type User struct {
	Username, Password, Email, Token string
}

type TokenEntry struct {
	Token TokenInfo
}

type TokenInfo struct {
	Username  string
	Token     string
	Timestamp time.Time
}

type StorageEngine interface {
	initialize() error
	StoreTarball() error
	RetrieveTarball() ([]byte, error)
	RetrieveUser() (User, error)
	StoreUserToken(username string, token string) error
}
