package config

import "os"

// Config holds application configuration
type Config struct {
	Editor string
	Picker string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Editor: getEnv("EDITOR", "nano"),
		Picker: getEnv("PROOMPT_PICKER", "fzf"),
	}
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
