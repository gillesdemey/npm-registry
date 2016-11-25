package routes

import (
	"errors"
	"github.com/gillesdemey/npm-registry/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"io/ioutil"
	"net/http"
)

type TarballStorageSuite struct {
  suite.Suite
	storage mocks.MockedStorage
}

func (s *TarballStorageSuite) SetupSuite() {
	httpmock.Activate()

	httpmock.RegisterResponder(
		"GET", "https://registry.npmjs.org/foo/-/foo-0.1.0.tgz",
		httpmock.NewBytesResponder(http.StatusOK, []byte("")),
	)

	packageNotFoundResponder, _ := httpmock.NewJsonResponder(
		http.StatusNotFound,
		`{"error":"package could not be found."}`,
	)
	httpmock.RegisterResponder(
		"GET", "https://registry.npmjs.org/foo/-/foo-0.0.0.tgz",
		packageNotFoundResponder,
	)

	httpmock.RegisterResponder(
		"GET", "https://registry.npmjs.org/fail/-/fail-0.1.0.tgz",
		httpmock.ConnectionFailure,
	)
}

func (s *TarballStorageSuite) TestTryUpstreamTarballFound() {
	_, err := tryUpstreamTarball("foo", "foo-0.1.0.tgz")
	s.Nil(err)
}

func (s *TarballStorageSuite) TestTryUpstreamTarballNotFound() {
	_, err := tryUpstreamTarball("foo", "foo-0.0.0.tgz")
	s.EqualError(err, "no such package available")
}

func (s *TarballStorageSuite) TestTryUpstreamTarballConnectionFailure() {
	_, err := tryUpstreamTarball("fail", "fail-0.1.0.tgz")
	s.NotNil(err)
}

func (s *TarballStorageSuite) TestTryStorageTarball() {
	storage := new(mocks.MockedStorage)
	storage.On("RetrieveTarball", "foo", "foo-0.1.0.tgz", ioutil.Discard).
		Return(nil)

	err := tryStorageTarball(storage, "foo", "foo-0.1.0.tgz", ioutil.Discard)
	s.Nil(err)
}

func (s *TarballStorageSuite) TestTryStorageTarballFailure() {
	storage := new(mocks.MockedStorage)
	storage.On("RetrieveTarball", "foo", "foo-0.1.0.tgz", ioutil.Discard).
		Return(errors.New("something happend"))

	err := tryStorageTarball(storage, "foo", "foo-0.1.0.tgz", ioutil.Discard)
	s.NotNil(err)
}

func (s *TarballStorageSuite) TestUpdateTarballStorage() {
	storage := new(mocks.MockedStorage)
	storage.On("StoreTarball", "foo/foo-0.1.0.tgz", strings.NewReader("")).
		Return(nil)

	err := updateTarballStorage(storage, "foo/foo-0.1.0.tgz", strings.NewReader(""))
	s.Nil(err)
}

func (s *TarballStorageSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestTarballStorageSuite(t *testing.T) {
  suite.Run(t, new(TarballStorageSuite))
}
