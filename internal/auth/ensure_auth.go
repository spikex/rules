// auth/ensure_auth.go
package auth

import "github.com/fatih/color"

// EnsureAuthenticated ensures the user is authenticated before proceeding
// Returns true if authentication is successful, false otherwise
func EnsureAuthenticated(requireAuth bool) (bool, error) {
	if IsAuthenticated() {
		return true, nil
	}

	if !requireAuth {
		return false, nil
	}

	color.Yellow("Authentication required.")

	_, err := Login(false)
	if err != nil {
		color.Red("Failed to authenticate: %v", err)
		return false, err
	}

	return true, nil
}