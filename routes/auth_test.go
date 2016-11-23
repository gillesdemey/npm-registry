package routes

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	assert.Equal(t, rec.Body.String(), `{"username":"foo"}`)
}
