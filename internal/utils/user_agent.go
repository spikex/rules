package utils

import (
	"fmt"
	"net/http"
	"runtime"
)

// GetUserAgent returns the user agent string for HTTP requests
// This includes the CLI name, version, and runtime information
func GetUserAgent() string {
	// For now, use a hardcoded version. In the future, this could be injected
	// during the build process via ldflags
	version := "1.0.0"
	
	// Include runtime information for better analytics
	return fmt.Sprintf("rules-cli/%s (%s; %s; %s)", 
		version, 
		runtime.GOOS, 
		runtime.GOARCH, 
		runtime.Version())
}

// SetUserAgent sets the User-Agent header on an HTTP request
func SetUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", GetUserAgent())
} 