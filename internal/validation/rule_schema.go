package validation

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// RuleFrontmatterSchema contains the embedded JSON schema for rule frontmatter
//
//go:embed rule_frontmatter_schema.json
var RuleFrontmatterSchema string

// RuleFrontmatter represents the structure of rule frontmatter
type RuleFrontmatter struct {
	AlwaysApply bool     `yaml:"alwaysApply,omitempty" json:"alwaysApply,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Globs       string   `yaml:"globs,omitempty" json:"globs,omitempty"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field      string      `json:"field"`
	Message    string      `json:"message"`
	Value      interface{} `json:"value,omitempty"`
	Type       string      `json:"type,omitempty"`
	SchemaPath string      `json:"schemaPath,omitempty"`
	Context    string      `json:"context,omitempty"`
}

// ValidationResult contains the validation results
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

var (
	// compiledSchema holds the compiled JSON schema for reuse
	compiledSchema *gojsonschema.Schema
)

// init compiles the schema once at startup
func init() {
	schemaLoader := gojsonschema.NewStringLoader(RuleFrontmatterSchema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		panic(fmt.Errorf("failed to compile frontmatter schema: %w", err))
	}
	compiledSchema = schema
}

// ValidateRuleFile validates a markdown rule file's frontmatter against the schema
func ValidateRuleFile(filePath string) (*ValidationResult, error) {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse frontmatter
	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "frontmatter",
					Message: fmt.Sprintf("Failed to parse frontmatter: %v", err),
				},
			},
		}, nil
	}

	// If no frontmatter found, consider it valid (rules without frontmatter are allowed)
	if frontmatter == nil {
		return &ValidationResult{Valid: true}, nil
	}

	return ValidateFrontmatter(frontmatter)
}

// ValidateFrontmatter validates frontmatter data against the JSON schema
func ValidateFrontmatter(frontmatter map[string]interface{}) (*ValidationResult, error) {
	// Create a document loader from the frontmatter data
	documentLoader := gojsonschema.NewGoLoader(frontmatter)

	// Validate against the compiled schema
	result, err := compiledSchema.Validate(documentLoader)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	validationResult := &ValidationResult{
		Valid:  result.Valid(),
		Errors: make([]ValidationError, 0),
	}

	// Convert validation errors to our format
	if !result.Valid() {
		for _, err := range result.Errors() {
			validationErr := ValidationError{
				Field:   extractFieldName(err.Field()),
				Message: err.Description(),
				Type:    err.Type(),
				Context: err.Context().String(),
				Value:   err.Value(),
			}

			// Safely extract schema path from details
			if details := err.Details(); details != nil {
				if property, ok := details["property"]; ok && property != nil {
					validationErr.SchemaPath = fmt.Sprintf("%v", property)
				}
			}

			validationResult.Errors = append(validationResult.Errors, validationErr)
		}
	}

	return validationResult, nil
}

// extractFieldName cleans up the field name for better display
func extractFieldName(field string) string {
	// Remove (root). prefix and return clean field name
	cleanField := strings.TrimPrefix(field, "(root).")
	if cleanField == "(root)" {
		return "root"
	}
	return cleanField
}

// extractFrontmatter extracts YAML frontmatter from markdown content
func extractFrontmatter(content []byte) (map[string]interface{}, error) {
	lines := string(content)

	// Check if content starts with frontmatter delimiter
	if len(lines) < 4 || lines[:4] != "---\n" {
		return nil, nil // No frontmatter
	}

	// Find the end of frontmatter
	endDelimiter := "\n---\n"
	endIndex := strings.Index(lines[4:], endDelimiter)
	if endIndex == -1 {
		return nil, fmt.Errorf("unterminated frontmatter: missing closing ---")
	}

	// Extract frontmatter content
	frontmatterContent := lines[4 : 4+endIndex]
	if strings.TrimSpace(frontmatterContent) == "" {
		return nil, nil // Empty frontmatter
	}

	// Parse YAML
	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatterContent), &frontmatter); err != nil {
		return nil, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}

	return frontmatter, nil
}

// GetSchema returns the JSON schema as a string for external use
func GetSchema() string {
	return RuleFrontmatterSchema
}

// ValidateJSON validates JSON data directly against the schema
func ValidateJSON(data []byte) (*ValidationResult, error) {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "json",
					Message: fmt.Sprintf("Invalid JSON: %v", err),
				},
			},
		}, nil
	}

	if jsonMap, ok := jsonData.(map[string]interface{}); ok {
		return ValidateFrontmatter(jsonMap)
	}

	return &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{
				Field:   "root",
				Message: "JSON data must be an object",
			},
		},
	}, nil
}

// ValidateYAML validates YAML data directly against the schema
func ValidateYAML(data []byte) (*ValidationResult, error) {
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "yaml",
					Message: fmt.Sprintf("Invalid YAML: %v", err),
				},
			},
		}, nil
	}

	return ValidateFrontmatter(yamlData)
}