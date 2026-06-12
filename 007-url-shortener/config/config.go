package config

import "os"

// Config holds everything the service reads from the environment.
type Config struct {
	Port        string
	DatabaseURL string // empty means "use the in-memory store"
}

// Load reads the environment and applies defaults.
func Load() Config {
	return Config{
		Port:        envOr("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

// envOr returns the env var named key, or fallback if it is unset/empty.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
