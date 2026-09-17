package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration, loaded from environment variables.
type Config struct {
	DatabaseURL   string
	Port          string
	AdminUser     string
	AdminPassword string
	// StaticDir, when set, makes the server also serve a static frontend from
	// this directory (single-origin deploy). Empty = API only.
	StaticDir string
}

// Load reads configuration from a .env file (if present) and the environment.
// Environment variables always win over the .env file.
func Load() Config {
	// Ignore the error: in production the .env file may legitimately be absent.
	_ = godotenv.Load()

	return Config{
		DatabaseURL:   env("DATABASE_URL", "postgres://travelcrm:travelcrm@localhost:5432/travel_crm?sslmode=disable"),
		Port:          env("PORT", "8080"),
		AdminUser:     env("ADMIN_USER", "admin"),
		AdminPassword: env("ADMIN_PASSWORD", "admin123"),
		StaticDir:     env("STATIC_DIR", ""),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
