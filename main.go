// Package main is the entry point for the bootstrap-cli application,
// handling initial configuration extraction and command execution setup.
package main

import (
	"log"
	"os"

	"github.com/YitzhakMizrahi/bootstrap-cli/cmd"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
)

func main() {
	// Create a temporary directory for extracted configs
	tempDir, err := os.MkdirTemp("", "bootstrap-cli-config-*")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up on exit

	// Create config loader and extract configs
	configLoader := config.NewLoader(tempDir)
	if err := configLoader.ExtractEmbeddedConfigs(); err != nil {
		log.Fatalf("Failed to extract embedded configs: %v", err)
	}

	// Set the config path in the environment
	os.Setenv("BOOTSTRAP_CLI_CONFIG", tempDir)

	// Execute the root command
	cmd.Execute()
} 