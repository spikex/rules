package formats

import (
	"fmt"
	"os"
	"path/filepath"
)

// RenderOptions holds configuration options for the rendering process
type RenderOptions struct {
	SourceDir    string
	TargetFormat Format
	Verbose      bool
}

// RenderRulesToFormat renders rules from the source directory to the target format
// This is a higher-level function that sets up the rendering process
func RenderRulesToFormat(sourceDir string, targetFormatName string, verbose bool) error {
	// Get the target format
	targetFormat := GetFormat(targetFormatName)
	
	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory %s does not exist", sourceDir)
	}
	
	// For single file formats, make sure the parent directory exists
	if targetFormat.IsSingleFile {
		parentDir := filepath.Dir(targetFormat.SingleFilePath)
		if parentDir != "." && parentDir != "" {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for single file format: %w", err)
			}
		}
	} else {
		// For directory-based formats, create the target directory
		if err := os.MkdirAll(targetFormat.DirectoryPrefix, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetFormat.DirectoryPrefix, err)
		}
	}
	
	// Process the rule files
	return ProcessRuleFiles(sourceDir, targetFormat)
}