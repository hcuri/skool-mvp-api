package config

import (
	"os"
)

// Config holds runtime configuration for the API server.
type Config struct {
	Port     string
	LogLevel string
}

// Load reads configuration from environment variables, supplying defaults when unset.
func Load() Config {
	return Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
