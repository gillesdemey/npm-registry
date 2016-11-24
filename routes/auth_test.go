package routes

import (
	"github.com/gillesdemey/npm-registry/auth"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

func TestLogin(t *testing.T) {
	req, err := http.NewRequest(
		"PUT", "/-/user/foo",
		strings.NewReader(`{"name":"foo","password":"bar"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	n := negroni.New()
	n.UseHandlerFunc(Login)

	storage := new(MockStorage)
	auth := auth.NewHtpasswdProvider("../test/htpasswd.test")

	ctx := NewRendererContext()
	ctx = context.WithValue(ctx, "storage", storage)
	ctx = context.WithValue(ctx, "auth", auth)

	n.ServeHTTP(rec, req.WithContext(ctx))

	assert.Equal(t, rec.Code, http.StatusCreated)
	assert.Contains(t, rec.Body.String(), `{"token":"`, `"}`)
}

func TestLoginInvalid(t *testing.T) {
	req, err := http.NewRequest(
		"PUT", "/-/user/bar",
		strings.NewReader(`{"name":"bar","password":"foo"}`),
	)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	n := negroni.New()
	n.UseHandlerFunc(Login)

	auth := auth.NewHtpasswdProvider("../test/htpasswd.test")

	ctx := NewRendererContext()
	ctx = context.WithValue(ctx, "auth", auth)

	n.ServeHTTP(rec, req.WithContext(ctx))

	assert.Equal(t, rec.Code, http.StatusUnauthorized)
	assert.JSONEq(t, rec.Body.String(), `{"error":"invalid credentials"}`)
}

func TestWhoami(t *testing.T) {
	req, err := http.NewRequest("GET", "/-/whoami", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	n := negroni.New(negroni.HandlerFunc(Whoami))

	ctx := NewRendererContext()
	ctx = context.WithValue(ctx, "user", "foo")
	n.ServeHTTP(rec, req.WithContext(ctx))

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.JSONEq(t, rec.Body.String(), `{"username":"foo"}`)
}
