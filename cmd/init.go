package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"rules-cli/internal/formats"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new rules directory",
	Long: `Initialize a new rules directory with the specified format.
This creates the necessary directory structure and an empty rules.json file.`,
	Example: `  rules init
  rules init --format cursor`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Initializing rules with format: %s\n", format)
		
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