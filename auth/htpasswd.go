package auth

import (
	"encoding/csv"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func NewHtpasswdFile(filename string) *HtpasswdFile {
	return &HtpasswdFile{Path: filename}
}

type HtpasswdFile struct {
	Path  string
	Info  os.FileInfo
	Users map[string]string
}

func (file *HtpasswdFile) Reload() {
	r, err := os.Open(file.Path)
	if err != nil {
		panic(err)
	}
	csv_reader := csv.NewReader(r)
	csv_reader.Comma = ':'
	csv_reader.Comment = '#'
	csv_reader.TrimLeadingSpace = true

	records, err := csv_reader.ReadAll()
	if err != nil {
		panic(err)
	}

	file.Users = make(map[string]string)
	for _, record := range records {
		file.Users[record[0]] = record[1]
	}
}

func (file *HtpasswdFile) ReloadIfNeeded() {
	info, err := os.Stat(file.Path)
	if err != nil {
		panic(err)
	}
	if file.Info == nil || file.Info.ModTime() != info.ModTime() {
		file.Info = info
		file.Reload()
	}
}

func (file *HtpasswdFile) GetPasswordForUser(user string) (hash string, ok bool) {
	file.ReloadIfNeeded()
	hash, ok = file.Users[user]
	return hash, ok
}

func (file *HtpasswdFile) CompareUsernameAndPassword(username, password string) error {
	hashed, ok := file.GetPasswordForUser(username)
	if !ok {
		return errors.New("username does not exist")
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
