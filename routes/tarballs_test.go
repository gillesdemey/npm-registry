package routes

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type MockedRoundTripper struct {
	mock.Mock
}

func (m *MockedRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	args := m.Called()
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestTryUpstreamTarballFound(t *testing.T) {
	// Overwrites the Roundtripper of http's `DefaultClient` to the mocked
	// `RoundTripper` so we can make assertions on it.
	mockedTransport := new(MockedRoundTripper)
	mockedTransport.On("RoundTrip").Return(
		&http.Response{
			Status:     "200 OK",
			StatusCode: 200,
		}, nil)
	http.DefaultClient.Transport = mockedTransport

	_, err := tryUpstreamTarball("hello", "world")
	assert.Nil(t, err)
}

func TestTryUpstreamTarballNotFound(t *testing.T) {
	mockedTransport := new(MockedRoundTripper)
	mockedTransport.On("RoundTrip").Return(&http.Response{}, errors.New("Nope"))
	http.DefaultClient.Transport = mockedTransport

	_, err := tryUpstreamTarball("hello", "world")
	assert.Error(t, err)
}
