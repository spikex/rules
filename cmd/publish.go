package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"rules-cli/internal/auth"
	"rules-cli/internal/registry"

	"github.com/spf13/cobra"
)

var (
	slug       string
	visibility string
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish <rule-file>",
	Short: "Publish a rule file to the registry",
	Long: `Publishes a rule file to the registry.

This command requires authentication and will prompt you to login if you're not already.
You must specify the organization/ruleset slug to publish to.

The visibility can be set to "public" (default) or "private".

Examples:
  rules publish my-rule.md --slug my-org/my-rules
  rules publish my-rule.md --slug my-org/my-rules --visibility private`,
	Args: cobra.ExactArgs(1),
	RunE: runPublishCommand,
}

// runPublishCommand implements the main logic for the publish command
func runPublishCommand(cmd *cobra.Command, args []string) error {
	// Ensure the user is authenticated
	authenticated, err := auth.EnsureAuthenticated(true)
	if err != nil || !authenticated {
		return fmt.Errorf("authentication required to publish rules")
	}

	// Validate the slug format
	if slug == "" {
		return fmt.Errorf("--slug is required and must be in the format 'organization/ruleset'")
	}

	if !strings.Contains(slug, "/") {
		return fmt.Errorf("slug must be in the format 'organization/ruleset'")
	}

	parts := strings.Split(slug, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("slug must be in the format 'organization/ruleset'")
	}

	ownerSlug := parts[0]
	ruleSlug := parts[1]

	// Validate the visibility
	if visibility != "public" && visibility != "private" {
		return fmt.Errorf("visibility must be either 'public' or 'private'")
	}

	// Get the rule file path
	ruleFilePath := args[0]

	// Check if the file exists
	if _, err := os.Stat(ruleFilePath); os.IsNotExist(err) {
		return fmt.Errorf("rule file '%s' does not exist", ruleFilePath)
	}

	// Read the file content
	content, err := ioutil.ReadFile(ruleFilePath)
	if err != nil {
		return fmt.Errorf("failed to read rule file: %w", err)
	}

	// Create registry client
	authConfig := auth.LoadAuthConfig()
	client := registry.NewClient(cfg.RegistryURL)
	client.SetAuthToken(authConfig.AccessToken)

	// Publish the rule
	fmt.Printf("Publishing rule to %s/%s with visibility: %s\n", ownerSlug, ruleSlug, visibility)
	err = client.PublishRule(ownerSlug, ruleSlug, string(content), visibility)
	if err != nil {
		return fmt.Errorf("failed to publish rule: %w", err)
	}

	ruleFileName := filepath.Base(ruleFilePath)
	fmt.Printf("Successfully published rule '%s' to %s/%s\n", ruleFileName, ownerSlug, ruleSlug)
	return nil
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// Add flags
	publishCmd.Flags().StringVar(&slug, "slug", "", "The organization/ruleset slug to publish to (required)")
	publishCmd.Flags().StringVar(&visibility, "visibility", "public", "Set the visibility of the rule to 'public' or 'private'")

	// Mark required flags
	publishCmd.MarkFlagRequired("slug")
}