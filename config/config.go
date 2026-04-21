// Package config loads application configuration from environment variables.
// A .env file is optional — when running in Docker, variables are injected directly.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration values.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Load reads environment variables, optionally from a .env file.
// Missing .env is silently ignored so Docker deployments work without one.
func Load() (*Config, error) {
	// ignore missing .env — env vars may be injected directly (e.g. Docker)
	_ = godotenv.Load()

	return &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}, nil
}
