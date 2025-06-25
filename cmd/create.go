package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"rules-cli/internal/ruleset"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	globs       string
	description string
	alwaysApply bool
	ruleName    string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [rule-name] [rule-body]",
	Short: "Create a new rule using Continue format",
	Long: `Create a new rule in the current directory using Continue format specification.
If rule parameters are not provided, they will be prompted for interactively.
This command does not modify the rules.json file.`,
	Example: `  rules create my-rule "This is the body of the rule"
  rules create --globs "**/*.{tsx,jsx}" --description "React style guide" my-rule
  rules create --alwaysApply "Always apply this rule"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get rule name
		if len(args) > 0 {
			ruleName = args[0]
		} else {
			prompt := promptui.Prompt{
				Label: "Rule name",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("rule name cannot be empty")
					}
					return nil
				},
			}
			var err error
			ruleName, err = prompt.Run()
			if err != nil {
				return fmt.Errorf("failed to get rule name: %w", err)
			}
		}

		// Create rule struct for Continue format
		rule := ruleset.Rule{
			Description: description,
			Globs:       globs,
			AlwaysApply: alwaysApply,
		}

		// Prompt for globs if not provided
		if rule.Globs == "" {
			prompt := promptui.Prompt{
				Label:   "Glob patterns",
				Default: "**/*",
			}
			var err error
			rule.Globs, err = prompt.Run()
			if err != nil {
				return fmt.Errorf("failed to get globs: %w", err)
			}
		}

		// Prompt for description if not provided
		if rule.Description == "" {
			prompt := promptui.Prompt{
				Label: "Rule description",
			}
			var err error
			rule.Description, err = prompt.Run()
			if err != nil {
				return fmt.Errorf("failed to get description: %w", err)
			}
		}

		// Note: Tags are not part of Continue format specification and are omitted

		// Get rule body
		if len(args) > 1 {
			rule.Body = args[1]
		} else {
			color.Cyan("Enter rule body (press Ctrl+D on a new line when done):")
			scanner := bufio.NewScanner(os.Stdin)
			var bodyLines []string
			for scanner.Scan() {
				bodyLines = append(bodyLines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}
			rule.Body = strings.Join(bodyLines, "\n")
		}

		// Create the rule file in Continue format
		filePath, err := ruleset.CreateRule(rule, "continue", ruleName)
		if err != nil {
			return fmt.Errorf("failed to create rule: %w", err)
		}

		color.Green("Rule '%s' created successfully at %s", ruleName, filePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&globs, "globs", "", "Glob patterns to match files (defaults to \"**/*\" if not specified)")
	createCmd.Flags().StringVar(&description, "description", "", "Short description of the rule")
	createCmd.Flags().BoolVar(&alwaysApply, "alwaysApply", false, "Whether to always apply the rule (creates alwaysApply: true in frontmatter)")
}
