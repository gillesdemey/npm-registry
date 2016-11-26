package routes

import (
	"context"
	"errors"
	"github.com/gillesdemey/npm-registry/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MetaStorageSuite struct {
	suite.Suite
	storage mocks.MockedStorage
}

func (s *MetaStorageSuite) SetupSuite() {
	httpmock.Activate()

	httpmock.RegisterResponder(
		"GET", "https://registry.npmjs.org/foo",
		httpmock.NewBytesResponder(http.StatusOK, []byte("")),
	)

	packageNotFoundResponder, _ := httpmock.NewJsonResponder(
		http.StatusNotFound,
		`{"error":"package could not be found."}`,
	)
	httpmock.RegisterResponder(
		"GET", "https://registry.npmjs.org/notfound",
		packageNotFoundResponder,
	)

	httpmock.RegisterResponder(
		"GET", "https://registry.npmjs.org/failure",
		httpmock.ConnectionFailure,
	)
}

func (s *MetaStorageSuite) TestTryUpstream() {
	_, err := tryUpstream("foo")
	s.Nil(err)
}

func (s *MetaStorageSuite) TestTryUpstreamNotFound() {
	_, err := tryUpstream("notfound")
	s.EqualError(err, "no such package available")
}

func (s *MetaStorageSuite) TestTryUpstreamFailure() {
	_, err := tryUpstream("failure")
	s.Error(err)
}

func (s *MetaStorageSuite) TestTryMetaStorage() {
	storage := new(mocks.MockedStorage)
	storage.On("RetrieveMetadata", "foo", ioutil.Discard).
		Return(nil)

	err := tryMetaStorage(storage, "foo", ioutil.Discard)
	s.Nil(err)
}

func (s *MetaStorageSuite) TestTryMetaStorageNotFound() {
	storage := new(mocks.MockedStorage)
	storage.On("RetrieveMetadata", "foo", ioutil.Discard).
		Return(errors.New("no such package available"))

	err := tryMetaStorage(storage, "foo", ioutil.Discard)
	s.Error(err)
}

func (s *MetaStorageSuite) TestUpdateMetaStorage() {
	storage := new(mocks.MockedStorage)
	storage.On("StoreMetadata", "foo", strings.NewReader("")).
		Return(nil)

	err := updateMetaStorage(storage, "foo", strings.NewReader(""))
	s.Nil(err)
}

func (s *MetaStorageSuite) TestGetPackageMetadata() {
	httpmock.RegisterNoResponder(httpmock.InitialTransport.RoundTrip)

	req, _ := http.NewRequest("GET", "/foo", nil)
	q := req.URL.Query()
	q.Set(":pkg", "foo")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "Bearer abc123")

	rec := httptest.NewRecorder()

	ctx := NewRendererContext()
	storage := new(mocks.MockedStorage)

	storage.On("StoreMetadata", "foo", mock.Anything).
		Return(nil)

	ctx = context.WithValue(ctx, "storage", storage)

	GetPackageMetadata(rec, req.WithContext(ctx), func(w http.ResponseWriter, req *http.Request) {
		s.Equal(rec.Code, http.StatusOK)
	})}

func (s *MetaStorageSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestMetaStorageSuite(t *testing.T) {
	suite.Run(t, new(MetaStorageSuite))
}

func TestQueryEscape(t *testing.T) {
	expected := "@foo%2Fbar"
	output := queryEscape("@foo/bar")
	assert.Equal(t, output, expected)
}
