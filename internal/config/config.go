package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment         string
	DataBackend         string
	HTTPAddr            string
	WebOrigin           string
	DatabaseURL         string
	UploadDir           string
	SessionCookieSecure bool
	SessionTTL          time.Duration
	DevSeedPassword     string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		Environment:     envOr("APP_ENV", "development"),
		DataBackend:     envOr("DATA_BACKEND", "mock"),
		HTTPAddr:        envOr("HTTP_ADDR", ":8080"),
		WebOrigin:       envOr("WEB_ORIGIN", "http://localhost:5173"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		UploadDir:       envOr("UPLOAD_DIR", "./data/uploads"),
		DevSeedPassword: os.Getenv("DEV_SEED_PASSWORD"),
	}

	if cfg.DataBackend != "mock" && cfg.DataBackend != "postgres" {
		return Config{}, fmt.Errorf("DATA_BACKEND must be mock or postgres")
	}
	if cfg.Environment == "production" && cfg.DataBackend != "postgres" {
		return Config{}, fmt.Errorf("DATA_BACKEND must be postgres in production")
	}
	if cfg.DataBackend == "postgres" && cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required when DATA_BACKEND=postgres")
	}
	if cfg.DataBackend == "mock" && cfg.DevSeedPassword == "" {
		return Config{}, fmt.Errorf("DEV_SEED_PASSWORD is required when DATA_BACKEND=mock")
	}

	secure, err := strconv.ParseBool(envOr("SESSION_COOKIE_SECURE", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("parse SESSION_COOKIE_SECURE: %w", err)
	}
	if cfg.Environment == "production" && !secure {
		return Config{}, fmt.Errorf("SESSION_COOKIE_SECURE must be true in production")
	}
	cfg.SessionCookieSecure = secure

	cfg.SessionTTL, err = time.ParseDuration(envOr("SESSION_TTL", "168h"))
	if err != nil || cfg.SessionTTL <= 0 {
		return Config{}, fmt.Errorf("SESSION_TTL must be a positive duration")
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
