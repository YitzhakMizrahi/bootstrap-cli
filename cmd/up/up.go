// Package up provides the up command for running the bootstrap-cli TUI
package up

import (
	"fmt"
	"os"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/app"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	logger *log.Logger
)

// NewUpCmd creates the up command
func NewUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Run the interactive TUI to configure and set up your development environment",
		Long: `Guides you through selecting and installing components
for your development environment, including:
- Core development tools
- Modern CLI utilities
- Programming languages
- Fonts
- Shell setup
- Dotfiles management`,
		RunE: runUp,
	}
	return cmd
}

func runUp(cmd *cobra.Command, _ []string) error {
	logger = log.New(log.InfoLevel)
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetLevel(log.DebugLevel)
	}
	logger.Info("Starting Bootstrap CLI TUI...")

	// Get config path from environment
	configPath := os.Getenv("BOOTSTRAP_CLI_CONFIG")
	if configPath == "" {
		return fmt.Errorf("BOOTSTRAP_CLI_CONFIG environment variable not set")
	}

	// Initialize config loader with the correct path
	configLoader := config.NewLoader(configPath)

	// Create and run the application
	model := app.New(configLoader)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("application error: %w", err)
	}

	// Get selections from the final model
	m := finalModel.(*app.Model)
	if len(m.SelectedTools()) > 0 {
		logger.Info("Installing selected tools...")
		// TODO: Install tools
	}

	if len(m.SelectedFonts()) > 0 {
		logger.Info("Installing selected fonts...")
		// TODO: Install fonts
	}

	if len(m.SelectedLanguages()) > 0 {
		logger.Info("Installing selected languages...")
		// TODO: Install languages and managers
	}

	return nil
} 