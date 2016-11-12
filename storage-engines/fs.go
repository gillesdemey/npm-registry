package storageengines

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gillesdemey/npm-registry/model"
)

type FSStorage struct {
	Folder string
}

func NewFSStorage(folder string) *FSStorage {
	engine := &FSStorage{Folder: folder}
	engine.initialize()
	return engine
}

func (s *FSStorage) initialize() {
	os.MkdirAll(s.Folder, os.ModePerm)
	return
}

func (s *FSStorage) StoreTarball(pkg string, filename string, reader io.Reader) error {
	tarballPath := filepath.Join(s.Folder, "tarballs", pkg)
	tarballLocation := filepath.Join(tarballPath, filename)

	if _, err := os.Stat(tarballLocation); err == nil {
		log.Printf("Tarball %s already exists, aborting", tarballLocation)
		return nil // file already exists
	}

	// create subdirectories we need, ignore errors
	os.MkdirAll(tarballPath, os.ModePerm)

	tarball, err := os.Create(tarballLocation)
	if err != nil {
		return err
	}
	defer tarball.Close()

	io.Copy(tarball, reader)
	return nil
}

func (s *FSStorage) RetrieveTarball(pkg string, filename string, writer io.Writer) error {
	tarballLocation := filepath.Join(s.Folder, "tarballs", pkg, filename)

	tarball, err := os.Open(tarballLocation)
	if err != nil {
		return err
	}
	defer tarball.Close()

	io.Copy(writer, tarball)
	return nil
}

func (s *FSStorage) RetrieveUsernameFromToken(token string) (string, error) {
	tokenEntries := make(map[string]model.Token)
	tokenFile := filepath.Join(s.Folder, "tokens.toml")

	if _, err := toml.DecodeFile(tokenFile, &tokenEntries); err != nil {
		return "", err
	}
	tokenEntry := tokenEntries[token]

	return tokenEntry.Username, nil
}

func (s *FSStorage) StoreUser(pkg string) error {
	return nil
}

func (s *FSStorage) RetrieveUser(string, io.Writer) error {
	return nil
}

func (s *FSStorage) StoreUserToken(token string, username string) error {
	tokenEntry := make(map[string]model.Token, 1)
	tokenEntry[token] = model.Token{
		Username:  username,
		Timestamp: time.Now(),
	}

	entry := new(bytes.Buffer)
	if err := toml.NewEncoder(entry).Encode(tokenEntry); err != nil {
		return err
	}

	tokensFile := filepath.Join(s.Folder, "tokens.toml")

	if err := ioutil.WriteFile(tokensFile, entry.Bytes(), 0666); err != nil {
		return err
	}

	return nil
}

func (s *FSStorage) RetrieveMetadata(pkg string, writer io.Writer) error {
	metaFileName := fmt.Sprintf("%s.json", pkg)
	metaFileLocation := filepath.Join(s.Folder, "meta", metaFileName)

	metaFile, err := os.Open(metaFileLocation)
	if err != nil {
		return err
	}
	defer metaFile.Close()
	io.Copy(writer, metaFile)

	return nil
}

func (s *FSStorage) StoreMetadata(pkg string, data io.Reader) error {
	metaFilePath := filepath.Join(s.Folder, "meta")
	metaFileLocation := filepath.Join(metaFilePath, pkg+".json")

	// create subdirectories we need, ignore errors
	os.MkdirAll(metaFilePath, os.ModePerm)

	metaFile, err := os.Create(metaFileLocation)
	if err != nil {
		return err
	}
	defer metaFile.Close()
	io.Copy(metaFile, data)

	return nil
}
