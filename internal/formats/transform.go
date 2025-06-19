package formats

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// EmptyYAMLValue represents an empty YAML value that renders as just a space after the colon
type EmptyYAMLValue struct{}

func (e EmptyYAMLValue) MarshalYAML() (interface{}, error) {
	var node yaml.Node
	node.Kind = yaml.ScalarNode
	node.Value = ""
	node.Tag = ""
	return &node, nil
}

// RuleMetadata represents the frontmatter metadata of a rule file
type RuleMetadata map[string]interface{}

// ParseFrontmatter extracts the frontmatter from a markdown file
func ParseFrontmatter(content []byte) (RuleMetadata, []byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var frontmatterLines []string
	var bodyLines []string
	
	// Check if the file starts with a frontmatter delimiter
	if !scanner.Scan() || scanner.Text() != "---" {
		// No frontmatter, return the content as is
		return RuleMetadata{}, content, nil
	}
	
	// Read frontmatter lines
	inFrontmatter := true
	for scanner.Scan() {
		line := scanner.Text()
		
		if line == "---" {
			inFrontmatter = false
			continue
		}
		
		if inFrontmatter {
			frontmatterLines = append(frontmatterLines, line)
		} else {
			bodyLines = append(bodyLines, line)
		}
	}
	
	// Parse frontmatter as YAML
	frontmatterStr := strings.Join(frontmatterLines, "\n")
	metadata := RuleMetadata{}
	
	if err := yaml.Unmarshal([]byte(frontmatterStr), &metadata); err != nil {
		return nil, nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}
	
	// Join body lines
	bodyContent := []byte(strings.Join(bodyLines, "\n"))
	
	return metadata, bodyContent, nil
}

// TransformRuleContent transforms the content of a rule file based on the target format
func TransformRuleContent(content []byte, format Format) ([]byte, error) {
	// Parse frontmatter and content
	metadata, bodyContent, err := ParseFrontmatter(content)
	if err != nil {
		return nil, err
	}

	var trimmedBodyContent = bytes.TrimSpace(bodyContent)
	
	// Transform metadata based on format
	transformedMetadata, err := TransformMetadata(metadata, format)
	if err != nil {
		return nil, err
	}
	
	// If no metadata (or empty), just return the trimmed body content
	if len(transformedMetadata) == 0 {
		return trimmedBodyContent, nil
	}
	
	// Clean up empty/nil values before serialization
	cleanedMetadata := cleanMetadataForYAML(transformedMetadata)

	// Serialize metadata to YAML
	metadataBytes, err := yaml.Marshal(cleanedMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize metadata: %w", err)
	}
	
	// Combine frontmatter and content
	var result bytes.Buffer
	result.WriteString("---\n")
	result.Write(metadataBytes)
	result.WriteString("---\n\n")
	result.Write(trimmedBodyContent)
	
	return result.Bytes(), nil
}

// cleanMetadataForYAML converts nil/empty values to EmptyYAMLValue for clean YAML output
func cleanMetadataForYAML(metadata RuleMetadata) RuleMetadata {
	cleaned := RuleMetadata{}

	for key, value := range metadata {
		switch v := value.(type) {
		case nil:
			cleaned[key] = EmptyYAMLValue{}
		case string:
			if v == "" {
				cleaned[key] = EmptyYAMLValue{}
			} else {
				cleaned[key] = v
			}
		case *string:
			if v == nil || *v == "" {
				cleaned[key] = EmptyYAMLValue{}
			} else {
				cleaned[key] = *v
			}
		default:
			cleaned[key] = v
		}
	}
		
	return cleaned
}

