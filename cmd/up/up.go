// Package up provides the up command for running the bootstrap-cli TUI
package up

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	base_iface "github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces" // Base interfaces (like for UI selections)
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/pipeline" // Pipeline interfaces defined in pipeline package itself
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
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
		// Try default location if env var is not set
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(home, ".config", "bootstrap-cli")
		logger.Debug("BOOTSTRAP_CLI_CONFIG not set, using default: %s", configPath)
	}

	// Ensure base config directory exists (optional, loader might handle it)
	// if err := os.MkdirAll(configPath, 0755); err != nil {
	// 	logger.Warn("Could not create config directory %s: %v", configPath, err)
	// }

	// Initialize config loader with the correct path
	configLoader := config.NewLoader(configPath)

	// --- Run the TUI Application --- 
	appModel := app.New(configLoader)
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	finalModelInterface, err := p.Run()
	if err != nil {
		// Ensure terminal state is reset even on error
		_ = tea.ClearScreen() // Attempt to clear screen
		return fmt.Errorf("application error: %w", err)
	}
	logger.Info("TUI finished. Processing selections...")

	// --- Process Selections and Run Installation --- 
	m, ok := finalModelInterface.(*app.Model)
	if !ok {
		return fmt.Errorf("internal error: could not cast final model to *app.Model")
	}

	// Gather selections (selectedTools is now []*pipeline.Tool)
	selectedPipelineTools := m.SelectedTools()      
	manageDotfiles := m.GetManageDotfiles() // Use getter
	dotfilesRepoURL := m.GetDotfilesRepoURL() // Use getter
	// selectedFontInterfaces := m.SelectedFonts()        
	// selectedLanguageInterfaces := m.SelectedLanguages() 
	selectedShellInterface := m.GetSelectedShell()     
	// TODO: Get selected dotfiles

	// Early exit if nothing was selected
	if len(selectedPipelineTools) == 0 && !manageDotfiles /* && other selections empty */ && selectedShellInterface == nil {
		logger.Info("No items selected for installation or configuration. Exiting.")
		return nil
	}

	// Tool definitions are now correctly loaded in selectedPipelineTools from the UI model.
	// No extra loading/filtering needed here.

	// Detect system platform and package manager
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system info for installation: %w", err)
	}
	pkgManagerFactory := factory.NewPackageManagerFactory()
	pkgManagerImpl, err := pkgManagerFactory.GetPackageManager() // base_iface.PackageManager
	if err != nil {
		return fmt.Errorf("failed to detect package manager for installation: %w", err)
	}

	// Adapt the base PackageManager to the pipeline's PackageManager interface
	var pipelinePackageManager pipeline.PackageManager // Use pipeline's interface
	pipelinePackageManager = &packageManagerAdapter{impl: pkgManagerImpl} 
	// fmt.Println("TODO: Verify and complete PackageManager adapter implementation for pipeline.") // Remove TODO Print

	pipelinePlatform := &pipeline.Platform{
		OS:             sysInfo.OS,
		Arch:           sysInfo.Arch,
		PackageManager: pkgManagerImpl.GetName(), // Get name from the base implementation before adapting
		Shell:          sysInfo.Shell,
	}
	// fmt.Println("TODO: Ensure this is the correct way to set PackageManager for the pipeline context") // Remove TODO Print

	// Create the installer
	installer, err := pipeline.NewInstaller(pipelinePlatform, pipelinePackageManager)
	if err != nil {
		return fmt.Errorf("failed to create installer: %w", err)
	}

	// TODO: Adapt Fonts, Languages, Shell, Dotfiles for installation

	if len(selectedPipelineTools) > 0 || manageDotfiles { // Use local var
		logger.Info("Starting installation process...")
		// Pass dotfiles selections to the installer
		if err := installer.InstallSelections(selectedPipelineTools, manageDotfiles, dotfilesRepoURL); err != nil { // Use local vars 
			return fmt.Errorf("installation failed: %w", err)
		}
		logger.Info("Installation phase complete.")
	} else {
		logger.Info("No tools or dotfiles selected for installation.")
	}

	// TODO: Add similar blocks for Fonts, Languages, Shell setup

	if selectedShellInterface != nil {
		logger.Info("Configuring selected shell: %s", selectedShellInterface.Name)
		// TODO: Implement shell configuration logic
	}

	logger.Info("Bootstrap setup process finished.")
	return nil
} 

// Placeholder adapter - NEEDS REAL IMPLEMENTATION and matching interfaces defined
// Adapter implementation to bridge interfaces.PackageManager and pipeline.PackageManager
type packageManagerAdapter struct {
	impl base_iface.PackageManager // The implementation from internal/packages
}

func (a *packageManagerAdapter) Install(pkg string) error { return a.impl.Install(pkg) }
func (a *packageManagerAdapter) Uninstall(pkg string) error { return a.impl.Uninstall(pkg) } // Use renamed Uninstall
func (a *packageManagerAdapter) IsInstalled(pkg string) (bool, error) {
	// Now directly call the method with the correct signature
	return a.impl.IsInstalled(pkg)
}
func (a *packageManagerAdapter) Update() error { return a.impl.Update() }
func (a *packageManagerAdapter) SetupSpecialPackage(pkg string) error { 
	// Assuming base interface now has this method (verify if needed)
	return a.impl.SetupSpecialPackage(pkg) 
}
func (a *packageManagerAdapter) IsPackageAvailable(pkg string) bool { 
	// Now call the method added to the base interface
	return a.impl.IsPackageAvailable(pkg) 
}
func (a *packageManagerAdapter) GetName() string { 
	// Now call the method from the base interface
	return a.impl.GetName() 
}

// mapUIToolToPipelineTool removed as we now load pipeline.Tool directly via configLoader 