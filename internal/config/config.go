package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds configuration for the rules CLI
type Config struct {
	RegistryURL string `mapstructure:"registry_url"`
	DefaultFormat string `mapstructure:"default_format"`
	Username string `mapstructure:"username"`
	Email string `mapstructure:"email"`
}

// Initialize sets up the configuration
func Initialize() (*Config, error) {
	// Set defaults
	viper.SetDefault("registry_url", "https://rules.example.com")
	viper.SetDefault("default_format", "default")
	
	// Set config name
	viper.SetConfigName("rules-cli")
	viper.SetConfigType("yaml")
	
	// Add config search paths
	configHome, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(configHome, "rules-cli"))
	}
	viper.AddConfigPath(".")
	
	// Read environment variables
	viper.SetEnvPrefix("RULES")
	viper.AutomaticEnv()
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// It's ok if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}
	
	// Parse config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	
	return &config, nil
}

// SaveConfig saves the current configuration
func SaveConfig(config *Config) error {
	viper.Set("registry_url", config.RegistryURL)
	viper.Set("default_format", config.DefaultFormat)
	viper.Set("username", config.Username)
	viper.Set("email", config.Email)
	
	configHome, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("unable to find user config dir: %w", err)
	}
	
	configDir := filepath.Join(configHome, "rules-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("unable to create config directory: %w", err)
	}
	
	configPath := filepath.Join(configDir, "rules-cli.yaml")
	return viper.WriteConfigAs(configPath)
}