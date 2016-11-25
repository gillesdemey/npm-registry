package storageengines

import (
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type FSStorageSuite struct {
	suite.Suite
	storage *FSStorage
}

func (s *FSStorageSuite) SetupSuite() {
	os.RemoveAll(s.storage.Folder)
	os.MkdirAll(s.storage.Folder, os.ModePerm)

	file, _ := os.Open("../test/fixtures/tokens.toml")
	defer file.Close()

	dest, _ := os.Create("../test/store/tokens.toml")
	defer dest.Close()

	_, err := io.Copy(dest, file)
	s.Nil(err)
}

func (s *FSStorageSuite) TestStoreTarballUnscoped() {
	var err error
	storage := s.storage
	tarball, _ := os.Open("../test/fixtures/tarballs/inherits/inherits-2.0.3.tgz")
	defer tarball.Close()

	err = storage.StoreTarball("inherits-2.0.3.tgz", tarball)
	s.Nil(err)

	_, err = os.Stat(path.Join(s.storage.Folder, "tarballs", "inherits-2.0.3.tgz"))
	s.Nil(err)
}

func (s *FSStorageSuite) TestStoreTarballScoped() {
	storage := s.storage
	tarball, _ := os.Open("../test/fixtures/tarballs/@request/qs/qs-0.1.0.tgz")
	defer tarball.Close()

	err := storage.StoreTarball("@request/qs-0.1.0.tgz", tarball)
	s.Nil(err)

	_, err = os.Stat(path.Join(s.storage.Folder, "tarballs", "@request", "qs-0.1.0.tgz"))
	s.Nil(err)
}

func (s *FSStorageSuite) TestRetrieveTarballUnscoped() {
	storage := s.storage
	tarball, _ := os.Open("../test/fixtures/tarballs/inherits/inherits-2.0.3.tgz")
	storage.StoreTarball("inherits-2.0.3.tgz", tarball)

	err := storage.RetrieveTarball("inherits", "inherits-2.0.3.tgz", ioutil.Discard)
	s.Nil(err)
}

func (s *FSStorageSuite) TestRetrieveTarballScoped() {
	storage := s.storage
	tarball, _ := os.Open("../test/fixtures/tarballs/@request/qs/qs-0.1.0.tgz")
	storage.StoreTarball("@request/qs-0.1.0.tgz", tarball)

	err := storage.RetrieveTarball("@request/qs", "qs-0.1.0.tgz", ioutil.Discard)
	s.Nil(err)
}

func (s *FSStorageSuite) TestRetrieveUsernameFromToken() {
	storage := s.storage
	username, err := storage.RetrieveUsernameFromToken("abc-123-def-456")
	s.Nil(err)
	s.Equal(username, "foo")
}

func (s *FSStorageSuite) TestRetrieveUsernameFromTokenUnknownUser() {
	storage := s.storage
	username, err := storage.RetrieveUsernameFromToken("nope")
	s.NotNil(err)
	s.Equal(username, "")
}

func (s *FSStorageSuite) TestStoreUserToken() {
	storage := s.storage

	err := storage.StoreUserToken("foo-123-bar-456", "bar")
	s.Nil(err)

	username, err := storage.RetrieveUsernameFromToken("foo-123-bar-456")
	s.Nil(err)
	s.Equal(username, "bar")
}

func (s *FSStorageSuite) TestStoreMetadataUnscoped() {
	meta, _ := os.Open("../test/fixtures/meta/qs.json")
	defer meta.Close()

	err := s.storage.StoreMetadata("qs", meta)
	s.Nil(err)

	file, err := os.Open(path.Join(s.storage.Folder, "meta", "qs.json"))
	defer file.Close()

	s.Nil(err)
	s.NotNil(file)
}

func (s *FSStorageSuite) TestStoreMetadataScoped() {
	meta, _ := os.Open("../test/fixtures/meta/@request/qs.json")
	defer meta.Close()

	err := s.storage.StoreMetadata("@request/qs", meta)
	s.Nil(err)

	file, err := os.Open(path.Join(s.storage.Folder, "meta", "@request", "qs.json"))
	defer file.Close()

	s.Nil(err)
	s.NotNil(file)
}

func (s *FSStorageSuite) TestRetrieveMetadataUnscoped() {
	meta, _ := os.Open("../test/fixtures/meta/qs.json")
	defer meta.Close()

	s.storage.StoreMetadata("qs", meta)
	err := s.storage.RetrieveMetadata("qs", ioutil.Discard)
	s.Nil(err)
}

func (s *FSStorageSuite) TestRetrieveMetadataScoped() {
	meta, _ := os.Open("../test/fixtures/meta/@request/qs.json")
	defer meta.Close()

	s.storage.StoreMetadata("@request/qs", meta)
	err := s.storage.RetrieveMetadata("@request/qs", ioutil.Discard)
	s.Nil(err)
}

func (s *FSStorageSuite) TearDownSuite() {
	os.RemoveAll(s.storage.Folder)
}

func TestFSStorageSuite(t *testing.T) {
	suite.Run(t, &FSStorageSuite{
		storage: NewFSStorage("../test/store"),
	})
}
