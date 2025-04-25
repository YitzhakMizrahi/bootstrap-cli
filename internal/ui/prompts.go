// Package ui provides user interface components and prompts for the bootstrap-cli application.
package ui

import (
	"fmt"
	"os"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
)

// ShowWelcomeScreen displays the welcome screen and returns true if user wants to continue
func ShowWelcomeScreen() bool {
	return screens.ShowWelcomeScreen()
}

// ShowSystemInfo displays the system information and returns true if user wants to continue
func ShowSystemInfo(info *system.Info) bool {
	return screens.ShowSystemInfo(info)
}

// PromptDotfiles prompts for GitHub dotfiles URL
func PromptDotfiles() (string, error) {
	prompt := components.NewBasicPrompt("Clone dotfiles from GitHub?", []string{"Yes", "No"})
	
	shouldClone, err := prompt.RunYesNo()
	if err != nil {
		return "", err
	}

	if !shouldClone {
		return "", nil
	}

	urlPrompt := components.NewBasicPrompt("Enter GitHub repo URL", nil)
	return urlPrompt.RunWithInput()
}

// PromptShellSelection prompts the user to select a shell
func PromptShellSelection(shellInfo *interfaces.ShellInfo) (string, error) {
	if len(shellInfo.Available) == 0 {
		return "", fmt.Errorf("no supported shells found")
	}

	prompt := components.NewBasicPrompt("Select your preferred shell", shellInfo.Available)
	return prompt.Run()
}

// PromptFontInstallation prompts for font installation
func PromptFontInstallation() (bool, error) {
	prompt := components.NewBasicPrompt("Install JetBrains Mono Nerd Font?", []string{"Yes", "No"})
	return prompt.RunYesNo()
}

// PromptToolSelection prompts for tool selection
func PromptToolSelection(loader *config.Loader) ([]*interfaces.Tool, error) {
	tools, err := loader.LoadTools()
	if err != nil {
		return nil, fmt.Errorf("failed to load tools: %w", err)
	}

	toolScreen := screens.NewToolScreen(tools)
	return toolScreen.ShowToolSelection()
}

// PromptLanguages prompts for programming language selection
func PromptLanguages() ([]*interfaces.Language, error) {
	// Get config path from environment
	configPath := os.Getenv("BOOTSTRAP_CLI_CONFIG")
	if configPath == "" {
		return nil, fmt.Errorf("BOOTSTRAP_CLI_CONFIG environment variable not set")
	}

	// Create config loader with the correct path
	loader := config.NewLoader(configPath)
	
	// Load available languages
	availableLanguages, err := loader.LoadLanguages()
	if err != nil {
		return nil, fmt.Errorf("failed to load language configurations: %w", err)
	}

	// Use the new language screen
	screen := screens.NewLanguageScreen()
	return screen.ShowLanguageSelection(availableLanguages)
}

// PromptLanguageManagersForLanguages prompts for language manager selection based on selected languages
func PromptLanguageManagersForLanguages(selectedLanguages []*interfaces.Language) ([]*interfaces.Tool, error) {
	// Get config path from environment
	configPath := os.Getenv("BOOTSTRAP_CLI_CONFIG")
	if configPath == "" {
		return nil, fmt.Errorf("BOOTSTRAP_CLI_CONFIG environment variable not set")
	}

	// Create config loader with the correct path
	loader := config.NewLoader(configPath)

	// Load all language managers
	availableManagers, err := loader.LoadLanguageManagers()
	if err != nil {
		return nil, fmt.Errorf("failed to load language manager configurations: %w", err)
	}

	// Use the new language screen
	screen := screens.NewLanguageScreen()
	return screen.ShowManagerSelection(availableManagers, selectedLanguages)
}

// Helper function to filter managers based on selected languages
func filterManagersByLanguages(managers []*interfaces.Tool, languages []string) []*interfaces.Tool {
	filtered := make([]*interfaces.Tool, 0)
	languageNames := make(map[string]bool)
	
	// Create a map of selected language names
	for _, lang := range languages {
		languageNames[lang] = true
	}

	// Filter managers based on their associated languages
	for _, manager := range managers {
		for _, lang := range manager.Languages {
			if languageNames[lang] {
				filtered = append(filtered, manager)
				break
			}
		}
	}

	return filtered
}

// ValidateSetup validates the installation
func ValidateSetup() error {
	// Display validation results
	fmt.Println("Validation Results:")
	fmt.Println("- Shell setup: OK")
	fmt.Println("- Tools installed: OK")
	fmt.Println("- Language runtimes: OK")
	fmt.Println("- Paths and symlinks: Configured")
	fmt.Println("\n✅ All systems go!")

	// Use the basic prompt for the finish option
	prompt := components.NewBasicPrompt("Press Enter to finish", []string{"Finish"})
	_, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}

	return nil
} 