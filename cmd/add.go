package cmd

import (
	"encoding/json"
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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <rulename>",
	Short: "Add a rule from the registry",
	Long: `Add a rule from the registry to the current ruleset.
The rule will be downloaded and added to the rules.json file.

Rules are downloaded from the registry API using the GET endpoint 
(e.g., api.continue.dev/v0/<owner-slug>/<rule-slug>/latest/download).

For GitHub repositories, use the gh: prefix followed by the owner/repo.
For example: gh:owner/repo

When importing from GitHub repositories, the tool will:
- Download all files in the repository
- Use the main branch of the repository by default
- Look for rules.json in the downloaded files to find the version`,
	Example: `  rules add vercel/nextjs
  rules add redis
  rules add gh:owner/repo`,
	Args: cobra.ExactArgs(1),
	RunE: runAddCommand,
}

// RuleIdentifier contains the parsed components of a rule identifier
type RuleIdentifier struct {
	OwnerSlug string
	RuleSlug  string
	Version   string
	FullName  string // The full name as it should appear in rules.json
}

// parseRuleIdentifier extracts the owner, rule slug, and version from the input argument
func parseRuleIdentifier(ruleArg string) (*RuleIdentifier, error) {
	identifier := &RuleIdentifier{
		Version: "latest", // Default version
	}
	
	// Handle GitHub repositories
	if strings.HasPrefix(ruleArg, "gh:") {
		// Format: gh:owner/repo or gh:owner/repo@version
		repoPath := ruleArg[3:] // Remove "gh:" prefix
		
		// Check for version
		if parts := strings.Split(repoPath, "@"); len(parts) > 1 {
			repoPath = parts[0]
			identifier.Version = parts[1]
		}
		
		// Split owner/repo
		repoParts := strings.Split(repoPath, "/")
		if len(repoParts) != 2 {
			return nil, fmt.Errorf("GitHub repository must be in format 'gh:owner/repo'")
		}
		
		identifier.OwnerSlug = "gh:" + repoParts[0]
		identifier.RuleSlug = repoParts[1]
		identifier.FullName = ruleArg
		
		return identifier, nil
	}
	
	// Handle registry rules
	ruleName := ruleArg
	
	// Check if version is specified
	if parts := strings.Split(ruleName, "@"); len(parts) > 1 {
		ruleName = parts[0]
		identifier.Version = parts[1]
	}
	
	// Check if owner/rule format
	if parts := strings.Split(ruleName, "/"); len(parts) == 2 {
		identifier.OwnerSlug = parts[0]
		identifier.RuleSlug = parts[1]
		identifier.FullName = ruleName
	} else if len(parts) == 1 {
		// Single name - might need a default owner or handle differently
		// For now, we'll assume the rule name is the owner and rule slug
		identifier.OwnerSlug = parts[0]
		identifier.RuleSlug = parts[0]
		identifier.FullName = ruleName
	} else {
		return nil, fmt.Errorf("rule name must be in format 'owner/rule' or 'rulename'")
	}
	
	return identifier, nil
}

// setupRulesDirectory ensures the rules directory exists and returns paths
func setupRulesDirectory(format string) (rulesDir string, rulesJSONPath string, err error) {
	// Get rules directory for the format
	rulesDir, err = formats.GetRulesDirectory(format)
	if err != nil {
		return "", "", fmt.Errorf("failed to get rules directory: %w", err)
	}
	
	// Get rules.json path
	rulesJSONPath, err = formats.GetRulesJSONPath(format)
	if err != nil {
		return "", "", fmt.Errorf("failed to get rules.json path: %w", err)
	}
	
	return rulesDir, rulesJSONPath, nil
}

// getFormatSuggestion checks for existing format folders and provides suggestion if needed
func getFormatSuggestion(rulesJSONPath string) string {
	if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
		// Use the new helper function
		suggestion, err := formats.GetFormatSuggestionMessage()
		if err == nil && suggestion != "" {
			return suggestion
		}
	}
	return ""
}

