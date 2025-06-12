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
	RunE: func(cmd *cobra.Command, args []string) error {
		ruleName := args[0]
		ruleVersion := "0.0.1" // Default version
		
		// Check if version is specified (but not for GitHub URLs)
		if !strings.HasPrefix(ruleName, "gh:") {
			if parts := strings.Split(ruleName, "@"); len(parts) > 1 {
				ruleName = parts[0]
				ruleVersion = parts[1]
			}
		}
		
		// Get rules directory for the format
		rulesDir, err := formats.GetRulesDirectory(format)
		if err != nil {
			return fmt.Errorf("failed to get rules directory: %w", err)
		}
		
		// Get rules.json path
		rulesJSONPath, err := formats.GetRulesJSONPath(format)
		if err != nil {
			return fmt.Errorf("failed to get rules.json path: %w", err)
		}
		
		// Load ruleset or create a new one if it doesn't exist
		var rs *ruleset.RuleSet
		if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
			// Create directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(rulesJSONPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for rules.json: %w", err)
			}
			
			// Create a default ruleset
			rs = ruleset.DefaultRuleSet(filepath.Base(filepath.Dir(rulesJSONPath)))
			fmt.Println("Creating new rules.json file with default structure")
		} else if err != nil {
			return fmt.Errorf("failed to check rules.json file: %w", err)
		} else {
			// Load existing ruleset
			rs, err = ruleset.LoadRuleSet(rulesJSONPath)
			if err != nil {
				return fmt.Errorf("failed to load ruleset: %w", err)
			}
		}
		
		// Check if rule already exists
		if rs.RuleExists(ruleName) {
			version, _ := rs.GetRuleVersion(ruleName)
			return fmt.Errorf("rule '%s' already exists with version %s", ruleName, version)
		}
		
		// Create registry client
		client := registry.NewClient(cfg.RegistryURL)
		
		// Download rule
		if strings.HasPrefix(ruleName, "gh:") {
			fmt.Printf("Downloading rules from GitHub repository '%s' (src/ directory)...\n", ruleName[3:])
		} else {
			fmt.Printf("Downloading rule '%s' (version %s)...\n", ruleName, ruleVersion)
		}
		
		if err := client.DownloadRule(ruleName, ruleVersion, rulesDir); err != nil {
			return fmt.Errorf("failed to download rule: %w", err)
		}
		
		// Add rule to ruleset
		rs.AddRule(ruleName, ruleVersion)
		if err := rs.SaveRuleSet(rulesJSONPath); err != nil {
			return fmt.Errorf("failed to save ruleset: %w", err)
		}
		
		fmt.Printf("Rule '%s' (version %s) added successfully\n", ruleName, ruleVersion)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}