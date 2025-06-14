package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rules-cli/internal/formats"
	"rules-cli/internal/ruleset"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	deleteFlag bool
	forceFlag  bool
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <rulename>",
	Short: "Remove a rule from the ruleset",
	Long: `Remove a rule from the ruleset.
By default, this only removes the rule from rules.json.
Use --delete to also remove rule files from disk.

For GitHub repositories, use the same gh: prefix as when adding.`,
	Example: `  rules remove vercel/nextjs
  rules remove redis --delete
  rules remove workos/authkit-nextjs --delete --force
  rules remove gh:owner/repo --delete`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ruleName := args[0]

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

		// Load ruleset
		rs, err := ruleset.LoadRuleSet(rulesJSONPath)
		if err != nil {
			return fmt.Errorf("failed to load ruleset: %w", err)
		}

		// Check if rule exists
		if !rs.RuleExists(ruleName) {
			return fmt.Errorf("rule '%s' does not exist in the ruleset", ruleName)
		}

		// Get rule version before removal
		version, _ := rs.GetRuleVersion(ruleName)

		// Handle deletion of rule files
		if deleteFlag {
			// Calculate rule directory path
			ruleDir := filepath.Join(rulesDir, ruleName)
			
			// Check if we should prompt for confirmation
			if !forceFlag {
				color.Yellow("This will remove rule '%s' (version %s) and delete its files. Continue? [y/N]: ", ruleName, version)
				var response string
				fmt.Scanln(&response)
				
				if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
					color.Yellow("Operation cancelled")
					return nil
				}
			}

			// Delete rule directory and files
			if err := os.RemoveAll(ruleDir); err != nil {
				color.Red("Warning: Failed to delete rule files: %v", err)
				// Continue anyway to remove from rules.json
			} else {
				color.Cyan("Deleted rule files from %s", ruleDir)
			}
		}

		// Remove rule from ruleset
		rs.RemoveRule(ruleName)
		if err := rs.SaveRuleSet(rulesJSONPath); err != nil {
			return fmt.Errorf("failed to save ruleset: %w", err)
		}

		color.Green("Rule '%s' (version %s) removed successfully", ruleName, version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	
	// Add flags specific to remove command
	removeCmd.Flags().BoolVar(&deleteFlag, "delete", false, "Also delete rule files from disk")
	removeCmd.Flags().BoolVar(&forceFlag, "force", false, "Skip confirmation prompts")
}