package routes

import (
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestPing(t *testing.T) {
  req, err := http.NewRequest("GET", "/-/ping", nil)
  if err != nil {
    t.Fatal(err)
  }

  rec := httptest.NewRecorder()
  handler := http.HandlerFunc(Ping)

  ctx := NewRendererContext()
  handler.ServeHTTP(rec, req.WithContext(ctx))

  if status := rec.Code; status != http.StatusOK {
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusOK)
  }

  expected := `{}`
  if rec.Body.String() != expected {
    t.Errorf("handler returned unexpected body: got %v want %v",
      rec.Body.String(), expected)
  }
}