// TransformMetadata transforms rule metadata based on the target format
func TransformMetadata(metadata RuleMetadata, format Format) (RuleMetadata, error) {
	// Create a copy of the metadata
	transformed := RuleMetadata{}
	for k, v := range metadata {
		transformed[k] = v
	}
	
	// Transform based on format
	switch format.Name {
	case "default", "continue", "cursor":
		// Keep alwaysApply, description, and globs
		for k := range transformed {
			if k != "alwaysApply" && k != "description" && k != "globs" {
				delete(transformed, k)
			}
		}
		
	case "windsurf":
		// Transform alwaysApply to trigger
		if alwaysApply, ok := metadata["alwaysApply"]; ok {
			delete(transformed, "alwaysApply")
			
			// If alwaysApply is true, set trigger to always_on
			if alwaysApply == true {
				transformed["trigger"] = "always_on"
			} else {
			// If alwaysApply is false, set trigger to manual
				transformed["trigger"] = "manual"
			}
		} else {
			// Default to manual trigger
			transformed["trigger"] = "manual"
		}
		
		// Keep description and globs
		for k := range transformed {
			if k != "trigger" && k != "description" && k != "globs" {
				delete(transformed, k)
			}
		}
		
	case "copilot":
		// Transform globs to applyTo
		if globs, ok := metadata["globs"]; ok {
			delete(transformed, "globs")
			transformed["applyTo"] = globs
		} else {
			// Default to all files
			transformed["applyTo"] = "**"
		}
		
		// Remove other fields except description
		for k := range transformed {
			if k != "applyTo" && k != "description" {
				delete(transformed, k)
			}
		}
		
	case "cline", "cody":
		// For these formats, we only keep description
		for k := range transformed {
			if k != "description" {
				delete(transformed, k)
			}
		}
		
	case "claude", "codex", "amp":
		// For single file formats, we don't include any frontmatter
		return RuleMetadata{}, nil
	}
	
	return transformed, nil
}

// IsRuleApplicable checks if a rule should be included in the target format
func IsRuleApplicable(content []byte, format Format) (bool, error) {
	// For single file formats, only include rules with alwaysApply: true
	if format.IsSingleFile {
		return isAlwaysApply(content), nil
	}
	
	// For other formats, include all rules
	return true, nil
}

// ProcessRuleFiles processes all rule files in the source directory and renders them to the target format
func ProcessRuleFiles(sourceDir string, targetFormat Format) error {
	// For single file formats, we need to gather all rules with alwaysApply: true
	if targetFormat.IsSingleFile {
		return renderToSingleFile(sourceDir, targetFormat)
	}
	
	// For directory-based formats, we need to create the target directory and process each file
	targetDir := targetFormat.DirectoryPrefix
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}
	
	// Walk through all files in the source directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Only process Markdown files
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		
		// Check if the rule should be included in the target format
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read source file %s: %w", path, err)
		}
		
		applicable, err := IsRuleApplicable(content, targetFormat)
		if err != nil {
			return err
		}
		
		if !applicable {
			// Skip this rule
			return nil
		}
		
		// Get relative path from source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		
		// Create target path with appropriate extension
		targetPath := filepath.Join(targetDir, strings.TrimSuffix(relPath, ".md") + targetFormat.FileExtension)
		
		// Ensure target directory exists
		targetFileDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetFileDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetFileDir, err)
		}
		
		// Transform the content based on the format
		transformedContent, err := TransformRuleContent(content, targetFormat)
		if err != nil {
			return fmt.Errorf("failed to transform content: %w", err)
		}
		
		// Write to target file
		if err := os.WriteFile(targetPath, transformedContent, 0644); err != nil {
			return fmt.Errorf("failed to write to target file %s: %w", targetPath, err)
		}
		
		return nil
	})
}

// GetRuleName extracts the rule name from a file path
func GetRuleName(filePath string, sourceDir string) (string, error) {
	// Get relative path from source directory
	relPath, err := filepath.Rel(sourceDir, filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}
	
	// Remove extension
	ruleName := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	
	// Replace directory separators with slashes for consistency
	ruleName = strings.ReplaceAll(ruleName, string(os.PathSeparator), "/")
	
	return ruleName, nil
}

// ExtractRuleTitle extracts the title from the content of a rule file
func ExtractRuleTitle(content []byte) string {
	// Look for a markdown heading (# Title)
	re := regexp.MustCompile(`(?m)^#\s+(.+)$`)
	match := re.FindSubmatch(content)
	
	if len(match) > 1 {
		return string(match[1])
	}
	
	// If no heading found, return empty string
	return ""
}