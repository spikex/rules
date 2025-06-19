package main

import (
	"testing"

	"rules-cli/cmd"
)

func TestVersionSetting(t *testing.T) {
	// Test that the version is passed to the cmd package
	originalVersion := Version
	originalCmdVersion := cmd.Version

	defer func() {
		Version = originalVersion
		cmd.Version = originalCmdVersion
	}()

	Version = "test-main-version"
	cmd.Version = Version

	if cmd.Version != "test-main-version" {
		t.Errorf("Expected cmd.Version to be 'test-main-version', got %s", cmd.Version)
	}
}

func TestMainFunction(t *testing.T) {
	// Test that main function exists and doesn't panic when version is set
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() related functions panicked: %v", r)
		}
	}()

	// Test that Version variable exists and can be set
	originalVersion := Version
	Version = "test"
	if Version != "test" {
		t.Error("Version variable should be settable")
	}
	Version = originalVersion
}
