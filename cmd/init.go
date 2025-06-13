package cmd

import (
	"fmt"
	"os"
	"strings"

	"rules-cli/internal/formats"

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
		fmt.Printf("Initializing rules with format: %s\n", format)
		
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
				fmt.Printf("Found existing rules folder(s): %s\n", strings.Join(formatFolders, ", "))
				fmt.Printf("Consider running 'rules render %s' to initialize rules.json from existing rules\n", formatFolders[0])
			}
		}
		
		if err := formats.InitializeFormat(format); err != nil {
			return fmt.Errorf("initialization failed: %w", err)
		}
		
		fmt.Printf("Rules initialized successfully. Format: %s\n", format)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}