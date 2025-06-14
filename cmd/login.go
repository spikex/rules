package cmd

import (
	"fmt"

	"rules-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the registry service",
	Long:  `Starts the authorization flow to authenticate with the registry service.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting login process...")
		
		// Call the login function from the auth package
		authConfig, err := auth.Login()
		if err != nil {
			color.Red("Login failed: %v", err)
			return
		}
		
		if authConfig.UserEmail != "" {
			color.Green("Successfully logged in as %s", authConfig.UserEmail)
		} else {
			color.Green("Successfully authenticated")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}