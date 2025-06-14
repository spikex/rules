package tests

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// CommandConfig represents a mapping between a command and its golden file
type CommandConfig struct {
	Command   string
	GoldenFile string
}

// readCommandConfigs reads the configuration file and returns a map of golden files to commands
func readCommandConfigs() (map[string]string, error) {
	configFile := "../tests/golden_commands.txt"
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Map golden files to their commands
	goldenToCommand := make(map[string]string)
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			cmd := strings.TrimSpace(parts[0])
			goldenFile := strings.TrimSpace(parts[1])
			// Store with the golden file as the key
			goldenToCommand[goldenFile] = cmd
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return goldenToCommand, nil
}

// matchWithPlaceholders checks if the output matches the expected text with placeholders
func matchWithPlaceholders(actual, expected string) bool {
	// First try exact match
	if actual == expected {
		return true
	}

	// Handle login.golden placeholders
	if strings.Contains(expected, "<STATE_PLACEHOLDER>") {
		// Extract parts before and after the placeholder
		beforeState := expected[:strings.Index(expected, "<STATE_PLACEHOLDER>")]
		afterState := expected[strings.Index(expected, "<STATE_PLACEHOLDER>")+len("<STATE_PLACEHOLDER>"):]

		// Check if the output contains the parts before and after
		if !strings.HasPrefix(actual, beforeState) {
			return false
		}
		
		// Handle the error placeholder if it exists
		if strings.Contains(afterState, "<ERROR_PLACEHOLDER>") {
			beforeError := afterState[:strings.Index(afterState, "<ERROR_PLACEHOLDER>")]
			afterError := afterState[strings.Index(afterState, "<ERROR_PLACEHOLDER>")+len("<ERROR_PLACEHOLDER>"):]

			// The middle part should be the UUID - find where the beforeError starts
			stateEnd := strings.Index(actual[len(beforeState):], beforeError) + len(beforeState)
			if stateEnd < len(beforeState) {
				return false
			}

			// Extract the UUID
			stateValue := actual[len(beforeState):stateEnd]
			
			// Validate UUID format (simple check)
			if !strings.Contains(stateValue, "-") || len(stateValue) < 30 {
				return false
			}

			// Remaining text after error message
			remainingText := actual[stateEnd+len(beforeError):]
			
			// Check if error message is valid (either "unexpected newline" or "EOF")
			if !strings.HasPrefix(remainingText, "unexpected newline") && 
			   !strings.HasPrefix(remainingText, "EOF") {
				return false
			}
			
			// Check the final part after the error message
			errorMsgEnd := 0
			if strings.HasPrefix(remainingText, "unexpected newline") {
				errorMsgEnd = len("unexpected newline")
			} else if strings.HasPrefix(remainingText, "EOF") {
				errorMsgEnd = len("EOF")
			}
			
			return strings.HasPrefix(remainingText[errorMsgEnd:], afterError)
		} else {
			// Just check if the remaining text appears after the state
			stateEnd := strings.Index(actual[len(beforeState):], afterState)
			return stateEnd >= 0
		}
	}

	return false
}

func TestGoldenFiles(t *testing.T) {
	// Read command configurations
	goldenToCommand, err := readCommandConfigs()
	if err != nil {
		t.Fatalf("Failed to read command configurations: %v", err)
	}
	
	if len(goldenToCommand) == 0 {
		t.Fatal("No commands configured. Check golden_commands.txt file.")
	}
	
	// Find all golden files
	goldenFiles, err := filepath.Glob("golden/**/*.golden")
	if err != nil {
		t.Fatalf("Failed to find golden files: %v", err)
	}
	
	if len(goldenFiles) == 0 {
		t.Fatal("No golden files found. Run scripts/generate_golden.sh first.")
	}
	
	// Make sure the CLI binary exists
	cliPath := "../rules-cli"
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		// Build the CLI if it doesn't exist
		cmd := exec.Command("go", "build", "-o", cliPath, "../main.go")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to build CLI: %v", err)
		}
	}
	
	for _, goldenFile := range goldenFiles {
		t.Run(goldenFile, func(t *testing.T) {
			// Look up the command for this golden file
			cmd, found := goldenToCommand["tests/"+goldenFile]
			if !found {
				t.Skipf("No command configured for golden file: %s", goldenFile)
				return
			}
			
			// Split the command into arguments
			var args []string
			if cmd != "" {
				args = strings.Fields(cmd)
			}
			
			// Run the command
			execCmd := exec.Command(cliPath, args...)
			output, err := execCmd.CombinedOutput()
			
			// We don't fail the test if the command returns non-zero
			// because some golden files might be testing error cases
			
			// Read expected output
			expectedBytes, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("Failed to read golden file: %v", err)
			}
			
			expected := string(expectedBytes)
			actual := string(output)

			// Normalize line endings for cross-platform compatibility
			expected = strings.ReplaceAll(expected, "\r\n", "\n")
			actual = strings.ReplaceAll(actual, "\r\n", "\n")

			// Special case for login test which has placeholders
			if strings.Contains(expected, "<STATE_PLACEHOLDER>") {
				if !matchWithPlaceholders(actual, expected) {
					t.Errorf("Output does not match golden file with placeholders.\nCommand: %s\nExpected:\n%s\n\nGot:\n%s",
						cmd, expected, actual)
				}
			} else {
				// Standard equality check for other files
				if actual != expected {
					t.Errorf("Output does not match golden file.\nCommand: %s\nExpected:\n%s\n\nGot:\n%s",
						cmd, expected, actual)
				}
			}
		})
	}
}

// TestHelp provides a simple example to ensure our testing framework works
func TestHelp(t *testing.T) {
	// This test just makes sure the CLI can run and produce some output
	// It's useful for initial verification of the test setup
	
	cmd := exec.Command("../rules-cli", "--help")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Command output: %s", string(output))
		t.Fatalf("Command failed: %v", err)
	}
	
	if len(output) == 0 {
		t.Error("Help command produced no output")
	}
}