package cmd

import (
	"fmt"
	"os"

	"rules-cli/internal/auth"
	"rules-cli/internal/registry"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display information about the currently authenticated user",
	Long:  `Displays information about the currently authenticated user, including username, email, and organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the user is authenticated
		if !auth.IsAuthenticated() {
			color.Yellow("You are not currently authenticated.")
			fmt.Println("Use 'rules login' to authenticate.")
			return
		}

		// Load auth config to get user information
		authConfig := auth.LoadAuthConfig()

		// Display user information
		color.Green("Authenticated User:")

		if authConfig.UserID != "" {
			fmt.Printf("User ID: %s\n", authConfig.UserID)
		}

		if authConfig.UserEmail != "" {
			fmt.Printf("Email: %s\n", authConfig.UserEmail)
		}

		// If using API key from environment variable
		if apiKey := fmt.Sprintf("%s", authConfig.AccessToken); apiKey != "" && authConfig.UserEmail == "" {
			if os.Getenv("CONTINUE_API_KEY") != "" {
				fmt.Println("Authentication method: Environment variable (CONTINUE_API_KEY)")
			} else {
				fmt.Println("Authentication method: Access token")
			}
		}

		// Attempt to get additional user info from the registry API
		apiBase := viper.GetString("api_base")
		client, err := registry.GetAuthenticatedClient(apiBase, false)
		if err == nil && client.IsLoggedIn {
			// Here we would make an API call to get more user information
			// This would typically include organization details
		}
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
