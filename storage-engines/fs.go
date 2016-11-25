package storageengines

import (
	"bytes"
	"errors"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gillesdemey/npm-registry/model"
	"github.com/gillesdemey/npm-registry/packages"
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

func (s *FSStorage) StoreTarball(filename string, reader io.Reader) error {
	scope, filename := packages.SplitPackageName(filename)
	tarballPath := filepath.Join(s.Folder, "tarballs")

	if scope != "" {
		tarballPath = filepath.Join(tarballPath, scope)
	}

	tarballLocation := filepath.Join(tarballPath, filename)

	if _, err := os.Stat(tarballLocation); err == nil {
		log.WithFields(log.Fields{
			"tarball": filename,
		}).Info("Tarball already exists, skipping")
		return nil
	}

	log.WithFields(log.Fields{"path": tarballLocation}).Info("Storing tarball")

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
	scope, _ := packages.SplitPackageName(pkg)
	tarballPath := filepath.Join(s.Folder, "tarballs")

	if scope != "" {
		tarballPath = filepath.Join(tarballPath, scope)
	}

	tarballLocation := filepath.Join(tarballPath, filename)

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
	if tokenEntry.Username == "" {
		return "", errors.New("no user found for token")
	}

	return tokenEntry.Username, nil
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
	scope, pkgName := packages.SplitPackageName(pkg)
	metaFilePath := filepath.Join(s.Folder, "meta")

	if scope != "" {
		metaFilePath = filepath.Join(metaFilePath, scope)
	}

	metaFileLocation := filepath.Join(metaFilePath, pkgName+".json")

	metaFile, err := os.Open(metaFileLocation)
	if err != nil {
		return err
	}
	defer metaFile.Close()
	io.Copy(writer, metaFile)

	return nil
}

func (s *FSStorage) StoreMetadata(pkg string, data io.Reader) error {
	scope, pkgName := packages.SplitPackageName(pkg)
	metaFilePath := filepath.Join(s.Folder, "meta")

	if scope != "" {
		metaFilePath = filepath.Join(metaFilePath, scope)
	}

	metaFileLocation := filepath.Join(metaFilePath, pkgName+".json")

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
