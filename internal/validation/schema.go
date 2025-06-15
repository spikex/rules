package validation

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema/rules-schema.json
var schemaFS embed.FS

// ValidateRulesJSON validates a rules.json file against the JSON schema
func ValidateRulesJSON(rulesData []byte) error {
	// Load the schema from embedded file
	schemaBytes, err := schemaFS.ReadFile("schema/rules-schema.json")
	if err != nil {
		return fmt.Errorf("failed to load embedded schema: %w", err)
	}

	// Create schema loader
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)

	// Create document loader with the rules data
	documentLoader := gojsonschema.NewBytesLoader(rulesData)

	// Validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed to validate JSON schema: %w", err)
	}

	if !result.Valid() {
		var errorMsg string
		for i, validationError := range result.Errors() {
			if i > 0 {
				errorMsg += "\n"
			}
			errorMsg += fmt.Sprintf("  - %s", validationError.String())
		}
		return fmt.Errorf("rules.json validation failed:\n%s", errorMsg)
	}

	return nil
}

// ValidateRulesJSONFromFile validates a rules.json file from disk
func ValidateRulesJSONFromFile(filePath string) error {
	// Read the file content
	rulesData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read rules file: %w", err)
	}

	return ValidateRulesJSON(rulesData)
}

// ValidateRulesObject validates a rules object directly
func ValidateRulesObject(rules interface{}) error {
	// Convert the rules object to JSON bytes
	rulesBytes, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("failed to marshal rules object: %w", err)
	}

	return ValidateRulesJSON(rulesBytes)
}
