package registry

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Client represents a registry client
type Client struct {
	BaseURL string
}

// RuleInfo contains information about a rule in the registry
type RuleInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
}

// RegistryResponse represents the response from the registry API
type RegistryResponse struct {
	Content string `json:"content"`
}

// NewClient creates a new registry client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

// GetRule fetches a rule from the registry
func (c *Client) GetRule(name, version string) (*RuleInfo, error) {
	url := fmt.Sprintf("%s/rules/%s/%s", c.BaseURL, name, version)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request rule: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch rule: status %d", resp.StatusCode)
	}
	
	var ruleInfo RuleInfo
	if err := json.NewDecoder(resp.Body).Decode(&ruleInfo); err != nil {
		return nil, fmt.Errorf("failed to decode rule info: %w", err)
	}
	
	return &ruleInfo, nil
}

// DownloadRule downloads a rule to the specified directory
func (c *Client) DownloadRule(name, version, formatDir string) error {
	// Check if this is a GitHub repository
	if strings.HasPrefix(name, "gh:") {
		return c.downloadFromGitHub(name[3:], formatDir)
	}
	
	// Use the registry API GET endpoint to fetch rule content
	url := fmt.Sprintf("%s/registry/v1/%s/latest", c.BaseURL, name)
	if version != "latest" && version != "" {
		url = fmt.Sprintf("%s/registry/v1/%s/%s", c.BaseURL, name, version)
	}
	
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to request rule from registry API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch rule from registry API: status %d", resp.StatusCode)
	}
	
	var registryResponse RegistryResponse
	if err := json.NewDecoder(resp.Body).Decode(&registryResponse); err != nil {
		return fmt.Errorf("failed to decode registry API response: %w", err)
	}
	
	// Create rule directory
	ruleDir := filepath.Join(formatDir, name)
	if err := os.MkdirAll(ruleDir, 0755); err != nil {
		return fmt.Errorf("failed to create rule directory: %w", err)
	}
	
	// Write the main rule file
	rulePath := filepath.Join(ruleDir, "index.md")
	if err := ioutil.WriteFile(rulePath, []byte(registryResponse.Content), 0644); err != nil {
		return fmt.Errorf("failed to write rule file: %w", err)
	}
	
	// Print summary of downloaded rule
	fmt.Printf("Successfully downloaded rule '%s' (version: %s)\n", name, version)
	fmt.Printf("Rule saved to: %s\n", rulePath)
	
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
	
	// Print summary of downloaded files
	fmt.Printf("Successfully downloaded rules from GitHub repository: %s\n", repoPath)
	fmt.Printf("Downloaded %d files to: %s\n", len(downloadedFiles), ruleDir)
	for _, file := range downloadedFiles {
		fmt.Printf("  - %s\n", file)
	}
	
	return nil
}