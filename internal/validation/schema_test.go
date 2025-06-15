package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateRuleFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "rule_validation_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		content     string
		expectValid bool
		expectError string
	}{
		{
			name: "Valid rule with all fields",
			content: `---
alwaysApply: true
description: "This is a test rule"
globs: "*.{js,ts}"
tags: ["frontend", "javascript"]
---

# Test Rule

This is the body of the rule.
`,
			expectValid: true,
		},
		{
			name: "Valid rule with minimal frontmatter",
			content: `---
description: "Simple rule"
---

# Simple Rule

Just a description.
`,
			expectValid: true,
		},
		{
			name: "Valid rule without frontmatter",
			content: `# No Frontmatter Rule

This rule has no frontmatter.
`,
			expectValid: true,
		},
		{
			name: "Invalid - unknown field",
			content: `---
description: "Test rule"
unknownField: "not allowed"
---

# Test Rule
`,
			expectValid: false,
			expectError: "Additional property unknownField is not allowed",
		},
		{
			name: "Invalid - alwaysApply not boolean",
			content: `---
alwaysApply: "yes"
description: "Test rule"
---

# Test Rule
`,
			expectValid: false,
			expectError: "Invalid type. Expected: boolean, given: string",
		},
		{
			name: "Invalid - empty description",
			content: `---
description: ""
---

# Test Rule
`,
			expectValid: false,
			expectError: "String length must be greater than or equal to 1",
		},
		{
			name: "Invalid - description too long",
			content: `---
description: "` + strings.Repeat("a", 501) + `"
---

# Test Rule
`,
			expectValid: false,
			expectError: "String length must be less than or equal to 500",
		},
		{
			name: "Invalid - empty globs",
			content: `---
globs: ""
---

# Test Rule
`,
			expectValid: false,
			expectError: "String length must be greater than or equal to 1",
		},
		{
			name: "Invalid - tags not array",
			content: `---
tags: "not an array"
---

# Test Rule
`,
			expectValid: false,
			expectError: "Invalid type. Expected: array, given: string",
		},
		{
			name: "Invalid - duplicate tags",
			content: `---
tags: ["frontend", "frontend"]
---

# Test Rule
`,
			expectValid: false,
			expectError: "array items[0,1] must be unique",
		},
		{
			name: "Invalid - too many tags",
			content: `---
tags: [` + generateTags(21) + `]
---

# Test Rule
`,
			expectValid: false,
			expectError: "Array must have at most 20 items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary rule file
			filePath := filepath.Join(tempDir, "test-rule.md")
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Validate the file
			result, err := ValidateRuleFile(filePath)
			if err != nil {
				t.Fatalf("ValidateRuleFile returned error: %v", err)
			}

			// Check if validation result matches expectation
			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
				if len(result.Errors) > 0 {
					t.Errorf("Validation errors: %+v", result.Errors)
				}
			}

			// Check for expected error message
			if !tt.expectValid && tt.expectError != "" {
				found := false
				for _, validationErr := range result.Errors {
					if validationErr.Message == tt.expectError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in validation errors: %+v", tt.expectError, result.Errors)
				}
			}
		})
	}
}

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectNil   bool
		expectError bool
	}{
		{
			name: "Valid frontmatter",
			content: `---
alwaysApply: true
description: "Test"
---

Content here
`,
			expectNil: false,
		},
		{
			name: "No frontmatter",
			content: `# Just Content

No frontmatter here.
`,
			expectNil: true,
		},
		{
			name: "Empty frontmatter",
			content: `---

---

Content here
`,
			expectNil: true,
		},
		{
			name: "Invalid YAML",
			content: `---
invalid: yaml: content
---

Content here
`,
			expectError: true,
		},
		{
			name: "Unterminated frontmatter",
			content: `---
description: "Test"

Content without closing delimiter
`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractFrontmatter([]byte(tt.content))

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result but got: %+v", result)
			}

			if !tt.expectNil && result == nil {
				t.Errorf("Expected non-nil result but got nil")
			}
		})
	}
}

// generateTags creates a comma-separated string of n tags for testing
func generateTags(n int) string {
	tags := make([]string, n)
	for i := 0; i < n; i++ {
		tags[i] = `"tag` + string(rune('0'+i%10)) + `"`
	}
	if len(tags) == 0 {
		return ""
	}

	result := tags[0]
	for i := 1; i < len(tags); i++ {
		result += ", " + tags[i]
	}
	return result
}
