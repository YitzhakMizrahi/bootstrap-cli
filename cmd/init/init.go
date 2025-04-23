package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/tools"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
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

func runInit(cmd *cobra.Command, args []string) error {
	logger = log.New(log.InfoLevel)
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetLevel(log.DebugLevel)
	}
	logger.Info("Starting Bootstrap CLI initialization...")

	// Initialize config loader
	configLoader := config.NewConfigLoader(filepath.Join(os.TempDir(), "bootstrap-cli"))

	// Detect system info and package manager early as we'll need it for multiple phases
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}
	
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}
	logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
	logger.Info("Package Manager: %s", pm.Name())

	// --- Phase 1: Tool Selection ---
	logger.Info("Step 1: Select tools to install")
	availableTools, err := configLoader.LoadTools()
	if err != nil {
		return fmt.Errorf("failed to load tool configurations: %w", err)
	}

	// Extract tool names for the selector
	var toolNames []string
	toolMap := make(map[string]*install.Tool)
	for _, tool := range availableTools {
		toolNames = append(toolNames, tool.Name)
		toolMap[tool.Name] = tool
	}

	selector := ui.NewToolSelector(toolNames)
	p := tea.NewProgram(selector)

	// Run the interactive UI
	model, err := p.Run()
	if err != nil {
		return fmt.Errorf("UI error: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*ui.ToolSelector); ok {
		selectedNames := selectorModel.GetSelectedTools()
		if len(selectedNames) == 0 && !selectorModel.Finished() {
			logger.Info("Initialization cancelled or no tools selected.")
			return nil
		}

		// Convert selected names back to Tool objects
		var selectedTools []*install.Tool
		for _, name := range selectedNames {
			if tool, exists := toolMap[name]; exists {
				selectedTools = append(selectedTools, tool)
			}
		}

		// Store selected tools for the installer
		tools.SetSelectedTools(selectedTools)

		if len(selectedTools) == 0 {
			logger.Info("No tools selected. Skipping tool installation.")
		} else {
			// --- Phase 2: Installation (Requires Sudo) ---
			logger.Info("Step 2: Installing selected tools...")

			// Check for root privileges before proceeding with installation
			if !system.IsRoot() {
				logger.Warn("Tool installation requires root privileges.")
				logger.Info("Attempting to relaunch with sudo...")
				fmt.Printf("\nPlease re-run the command with sudo: sudo %s\n\n", strings.Join(os.Args, " "))
				return fmt.Errorf("root privileges required for tool installation")
			}

			// Configure needrestart automatically
			if err := configureNeedrestart(); err != nil {
				logger.Warn("Failed to configure needrestart for automatic restarts: %v", err)
			}

			opts := &tools.InstallOptions{
				Logger:           logger,
				PackageManager:   pm,
				Tools:            selectedTools,
				SkipVerification: false,
				AdditionalPaths:  []string{"/usr/bin", "/usr/local/bin", "/opt/homebrew/bin"},
			}

			if err := tools.InstallCoreTools(opts); err != nil {
				return fmt.Errorf("failed to install selected tools: %w", err)
			}

			logger.Success("Tool installation completed successfully!")
		}
	} else {
		return fmt.Errorf("failed to get tool selector model")
	}

	// --- Phase 3: Font Selection ---
	logger.Info("Step 3: Select fonts to install")
	_, err = configLoader.LoadFonts() // We'll use the fonts later when we implement the UI
	if err != nil {
		return fmt.Errorf("failed to load font configurations: %w", err)
	}

	// TODO: Add font selection UI similar to tools
	// For now, just install JetBrains Mono if user confirms
	if installFonts, err := ui.PromptFontInstallation(); err == nil && installFonts {
		fontInstaller := install.NewFontInstaller(logger)
		if err := fontInstaller.InstallJetBrainsMono(); err != nil {
			logger.Error("Failed to install JetBrains Mono font: %v", err)
		} else {
			logger.Success("Font installation completed successfully!")
		}
	}

	// --- Phase 4: Language Selection ---
	logger.Info("Step 4: Select programming languages to install")
	_, err = configLoader.LoadLanguages() // We'll use the languages later when we implement the UI
	if err != nil {
		return fmt.Errorf("failed to load language configurations: %w", err)
	}

	// TODO: Add language selection UI similar to tools
	// For now, prompt for common runtimes
	if selectedRuntimes, err := ui.PromptLanguageRuntimes(); err == nil && len(selectedRuntimes) > 0 {
		runtimeInstaller := install.NewRuntimeInstaller(pm, logger)
		for _, runtime := range selectedRuntimes {
			if err := runtimeInstaller.Install(runtime); err != nil {
				logger.Error("Failed to install runtime %s: %v", runtime, err)
			}
		}
		logger.Success("Language runtime installation completed successfully!")
	}

	// --- Phase 5: Dotfiles Setup ---
	logger.Info("Step 5: Configure dotfiles")
	_, err = configLoader.LoadDotfiles() // We'll use the dotfiles later when we implement the UI
	if err != nil {
		return fmt.Errorf("failed to load dotfile configurations: %w", err)
	}

	// TODO: Add dotfiles selection and setup UI
	// For now, just log that it's coming soon
	logger.Info("Dotfiles setup coming soon!")

	logger.Success("Bootstrap CLI initialization finished!")
	return nil
}

// configureNeedrestart sets needrestart to automatic mode (requires root)
func configureNeedrestart() error {
	logger.Debug("Configuring needrestart for automatic restarts...")
	content := []byte(`$nrconf{restart} = 'a';`)
	filePath := "/etc/needrestart/conf.d/50-autorestart.conf"

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		// Don't fail the whole process if this write fails, just log it.
		logger.Warn("Failed to write needrestart config %s: %v", filePath, err)
		return nil // Return nil to indicate non-fatal error
	}
	logger.Debug("Successfully configured needrestart.")
	return nil
}

// Note: The Quitting() method needs to be implemented or accessible
// on the ui.ToolSelector struct for the cancellation check to work correctly.
// This might involve adding a method like `Finished() bool` to ui.ToolSelector
// that returns true if the UI exited normally (Enter on last category) vs. quit (Ctrl+C). 