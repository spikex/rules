package main

import (
	"fmt"
	"os"

	"rules-cli/cmd"
)

// The Version variable will be set during build with ldflags
var Version = "dev"

func init() {
	// Pass the version to the cmd package
	cmd.Version = Version
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
