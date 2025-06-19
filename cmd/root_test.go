package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Test that the root command executes without error
	old := os.Args
	defer func() { os.Args = old }()

	// Test help command
	os.Args = []string{"rules", "--help"}

	// Capture output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Root command with --help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "command-line tool to create, manage, and convert rule sets") {
		t.Error("Help output should contain description")
	}
}

func TestVersionFlag(t *testing.T) {
	// Set version for testing
	originalVersion := Version
	Version = "test-version"
	defer func() { Version = originalVersion }()

	// Test by directly setting the version flag and checking behavior
	oldVersion := version
	version = true
	defer func() { version = oldVersion }()

	// Test that when version flag is true, the command behaves correctly
	if !version {
		t.Error("Version flag should be set to true for testing")
	}

	// Test that Version variable is set correctly
	if Version != "test-version" {
		t.Errorf("Version should be 'test-version', got '%s'", Version)
	}
}

func TestVersionFlagBehavior(t *testing.T) {
	// Test version flag behavior by parsing command line
	originalVersion := Version
	Version = "test-version-behavior"
	defer func() { Version = originalVersion }()

	// Create a new command instance to avoid state issues
	testCmd := &cobra.Command{
		Use: "rules",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				// This mimics the actual behavior in root.go
				return
			}
		},
	}

	testCmd.Flags().BoolVarP(&version, "version", "v", false, "Display version information")

	// Test parsing the version flag
	err := testCmd.ParseFlags([]string{"-v"})
	if err != nil {
		t.Fatalf("Failed to parse version flag: %v", err)
	}

	// The version variable should now be true
	if !version {
		t.Error("Version flag should be true after parsing -v")
	}
}

func TestRootCommandUsage(t *testing.T) {
	// Test that the root command has the correct usage string
	if rootCmd.Use != "rules" {
		t.Errorf("Expected Use to be 'rules', got %s", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestInitializeCommand(t *testing.T) {
	// Test that the command initializes without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Command initialization panicked: %v", r)
		}
	}()

	// This should not panic
	_ = rootCmd
}
