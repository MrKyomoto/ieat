package config

import (
	"strings"
	"testing"
)

func TestProductionRequiresSecureSessionCookie(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATA_BACKEND", "postgres")
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("SESSION_COOKIE_SECURE", "false")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "SESSION_COOKIE_SECURE") {
		t.Fatalf("Load() error = %v, want SESSION_COOKIE_SECURE validation error", err)
	}
}
