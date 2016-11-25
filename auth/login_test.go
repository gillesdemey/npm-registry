package auth

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestLogin(t *testing.T) {
  htpasswdProvider := NewHtpasswdProvider("../test/htpasswd.test")
  token, err := htpasswdProvider.Login("foo", "bar")
  assert.Nil(t, err)
  assert.NotNil(t, token)
}

func TestLoginInvalid(t *testing.T) {
  htpasswdProvider := &HtpasswdProvider{
    File:NewHtpasswdFile("../test/htpasswd.test"),
  }
  token, err := htpasswdProvider.Login("bar", "foo")
  assert.Error(t, err)
  assert.Equal(t, token, "")
}
