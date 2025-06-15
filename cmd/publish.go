package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"rules-cli/internal/auth"
	"rules-cli/internal/registry"
	"rules-cli/internal/ruleset"
	"rules-cli/internal/validation"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	visibility string
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish [path]",
	Short: "Publish a rule file to the registry",
	Long: `Publishes a rule file to the registry.

This command requires authentication and will prompt you to login if you're not already.
The slug is automatically determined from rules.json in the current directory or specified path.

The visibility can be set to "public" (default) or "private".

Examples:
  rules publish                    # Publish from current directory
  rules publish ./my-rules         # Publish from specified directory
  rules publish --visibility private`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPublishCommand,
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

	// Get user information to determine the organization slug
	userInfo, err := client.GetUserInfo()
	if err != nil {
		return fmt.Errorf("failed to get user information: %w", err)
	}

	// Determine the organization slug
	ownerSlug := userInfo.OrgSlug
	if ownerSlug == "" {
		// Fallback to username if orgSlug is not available
		ownerSlug = userInfo.Username
		if ownerSlug == "" {
			// Fallback to email prefix if username is not available
			if userInfo.Email != "" {
				ownerSlug = strings.Split(userInfo.Email, "@")[0]
			} else {
				return fmt.Errorf("could not determine organization slug from user information")
			}
		}
	}

	// Use the ruleset name as the rule slug
	ruleSlug := rs.Name

	// Validate the slug format
	if !isValidSlug(ownerSlug) || !isValidSlug(ruleSlug) {
		return fmt.Errorf("invalid slug format: %s/%s", ownerSlug, ruleSlug)
	}

	// Find the main rule file to publish
	// Look for index.md in the rules directory
	var ruleFilePath string
	if rulesPath != "" {
		// If a path was specified, look for index.md in that directory
		stat, err := os.Stat(rulesPath)
		if err == nil && stat.IsDir() {
			ruleFilePath = filepath.Join(rulesPath, "index.md")
		}
	} else {
		// Look in current directory for index.md
		ruleFilePath = "index.md"
	}

	// If index.md doesn't exist, look for any .md file in the rules directory
	if _, err := os.Stat(ruleFilePath); os.IsNotExist(err) {
		// Look for any .md file in the current directory
		files, err := filepath.Glob("*.md")
		if err == nil && len(files) > 0 {
			ruleFilePath = files[0]
		} else {
			return fmt.Errorf("no rule file found to publish. Please create an index.md file or specify a rule file")
		}
	}

	// Check if the file exists
	if _, err := os.Stat(ruleFilePath); os.IsNotExist(err) {
		return fmt.Errorf("rule file '%s' does not exist", ruleFilePath)
	}

	// Read the file content
	content, err := ioutil.ReadFile(ruleFilePath)
	if err != nil {
		return fmt.Errorf("failed to read rule file: %w", err)
	}

	// Publish the rule
	color.Cyan("Publishing rule to %s/%s with visibility: %s", ownerSlug, ruleSlug, visibility)
	err = client.PublishRule(ownerSlug, ruleSlug, string(content), visibility)
	if err != nil {
		return fmt.Errorf("failed to publish rule: %w", err)
	}

	ruleFileName := filepath.Base(ruleFilePath)
	color.Green("Successfully published rule '%s' to %s/%s", ruleFileName, ownerSlug, ruleSlug)
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
}
