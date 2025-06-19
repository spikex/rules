package cmd

import (
	"fmt"
	"os"
	"rules-cli/internal/formats"

	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render [format]",
	Short: "Render rules to a specific format",
	Long: `Renders existing rules to a specified format.
Copies all rules from the default location (.rules/) to the target format
as described in render-formats.md.

Supported formats: continue, cursor, windsurf, claude, copilot, codex, cline, cody, amp`,
	Example: `  rules render cursor
  rules render continue`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("format is required")
		}
		
		formatName := args[0]
		
		if formatName == "default" {
			return fmt.Errorf("cannot render to default format as it is the source")
		}

		// Get the source directory
		sourceDir, err := formats.GetRulesDirectory("default")
		if err != nil {
			return fmt.Errorf("failed to get source directory: %w", err)
		}

		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			return fmt.Errorf("source directory %s does not exist", sourceDir)
		}

		fmt.Printf("Rendering rules to %s format...\n", formatName)
		
		// Use the formats package to handle the rendering based on the target format
		verbose, _ := cmd.Flags().GetBool("verbose")
		if err := formats.RenderRulesToFormat(sourceDir, formatName, verbose); err != nil {
			return fmt.Errorf("failed to render rules to %s format: %w", formatName, err)
		}

		fmt.Printf("Successfully rendered rules to %s format\n", formatName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}