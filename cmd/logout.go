package cmd

import (
	"rules-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the registry service",
	Long:  `Logs the user out by removing stored authentication information.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Logging out...")
		// Call the logout function from the auth package
		auth.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}