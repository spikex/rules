package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds configuration for the rules CLI
type Config struct {
	RegistryURL    string
	DefaultFormat  string
	Username       string
	Email          string
	Formats        []string
}

// Initialize sets up the configuration from environment variables and Viper
func Initialize() (*Config, error) {
	// Set up Viper
	viper.SetConfigName("rules-cli")
	viper.SetConfigType("yaml")
	
	// Look for config in the user's home directory
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(home, ".rules-cli"))
	}
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("registry_url", "https://api.continue.dev")
	viper.SetDefault("default_format", "default")
	viper.SetDefault("username", "")
	viper.SetDefault("email", "")
	viper.SetDefault("formats", []string{"default", "cursor"})
	
	// Set default values for auth-related configurations
	viper.SetDefault("workos_client_id", "client_01J0FW6XN8N2XJAECF7NE0Y65J")
	viper.SetDefault("app_url", "https://hub.continue.dev/")
	viper.SetDefault("api_base", "https://api.continue.dev/")

	// Bind environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("RULES")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	_ = viper.ReadInConfig() // Ignore error if config file is not found

	config := Config{
		RegistryURL:   viper.GetString("registry_url"),
		DefaultFormat: viper.GetString("default_format"),
		Username:      viper.GetString("username"),
		Email:         viper.GetString("email"),
		Formats:       viper.GetStringSlice("formats"),
	}
	
	return &config, nil
}

// LoadConfig loads the configuration
func LoadConfig() (*Config, error) {
	return Initialize()
}