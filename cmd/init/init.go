// Package init provides initialization functionality for the bootstrap-cli,
// handling initial system setup, tool installation, and configuration.
package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/utils"

	"github.com/spf13/cobra"
)

var (
	logger *log.Logger
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new development environment",
		Long: `Guides you through selecting and installing components
for your development environment, including:
- Core development tools
- Modern CLI utilities
- Shell setup (coming soon)
- Dotfiles management (coming soon)`,
		RunE: runInit,
	}
	// Add flags here if needed later
	return cmd
}

func runInit(cmd *cobra.Command, _ []string) error {
	logger = log.New(log.InfoLevel)
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetLevel(log.DebugLevel)
	}
	logger.Info("Starting Bootstrap CLI initialization...")

	// Define installation sections
	sections := []struct {
		Name   string
		Status string
	}{
		{Name: "System Detection", Status: "pending"},
		{Name: "Tool Selection", Status: "pending"},
		{Name: "Font Selection", Status: "pending"},
		{Name: "Language Selection", Status: "pending"},
		{Name: "Dotfiles Setup", Status: "pending"},
	}

	// Get config path from environment
	configPath := os.Getenv("BOOTSTRAP_CLI_CONFIG")
	if configPath == "" {
		return fmt.Errorf("BOOTSTRAP_CLI_CONFIG environment variable not set")
	}

	// Initialize config loader with the correct path
	configLoader := config.NewLoader(configPath)

	// Update progress for system detection
	sections[0].Status = "current"
	fmt.Println(utils.RenderProgressBar(sections, 80))

	// Detect system info and package manager early as we'll need it for multiple phases
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}
	
	// Initialize package manager
	pmFactory := factory.NewPackageManagerFactory()
	packageManager, err := pmFactory.GetPackageManager()
	if err != nil {
		return fmt.Errorf("failed to initialize package manager: %w", err)
	}
	logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
	logger.Info("Package Manager: %s", packageManager.GetName())

	// Mark system detection as completed and update tool selection
	sections[0].Status = "completed"
	sections[1].Status = "current"
	fmt.Println(utils.RenderProgressBar(sections, 80))

	// --- Phase 1: Tool Selection ---
	logger.Info("Step 1: Select tools to install")
	
	// Load and select tools
	selectedTools, err := ui.PromptToolSelection(configLoader)
	if err != nil {
		return fmt.Errorf("failed to select tools: %w", err)
	}

	// Check if user cancelled
	if len(selectedTools) == 0 {
		fmt.Println("No tools selected. Skipping tool installation.")
		return nil
	}

	// Install selected tools
	if err := installTools(selectedTools, packageManager); err != nil {
		return fmt.Errorf("failed to install tools: %w", err)
	}

	// Configure automatic restarts for services
	if err := configureServices(selectedTools); err != nil {
		return fmt.Errorf("failed to configure services: %w", err)
	}

	// Mark tool selection as completed and update font selection
	sections[1].Status = "completed"
	sections[2].Status = "current"
	fmt.Println(utils.RenderProgressBar(sections, 80))

	// --- Phase 3: Font Selection ---
	logger.Info("Step 3: Select fonts to install")
	
	// Use the new font screen
	fontScreen := screens.NewFontScreen(configLoader)
	selectedFonts, err := fontScreen.ShowFontSelection()
	if err != nil {
		return fmt.Errorf("font selection failed: %w", err)
	}

	// Install selected fonts
	if len(selectedFonts) > 0 {
		logger.Info("Installing selected fonts...")
		fontInstaller := install.NewFontInstaller(logger)
		for _, font := range selectedFonts {
			if err := fontInstaller.InstallFont(font); err != nil {
				logger.Error("Failed to install font %s: %v", font.Name, err)
			}
		}
	}

	// Mark font selection as completed and update language selection
	sections[2].Status = "completed"
	sections[3].Status = "current"
	fmt.Println(utils.RenderProgressBar(sections, 80))

	// --- Phase 4: Language Selection ---
	logger.Info("Step 4: Select languages to install")
	
	// Load available languages
	availableLanguages, err := configLoader.LoadLanguages()
	if err != nil {
		return fmt.Errorf("failed to load language configurations: %w", err)
	}

	// Use the new language screen
	languageScreen := screens.NewLanguageScreen()
	selectedLanguages, err := languageScreen.ShowLanguageSelection(availableLanguages)
	if err != nil {
		logger.Error("Failed to get language selection: %v", err)
	} else if len(selectedLanguages) > 0 {
		// Then, select language managers based on selected languages
		availableManagers, err := configLoader.LoadLanguageManagers()
		if err != nil {
			logger.Error("Failed to load language manager configurations: %v", err)
		} else {
			selectedManagers, err := languageScreen.ShowManagerSelection(availableManagers, selectedLanguages)
			if err != nil {
				logger.Error("Failed to get language manager selection: %v", err)
			} else if len(selectedManagers) > 0 {
				// Install language managers first
				logger.Info("Installing selected language managers...")
				for _, manager := range selectedManagers {
					if err := install.NewInstaller(packageManager).Install(manager); err != nil {
						logger.Error("Failed to install language manager %s: %v", manager.Name, err)
					}
				}

				// Install selected languages
				logger.Info("Installing selected languages...")
				for _, lang := range selectedLanguages {
					if err := install.NewInstaller(packageManager).Install(lang.ToTool()); err != nil {
						logger.Error("Failed to install language %s: %v", lang.Name, err)
					}
				}
			}
		}
	}

	// Mark language selection as completed and update dotfiles setup
	sections[3].Status = "completed"
	sections[4].Status = "current"
	fmt.Println(utils.RenderProgressBar(sections, 80))

	// --- Phase 5: Dotfiles Setup ---
	logger.Info("Step 5: Setting up dotfiles")
	// TODO: Implement dotfiles setup

	// Mark dotfiles setup as completed
	sections[4].Status = "completed"
	fmt.Println(utils.RenderProgressBar(sections, 80))

	logger.Success("Bootstrap CLI initialization completed successfully!")
	return nil
}

