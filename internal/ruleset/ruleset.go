package ruleset

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RuleSet represents the rules.json file structure
type RuleSet struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	License     string            `json:"license"`
	Version     string            `json:"version"`
	Rules       map[string]string `json:"rules"`
}

// Rule represents a single rule with front matter and content
type Rule struct {
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Globs       string   `json:"globs,omitempty" yaml:"globs,omitempty"`
	AlwaysApply bool     `json:"alwaysApply,omitempty" yaml:"alwaysApply,omitempty"`
	Body        string   `json:"-" yaml:"-"`
}

// DefaultRuleSet creates a new ruleset with default values
func DefaultRuleSet(name string) *RuleSet {
	return &RuleSet{
		Name:        name,
		Description: "A ruleset for AI code assistants",
		Author:      "",
		License:     "CC0-1.0",
		Version:     "1.0.0",
		Rules:       make(map[string]string),
	}
}

// LoadRuleSet loads a ruleset from a specified path
func LoadRuleSet(path string) (*RuleSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rs RuleSet
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, err
	}

	return &rs, nil
}

// SaveRuleSet saves a ruleset to the specified path
func (rs *RuleSet) SaveRuleSet(path string) error {
	data, err := json.MarshalIndent(rs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// AddRule adds a rule to the ruleset with the given name and version
func (rs *RuleSet) AddRule(name, version string) {
	rs.Rules[name] = version
}

// RemoveRule removes a rule from the ruleset if it exists
func (rs *RuleSet) RemoveRule(name string) bool {
	if _, exists := rs.Rules[name]; exists {
		delete(rs.Rules, name)
		return true
	}
	return false
}

// RuleExists checks if a rule exists in the ruleset
func (rs *RuleSet) RuleExists(name string) bool {
	_, exists := rs.Rules[name]
	return exists
}

// GetRuleVersion gets the version of a rule if it exists
func (rs *RuleSet) GetRuleVersion(name string) (string, bool) {
	version, exists := rs.Rules[name]
	return version, exists
}

// CreateRule writes a rule to a markdown file and returns the file path
func CreateRule(rule Rule, format, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", errors.New("rule name cannot be empty")
	}

	// Create content with Continue format front matter
	content := "---\n"

	// Continue format: alwaysApply field (required in Continue format)
	if rule.AlwaysApply {
		content += "alwaysApply: true\n"
	} else {
		content += "alwaysApply: false\n"
	}

	// Continue format: description field
	if rule.Description != "" {
		content += "description: " + rule.Description + "\n"
	}

	// Continue format: globs field
	if rule.Globs != "" {
		content += "globs: \"" + rule.Globs + "\"\n"
	}

	// Note: Tags are not part of Continue format specification

	content += "---\n\n" + rule.Body

	// Create the rule file in the current directory
	fileName := name + ".md"
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// FindRuleSetFile looks for rules.json in the current directory or specified path
func FindRuleSetFile(path string) (string, error) {
	if path != "" {
		// If path is specified, check if it's a directory or file
		stat, err := os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("path does not exist: %w", err)
		}

		if stat.IsDir() {
			// If it's a directory, look for rules.json inside it
			rulesPath := filepath.Join(path, "rules.json")
			if _, err := os.Stat(rulesPath); err == nil {
				return rulesPath, nil
			}
			return "", fmt.Errorf("rules.json not found in directory: %s", path)
		} else {
			// If it's a file, check if it's rules.json
			if filepath.Base(path) == "rules.json" {
				return path, nil
			}
			return "", fmt.Errorf("specified file is not rules.json")
		}
	}

	// Look for rules.json in current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	rulesPath := filepath.Join(currentDir, "rules.json")
	if _, err := os.Stat(rulesPath); err == nil {
		return rulesPath, nil
	}

	return "", fmt.Errorf("rules.json not found in current directory")
}

// LoadRuleSetFromPath loads a ruleset from the current directory or specified path
func LoadRuleSetFromPath(path string) (*RuleSet, error) {
	rulesPath, err := FindRuleSetFile(path)
	if err != nil {
		return nil, err
	}

	return LoadRuleSet(rulesPath)
}
