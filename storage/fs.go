package storage

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/gillesdemey/npm-registry/model"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type FSStorage struct {
	StorageEngine
	folder string
}

func NewFSStorage() *FSStorage {
	engine := &FSStorage{}
	engine.initialize()
	return engine
}

func (s *FSStorage) initialize() error {
	return nil
}

func (s *FSStorage) StoreTarball() error {
	return nil
}

func (s *FSStorage) RetrieveTarball() ([]byte, error) {
	return nil, nil
}

func (s *FSStorage) RetrieveUser() (model.User, error) {
	return model.User{}, nil
}

func (s *FSStorage) StoreUserToken(token string, username string) error {
	tokenEntry := make(map[string]model.Token, 1)
	tokenEntry[token] = model.Token{
		Username:  username,
		Timestamp: time.Now(),
	}

	entry := new(bytes.Buffer)
	if err := toml.NewEncoder(entry).Encode(tokenEntry); err != nil {
		log.Fatal(err)
	}

	tokensFile := filepath.Join(s.folder, "tokens.toml")

	if err := ioutil.WriteFile(tokensFile, entry.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}

	return nil
}