// loadOrCreateRuleSet loads an existing ruleset or creates a new one
func loadOrCreateRuleSet(rulesJSONPath string) (*ruleset.RuleSet, error) {
	if _, err := os.Stat(rulesJSONPath); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(rulesJSONPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for rules.json: %w", err)
		}
		
		// Create a default ruleset
		rs := ruleset.DefaultRuleSet(filepath.Base(filepath.Dir(rulesJSONPath)))
		color.Cyan("Creating new rules.json file with default structure")
		return rs, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to check rules.json file: %w", err)
	}
	
	// Load existing ruleset
	rs, err := ruleset.LoadRuleSet(rulesJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ruleset: %w", err)
	}
	
	return rs, nil
}

// downloadRule downloads a rule from the registry
func downloadRule(client *registry.Client, identifier *RuleIdentifier, rulesDir string) (string, error) {
	if strings.HasPrefix(identifier.FullName, "gh:") {
		color.Cyan("Downloading rules from GitHub repository '%s'...", identifier.FullName[3:])
	} else {
		color.Cyan("Downloading rule '%s/%s' (version %s) from registry API...", identifier.OwnerSlug, identifier.RuleSlug, identifier.Version)
	}
	
	if err := client.DownloadRule(identifier.OwnerSlug, identifier.RuleSlug, identifier.Version, rulesDir); err != nil {
		return "", fmt.Errorf("failed to download rule: %w", err)
	}

	// Check for the actual version in the downloaded rule.json file
	var ruleDir string
	if strings.HasPrefix(identifier.FullName, "gh:") {
		ruleDir = filepath.Join(rulesDir, "gh:"+identifier.OwnerSlug[3:]+"/"+identifier.RuleSlug)
	} else {
		ruleDir = filepath.Join(rulesDir, identifier.OwnerSlug, identifier.RuleSlug)
	}

	ruleInfoPath := filepath.Join(ruleDir, "rules.json")

	if _, err := os.Stat(ruleInfoPath); err == nil {
		// Rule info file exists, try to get the actual version
		data, err := os.ReadFile(ruleInfoPath)
		if err == nil {
			var ruleInfo map[string]interface{}
			if err := json.Unmarshal(data, &ruleInfo); err == nil {
				if version, ok := ruleInfo["version"].(string); ok && version != "" {
					return version, nil
				}
			}
		}
	}
	
	// If we can't find or parse the version, return the requested version
	return identifier.Version, nil
}
	
// runAddCommand implements the main logic for the add command
func runAddCommand(cmd *cobra.Command, args []string) error {
	// Parse rule identifier
	identifier, err := parseRuleIdentifier(args[0])
	if err != nil {
		return fmt.Errorf("invalid rule identifier: %w", err)
	}
	
	// Setup rules directory
	rulesDir, rulesJSONPath, err := setupRulesDirectory(format)
	if err != nil {
		return err
	}
	
	// Check for format suggestions
	formatSuggestion := getFormatSuggestion(rulesJSONPath)
	
	// Load or create ruleset
	rs, err := loadOrCreateRuleSet(rulesJSONPath)
	if err != nil {
		return err
	}
	
	// Check if rule already exists (using the full name for consistency)
	if rs.RuleExists(identifier.FullName) {
		version, _ := rs.GetRuleVersion(identifier.FullName)
		return fmt.Errorf("rule '%s' already exists with version %s", identifier.FullName, version)
	}
	
	// Create registry client
	authConfig := auth.LoadAuthConfig()
	client := registry.NewClient(cfg.RegistryURL)
	client.SetAuthToken(authConfig.AccessToken)
	
	// Download rule and get the actual version
	actualVersion, err := downloadRule(client, identifier, rulesDir)
	if err != nil {
		return fmt.Errorf("failed to download rule: %w", err)
	}
	
	// Add rule to ruleset using the full name and actual version
	rs.AddRule(identifier.FullName, actualVersion)
	if err := rs.SaveRuleSet(rulesJSONPath); err != nil {
		return fmt.Errorf("failed to save ruleset: %w", err)
	}
	
	color.Green("Rule '%s' (version %s) added successfully", identifier.FullName, actualVersion)
	
	// Print format suggestion at the very end if applicable
	if formatSuggestion != "" {
		fmt.Println(formatSuggestion)
	}
	
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}