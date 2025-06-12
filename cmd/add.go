package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rules-cli/internal/formats"
	"rules-cli/internal/registry"
	"rules-cli/internal/ruleset"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <rulename>",
	Short: "Add a rule from the registry",
	Long: `Add a rule from the registry to the current ruleset.
The rule will be downloaded and added to the rules.json file.

For GitHub repositories, use the gh: prefix followed by the owner/repo.
For example: gh:owner/repo

When importing from GitHub repositories, the tool will:
- Download all files from the 'src/' directory in the repository
- Use the main branch of the repository by default`,
	Example: `  rules add vercel/nextjs
  rules add redis
  rules add gh:owner/repo`,
	Args: cobra.ExactArgs(1),
	RunE: runAddCommand,
}

// parseRuleNameAndVersion extracts the rule name and version from the input argument
func parseRuleNameAndVersion(ruleArg string) (name string, version string) {
	version = "0.0.1" // Default version
	name = ruleArg
	
	// Check if version is specified (but not for GitHub URLs)
	if !strings.HasPrefix(name, "gh:") {
		if parts := strings.Split(name, "@"); len(parts) > 1 {
			name = parts[0]
			version = parts[1]
		}
	}
	
	return name, version
}

// setupRulesDirectory ensures the rules directory exists and returns paths
func setupRulesDirectory(format string) (rulesDir string, rulesJSONPath string, err error) {
	// Get rules directory for the format
	rulesDir, err = formats.GetRulesDirectory(format)
	if err != nil {
		return "", "", fmt.Errorf("failed to get rules directory: %w", err)
	}
	
	// Get rules.json path
	rulesJSONPath, err = formats.GetRulesJSONPath(format)
	if err != nil {
		return "", "", fmt.Errorf("failed to get rules.json path: %w", err)
	}
	
	return rulesDir, rulesJSONPath, nil
}

// getFormatSuggestion checks for existing format folders and provides suggestion if needed
func getFormatSuggestion(rulesJSONPath string) string {
	if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
		// Use the new helper function
		suggestion, err := formats.GetFormatSuggestionMessage()
		if err == nil && suggestion != "" {
			return suggestion
		}
	}
	return ""
}

// loadOrCreateRuleSet loads an existing ruleset or creates a new one
func loadOrCreateRuleSet(rulesJSONPath string) (*ruleset.RuleSet, error) {
	if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(rulesJSONPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for rules.json: %w", err)
		}
		
		// Create a default ruleset
		rs := ruleset.DefaultRuleSet(filepath.Base(filepath.Dir(rulesJSONPath)))
		fmt.Println("Creating new rules.json file with default structure")
		return rs, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to check rules.json file: %w", err)
	}
	
	// Load existing ruleset
	rs, err := ruleset.LoadRuleSet(rulesJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ruleset: %w", err)
	}
	
	return rs, nil
}

// downloadRule downloads a rule from the registry
func downloadRule(client *registry.Client, ruleName, ruleVersion, rulesDir string) error {
	if strings.HasPrefix(ruleName, "gh:") {
		fmt.Printf("Downloading rules from GitHub repository '%s' (src/ directory)...\n", ruleName[3:])
	} else {
		fmt.Printf("Downloading rule '%s' (version %s)...\n", ruleName, ruleVersion)
	}
	
	return client.DownloadRule(ruleName, ruleVersion, rulesDir)
}

// runAddCommand implements the main logic for the add command
func runAddCommand(cmd *cobra.Command, args []string) error {
	// Parse rule name and version
	ruleName, ruleVersion := parseRuleNameAndVersion(args[0])
	
	// Setup rules directory
	rulesDir, rulesJSONPath, err := setupRulesDirectory(format)
	if err != nil {
		return err
	}
	
	// Check for format suggestions
	formatSuggestion := getFormatSuggestion(rulesJSONPath)
	
	// Load or create ruleset
	rs, err := loadOrCreateRuleSet(rulesJSONPath)
	if err != nil {
		return err
	}
	
	// Check if rule already exists
	if rs.RuleExists(ruleName) {
		version, _ := rs.GetRuleVersion(ruleName)
		return fmt.Errorf("rule '%s' already exists with version %s", ruleName, version)
	}
	
	// Create registry client
	client := registry.NewClient(cfg.RegistryURL)
	
	// Download rule
	if err := downloadRule(client, ruleName, ruleVersion, rulesDir); err != nil {
		return fmt.Errorf("failed to download rule: %w", err)
	}
	
	// Add rule to ruleset
	rs.AddRule(ruleName, ruleVersion)
	if err := rs.SaveRuleSet(rulesJSONPath); err != nil {
		return fmt.Errorf("failed to save ruleset: %w", err)
	}
	
	fmt.Printf("Rule '%s' (version %s) added successfully\n", ruleName, ruleVersion)
	
	// Print format suggestion at the very end if applicable
	if formatSuggestion != "" {
		fmt.Println(formatSuggestion)
	}
	
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}