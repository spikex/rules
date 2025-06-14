package formats

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRenderRulesToFormat_DirectoryBasedFormats tests all directory-based formats
func TestRenderRulesToFormat_DirectoryBasedFormats(t *testing.T) {
	tests := []struct {
		name                string
		format              string
		expectedDir         string
		expectedExtension   string
		expectedFrontmatter string
	}{
		{
			name:              "Continue format",
			format:            "continue",
			expectedDir:       ".continue/rules",
			expectedExtension: ".md",
			expectedFrontmatter: `---
alwaysApply: true
description: Test rule description
globs: "**/*.tsx"
---`,
		},
		{
			name:              "Cursor format",
			format:            "cursor",
			expectedDir:       ".cursor/rules",
			expectedExtension: ".mdc",
			expectedFrontmatter: `---
alwaysApply: true
description: Test rule description
globs: "**/*.tsx"
---`,
		},
		{
			name:              "Windsurf format",
			format:            "windsurf",
			expectedDir:       ".windsurf/rules",
			expectedExtension: ".md",
			expectedFrontmatter: `---
description: Test rule description
globs: "**/*.tsx"
trigger: always_on
---`,
		},
		{
			name:              "Copilot format",
			format:            "copilot",
			expectedDir:       ".github/instructions",
			expectedExtension: ".instructions.md",
			expectedFrontmatter: `---
applyTo: "**/*.tsx"
description: Test rule description
---`,
		},
		{
			name:              "Cline format",
			format:            "cline",
			expectedDir:       ".clinerules",
			expectedExtension: ".md",
			expectedFrontmatter: `---
description: Test rule description
---`,
		},
		{
			name:              "Cody format",
			format:            "cody",
			expectedDir:       ".sourcegraph",
			expectedExtension: ".rule.md",
			expectedFrontmatter: `---
description: Test rule description
---`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary source directory
			sourceDir := createTempSourceDir(t)
			defer os.RemoveAll(sourceDir)

			// Create test rule file
			testRuleContent := `---
alwaysApply: true
description: Test rule description
globs: "**/*.tsx"
tags: [frontend, react]
---

# Test Rule

This is a test rule for React components.
`
			testRuleFile := filepath.Join(sourceDir, "test-rule.md")
			if err := os.WriteFile(testRuleFile, []byte(testRuleContent), 0644); err != nil {
				t.Fatalf("Failed to create test rule file: %v", err)
			}

			// Change to temporary directory for the test
			origDir := getCurrentDir(t)
			tmpTestDir := createTempTestDir(t)
			defer func() {
				os.Chdir(origDir)
				os.RemoveAll(tmpTestDir)
			}()
			os.Chdir(tmpTestDir)

			// Run RenderRulesToFormat
			err := RenderRulesToFormat(sourceDir, tt.format, false)
			if err != nil {
				t.Fatalf("RenderRulesToFormat failed: %v", err)
			}

			// Verify directory was created
			if _, err := os.Stat(tt.expectedDir); os.IsNotExist(err) {
				t.Errorf("Expected directory %s was not created", tt.expectedDir)
			}

			// Verify file was created with correct extension
			expectedFile := filepath.Join(tt.expectedDir, "test-rule"+tt.expectedExtension)
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not created", expectedFile)
			}

			// Verify file content and frontmatter transformation
			content, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			contentStr := string(content)
			// Use more flexible matching for YAML quotes
			expectedFrontmatterFlexible := strings.ReplaceAll(tt.expectedFrontmatter, `"`, `'`)
			if !strings.Contains(contentStr, expectedFrontmatterFlexible) {
				t.Errorf("Expected frontmatter not found.\nExpected:\n%s\nActual content:\n%s", expectedFrontmatterFlexible, contentStr)
			}

			// Verify body content is preserved
			if !strings.Contains(contentStr, "# Test Rule") {
				t.Errorf("Rule body content was not preserved")
			}
		})
	}
}

