package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"rules-cli/internal/formats"

	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render [format]",
	Short: "Render rules to a specific format",
	Long: `Renders existing rules to a specified format.
Creates .{format}/rules/ directory and copies all rules from the default location
(.rules/) to the target format location. Preserves directory structure of rule sets.`,
	Example: `  rules render cursor
  rules render continue`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("format is required")
		}
		
		formatName := args[0]
		
		if formatName == "default" {
			return fmt.Errorf("cannot render to default format as it is the source")
		}

		// Get the source directory
		sourceDir, err := formats.GetRulesDirectory("default")
		if err != nil {
			return fmt.Errorf("failed to get source directory: %w", err)
		}

		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			return fmt.Errorf("source directory %s does not exist", sourceDir)
		}

		fmt.Printf("Rendering rules to %s format...\n", formatName)
		
		// Get target directory
		targetDir, err := formats.GetRulesDirectory(formatName)
		if err != nil {
			return fmt.Errorf("failed to get target directory for format %s: %w", formatName, err)
		}

		// Create target directory if it doesn't exist
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
		}

		// Copy rules from source to target
		if err := copyDir(sourceDir, targetDir); err != nil {
			return fmt.Errorf("failed to copy rules to target directory: %w", err)
		}

		fmt.Printf("Successfully rendered rules to %s format\n", formatName)

		return nil
	},
}

// copyDir recursively copies a directory tree, preserving directory structure
func copyDir(src, dst string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Get directory contents
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err = copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy the file
			if err = copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	// Skip if not a markdown file
	if !strings.HasSuffix(src, ".md") {
		return nil
	}

	// Open source file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy contents
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Flush to disk
	return out.Sync()
}

func init() {
	rootCmd.AddCommand(renderCmd)
}