func configureNeedrestart() error {
	configDir := "/etc/needrestart/conf.d"
	configFile := filepath.Join(configDir, "50-autorestart.conf")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create needrestart config directory: %w", err)
	}

	// Write configuration to automatically restart services
	config := []byte(`# Automatically restart services
$nrconf{restart} = 'a';
`)

	if err := os.WriteFile(configFile, config, 0644); err != nil {
		return fmt.Errorf("failed to write needrestart config: %w", err)
	}

	return nil
}

// installTools installs the selected tools using the package manager
func installTools(tools []*interfaces.Tool, packageManager interfaces.PackageManager) error {
	logger.Info("Installing selected tools...")

	// Check for root privileges before proceeding with installation
	if !system.IsRoot() {
		logger.Warn("Tool installation requires root privileges.")
		logger.Info("Attempting to relaunch with sudo...")
		fmt.Printf("\nPlease re-run the command with sudo: sudo %s\n\n", strings.Join(os.Args, " "))
		return fmt.Errorf("root privileges required for tool installation")
	}

	// Create tool installer
	installer := install.NewInstaller(packageManager)

	// Install each selected tool
	for _, tool := range tools {
		logger.Info("Installing %s...", tool.Name)
		installer.Logger = logger
		if err := installer.Install(tool); err != nil {
			logger.Error("Failed to install %s: %v", tool.Name, err)
			return err
		}
	}

	logger.Success("Tool installation completed successfully!")
	return nil
}

// configureServices sets up automatic restarts for services
func configureServices(tools []*interfaces.Tool) error {
	// Configure needrestart automatically
	if err := configureNeedrestart(); err != nil {
		logger.Warn("Failed to configure needrestart for automatic restarts: %v", err)
		return err
	}
	return nil
}

// Note: The Quitting() method needs to be implemented or accessible
// on the ui.ToolSelector struct for the cancellation check to work correctly.
// This might involve adding a method like `Finished() bool` to ui.ToolSelector
// that returns true if the UI exited normally (Enter on last category) vs. quit (Ctrl+C). 