// TestRenderRulesToFormat_SingleFileFormats tests all single file formats
func TestRenderRulesToFormat_SingleFileFormats(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		expectedFile string
	}{
		{
			name:         "Claude format",
			format:       "claude",
			expectedFile: "CLAUDE.md",
		},
		{
			name:         "Codex format",
			format:       "codex",
			expectedFile: "AGENT.md",
		},
		{
			name:         "Amp format",
			format:       "amp",
			expectedFile: "AGENT.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary source directory
			sourceDir := createTempSourceDir(t)
			defer os.RemoveAll(sourceDir)

			// Create test rule files - one with alwaysApply: true, one without
			alwaysApplyRuleContent := `---
alwaysApply: true
description: Always apply rule
---

# Always Apply Rule

This rule should be included in single file formats.
`

			nonAlwaysApplyRuleContent := `---
alwaysApply: false
description: Manual rule
---

# Manual Rule

This rule should NOT be included in single file formats.  
`

			noFrontmatterRuleContent := `# No Frontmatter Rule

This rule has no frontmatter and should be included (treated as alwaysApply: true).
`

			// Write test files
			if err := os.WriteFile(filepath.Join(sourceDir, "always-apply.md"), []byte(alwaysApplyRuleContent), 0644); err != nil {
				t.Fatalf("Failed to create always-apply rule file: %v", err)
			}
			if err := os.WriteFile(filepath.Join(sourceDir, "manual.md"), []byte(nonAlwaysApplyRuleContent), 0644); err != nil {
				t.Fatalf("Failed to create manual rule file: %v", err)
			}
			if err := os.WriteFile(filepath.Join(sourceDir, "no-frontmatter.md"), []byte(noFrontmatterRuleContent), 0644); err != nil {
				t.Fatalf("Failed to create no-frontmatter rule file: %v", err)
			}

			// Change to temporary directory for the test
			origDir := getCurrentDir(t)
			tmpTestDir := createTempTestDir(t)
			defer func() {
				os.Chdir(origDir)
				os.RemoveAll(tmpTestDir)
			}()
			os.Chdir(tmpTestDir)

			// Run RenderRulesToFormat
			err := RenderRulesToFormat(sourceDir, tt.format, false)
			if err != nil {
				t.Fatalf("RenderRulesToFormat failed: %v", err)
			}

			// Verify single file was created
			if _, err := os.Stat(tt.expectedFile); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not created", tt.expectedFile)
			}

			// Verify file content
			content, err := os.ReadFile(tt.expectedFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			contentStr := string(content)

			// Should contain the always-apply rule
			if !strings.Contains(contentStr, "Always Apply Rule") {
				t.Errorf("Single file should contain always-apply rule")
			}

			// Should contain the no-frontmatter rule (treated as alwaysApply: true)
			if !strings.Contains(contentStr, "No Frontmatter Rule") {
				t.Errorf("Single file should contain no-frontmatter rule")
			}

			// Should NOT contain the manual rule (alwaysApply: false)
			if strings.Contains(contentStr, "Manual Rule") {
				t.Errorf("Single file should NOT contain manual rule")
			}

			// Verify it starts with "# Rules" header
			if !strings.HasPrefix(contentStr, "# Rules\n") {
				t.Errorf("Single file should start with '# Rules' header")
			}

			// Verify no frontmatter is included in single file formats
			if strings.Contains(contentStr, "---") {
				t.Errorf("Single file formats should not contain frontmatter")
			}
		})
	}
}

// TestRenderRulesToFormat_SubdirectoryStructure tests that subdirectory structure is preserved
func TestRenderRulesToFormat_SubdirectoryStructure(t *testing.T) {
	// Create temporary source directory with subdirectories
	sourceDir := createTempSourceDir(t)
	defer os.RemoveAll(sourceDir)

	// Create subdirectory structure
	subDir := filepath.Join(sourceDir, "frontend", "react")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test rule in subdirectory
	testRuleContent := `---
alwaysApply: true
description: React component rule
---

# React Component Rule

Rules for React components.
`
	testRuleFile := filepath.Join(subDir, "components.md")
	if err := os.WriteFile(testRuleFile, []byte(testRuleContent), 0644); err != nil {
		t.Fatalf("Failed to create test rule file: %v", err)
	}

	// Change to temporary directory for the test
	origDir := getCurrentDir(t)
	tmpTestDir := createTempTestDir(t)
	defer func() {
		os.Chdir(origDir)
		os.RemoveAll(tmpTestDir)
	}()
	os.Chdir(tmpTestDir)

	// Run RenderRulesToFormat for continue format
	err := RenderRulesToFormat(sourceDir, "continue", false)
	if err != nil {
		t.Fatalf("RenderRulesToFormat failed: %v", err)
	}

	// Verify subdirectory structure is preserved
	expectedFile := filepath.Join(".continue", "rules", "frontend", "react", "components.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file with subdirectory structure %s was not created", expectedFile)
	}

	// Verify file content
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	if !strings.Contains(string(content), "React Component Rule") {
		t.Errorf("File content was not preserved correctly")
	}
}

