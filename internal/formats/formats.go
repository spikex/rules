package formats

import (
	"fmt"
	"os"
	"path/filepath"
)

// Format represents a rules format
type Format struct {
	Name            string
	DirectoryPrefix string
	ConfigTemplate  string
}

// Known formats
var (
	DefaultFormat = Format{
		Name:            "default",
		DirectoryPrefix: ".rules",
		ConfigTemplate:  "",
	}
	
	CursorFormat = Format{
		Name:            "cursor",
		DirectoryPrefix: ".cursor/rules",
		ConfigTemplate:  "",
	}
	
	// Add more formats as needed
	Formats = map[string]Format{
		"default": DefaultFormat,
		"cursor":  CursorFormat,
	}
)

// InitializeFormat creates the directory structure for a format
func InitializeFormat(formatName string) error {
	format, exists := Formats[formatName]
	if !exists {
		return fmt.Errorf("unknown format: %s", formatName)
	}
	
	dirPath := format.DirectoryPrefix
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}
	
	// Create rules.json in the root directory if it doesn't exist
	rulesJSONPath := "rules.json"
	if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
		defaultRulesJSON := `{
  "name": "ruleset-name",
  "description": "Description of the ruleset",
  "author": "Author Name",
  "license": "Apache-2.0",
  "version": "1.0.0",
  "rules": {}
}`
		
		if err := os.WriteFile(rulesJSONPath, []byte(defaultRulesJSON), 0644); err != nil {
			return fmt.Errorf("failed to create rules.json: %w", err)
		}
	}
	
	return nil
}

// GetRulesDirectory returns the rules directory for a given format
func GetRulesDirectory(formatName string) (string, error) {
	format, exists := Formats[formatName]
	if !exists {
		return "", fmt.Errorf("unknown format: %s", formatName)
	}
	
	return format.DirectoryPrefix, nil
}

// GetRulesJSONPath returns the path to rules.json file in the root directory
func GetRulesJSONPath(formatName string) (string, error) {
	// Check if format exists
	if _, exists := Formats[formatName]; !exists {
		return "", fmt.Errorf("unknown format: %s", formatName)
	}
	
	// Return the path to rules.json in the root directory
	return "rules.json", nil
}