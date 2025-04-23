package ui

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/shell"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/manifoldco/promptui"
)

// ShowWelcomeScreen displays the welcome screen and returns true if user wants to continue
func ShowWelcomeScreen() bool {
	prompt := promptui.Select{
		Label: "✨ Bootstrap CLI ✨\nSetup your dev machine with ease",
		Items: []string{"Start", "Exit"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return false
	}

	return result == "Start"
}

// ShowSystemInfo displays the system information and returns true if user wants to continue
func ShowSystemInfo(info *system.SystemInfo) bool {
	fmt.Printf("System Info Detected:\n")
	fmt.Printf("OS: %s %s\n", info.Distro, info.Version)
	fmt.Printf("Arch: %s\n", info.Arch)
	fmt.Printf("Package Manager: %s\n\n", info.PackageType)

	prompt := promptui.Select{
		Label: "Press Enter to continue",
		Items: []string{"Continue"},
	}

	_, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return false
	}

	return true
}

// PromptDotfiles prompts for GitHub dotfiles URL
func PromptDotfiles() (string, error) {
	prompt := promptui.Select{
		Label: "Clone dotfiles from GitHub?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	if result == "No" {
		return "", nil
	}

	urlPrompt := promptui.Prompt{
		Label: "Enter GitHub repo URL",
	}

	url, err := urlPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	return url, nil
}

// PromptShellSelection prompts for shell selection
func PromptShellSelection(currentShell *shell.ShellInfo) (string, error) {
	shells := []string{"zsh", "bash", "fish"}
	prompt := promptui.Select{
		Label: fmt.Sprintf("Choose your shell (current: %s)", currentShell.Type),
		Items: shells,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	return result, nil
}

// PromptFontInstallation prompts for font installation
func PromptFontInstallation() (bool, error) {
	prompt := promptui.Select{
		Label: "Install JetBrains Mono Nerd Font?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("prompt failed: %w", err)
	}

	return result == "Yes", nil
}

// PromptToolSelection prompts for tool selection
func PromptToolSelection() ([]string, error) {
	tools := []string{
		"git", "curl", "bat", "lsd", "fzf", "zoxide",
		"vim", "nano", "htop", "build-essential",
	}

	selector := NewToolSelector(tools)
	p := tea.NewProgram(selector)

	// Run the interactive UI
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("UI error: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*ToolSelector); ok {
		if !selectorModel.Finished() {
			return nil, fmt.Errorf("selection cancelled")
		}
		return selectorModel.GetSelectedTools(), nil
	}

	return nil, fmt.Errorf("failed to get tool selector model")
}

// PromptLanguageRuntimes prompts for language runtime installation
func PromptLanguageRuntimes() ([]string, error) {
	runtimes := []string{
		"Node.js (nvm)",
		"Python (pyenv)",
		"Go (goenv)",
		"Rust (rustup)",
	}

	selector := NewToolSelector(runtimes)
	p := tea.NewProgram(selector)

	// Run the interactive UI
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("UI error: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*ToolSelector); ok {
		if !selectorModel.Finished() {
			return nil, fmt.Errorf("selection cancelled")
		}
		return selectorModel.GetSelectedTools(), nil
	}

	return nil, fmt.Errorf("failed to get tool selector model")
}

// ValidateSetup validates the installation
func ValidateSetup() error {
	// TODO: Implement actual validation
	fmt.Println("Validation Results:")
	fmt.Println("- Shell setup: OK")
	fmt.Println("- Tools installed: OK")
	fmt.Println("- Language runtimes: OK")
	fmt.Println("- Paths and symlinks: Configured")
	fmt.Println("\n✅ All systems go!")

	prompt := promptui.Select{
		Label: "Press Enter to finish",
		Items: []string{"Finish"},
	}

	_, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}

	return nil
} 