package registry

import (
	"fmt"

	"github.com/fatih/color"

	"rules-cli/internal/auth"
)

// InitClientWithAuth initializes a registry client with authentication
func InitClientWithAuth(baseURL string) (*Client, error) {
	client := NewClient(baseURL)
	
	// Load auth config
	authConfig := auth.LoadAuthConfig()
	
	// If we have a valid token, use it
	if authConfig.AccessToken != "" {
		client.SetAuthToken(authConfig.AccessToken)
		return client, nil
	}
	
	return client, nil
}

// EnsureClientAuth ensures the client has authentication
// Returns a boolean indicating if authentication was successful
func EnsureClientAuth(client *Client, requireAuth bool) (bool, error) {
	if client.IsLoggedIn {
		return true, nil
	}
	
	if !requireAuth {
		return false, nil
	}
	
	// Start login flow
	color.Yellow("Authentication required to perform this operation.")
	authConfig, err := auth.Login(false)
	if err != nil {
		color.Red("Authentication failed: %v", err)
		return false, err
	}
	
	// Set the token on the client
	client.SetAuthToken(authConfig.AccessToken)
	return true, nil
}

// GetAuthenticatedClient gets a registry client with authentication
func GetAuthenticatedClient(baseURL string, requireAuth bool) (*Client, error) {
	client, err := InitClientWithAuth(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize registry client: %w", err)
	}
	
	_, err = EnsureClientAuth(client, requireAuth)
	if err != nil {
		return nil, err
	}
	
	return client, nil
}