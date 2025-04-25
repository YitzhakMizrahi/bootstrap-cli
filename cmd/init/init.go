// Package init provides initialization functionality for the bootstrap-cli,
// handling initial system setup and configuration extraction.
package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/spf13/cobra"
)

var (
	logger *log.Logger
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize bootstrap-cli configuration",
		Long: `Initialize bootstrap-cli by:
- Creating configuration directory
- Extracting default configurations
- Setting up environment variables`,
		RunE: runInit,
	}
	return cmd
}

func runInit(cmd *cobra.Command, _ []string) error {
	logger = log.New(log.InfoLevel)
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetLevel(log.DebugLevel)
	}
	logger.Info("Initializing Bootstrap CLI...")

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create config directory
	configDir := filepath.Join(home, ".config", "bootstrap-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set environment variable
	if err := os.Setenv("BOOTSTRAP_CLI_CONFIG", configDir); err != nil {
		return fmt.Errorf("failed to set BOOTSTRAP_CLI_CONFIG: %w", err)
	}

	// Extract default configurations
	configLoader := config.NewLoader(configDir)
	if err := configLoader.ExtractDefaults(); err != nil {
		return fmt.Errorf("failed to extract default configurations: %w", err)
	}

	logger.Success("Bootstrap CLI initialized successfully!")
	logger.Info("Run 'bootstrap-cli up' to start configuring your development environment")

	return nil
} 