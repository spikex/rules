package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"rules-cli/internal/formats"
	"rules-cli/internal/registry"
	"rules-cli/internal/ruleset"
)

var (
	forceInstall bool
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Synchronize rules directory with rules.json",
	Long: `Synchronizes the .rules directory with the contents of rules.json.
Performs a clean installation by:
- Removing all existing rule files and directories first
- Re-downloading and installing all rules specified in rules.json

This ensures the rules directory matches exactly what's defined in rules.json.`,
	Example: `  rules install
  rules install --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		// Check if rules directory exists
		if _, err := os.Stat(rulesDir); err == nil {
			// Directory exists, check if we should clean it
			if !forceInstall {
				fmt.Printf("This will remove all existing rules in '%s' and reinstall them from rules.json.\n", rulesDir)
				fmt.Print("Continue? (y/N): ")
				
				var response string
				fmt.Scanln(&response)
				
				if strings.ToLower(response) != "y" {
					fmt.Println("Installation cancelled.")
					return nil
				}
			}
			
			// Remove all files and directories in the rules directory
			fmt.Printf("Removing existing rules from '%s'...\n", rulesDir)
			if err := removeContents(rulesDir); err != nil {
				return fmt.Errorf("failed to clean rules directory: %w", err)
			}
		} else if os.IsNotExist(err) {
			// Directory doesn't exist, create it
			if err := os.MkdirAll(rulesDir, 0755); err != nil {
				return fmt.Errorf("failed to create rules directory: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check rules directory: %w", err)
		}

		// Create registry client
		client := registry.NewClient(cfg.RegistryURL)

		// Install each rule from rules.json
		fmt.Println("Installing rules from rules.json...")
		if len(rs.Rules) == 0 {
			fmt.Println("No rules found in rules.json.")
			return nil
		}

		successCount := 0
		errorCount := 0

		for ruleName, ruleVersion := range rs.Rules {
			fmt.Printf("Installing rule '%s' (version: %s)...\n", ruleName, ruleVersion)
			
			if err := client.DownloadRule(ruleName, ruleVersion, rulesDir); err != nil {
				fmt.Printf("Error installing rule '%s': %v\n", ruleName, err)
				errorCount++
			} else {
				successCount++
			}
		}

		// Print summary
		fmt.Printf("\nInstallation complete: %d rules installed, %d failed\n", successCount, errorCount)
		if errorCount > 0 {
			return fmt.Errorf("%d rules failed to install", errorCount)
		}
		
		return nil
	},
}

// removeContents removes all files and directories inside the specified directory
// but keeps the directory itself
func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	
	entries, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(dir, entry))
		if err != nil {
			return err
		}
	}
	
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
	
	// Add flags
	installCmd.Flags().BoolVar(&forceInstall, "force", false, "Skip confirmation prompts")
}