// TestRenderRulesToFormat_WindsurfTriggerTransformation tests Windsurf's specific trigger transformation
func TestRenderRulesToFormat_WindsurfTriggerTransformation(t *testing.T) {
	tests := []struct {
		name            string
		alwaysApply     interface{}
		expectedTrigger string
	}{
		{
			name:            "alwaysApply true becomes trigger always_on",
			alwaysApply:     true,
			expectedTrigger: "always_on",
		},
		{
			name:            "alwaysApply false becomes trigger manual",
			alwaysApply:     false,
			expectedTrigger: "manual",
		},
		{
			name:            "no alwaysApply becomes trigger manual",
			alwaysApply:     nil,
			expectedTrigger: "manual",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary source directory
			sourceDir := createTempSourceDir(t)
			defer os.RemoveAll(sourceDir)

			// Create test rule content based on test case
			var testRuleContent string
			if tt.alwaysApply != nil {
				testRuleContent = fmt.Sprintf(`---
alwaysApply: %v
description: Test rule
globs: "**/*.tsx"
---

# Test Rule

Test content.
`, tt.alwaysApply)
			} else {
				testRuleContent = `---
description: Test rule  
globs: "**/*.tsx"
---

# Test Rule

Test content.
`
			}

			testRuleFile := filepath.Join(sourceDir, "test-rule.md")
			if err := os.WriteFile(testRuleFile, []byte(testRuleContent), 0644); err != nil {
				t.Fatalf("Failed to create test rule file: %v", err)
			}

			// Change to temporary directory for the test
			origDir := getCurrentDir(t)
			tmpTestDir := createTempTestDir(t)
			defer func() {
				os.Chdir(origDir)
				os.RemoveAll(tmpTestDir)
			}()
			os.Chdir(tmpTestDir)

			// Run RenderRulesToFormat for windsurf
			err := RenderRulesToFormat(sourceDir, "windsurf", false)
			if err != nil {
				t.Fatalf("RenderRulesToFormat failed: %v", err)
			}

			// Verify file content has correct trigger value
			expectedFile := filepath.Join(".windsurf", "rules", "test-rule.md")
			content, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			contentStr := string(content)
			expectedTriggerLine := fmt.Sprintf("trigger: %s", tt.expectedTrigger)
			if !strings.Contains(contentStr, expectedTriggerLine) {
				t.Errorf("Expected trigger transformation not found.\nExpected: %s\nActual content:\n%s", expectedTriggerLine, contentStr)
			}

			// Verify alwaysApply is not present
			if strings.Contains(contentStr, "alwaysApply:") {
				t.Errorf("alwaysApply should be transformed to trigger, not preserved")
			}
		})
	}
}

