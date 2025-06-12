package cmd

import (
	"github.com/spf13/cobra"
	"rules-cli/internal/auth"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the registry service",
	Long:  `Logs the user out by removing stored authentication information.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the logout function from the auth package
		auth.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}