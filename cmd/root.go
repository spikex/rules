package cmd

import (
	"os"

	"rules-cli/internal/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version represents the current version of the CLI
// This will be set during build via ldflags
var Version = "dev"

var (
	cfgFile string
	cfg     *config.Config
	format  string
	version bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage rule sets for AI code assistants",
	Long: `A command-line tool to create, manage, and convert rule sets 
for code guidance across different AI assistant platforms 
(Continue, Cursor, Windsurf, Copilot, etc.).`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			color.Cyan("rules version %s", Version)
			return
		}
		// If no subcommand is specified and no version flag, show help
		cmd.Help()
	},
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
	
	// Version flag
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Display version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	cfg, err = config.Initialize()
	if err != nil {
		color.Red("Error initializing config: %v", err)
		os.Exit(1)
	}

	// If format not specified, use default from config
	if format == "" {
		format = cfg.DefaultFormat
	}
}