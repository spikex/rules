package utils

import (
	"fmt"
	"net/http"
	"runtime"
)

// Version will be set during build with ldflags
var Version = "dev"

// GetUserAgent returns the user agent string for HTTP requests
// This includes the CLI name, version, and runtime information
func GetUserAgent() string {
	// Include runtime information for better analytics
	return fmt.Sprintf("rules-cli/%s (%s; %s; %s)", 
		Version, 
		runtime.GOOS, 
		runtime.GOARCH, 
		runtime.Version())
}

// SetUserAgent sets the User-Agent header on an HTTP request
func SetUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", GetUserAgent())
} 