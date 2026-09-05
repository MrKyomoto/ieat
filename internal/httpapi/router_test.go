package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MrKyomoto/ieat/internal/config"
)

func TestLoginRejectsUntrustedOrigin(t *testing.T) {
	router := NewRouter(config.Config{WebOrigin: "https://ieat.example.edu"}, Dependencies{})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/session", strings.NewReader(`{}`))
	request.Header.Set("Origin", "https://attacker.example")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}
