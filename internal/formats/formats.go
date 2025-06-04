package formats

import (
	"fmt"
	"os"
)

// Format represents a rules format
type Format struct {
	Name            string
	DirectoryPrefix string
}

// GetFormat returns a Format for the given format name
func GetFormat(formatName string) Format {
	// Default to ".rules" if no format is specified or format is "default"
	if formatName == "" || formatName == "default" {
		return Format{
			Name:            "default",
			DirectoryPrefix: ".rules",
		}
	}
	
	// Otherwise use .<format>/rules
	return Format{
		Name:            formatName,
		DirectoryPrefix: fmt.Sprintf(".%s/rules", formatName),
	}
}
	
// InitializeFormat creates the directory structure for a format
func InitializeFormat(formatName string) error {
	format := GetFormat(formatName)
	
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
	format := GetFormat(formatName)
	return format.DirectoryPrefix, nil
}

// GetRulesJSONPath returns the path to rules.json file in the root directory
func GetRulesJSONPath(formatName string) (string, error) {
	// Return the path to rules.json in the root directory
	return "rules.json", nil
}