// TestRenderRulesToFormat_CopilotGlobsTransformation tests Copilot's globs to applyTo transformation
func TestRenderRulesToFormat_CopilotGlobsTransformation(t *testing.T) {
	tests := []struct {
		name            string
		globs           string
		expectedApplyTo string
	}{
		{
			name:            "globs becomes applyTo",
			globs:           "**/*.tsx",
			expectedApplyTo: "**/*.tsx",
		},
		{
			name:            "no globs defaults to all files",
			globs:           "",
			expectedApplyTo: "**",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary source directory
			sourceDir := createTempSourceDir(t)
			defer os.RemoveAll(sourceDir)

			// Create test rule content based on test case
			var testRuleContent string
			if tt.globs != "" {
				testRuleContent = fmt.Sprintf(`---
description: Test rule
globs: "%s"
---

# Test Rule

Test content.
`, tt.globs)
			} else {
				testRuleContent = `---
description: Test rule
---

# Test Rule

Test content.
`
			}

			testRuleFile := filepath.Join(sourceDir, "test-rule.md")
			if err := os.WriteFile(testRuleFile, []byte(testRuleContent), 0644); err != nil {
				t.Fatalf("Failed to create test rule file: %v", err)
			}

			// Change to temporary directory for the test
			origDir := getCurrentDir(t)
			tmpTestDir := createTempTestDir(t)
			defer func() {
				os.Chdir(origDir)
				os.RemoveAll(tmpTestDir)
			}()
			os.Chdir(tmpTestDir)

			// Run RenderRulesToFormat for copilot
			err := RenderRulesToFormat(sourceDir, "copilot", false)
			if err != nil {
				t.Fatalf("RenderRulesToFormat failed: %v", err)
			}

			// Verify file content has correct applyTo value
			expectedFile := filepath.Join(".github", "instructions", "test-rule.instructions.md")
			content, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			contentStr := string(content)
			expectedApplyToLine := fmt.Sprintf("applyTo: '%s'", tt.expectedApplyTo)
			if !strings.Contains(contentStr, expectedApplyToLine) {
				t.Errorf("Expected applyTo transformation not found.\nExpected: %s\nActual content:\n%s", expectedApplyToLine, contentStr)
			}

			// Verify globs is not present
			if strings.Contains(contentStr, "globs:") {
				t.Errorf("globs should be transformed to applyTo, not preserved")
			}
		})
	}
}

// TestRenderRulesToFormat_ErrorConditions tests error conditions
func TestRenderRulesToFormat_ErrorConditions(t *testing.T) {
	t.Run("Non-existent source directory", func(t *testing.T) {
		nonExistentDir := "/tmp/non-existent-dir-" + t.Name()

		// Change to temporary directory for the test
		origDir := getCurrentDir(t)
		tmpTestDir := createTempTestDir(t)
		defer func() {
			os.Chdir(origDir)
			os.RemoveAll(tmpTestDir)
		}()
		os.Chdir(tmpTestDir)

		err := RenderRulesToFormat(nonExistentDir, "continue", false)
		if err == nil {
			t.Errorf("Expected error for non-existent source directory, but got nil")
		}

		expectedErrorMsg := "source directory"
		if !strings.Contains(err.Error(), expectedErrorMsg) {
			t.Errorf("Expected error message to contain '%s', got: %v", expectedErrorMsg, err)
		}
	})
}

// TestRenderRulesToFormat_OnlyMarkdownFiles tests that only .md files are processed
func TestRenderRulesToFormat_OnlyMarkdownFiles(t *testing.T) {
	// Create temporary source directory
	sourceDir := createTempSourceDir(t)
	defer os.RemoveAll(sourceDir)

	// Create various file types
	files := map[string]string{
		"rule.md":     "# Valid Rule\nThis should be processed.",
		"readme.txt":  "This should be ignored.",
		"config.json": `{"key": "value"}`,
		"script.sh":   "#!/bin/bash\necho 'ignored'",
	}

	for filename, content := range files {
		filePath := filepath.Join(sourceDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Change to temporary directory for the test
	origDir := getCurrentDir(t)
	tmpTestDir := createTempTestDir(t)
	defer func() {
		os.Chdir(origDir)
		os.RemoveAll(tmpTestDir)
	}()
	os.Chdir(tmpTestDir)

	// Run RenderRulesToFormat
	err := RenderRulesToFormat(sourceDir, "continue", false)
	if err != nil {
		t.Fatalf("RenderRulesToFormat failed: %v", err)
	}

	// Verify only .md file was processed
	expectedFile := filepath.Join(".continue", "rules", "rule.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected markdown file %s was not created", expectedFile)
	}

	// Verify non-markdown files were not processed
	unexpectedFiles := []string{
		filepath.Join(".continue", "rules", "readme.txt"),
		filepath.Join(".continue", "rules", "config.json"),
		filepath.Join(".continue", "rules", "script.sh"),
	}

	for _, unexpectedFile := range unexpectedFiles {
		if _, err := os.Stat(unexpectedFile); !os.IsNotExist(err) {
			t.Errorf("Non-markdown file %s should not have been processed", unexpectedFile)
		}
	}
}

// Helper functions

func createTempSourceDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "rules-test-src-")
	if err != nil {
		t.Fatalf("Failed to create temporary source directory: %v", err)
	}
	return tmpDir
}

func createTempTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "rules-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary test directory: %v", err)
	}
	return tmpDir
}

func getCurrentDir(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	return dir
}
