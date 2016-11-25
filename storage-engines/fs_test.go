package storageengines

import (
  "github.com/stretchr/testify/suite"
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

func (s *FSStorageSuite) TearDownSuite() {
  os.RemoveAll(s.storage.Folder)
}

func TestFSStorageSuite(t *testing.T) {
  suite.Run(t, &FSStorageSuite{
    storage: NewFSStorage("../test/store"),
  })
}
