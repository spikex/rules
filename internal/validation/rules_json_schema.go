package validation

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
			friendlyMsg := getFriendlyErrorMessage(validationError)
			errorMsg += fmt.Sprintf("  - %s", friendlyMsg)
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

// getFriendlyErrorMessage converts technical JSON schema errors into user-friendly messages
func getFriendlyErrorMessage(validationError gojsonschema.ResultError) string {
	field := validationError.Field()
	errorType := validationError.Type()
	description := validationError.Description()
	
	// Handle specific pattern validation errors for well-known fields
	if errorType == "pattern" {
		if strings.Contains(field, "name") {
			// This is the ruleset name pattern validation
			if strings.Contains(description, "^[a-zA-Z0-9]([a-zA-Z0-9-_]*[a-zA-Z0-9])?/[a-zA-Z0-9]([a-zA-Z0-9-_]*[a-zA-Z0-9])?$") {
				return fmt.Sprintf("name: Invalid format. Expected format: 'owner/ruleset' (e.g., 'acme/web-security'). "+
					"Both owner and ruleset must start and end with alphanumeric characters and may contain hyphens or underscores in the middle")
			}
		}
		
		if strings.Contains(field, "rules") {
			// This is for rule names in the rules object
			if strings.Contains(description, "gh:[a-zA-Z0-9-._]+/[a-zA-Z0-9-._]+") {
				return fmt.Sprintf("rules: Invalid rule name format. Expected format: "+
					"'owner/rule' (e.g., 'acme/security-check')")
			}
		}
		
		if strings.Contains(field, "version") {
			// This is semantic version validation
			if strings.Contains(description, "semantic version") || strings.Contains(description, "0|[1-9]") {
				return fmt.Sprintf("version: Invalid semantic version format. Expected format: 'MAJOR.MINOR.PATCH' "+
					"(e.g., '1.0.0', '2.1.3') with optional pre-release and build metadata")
			}
		}
	}
	
	// Handle missing required fields
	if errorType == "required" {
		return fmt.Sprintf("Missing required field: %s", strings.Replace(description, "property", "field", 1))
	}
	
	// Handle type mismatches
	if errorType == "invalid_type" {
		return fmt.Sprintf("%s: %s", field, strings.Replace(description, "Invalid type", "Expected different type", 1))
	}
	
	// Handle additional properties not allowed (often means pattern didn't match)
	if errorType == "additional_property_not_allowed" && strings.Contains(field, "rules") {
		return fmt.Sprintf("rules: Invalid rule name format. Expected format: "+
			"'owner/rule' (e.g., 'acme/security-check'). "+
			"Rule names cannot contain spaces or special characters except hyphens, underscores, and dots")
	}
	
	// For any other errors, return the original message but clean it up slightly
	return fmt.Sprintf("%s: %s", strings.TrimPrefix(field, "(root)."), description)
}
