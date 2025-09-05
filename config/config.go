package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port         string
	BucketName   string
	ProjectID    string
	CredentialsPath string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		BucketName:   getEnv("GCS_BUCKET_NAME", ""),
		ProjectID:    getEnv("GOOGLE_CLOUD_PROJECT", ""),
		CredentialsPath: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}