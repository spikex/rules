package config

import (
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	// Test basic initialization
	cfg, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Initialize() returned nil config")
	}

	// Test default values
	if cfg.RegistryURL == "" {
		t.Error("RegistryURL should not be empty")
	}

	if cfg.DefaultFormat == "" {
		t.Error("DefaultFormat should not be empty")
	}

	// Test that default values are set correctly
	expectedRegistryURL := "https://api.continue.dev"
	if cfg.RegistryURL != expectedRegistryURL {
		t.Errorf("Expected RegistryURL to be %s, got %s", expectedRegistryURL, cfg.RegistryURL)
	}

	expectedDefaultFormat := "default"
	if cfg.DefaultFormat != expectedDefaultFormat {
		t.Errorf("Expected DefaultFormat to be %s, got %s", expectedDefaultFormat, cfg.DefaultFormat)
	}
}

func TestLoadConfig(t *testing.T) {
	// Test LoadConfig function
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}
}

func TestInitializeWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("RULES_REGISTRY_URL", "https://test.example.com")
	os.Setenv("RULES_DEFAULT_FORMAT", "test-format")
	os.Setenv("RULES_USERNAME", "testuser")
	os.Setenv("RULES_EMAIL", "test@example.com")

	defer func() {
		// Clean up environment variables
		os.Unsetenv("RULES_REGISTRY_URL")
		os.Unsetenv("RULES_DEFAULT_FORMAT")
		os.Unsetenv("RULES_USERNAME")
		os.Unsetenv("RULES_EMAIL")
	}()

	cfg, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() with env vars failed: %v", err)
	}

	// Test that environment variables override defaults
	if cfg.RegistryURL != "https://test.example.com" {
		t.Errorf("Expected RegistryURL to be overridden by env var, got %s", cfg.RegistryURL)
	}

	if cfg.DefaultFormat != "test-format" {
		t.Errorf("Expected DefaultFormat to be overridden by env var, got %s", cfg.DefaultFormat)
	}

	if cfg.Username != "testuser" {
		t.Errorf("Expected Username to be set by env var, got %s", cfg.Username)
	}

	if cfg.Email != "test@example.com" {
		t.Errorf("Expected Email to be set by env var, got %s", cfg.Email)
	}
}
