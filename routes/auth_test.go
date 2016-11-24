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
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"PUT", "/-/user/foo",
		strings.NewReader(`{"name":"foo","password":"bar"}`),
	)
	serveTestRequest(rec, req)

	assert.Equal(t, rec.Code, http.StatusCreated)
	assert.Contains(t, rec.Body.String(), `{"token":"`, `"}`)
}

func TestLoginInvalid(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"PUT", "/-/user/bar",
		strings.NewReader(`{"name":"bar","password":"foo"}`),
	)
	serveTestRequest(rec, req)

	assert.Equal(t, rec.Code, http.StatusUnauthorized)
	assert.JSONEq(t, rec.Body.String(), `{"error":"invalid credentials"}`)
}

func TestLoginBadRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/-/user/foo", strings.NewReader(``))
	serveTestRequest(rec, req)

	assert.Equal(t, rec.Code, http.StatusBadRequest)
	assert.JSONEq(t, rec.Body.String(), `{"error":"bad request"}`)
}

func serveTestRequest(rec *httptest.ResponseRecorder, req *http.Request) {
	n := negroni.New()
	n.UseHandlerFunc(Login)

	auth := auth.NewHtpasswdProvider("../test/htpasswd.test")
	storage := new(MockStorage)

	ctx := NewRendererContext()
	ctx = context.WithValue(ctx, "storage", storage)
	ctx = context.WithValue(ctx, "auth", auth)
	n.ServeHTTP(rec, req.WithContext(ctx))
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
