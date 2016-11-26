package routes

import (
	"context"
	"errors"
	"github.com/gillesdemey/npm-registry/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	_ "fmt"
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
		Return(errors.New("something happened"))

	err := tryStorageTarball(storage, "foo", "foo-0.1.0.tgz", ioutil.Discard)
	s.NotNil(err)
}

func (s *TarballStorageSuite) TestUpdateTarballStorage() {
	storage := new(mocks.MockedStorage)
	storage.On("StoreTarball", "foo-0.1.0.tgz", mock.Anything).
		Return(nil)

	err := updateTarballStorage(storage, "foo-0.1.0.tgz", strings.NewReader(""))
	s.Nil(err)
}

func (s *TarballStorageSuite) TestGetTarball() {
	httpmock.RegisterNoResponder(httpmock.InitialTransport.RoundTrip)

	req, _ := http.NewRequest("GET", "/foo/-/foo-0.1.0.tgz", nil)
	q := req.URL.Query()
	q.Set(":pkg", "foo")
	q.Set(":filename", "foo-0.1.0.tgz")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "Bearer abc123")

	rec := httptest.NewRecorder()

	ctx := req.Context()
	storage := new(mocks.MockedStorage)

	// storage.On("RetrieveUsernameFromToken", "abc123").Return("foo", nil)
	storage.On("RetrieveTarball", "foo", "foo-0.1.0.tgz", mock.Anything).
		Return(nil)
	storage.On("StoreTarball", "foo-0.1.0.tgz", mock.Anything).
		Return(nil)

	ctx = context.WithValue(ctx, "storage", storage)

	GetTarball(rec, req.WithContext(ctx), func(w http.ResponseWriter, req *http.Request) {
		s.Equal(rec.Code, http.StatusOK)
	})
}

func (s *TarballStorageSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestTarballStorageSuite(t *testing.T) {
	suite.Run(t, new(TarballStorageSuite))
}
