package registry

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"rules-cli/internal/utils"
	"strings"
)

// Client represents a registry client
type Client struct {
	BaseURL    string
	AuthToken  string
	IsLoggedIn bool
}

// RuleInfo contains information about a rule in the registry
type RuleInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
}

// PublishMetadata represents the metadata for publishing a rule
type PublishMetadata struct {
	Visibility string `json:"visibility"`
}

// UserInfo represents user information from the registry
type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	OrgSlug  string `json:"orgSlug"`
}

// NewClient creates a new registry client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		AuthToken:  "",
		IsLoggedIn: false,
	}
}

// SetAuthToken sets the auth token for API requests
func (c *Client) SetAuthToken(token string) {
	c.AuthToken = token
	c.IsLoggedIn = token != ""
}

// DownloadRule downloads a rule to the specified directory
func (c *Client) DownloadRule(ownerSlug, ruleSlug, version, formatDir string) error {
	// Check if this is a GitHub repository
	if strings.HasPrefix(ownerSlug, "gh:") {
		return c.downloadFromGitHub(ownerSlug[3:]+"/"+ruleSlug, formatDir)
	}
	
	// Use the registry API download endpoint
	url := fmt.Sprintf("%s/v0/%s/%s/latest/download", c.BaseURL, ownerSlug, ruleSlug)
	if version != "latest" && version != "" {
		url = fmt.Sprintf("%s/v0/%s/%s/%s/download", c.BaseURL, ownerSlug, ruleSlug, version)
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add auth header if logged in
	if c.IsLoggedIn {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	}
	
	utils.SetUserAgent(req)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request rule from registry API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch rule from registry API: status %d", resp.StatusCode)
	}
	
	// Read the zip file
	zipData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read zip data: %w", err)
	}
	
	// Create a reader for the zip file
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to parse zip archive: %w", err)
	}
	
	// Create rule directory
	ruleDir := filepath.Join(formatDir, ownerSlug, ruleSlug)
	if err := os.MkdirAll(ruleDir, 0755); err != nil {
		return fmt.Errorf("failed to create rule directory: %w", err)
	}
	
	// Extract files from the zip
	for _, file := range zipReader.File {
		// Skip directories
		if file.FileInfo().IsDir() {
			continue
		}
		
		// Open the file
		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file from archive: %w", err)
		}
		
		// Get destination path
		destPath := filepath.Join(ruleDir, file.Name)
		
		// Create directory for file if needed
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			src.Close()
			return fmt.Errorf("failed to create directory: %w", err)
		}
		
		// Create the file
		dest, err := os.Create(destPath)
		if err != nil {
			src.Close()
			return fmt.Errorf("failed to create file: %w", err)
		}
		
		// Copy the content
		_, err = io.Copy(dest, src)
		src.Close()
		dest.Close()
		
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}
	
	return nil
}

// PublishRule publishes a new version of a rule to the registry
func (c *Client) PublishRule(ruleSlug, version, zipFilePath string, visibility string) error {
	if !c.IsLoggedIn {
		return fmt.Errorf("you must be logged in to publish a rule")
	}
	
	url := fmt.Sprintf("%s/v0/%s/%s", c.BaseURL, ruleSlug, version)
	
	// Read the zip file
	zipData, err := ioutil.ReadFile(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}
	
	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Add the zip file
	fileWriter, err := writer.CreateFormFile("file", filepath.Base(zipFilePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	
	if _, err := fileWriter.Write(zipData); err != nil {
		return fmt.Errorf("failed to write zip data to form: %w", err)
	}
	
	// Add metadata
	metadata := PublishMetadata{
		Visibility: visibility,
	}
	
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	
	if err := writer.WriteField("metadata", string(metadataJSON)); err != nil {
		return fmt.Errorf("failed to write metadata field: %w", err)
	}
	
	writer.Close()
	
	// Create request
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return fmt.Errorf("failed to create publish request: %w", err)
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	
	utils.SetUserAgent(req)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to publish rule: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to publish rule: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}
	
	return nil
}

// downloadFromGitHub downloads rules from a GitHub repository
func (c *Client) downloadFromGitHub(repoPath string, formatDir string) error {
	// Construct GitHub API URL to download zip of the main branch
	url := fmt.Sprintf("https://api.github.com/repos/%s/zipball/main", repoPath)
	
	// Create HTTP request with appropriate headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	utils.SetUserAgent(req)
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download GitHub repository: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download GitHub repository: status %d", resp.StatusCode)
	}
	
	// Read the response body
	zipData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read repository data: %w", err)
	}
	
	// Create a reader for the zip file
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to parse repository archive: %w", err)
	}
	
	// Create rule directory
	ruleDir := filepath.Join(formatDir, "gh:"+repoPath)
	if err := os.MkdirAll(ruleDir, 0755); err != nil {
		return fmt.Errorf("failed to create rule directory: %w", err)
	}
	
	// Extract rules from the zip
	foundRules := false
	repoPrefix := ""
	downloadedFiles := []string{}
	
	// First, determine the repository root directory name (it usually includes a commit hash)
	for _, file := range zipReader.File {
		parts := strings.Split(file.Name, "/")
		if len(parts) > 0 {
			repoPrefix = parts[0]
			break
		}
	}
	
	if repoPrefix == "" {
		return fmt.Errorf("could not determine repository structure")
	}
	
	// Look for files in the src/ directory
	srcPrefix := repoPrefix + "/src/"
	
	for _, file := range zipReader.File {
		// Check if file is in the src/ directory
		if strings.HasPrefix(file.Name, srcPrefix) {
			foundRules = true
			
			// Skip directories, we'll create them as needed
			if file.FileInfo().IsDir() {
				continue
			}
			
			// Open the file
			src, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file from archive: %w", err)
			}
			
			// Get destination path without the repository and src/ prefix
			relativePath := strings.TrimPrefix(file.Name, srcPrefix)
			destPath := filepath.Join(ruleDir, relativePath)
			
			// Create directory for file if needed
			destDir := filepath.Dir(destPath)
			if err := os.MkdirAll(destDir, 0755); err != nil {
				src.Close()
				return fmt.Errorf("failed to create directory: %w", err)
			}
			
			// Create the file
			dest, err := os.Create(destPath)
			if err != nil {
				src.Close()
				return fmt.Errorf("failed to create file: %w", err)
			}
			
			// Copy the content
			_, err = io.Copy(dest, src)
			src.Close()
			dest.Close()
			
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			
			downloadedFiles = append(downloadedFiles, relativePath)
		}
	}
	
	if !foundRules {
		// List all files in the repository to help with debugging
		var fileList strings.Builder
		fileList.WriteString("Files found in the repository:\n")
		for _, file := range zipReader.File {
			fileList.WriteString("  - " + file.Name + "\n")
		}
		return fmt.Errorf("no rules found in the src/ directory of the GitHub repository.\n%s", fileList.String())
	}
	
	return nil
}
