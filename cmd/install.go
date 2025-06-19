package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rules-cli/internal/auth"
	"rules-cli/internal/formats"
	"rules-cli/internal/registry"
	"rules-cli/internal/ruleset"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
	Example: `  rules install`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if arguments were provided
		if len(args) > 0 {
			color.Yellow("Did you mean 'rules add %s' instead?", strings.Join(args, " "))
			return fmt.Errorf("install command does not accept arguments")
		}

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

		// Check for existing format folders, but save this information for later
		var formatSuggestion string
		if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
			// Use the new helper function
			suggestion, err := formats.GetFormatSuggestionMessage()
			if err == nil && suggestion != "" {
				formatSuggestion = suggestion
		}
		}

		// Check if rules.json exists, create it if it doesn't
		var rs *ruleset.RuleSet
		if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
			// Create directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(rulesJSONPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for rules.json: %w", err)
			}
			
			// Create a default ruleset
			rs = ruleset.DefaultRuleSet(filepath.Base(filepath.Dir(rulesJSONPath)))
			color.Cyan("Creating new rules.json file with default structure")
			
			// Save the new ruleset
			if err := rs.SaveRuleSet(rulesJSONPath); err != nil {
				return fmt.Errorf("failed to create rules.json: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check rules.json file: %w", err)
		} else {
			// Load existing ruleset
			rs, err = ruleset.LoadRuleSet(rulesJSONPath)
			if err != nil {
				return fmt.Errorf("failed to load ruleset: %w", err)
			}
		}

		// Check if rules directory exists
		if _, err := os.Stat(rulesDir); err == nil {
			// Directory exists, clean it without confirmation
			color.Cyan("Removing existing rules from '%s'...", rulesDir)
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
		authConfig := auth.LoadAuthConfig()
		client := registry.NewClient(cfg.RegistryURL)
		client.SetAuthToken(authConfig.AccessToken)

		// Install each rule from rules.json
		color.Cyan("Installing rules from rules.json...")
		if len(rs.Rules) == 0 {
			color.Yellow("No rules found in rules.json.")
			
			// Print format suggestion at the very end if applicable
			if formatSuggestion != "" {
				fmt.Println(formatSuggestion)
			}
	return nil
}

		successCount := 0
		errorCount := 0

		for ruleName, ruleVersion := range rs.Rules {
			color.Cyan("Installing rule '%s' (version: %s)...", ruleName, ruleVersion)
			
			// Split the rule name into owner and rule parts
			parts := strings.Split(ruleName, "/")
			ownerSlug := parts[0]
			ruleSlug := ""
			if len(parts) > 1 {
				ruleSlug = parts[1]
			}
			if err := client.DownloadRule(ownerSlug, ruleSlug, ruleVersion, rulesDir); err != nil {
				color.Red("Error installing rule '%s': %v", ruleName, err)
				errorCount++
			} else {
				successCount++
			}
		}

		// Print summary
		color.Green("\nInstallation complete: %d rules installed, %d failed", successCount, errorCount)
		
		// Print format suggestion at the very end if applicable
		if formatSuggestion != "" {
			fmt.Println(formatSuggestion)
		}
		
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
}
