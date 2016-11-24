package routes

import (
	"io"
	"github.com/stretchr/testify/mock"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

func NewRendererContext() context.Context {
	ctx := context.Background()
	render := render.New()
	return context.WithValue(ctx, "renderer", render)
}

type MockStorage struct {
	mock.Mock
}

func (s *MockStorage) RetrieveUsernameFromToken(token string) (string, error) {
	args := s.Called(token)
	return args.String(0), args.Error(1)
}

func (s *MockStorage) StoreUserToken(token, username string) error {
	return nil
}

func (s *MockStorage) RetrieveUser(token string, writer io.Writer) error {
	return nil
}
