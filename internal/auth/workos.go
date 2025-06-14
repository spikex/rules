// auth/workos.go
package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

// AuthConfig holds authentication information
type AuthConfig struct {
	UserID       string `json:"userId,omitempty"`
	UserEmail    string `json:"userEmail,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresAt    int64  `json:"expiresAt,omitempty"`
}

// Config constants
var (
	homedir, _     = os.UserHomeDir()
	authConfigPath = filepath.Join(homedir, ".continue", "auth.json")
)

// LoadAuthConfig loads the authentication configuration from disk
func LoadAuthConfig() AuthConfig {
	// If CONTINUE_API_KEY environment variable exists, use that instead
	if apiKey := os.Getenv("CONTINUE_API_KEY"); apiKey != "" {
		return AuthConfig{
			AccessToken: apiKey,
		}
	}

	// Check if auth config file exists
	data, err := ioutil.ReadFile(authConfigPath)
	if err != nil {
		return AuthConfig{}
	}

	var config AuthConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error loading auth config: %v\n", err)
		return AuthConfig{}
	}

	return config
}

// SaveAuthConfig saves the authentication configuration to disk
func SaveAuthConfig(config AuthConfig) {
	// If using CONTINUE_API_KEY environment variable, don't save anything
	if os.Getenv("CONTINUE_API_KEY") != "" {
		return
	}

	// Make sure the directory exists
	dir := filepath.Dir(authConfigPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// Marshal the config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error saving auth config: %v\n", err)
		return
	}

	// Write to file
	if err := ioutil.WriteFile(authConfigPath, data, 0644); err != nil {
		fmt.Printf("Error saving auth config: %v\n", err)
	}
}

// IsAuthenticated checks if the user is authenticated and the token is valid
func IsAuthenticated() bool {
	// If CONTINUE_API_KEY environment variable exists, user is authenticated
	if os.Getenv("CONTINUE_API_KEY") != "" {
		return true
	}

	config := LoadAuthConfig()

	if config.UserID == "" || config.AccessToken == "" {
		return false
	}

	// Check if token is expired (if we have an expiration)
	if config.ExpiresAt > 0 && time.Now().UnixMilli() > config.ExpiresAt {
		// Try refreshing the token
		_, err := RefreshToken(config.RefreshToken)
		if err != nil {
			// If refresh fails, we're not authenticated
			return false
		}
	}

	return true
}

// Prompt asks the user for input
func Prompt(question string) (string, error) {
	fmt.Print(question)
	var answer string
	_, err := fmt.Scanln(&answer)
	return answer, err
}

// GetAuthUrlForTokenPage returns the auth URL for the token page
func GetAuthUrlForTokenPage() string {
	baseURL := "https://api.workos.com/user_management/authorize"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", viper.GetString("workos_client_id"))
	
	redirectPath := "tokens/callback/rules"
	params.Add("redirect_uri", fmt.Sprintf("%s%s", viper.GetString("app_url"), redirectPath))
	
	params.Add("state", uuid.New().String())
	params.Add("provider", "authkit")

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// RefreshToken refreshes the access token using a refresh token
func RefreshToken(refreshToken string) (AuthConfig, error) {
	type refreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}

	type user struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	type refreshResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		User         user   `json:"user"`
	}

	apiBase := viper.GetString("api_base")
	
	reqBody, err := json.Marshal(refreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return AuthConfig{}, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%sauth/refresh", apiBase),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return AuthConfig{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AuthConfig{}, fmt.Errorf("refresh token failed with status: %s", resp.Status)
	}

	var response refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return AuthConfig{}, err
	}

	// Calculate token expiration (assuming 1 hour validity)
	tokenExpiresAt := time.Now().Add(1 * time.Hour).UnixMilli()

	authConfig := AuthConfig{
		UserID:       response.User.ID,
		UserEmail:    response.User.Email,
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresAt:    tokenExpiresAt,
	}

	// Save the config
	SaveAuthConfig(authConfig)

	return authConfig, nil
}

// Login authenticates using the Continue web flow
func Login() (AuthConfig, error) {
	// If CONTINUE_API_KEY environment variable exists, use that instead
	if apiKey := os.Getenv("CONTINUE_API_KEY"); apiKey != "" {
		color.Green("Using CONTINUE_API_KEY from environment variables")
		return AuthConfig{
			AccessToken: apiKey,
		}, nil
	}

	color.Cyan("\nStarting authentication with Continue...")

	// Get auth URL
	authURL := GetAuthUrlForTokenPage()
	color.Green("Opening browser to sign in at: %s", authURL)
	if err := browser.OpenURL(authURL); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
		fmt.Printf("Please manually open: %s\n", authURL)
	}

	color.Yellow("\nAfter signing in, you'll receive a token.")

	// Get token from user
	token, err := Prompt(color.YellowString("Paste your sign-in token here: "))
	if err != nil {
		return AuthConfig{}, err
	}

	color.Cyan("Verifying token...")

	// Exchange token for session
	response, err := RefreshToken(token)
	if err != nil {
		return AuthConfig{}, errors.New("authentication failed: " + err.Error())
	}

	color.Green("\nAuthentication successful!")

	return response, nil
}

// Logout logs the user out by clearing saved credentials
func Logout() {
	if os.Getenv("CONTINUE_API_KEY") != "" {
		color.Yellow("Using CONTINUE_API_KEY from environment variables, nothing to log out")
		return
	}

	if _, err := os.Stat(authConfigPath); err == nil {
		os.Remove(authConfigPath)
		color.Green("Successfully logged out")
	} else {
		color.Yellow("No active session found")
	}
}