package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	// Get rule info
	ruleInfo, err := c.GetRule(name, version)
	if err != nil {
		return err
	}
	
	// Create rule directory
	ruleDir := filepath.Join(formatDir, name)
	if err := os.MkdirAll(ruleDir, 0755); err != nil {
		return fmt.Errorf("failed to create rule directory: %w", err)
	}
	
	// Download each file
	for _, file := range ruleInfo.Files {
		fileURL := fmt.Sprintf("%s/rules/%s/%s/files/%s", c.BaseURL, name, version, file)
		
		// Create directory for file if needed
		fileDir := filepath.Dir(filepath.Join(ruleDir, file))
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for file %s: %w", file, err)
		}
		
		// Download file
		resp, err := http.Get(fileURL)
		if err != nil {
			return fmt.Errorf("failed to download file %s: %w", file, err)
		}
		
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("failed to download file %s: status %d", file, resp.StatusCode)
		}
		
		// Create file
		outFile, err := os.Create(filepath.Join(ruleDir, file))
		if err != nil {
			resp.Body.Close()
			return fmt.Errorf("failed to create file %s: %w", file, err)
		}
		
		// Copy contents
		_, err = io.Copy(outFile, resp.Body)
		outFile.Close()
		resp.Body.Close()
		
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", file, err)
		}
	}
	
	return nil
}