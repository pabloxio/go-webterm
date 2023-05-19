package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pabloxio/go-webterm/handlers"
)

func TestHealthz(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HealthzHandler)
	handler.ServeHTTP(rr, req)

	// Check status code is 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("Wrong body: got %v want %v", rr.Body.String(), expected)
	}
}
