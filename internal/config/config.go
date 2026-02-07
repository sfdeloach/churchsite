package config

import (
	"fmt"
	"os"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	AppEnv  string
	AppURL  string
	AppPort string

	DatabaseURL string
	RedisURL    string

	JWTSecret     string
	JWTExpiration string

	SMTPHost  string
	SMTPPort  string
	SMTPUser  string
	SMTPPass  string
	FromEmail string
	FromName  string

	MaxUploadSize string
}

// Load reads configuration from environment variables and returns a Config.
// Returns an error if any required variable is missing.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppURL:  getEnv("APP_URL", "http://localhost:3000"),
		AppPort: getEnv("APP_PORT", "3000"),

		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),

		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		SMTPHost:  os.Getenv("SMTP_HOST"),
		SMTPPort:  os.Getenv("SMTP_PORT"),
		SMTPUser:  os.Getenv("SMTP_USER"),
		SMTPPass:  os.Getenv("SMTP_PASS"),
		FromEmail: os.Getenv("FROM_EMAIL"),
		FromName:  os.Getenv("FROM_NAME"),

		MaxUploadSize: getEnv("MAX_UPLOAD_SIZE", "10485760"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}

	return cfg, nil
}

// IsDevelopment returns true if the app is running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
