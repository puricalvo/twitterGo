package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeAPIOptions(t *testing.T) {
	t.Setenv("CORS_ORIGIN", "http://localhost:3000")
	request := httptest.NewRequest(http.MethodOptions, "/tweet", nil)
	response := httptest.NewRecorder()

	serveAPI(response, request, "test-secret")

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("allow-origin = %q, want %q", got, "http://localhost:3000")
	}
}
