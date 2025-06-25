package cmd

import (
	"fmt"
	"rules-cli/internal/formats"
	"rules-cli/internal/ruleset"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rules currently installed in the project",
	Long: `List all rules currently installed in the project, similar to 'npm list'.
This command reads from the rules.json file and displays the installed rules
with their versions in a tree-like format.

The command only reads local files and does not require network access.`,
	Example: `  rules list`,
	Args:    cobra.NoArgs,
	RunE:    runListCommand,
}

// runListCommand implements the main logic for the list command
func runListCommand(cmd *cobra.Command, args []string) error {
	// Get rules.json path
	rulesJSONPath, err := formats.GetRulesJSONPath(format)
	if err != nil {
		return fmt.Errorf("failed to get rules.json path: %w", err)
	}

	// Load the ruleset
	rs, err := ruleset.LoadRuleSet(rulesJSONPath)
	if err != nil {
		return fmt.Errorf("No rules.json file found in current directory\nRun 'rules init' to initialize a new project")
	}

	// Display the project name and version
	color.Cyan("%s@%s", rs.Name, rs.Version)

	// Check if rules are empty
	if len(rs.Rules) == 0 {
		color.Yellow("(empty)")
		return nil
	}

	// Sort rules for consistent output
	ruleNames := make([]string, 0, len(rs.Rules))
	for name := range rs.Rules {
		ruleNames = append(ruleNames, name)
	}

	// Simple alphabetical sort
	for i := 0; i < len(ruleNames); i++ {
		for j := i + 1; j < len(ruleNames); j++ {
			if ruleNames[i] > ruleNames[j] {
				ruleNames[i], ruleNames[j] = ruleNames[j], ruleNames[i]
			}
		}
	}

	// Display rules in tree format
	for i, name := range ruleNames {
		version := rs.Rules[name]

		// Use appropriate tree characters
		var prefix string
		if i == len(ruleNames)-1 {
			prefix = "└── "
		} else {
			prefix = "├── "
		}

		fmt.Printf("%s%s@%s\n", prefix, name, version)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
