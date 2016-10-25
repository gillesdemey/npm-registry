package storage

import (
  "io/ioutil"
  "path/filepath"
  "log"
)

type FSStorage struct {
  StorageEngine
  folder string
}

func NewFSStorage () *FSStorage {
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

func (s *FSStorage) RetrieveUser() (User, error) {
  return User{}, nil
}

func (s *FSStorage) StoreUserToken(token string, username string) error {
  tokensFile := filepath.Join(s.folder, "tokens")
  log.Printf("Writing to file %s", tokensFile)

	if err := ioutil.WriteFile(tokensFile, []byte(token), 0666); err != nil {
		log.Fatal(err)
	}

  return nil
}
