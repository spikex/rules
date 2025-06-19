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
	tags        []string
	globs       string
	description string
	alwaysApply bool
	ruleName    string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [rule-name] [rule-body]",
	Short: "Create a new rule",
	Long: `Create a new rule in the current ruleset.
If rule parameters are not provided, they will be prompted for interactively.
This command does not modify the rules.json file.`,
	Example: `  rules create my-rule "This is the body of the rule"
  rules create --tags frontend --globs "*.{tsx,jsx}" --description "React style guide" my-rule`,
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

		// Create rule struct
		rule := ruleset.Rule{
			Description: description,
			Globs:       globs,
			AlwaysApply: alwaysApply,
		}

		// Process tags
		if len(tags) > 0 {
			rule.Tags = tags
		} else if rule.Tags == nil {
			prompt := promptui.Prompt{
				Label: "Tags (comma-separated)",
			}
			tagsInput, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("failed to get tags: %w", err)
			}

			if tagsInput != "" {
				rule.Tags = strings.Split(tagsInput, ",")
				// Trim whitespace from tags
				for i, tag := range rule.Tags {
					rule.Tags[i] = strings.TrimSpace(tag)
				}
			}
		}

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

		// Create the rule file
		if err := ruleset.CreateRule(rule, format, ruleName); err != nil {
			return fmt.Errorf("failed to create rule: %w", err)
		}

		color.Green("Rule '%s' created successfully", ruleName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags for the rule (comma-separated)")
	createCmd.Flags().StringVar(&globs, "globs", "", "Glob patterns to match files")
	createCmd.Flags().StringVar(&description, "description", "", "Short description of the rule")
	createCmd.Flags().BoolVar(&alwaysApply, "alwaysApply", false, "Whether to always apply the rule")
}
