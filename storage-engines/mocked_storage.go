package storageengines

import (
	"github.com/stretchr/testify/mock"
	"io"
)

type MockedStorage struct {
	mock.Mock
}

func (s *MockedStorage) StoreTarball(filename string, reader io.Reader) error {
	args := s.Called(filename, reader)
	return args.Error(0)
}
func (s *MockedStorage) RetrieveTarball(pkg string, filename string, writer io.Writer) error {
	args := s.Called(pkg, filename, writer)
	return args.Error(0)
}

func (s *MockedStorage) StoreMetadata(pkg string, data io.Reader) error {
	args := s.Called(pkg, data)
	return args.Error(0)
}
func (s *MockedStorage) RetrieveMetadata(pkg string, writer io.Writer) error {
	args := s.Called(pkg, writer)
	return args.Error(0)
}

func (s *MockedStorage) StoreUser(pkg string) error {
	args := s.Called(pkg)
	return args.Error(0)
}
func (s *MockedStorage) StoreUserToken(token, username string) error {
	args := s.Called(token, username)
	return args.Error(0)
}

func (s *MockedStorage) RetrieveUsernameFromToken(token string) (string, error) {
	args := s.Called(token)
	return args.String(0), args.Error(1)
}
