// Package up provides the command for running the complete bootstrap process,
// orchestrating system detection, shell configuration, tool installation,
// and environment setup in a single unified flow.
package up

import (
	"fmt"
	"os"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/app"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// NewUpCmd creates a new up command
func NewUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run the complete bootstrap process",
		Long: `Run the complete bootstrap process including:
- System detection
- Shell detection and configuration
- Core tool installation
- Package management setup`,
		RunE: runUp,
	}
}

func runUp(_ *cobra.Command, _ []string) error {
	// Get config path from environment
	configPath := os.Getenv("BOOTSTRAP_CLI_CONFIG")
	if configPath == "" {
		return fmt.Errorf("BOOTSTRAP_CLI_CONFIG environment variable not set")
	}

	// Initialize config loader with the correct path
	configLoader := config.NewLoader(configPath)

	// Create and run the Bubble Tea application
	model := app.New(configLoader)
	p := tea.NewProgram(model)
	
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("UI error: %w", err)
	}

	// Get the final model state
	m, ok := finalModel.(*app.Model)
	if !ok {
		return fmt.Errorf("invalid model type")
	}

	// Check if user quit early
	if m.CurrentScreen() != app.FinishScreen {
		fmt.Println("Setup cancelled by user")
		return nil
	}

	// Get selected items
	selectedTools := m.SelectedTools()
	selectedFonts := m.SelectedFonts()
	selectedLanguages := m.SelectedLanguages()
	selectedManagers := m.SelectedManagers()

	// Step 9: Install Progress
	logger := log.New(log.InfoLevel)

	// Check for root privileges before proceeding with installation
	if !system.IsRoot() {
		logger.Warn("Tool installation requires root privileges.")
		fmt.Printf("\nPlease re-run with sudo: sudo %s\n\n", strings.Join(os.Args, " "))
		return fmt.Errorf("root privileges required for installation")
	}

	// Use the factory to get the package manager
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	installer := install.NewInstaller(pm)
	installer.Logger = logger

	// First install language managers
	logger.Info("Installing selected language managers...")
	for _, manager := range selectedManagers {
		logger.Info("Installing %s...", manager.Name)
		if err := installer.Install(&interfaces.Tool{
			Name: manager.Name,
			Description: manager.Description,
		}); err != nil {
			logger.Error("Failed to install language manager %s: %v. Continuing...", manager.Name, err)
		}
	}
	logger.Success("Language manager installation process finished.")

	// Then install languages
	logger.Info("Installing selected languages...")
	for _, lang := range selectedLanguages {
		logger.Info("Installing %s...", lang.Name)
		if err := installer.Install(&interfaces.Tool{
			Name: lang.Name,
			Description: lang.Description,
		}); err != nil {
			logger.Error("Failed to install language %s: %v. Continuing...", lang.Name, err)
		}
	}
	logger.Success("Language installation process finished.")

	// Finally install other selected tools
	logger.Info("Installing selected tools...")
	for _, tool := range selectedTools {
		logger.Info("Installing %s...", tool.Name)
		if err := installer.Install(&interfaces.Tool{
			Name: tool.Name,
			Description: tool.Description,
		}); err != nil {
			logger.Error("Failed to install tool %s: %v. Continuing...", tool.Name, err)
		}
	}
	logger.Success("Tool installation process finished.")

	// Install selected fonts
	if len(selectedFonts) > 0 {
		logger.Info("Installing selected fonts...")
		fontInstaller := install.NewFontInstaller(logger)
		for _, font := range selectedFonts {
			logger.Info("Installing %s...", font.Name)
			if err := fontInstaller.InstallFont(font); err != nil {
				logger.Error("Failed to install font %s: %v", font.Name, err)
			} else {
				logger.Success("%s installed successfully.", font.Name)
			}
		}
	}

	fmt.Println("Bootstrap process completed successfully!")
	return nil
} 