package up

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/shell"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/tools"
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
	shellMgr := shell.NewManager()
	shellInfo, err := shellMgr.DetectCurrent()
	if err != nil {
		return fmt.Errorf("failed to detect shell: %w", err)
	}
	selectedShell, err := ui.PromptShellSelection(shellInfo)
	if err != nil {
		return fmt.Errorf("failed to handle shell selection: %w", err)
	}
	if selectedShell != "" {
		// TODO: Implement shell configuration
		fmt.Printf("Configuring shell: %s\n", selectedShell)
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

	// Step 7: Language Runtimes
	runtimes, err := ui.PromptLanguageRuntimes()
	if err != nil {
		return fmt.Errorf("failed to handle language runtimes: %w", err)
	}

	// Step 8: Install Progress
	logger := log.New(log.InfoLevel)
	pm, err := packages.DetectPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	installer := install.NewInstaller(pm)
	installer.Logger = logger // Assign the existing logger

	// Convert selected tool names back to install.Tool objects
	var toolsToInstall []*install.Tool
	allToolCategories := tools.GetToolCategories() // Get all defined tools
	for _, selectedName := range selectedTools {
		found := false
		for _, category := range allToolCategories {
			for _, tool := range category.Tools {
				if tool.Name == selectedName {
					toolRef := tool // Create a local copy for the pointer
					toolsToInstall = append(toolsToInstall, &toolRef)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			// Handle case where a selected tool name doesn't match any known tool
			// For now, just log a warning, but maybe create a basic Tool object
			logger.Warn("Selected tool '%s' not found in defined categories.", selectedName)
			// Optionally, create a basic install.Tool if you want to attempt installation anyway
			// toolsToInstall = append(toolsToInstall, &install.Tool{Name: selectedName, PackageName: selectedName})
		}
	}

	// Install selected tools
	logger.Info("Installing selected tools...")
	for _, tool := range toolsToInstall {
		logger.Info("Installing %s...", tool.Name)
		if err := installer.Install(tool); err != nil {
			// Decide if one tool failing should stop the whole process
			logger.Error("Failed to install tool %s: %v. Continuing...", tool.Name, err)
			// return fmt.Errorf("failed to install tool %s: %w", tool.Name, err) // Uncomment to stop on first error
		}
	}
	logger.Success("Selected tools installation process finished.")

	// TODO: Install selected language runtimes
	logger.Info("Installing selected language runtimes...")
	for _, runtime := range runtimes {
		// Placeholder: Need actual installation logic for runtimes
		// This might involve a similar pattern: looking up runtime details
		// and calling an appropriate install function.
		logger.Info("Attempting to install %s...", runtime)
		// runtimeInstaller.Install(runtimeInfo) // Example
	}
	logger.Success("Language runtime installation process finished.")

	// Step 9: Validation & Finish
	if err := ui.ValidateSetup(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println("Bootstrap process completed successfully!")
	return nil
} 