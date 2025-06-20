package cmd

import (
	"fmt"
	"rules-cli/internal/formats"

	"github.com/spf13/cobra"
)

// formatsCmd represents the formats command
var formatsCmd = &cobra.Command{
	Use:   "formats",
	Short: "List all available render formats",
	Long: `List all supported render formats for the 'rules render' command.
Each format represents a different AI code assistant platform with specific
folder structures and file extensions.`,
	Example: `  rules formats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available render formats:")
		fmt.Println()

		formatList := formats.GetAllFormats()
		for _, format := range formatList {
			if format.IsSingleFile {
				fmt.Printf("%-10s - %s (%s)\n", format.Name, format.SingleFilePath, format.Description)
			} else {
				fmt.Printf("%-10s - %s/*%s (%s)\n", format.Name, format.DirectoryPrefix, format.FileExtension, format.Description)
			}
		}

		fmt.Println()
		fmt.Println("Usage: rules render <format>")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatsCmd)
}