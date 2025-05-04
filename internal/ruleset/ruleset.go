package ruleset

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
		Author:      "Anonymous",
		License:     "Apache-2.0",
		Version:     "1.0.0",
		Rules:       make(map[string]string),
	}
}

// LoadRuleSet loads a ruleset from a specified path
func LoadRuleSet(path string) (*RuleSet, error) {
	data, err := ioutil.ReadFile(path)
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

	return ioutil.WriteFile(path, data, 0644)
}

// AddRule adds a rule to the ruleset with the given name and version
func (rs *RuleSet) AddRule(name, version string) {
	rs.Rules[name] = version
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

// CreateRule writes a rule to a markdown file
func CreateRule(rule Rule, format, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("rule name cannot be empty")
	}

	// Create content with front matter
	content := "---\n"
	if rule.Description != "" {
		content += "description: " + rule.Description + "\n"
	}
	if len(rule.Tags) > 0 {
		content += "tags: [" + strings.Join(rule.Tags, ", ") + "]\n"
	}
	if rule.Globs != "" {
		content += "globs: " + rule.Globs + "\n"
	}
	if rule.AlwaysApply {
		content += "alwaysApply: true\n"
	}
	content += "---\n\n" + rule.Body

	// Determine directory based on format
	var ruleDir string
	if format == "default" {
		ruleDir = ".rules"
	} else {
		ruleDir = "." + format + "/rules"
	}

	// Ensure the directory exists
	if err := os.MkdirAll(ruleDir, 0755); err != nil {
		return err
	}

	// Create the rule file
	fileName := filepath.Join(ruleDir, name+".md")
	return ioutil.WriteFile(fileName, []byte(content), 0644)
}