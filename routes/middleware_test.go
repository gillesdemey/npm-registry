package routes

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateToken(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer abc123")

	rec := httptest.NewRecorder()

	ctx := req.Context()
	storage := new(MockStorage)

	storage.On("RetrieveUsernameFromToken", "abc123").Return("foo", nil)

	ctx = context.WithValue(ctx, "storage", storage)

	ValidateToken(rec, req.WithContext(ctx), func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, req.Context().Value("user").(string), "foo")
		assert.Equal(t, req.Context().Value("token").(string), "abc123")
	})
}

func TestValidateTokenWithInvalidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer abc123")

	rec := httptest.NewRecorder()

	ctx := req.Context()
	storage := new(MockStorage)

	storage.
		On("RetrieveUsernameFromToken", "abc123").
		Return("", errors.New("oh no!"))

	ctx = context.WithValue(ctx, "storage", storage)

	ValidateToken(rec, req.WithContext(ctx), func(w http.ResponseWriter, req *http.Request) {
		t.Fail()
	})

	assert.Equal(t, rec.Code, http.StatusUnauthorized)
}

type MockStorage struct {
	mock.Mock
}

func (s *MockStorage) RetrieveUsernameFromToken(token string) (string, error) {
	args := s.Called(token)
	return args.String(0), args.Error(1)
}
