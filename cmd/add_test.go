package cmd

import (
	"testing"
)

func TestParseRuleIdentifier(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *RuleIdentifier
		hasError bool
	}{
		{
			name:  "GitHub repo root",
			input: "gh:owner/repo",
			expected: &RuleIdentifier{
				OwnerSlug: "gh:owner",
				RepoName:  "repo",
				RuleSlug:  "repo",
				SubPath:   "",
				Version:   "latest",
				FullName:  "gh:owner/repo",
			},
			hasError: false,
		},
		{
			name:  "GitHub repo with subfolder",
			input: "gh:owner/repo/path/to/folder",
			expected: &RuleIdentifier{
				OwnerSlug: "gh:owner",
				RepoName:  "repo",
				RuleSlug:  "folder",
				SubPath:   "path/to/folder",
				Version:   "latest",
				FullName:  "gh:owner/repo/path/to/folder",
			},
			hasError: false,
		},
		{
			name:  "GitHub repo with single subfolder",
			input: "gh:owner/repo/rules",
			expected: &RuleIdentifier{
				OwnerSlug: "gh:owner",
				RepoName:  "repo",
				RuleSlug:  "rules",
				SubPath:   "rules",
				Version:   "latest",
				FullName:  "gh:owner/repo/rules",
			},
			hasError: false,
		},
		{
			name:  "GitHub repo with version",
			input: "gh:owner/repo@v1.0.0",
			expected: &RuleIdentifier{
				OwnerSlug: "gh:owner",
				RepoName:  "repo",
				RuleSlug:  "repo",
				SubPath:   "",
				Version:   "v1.0.0",
				FullName:  "gh:owner/repo@v1.0.0",
			},
			hasError: false,
		},
		{
			name:  "GitHub repo with subfolder and version",
			input: "gh:owner/repo/path/to/folder@v1.0.0",
			expected: &RuleIdentifier{
				OwnerSlug: "gh:owner",
				RepoName:  "repo",
				RuleSlug:  "folder",
				SubPath:   "path/to/folder",
				Version:   "v1.0.0",
				FullName:  "gh:owner/repo/path/to/folder@v1.0.0",
			},
			hasError: false,
		},
		{
			name:     "Invalid GitHub format",
			input:    "gh:owner",
			expected: nil,
			hasError: true,
		},
		{
			name:  "Registry rule",
			input: "owner/rule",
			expected: &RuleIdentifier{
				OwnerSlug: "owner",
				RuleSlug:  "rule",
				Version:   "latest",
				FullName:  "owner/rule",
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseRuleIdentifier(tc.input)
			
			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tc.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tc.input, err)
				return
			}
			
			if result.OwnerSlug != tc.expected.OwnerSlug {
				t.Errorf("OwnerSlug mismatch for %s: got %s, expected %s", tc.input, result.OwnerSlug, tc.expected.OwnerSlug)
			}
			
			if result.RuleSlug != tc.expected.RuleSlug {
				t.Errorf("RuleSlug mismatch for %s: got %s, expected %s", tc.input, result.RuleSlug, tc.expected.RuleSlug)
			}
			
			if result.Version != tc.expected.Version {
				t.Errorf("Version mismatch for %s: got %s, expected %s", tc.input, result.Version, tc.expected.Version)
			}
			
			if result.FullName != tc.expected.FullName {
				t.Errorf("FullName mismatch for %s: got %s, expected %s", tc.input, result.FullName, tc.expected.FullName)
			}
			
			if result.SubPath != tc.expected.SubPath {
				t.Errorf("SubPath mismatch for %s: got %s, expected %s", tc.input, result.SubPath, tc.expected.SubPath)
			}
			
			if result.RepoName != tc.expected.RepoName {
				t.Errorf("RepoName mismatch for %s: got %s, expected %s", tc.input, result.RepoName, tc.expected.RepoName)  
			}
		})
	}
}