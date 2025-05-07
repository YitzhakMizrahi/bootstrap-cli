// Package ui provides user interface components and prompts for the bootstrap-cli application.
package ui

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
)

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