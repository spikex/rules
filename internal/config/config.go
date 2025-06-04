package config

import (
	"os"
	"strings"
)

// Config holds configuration for the rules CLI
type Config struct {
	RegistryURL    string
	DefaultFormat  string
	Username       string
	Email          string
	Formats        []string
}

// Initialize sets up the configuration from environment variables
func Initialize() (*Config, error) {
	config := Config{
	// Set defaults
		RegistryURL:    "https://rules.example.com",
		DefaultFormat:  "default",
		Username:       "",
		Email:          "",
		Formats:        []string{"default", "cursor"},
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

	if envFormats := os.Getenv("RULES_FORMATS"); envFormats != "" {
		config.Formats = strings.Split(envFormats, ",")
		// Trim whitespace from format names
		for i, format := range config.Formats {
			config.Formats[i] = strings.TrimSpace(format)
		}
	}
	
	return &config, nil
}

// LoadConfig loads the configuration
func LoadConfig() (*Config, error) {
	return Initialize()
}