package cmd

import (
	"fmt"
	"os"
	"strings"

	"rules-cli/internal/formats"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new rules directory",
	Long: `Initialize a new rules directory with the specified format.
This creates the necessary directory structure and an empty rules.json file.`,
	Example: `  rules init`,
	RunE: func(cmd *cobra.Command, args []string) error {
		color.Cyan("Initializing rules with format: %s", format)

		// Check if rules.json doesn't exist
		rulesJSONPath, err := formats.GetRulesJSONPath(format)
		if err != nil {
			return fmt.Errorf("failed to get rules.json path: %w", err)
		}

		if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
			// Check for any top-level folder of the structure ".{folder-name}/rules"
			formatFolders, err := formats.FindRulesFormats()
			if err == nil && len(formatFolders) > 0 {
				// Suggest rendering to the user
				color.Yellow("Found existing rules folder(s): %s", strings.Join(formatFolders, ", "))
				color.Yellow("Consider running 'rules render %s' to initialize rules.json from existing rules", formatFolders[0])
			}
		}

		if err := formats.InitializeFormat(format); err != nil {
			return fmt.Errorf("initialization failed: %w", err)
		}

		color.Green("Rules initialized successfully. Format: %s", format)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
