package config

import (
	"os"
)

// Config holds configuration for the rules CLI
type Config struct {
	RegistryURL    string
	DefaultFormat  string
	Username       string
	Email          string
}

// Initialize sets up the configuration from environment variables
func Initialize() (*Config, error) {
	config := Config{
	// Set defaults
		RegistryURL:    "https://rules.example.com",
		DefaultFormat:  "default",
		Username:       "",
		Email:          "",
	}

	// Override from environment variables
	if envURL := os.Getenv("RULES_REGISTRY_URL"); envURL != "" {
		config.RegistryURL = envURL
	}

	if envFormat := os.Getenv("RULES_DEFAULT_FORMAT"); envFormat != "" {
		config.DefaultFormat = envFormat
	}
	
	if envUsername := os.Getenv("RULES_USERNAME"); envUsername != "" {
		config.Username = envUsername
	}

	if envEmail := os.Getenv("RULES_EMAIL"); envEmail != "" {
		config.Email = envEmail
	}
	return &config, nil
}
