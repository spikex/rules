package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"rules-cli/internal/config"
)

var (
	cfgFile string
	cfg     *config.Config
	format  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage rule sets for AI code assistants",
	Long: `A command-line tool to create, manage, and convert rule sets 
for code guidance across different AI assistant platforms 
(Continue, Cursor, Claude Code, Copilot, etc.).`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rules-cli/rules-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&format, "format", "", "rule format (default is set in config)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	cfg, err = config.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// If format not specified, use default from config
	if format == "" {
		format = cfg.DefaultFormat
	}
}