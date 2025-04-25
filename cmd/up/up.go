package up

import (
	"fmt"
	"os"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/shell"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui"
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

func runUp(cmd *cobra.Command, args []string) error {
	// Step 1: Welcome Screen
	if !ui.ShowWelcomeScreen() {
		fmt.Println("Setup cancelled by user")
		return nil
	}

	// Step 2: System Detection
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system: %w", err)
	}
	if !ui.ShowSystemInfo(sysInfo) {
		fmt.Println("Setup cancelled by user")
		return nil
	}

	// Step 3: Dotfiles GitHub Clone (Optional)
	dotfilesURL, err := ui.PromptDotfiles()
	if err != nil {
		return fmt.Errorf("failed to handle dotfiles: %w", err)
	}
	if dotfilesURL != "" {
		// TODO: Implement dotfiles cloning
		fmt.Printf("Cloning dotfiles from: %s\n", dotfilesURL)
	}

	// Step 4: Shell Selection
	shellMgr := shell.NewDefaultManager()
	shellInfo, err := shellMgr.DetectCurrent()
	if err != nil {
		return fmt.Errorf("failed to detect shell: %w", err)
	}

	// Get list of available shells
	availableShells, err := shellMgr.ListAvailable()
	if err != nil {
		return fmt.Errorf("failed to list available shells: %w", err)
	}

	// Update shellInfo with the complete list of available shells
	shellInfo.Available = make([]string, 0)
	for _, shell := range availableShells {
		if !contains(shellInfo.Available, shell.Current) {
			shellInfo.Available = append(shellInfo.Available, shell.Current)
		}
	}

	selectedShell, err := ui.PromptShellSelection(shellInfo)
	if err != nil {
		return fmt.Errorf("failed to handle shell selection: %w", err)
	}
	if selectedShell != "" {
		// Create shell config from selection
		config := &interfaces.ShellConfig{
			// Add default configuration for the selected shell
			Aliases: make(map[string]string),
			Exports: make(map[string]string),
			Functions: make(map[string]string),
			Path: []string{},
			Source: []string{},
		}
		
		// Configure the selected shell
		if err := shellMgr.ConfigureShell(config); err != nil {
			return fmt.Errorf("failed to configure shell: %w", err)
		}
		fmt.Printf("Shell configured: %s\n", selectedShell)
	}

	// Step 5: Font Installer (Optional)
	installFonts, err := ui.PromptFontInstallation()
	if err != nil {
		return fmt.Errorf("failed to handle font installation: %w", err)
	}
	if installFonts {
		// TODO: Implement font installation
		fmt.Println("Installing JetBrains Mono Nerd Font...")
	}

	// Step 6: Tool Selection
	selectedTools, err := ui.PromptToolSelection()
	if err != nil {
		return fmt.Errorf("failed to handle tool selection: %w", err)
	}

	// Step 7: Language Selection
	selectedLanguages, err := ui.PromptLanguages()
	if err != nil {
		return fmt.Errorf("failed to handle language selection: %w", err)
	}

	// Step 8: Language Manager Selection
	selectedManagers, err := ui.PromptLanguageManagersForLanguages(selectedLanguages)
	if err != nil {
		return fmt.Errorf("failed to handle language manager selection: %w", err)
	}

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
		logger.Info("Installing %s...", manager)
		if err := installer.Install(&install.Tool{Name: manager}); err != nil {
			logger.Error("Failed to install language manager %s: %v. Continuing...", manager, err)
		}
	}
	logger.Success("Language manager installation process finished.")

	// Then install languages
	logger.Info("Installing selected languages...")
	for _, lang := range selectedLanguages {
		logger.Info("Installing %s...", lang)
		if err := installer.Install(&install.Tool{Name: lang}); err != nil {
			logger.Error("Failed to install language %s: %v. Continuing...", lang, err)
		}
	}
	logger.Success("Language installation process finished.")

	// Finally install other selected tools
	logger.Info("Installing selected tools...")
	for _, tool := range selectedTools {
		logger.Info("Installing %s...", tool)
		if err := installer.Install(&install.Tool{Name: tool}); err != nil {
			logger.Error("Failed to install tool %s: %v. Continuing...", tool, err)
		}
	}
	logger.Success("Tool installation process finished.")

	// Install JetBrains Mono Nerd Font if selected
	if installFonts {
		logger.Info("Installing JetBrains Mono Nerd Font...")
		fontInstaller := install.NewFontInstaller(logger)
		if err := fontInstaller.InstallJetBrainsMono(); err != nil {
			logger.Error("Failed to install JetBrains Mono font: %v", err)
		} else {
			logger.Success("JetBrains Mono font installed successfully.")
		}
	}

	// Step 9: Validation & Finish
	if err := ui.ValidateSetup(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println("Bootstrap process completed successfully!")
	return nil
}

// contains checks if a string is present in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
} 