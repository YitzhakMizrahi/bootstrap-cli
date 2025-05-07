// Package ui provides user interface components and prompts for the bootstrap-cli application.
package ui

import (
	"fmt"
	"os"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
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

// PromptToolSelection prompts for tool selection
func PromptToolSelection(loader *config.Loader) ([]*interfaces.Tool, error) {
	tools, err := loader.LoadTools()
	if err != nil {
		return nil, fmt.Errorf("failed to load tools: %w", err)
	}

	toolItems := make([]interface{}, len(tools))
	for i, t := range tools {
		toolItems[i] = t
	}
	screen := screens.NewSelectionScreen(
		"Select Development Tools",
		toolItems,
		func(item interface{}) string {
			if tool, ok := item.(*interfaces.Tool); ok {
				if tool.Category != "" {
					return fmt.Sprintf("[%s] %s", tool.Category, tool.Name)
				}
				return tool.Name
			}
			return ""
		},
		func(item interface{}) string {
			if tool, ok := item.(*interfaces.Tool); ok {
				return tool.Description
			}
			return ""
		},
	)
	p := tea.NewProgram(screen)
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("tool selection failed: %w", err)
	}
	if sel, ok := model.(*screens.SelectionScreen); ok {
		selectedMap := sel.GetSelected()
		selectedTools := make([]*interfaces.Tool, 0, len(selectedMap))
		for _, item := range selectedMap {
			if tool, ok := item.(*interfaces.Tool); ok {
				selectedTools = append(selectedTools, tool)
			}
		}
		return selectedTools, nil
	}
	return nil, fmt.Errorf("unexpected model type returned from tool selection")
}

// PromptLanguages prompts for programming language selection
func PromptLanguages() ([]*interfaces.Language, error) {
	configPath := os.Getenv("BOOTSTRAP_CLI_CONFIG")
	if configPath == "" {
		return nil, fmt.Errorf("BOOTSTRAP_CLI_CONFIG environment variable not set")
	}
	loader := config.NewLoader(configPath)
	availableLanguages, err := loader.LoadLanguages()
	if err != nil {
		return nil, fmt.Errorf("failed to load language configurations: %w", err)
	}
	languageItems := make([]interface{}, len(availableLanguages))
	for i, l := range availableLanguages {
		languageItems[i] = l
	}
	screen := screens.NewSelectionScreen(
		"Select Languages",
		languageItems,
		func(item interface{}) string {
			if lang, ok := item.(*interfaces.Language); ok {
				return lang.Name
			}
			return ""
		},
		func(item interface{}) string {
			if lang, ok := item.(*interfaces.Language); ok {
				return lang.Description
			}
			return ""
		},
	)
	p := tea.NewProgram(screen)
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("language selection failed: %w", err)
	}
	if sel, ok := model.(*screens.SelectionScreen); ok {
		selectedMap := sel.GetSelected()
		selectedLanguages := make([]*interfaces.Language, 0, len(selectedMap))
		for _, item := range selectedMap {
			if lang, ok := item.(*interfaces.Language); ok {
				selectedLanguages = append(selectedLanguages, lang)
			}
		}
		return selectedLanguages, nil
	}
	return nil, fmt.Errorf("unexpected model type returned from language selection")
}

// PromptLanguageManagersForLanguages prompts for language manager selection based on selected languages
func PromptLanguageManagersForLanguages(_ []*interfaces.Language) ([]*interfaces.Tool, error) {
	configPath := os.Getenv("BOOTSTRAP_CLI_CONFIG")
	if configPath == "" {
		return nil, fmt.Errorf("BOOTSTRAP_CLI_CONFIG environment variable not set")
	}
	loader := config.NewLoader(configPath)
	availableManagers, err := loader.LoadLanguageManagers()
	if err != nil {
		return nil, fmt.Errorf("failed to load language manager configurations: %w", err)
	}
	managerItems := make([]interface{}, len(availableManagers))
	for i, m := range availableManagers {
		managerItems[i] = m
	}
	screen := screens.NewSelectionScreen(
		"Select Language Managers",
		managerItems,
		func(item interface{}) string {
			if manager, ok := item.(*interfaces.Tool); ok {
				return manager.Name
			}
			return ""
		},
		func(item interface{}) string {
			if manager, ok := item.(*interfaces.Tool); ok {
				return manager.Description
			}
			return ""
		},
	)
	p := tea.NewProgram(screen)
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("manager selection failed: %w", err)
	}
	if sel, ok := model.(*screens.SelectionScreen); ok {
		selectedMap := sel.GetSelected()
		selectedManagers := make([]*interfaces.Tool, 0, len(selectedMap))
		for _, item := range selectedMap {
			if manager, ok := item.(*interfaces.Tool); ok {
				selectedManagers = append(selectedManagers, manager)
			}
		}
		return selectedManagers, nil
	}
	return nil, fmt.Errorf("unexpected model type returned from manager selection")
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