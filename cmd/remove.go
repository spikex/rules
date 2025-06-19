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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <rulename>",
	Short: "Remove a rule from the ruleset",
	Long: `Remove a rule from the ruleset.
This will remove the rule from rules.json and delete its files from the .rules folder.

For GitHub repositories, use the same gh: prefix as when adding.`,
	Example: `  rules remove vercel/nextjs
  rules remove redis
  rules remove workos/authkit-nextjs
  rules remove gh:owner/repo`,
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

		// Calculate rule directory path
		ruleDir := filepath.Join(rulesDir, ruleName)

		// Delete rule directory and files
		if err := os.RemoveAll(ruleDir); err != nil {
			color.Red("Warning: Failed to delete rule files: %v", err)
			// Continue anyway to remove from rules.json
		} else {
			color.Cyan("Deleted rule files from %s", ruleDir)

			// Check if parent directory is empty and delete it if so
			if strings.Contains(ruleName, "/") {
				// Get the parent directory name (e.g., "starter" from "starter/nextjs-rules")
				parentName := ruleName[:strings.LastIndex(ruleName, "/")]
				parentDir := filepath.Join(rulesDir, parentName)

				// Check if parent directory exists and is empty
				if entries, err := os.ReadDir(parentDir); err == nil && len(entries) == 0 {
					os.Remove(parentDir)
				}
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
}
