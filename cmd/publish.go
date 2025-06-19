package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rules-cli/internal/auth"
	"rules-cli/internal/registry"
	"rules-cli/internal/ruleset"
	"rules-cli/internal/validation"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	visibility     string
	packageVersion string
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish [path]",
	Short: "Publish a rule package to the registry",
	Long: `Publishes a rule package to the registry as a zip file.

This command requires authentication and will prompt you to login if you're not already.
The slug is automatically determined from rules.json in the current directory or specified path.

The visibility can be set to "public" (default) or "private".
The version can be specified with --version flag, defaults to timestamp-based version.

Examples:
  rules publish                           # Publish from current directory
  rules publish ./my-rules                # Publish from specified directory
  rules publish --visibility private      # Publish as private
  rules publish --version 1.0.0          # Publish with specific version`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPublishCommand,
}

// createPackageZip creates a zip file containing all the rule files
func createPackageZip(rulesPath string, tempDir string) (string, error) {
	// Create temporary zip file
	zipFileName := fmt.Sprintf("package-%d.zip", time.Now().Unix())
	zipPath := filepath.Join(tempDir, zipFileName)
	
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// Determine the source directory
	sourceDir := "."
	if rulesPath != "" {
		if stat, err := os.Stat(rulesPath); err == nil && stat.IsDir() {
			sourceDir = rulesPath
		} else {
			sourceDir = filepath.Dir(rulesPath)
		}
	}
	
	// Walk through the directory and add files to zip
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Skip hidden files and directories
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		
		// Skip temporary files
		if strings.HasSuffix(info.Name(), ".tmp") || strings.HasSuffix(info.Name(), "~") {
			return nil
		}
		
		// Get relative path from source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		
		// Create file in zip
		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		
		// Open source file
		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()
		
		// Copy file content to zip
		_, err = io.Copy(zipFileWriter, sourceFile)
		return err
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to create package zip: %w", err)
	}
	
	return zipPath, nil
}

// runPublishCommand implements the main logic for the publish command
func runPublishCommand(cmd *cobra.Command, args []string) error {
	// Validate the visibility
	if visibility != "public" && visibility != "private" {
		return fmt.Errorf("visibility must be either 'public' or 'private'")
	}

	// Determine the path to look for rules.json
	var rulesPath string
	if len(args) > 0 {
		rulesPath = args[0]
	}

	// Determine the rules.json file path
	var rulesJSONPath string
	if rulesPath != "" {
		// If a path was specified, look for rules.json in that directory
		stat, err := os.Stat(rulesPath)
		if err == nil && stat.IsDir() {
			rulesJSONPath = filepath.Join(rulesPath, "rules.json")
		} else {
			rulesJSONPath = rulesPath
		}
	} else {
		// Look in current directory for rules.json
		rulesJSONPath = "rules.json"
	}

	// Validate the rules.json file against the schema FIRST
	color.Cyan("Validating rules.json against schema...")
	if err := validation.ValidateRulesJSONFromFile(rulesJSONPath); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}
	color.Green("✓ rules.json is valid")

	// Load ruleset from the specified path or current directory
	rs, err := ruleset.LoadRuleSetFromPath(rulesPath)
	if err != nil {
		return fmt.Errorf("failed to load rules.json: %w", err)
	}

	// Validate that the ruleset has a name
	if rs.Name == "" {
		return fmt.Errorf("rules.json must have a 'name' field")
	}

	// NOW ensure the user is authenticated (after validation passes)
	authenticated, err := auth.EnsureAuthenticated(true)
	if err != nil || !authenticated {
		return fmt.Errorf("authentication required to publish rules")
	}

	// Create registry client and get user info
	authConfig := auth.LoadAuthConfig()
	client := registry.NewClient(cfg.RegistryURL)
	client.SetAuthToken(authConfig.AccessToken)

	// Generate version if not specified
	packageVersion = rs.Version;

	// Create temporary directory for package creation
	tempDir, err := ioutil.TempDir("", "rules-publish-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create package zip file
	color.Cyan("Creating package zip file...")
	zipPath, err := createPackageZip(rulesPath, tempDir)
	if err != nil {
		return fmt.Errorf("failed to create package: %w", err)
	}
	defer os.Remove(zipPath)

	// Publish the rule package
	color.Cyan("Publishing package to %s (version %s) with visibility: %s", rs.Name, packageVersion, visibility)
	err = client.PublishRule(rs.Name, packageVersion, zipPath, visibility)
	if err != nil {
		return fmt.Errorf("failed to publish rule: %w", err)
	}

	color.Green("Successfully published package '%s/%s' (version %s)", rs.Name, packageVersion)
	return nil
}

// isValidSlug checks if a slug is valid (alphanumeric, hyphens, underscores only)
func isValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	// Check if slug contains only valid characters
	for _, char := range slug {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}

	return true
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// Add flags
	publishCmd.Flags().StringVar(&visibility, "visibility", "public", "Set the visibility of the rule to 'public' or 'private'")
	publishCmd.Flags().StringVar(&packageVersion, "version", "", "Version for the package (defaults to timestamp-based version)")
}