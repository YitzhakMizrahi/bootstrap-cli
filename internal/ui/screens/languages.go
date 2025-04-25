package screens

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// LanguageScreen represents the language selection screen
type LanguageScreen struct {
	title string
}

// NewLanguageScreen creates a new language selection screen
func NewLanguageScreen() *LanguageScreen {
	return &LanguageScreen{
		title: "Language Runtimes",
	}
}

// ShowLanguageSelection shows the language selection screen
func (s *LanguageScreen) ShowLanguageSelection(availableLanguages []*interfaces.Language) ([]*interfaces.Language, error) {
	// Create a selector for languages
	selector := components.NewBaseSelector(s.title)
	
	// Set up the selector with language items
	selector.SetItems(
		utils.ConvertToInterfaceSlice(availableLanguages),
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

	// Run the selector
	p := tea.NewProgram(selector)
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("language selection failed: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*components.BaseSelector); ok {
		if !selectorModel.Finished() {
			return nil, fmt.Errorf("selection cancelled")
		}

		// Convert selected items to languages
		selected := selectorModel.GetSelected()
		languages := make([]*interfaces.Language, 0, len(selected))
		for _, item := range selected {
			if lang, ok := item.(*interfaces.Language); ok {
				languages = append(languages, lang)
			}
		}
		return languages, nil
	}

	return nil, fmt.Errorf("failed to get language selector model")
}

// ShowManagerSelection prompts the user to select language managers based on selected languages
func (s *LanguageScreen) ShowManagerSelection(availableManagers []*interfaces.Tool, selectedLanguages []*interfaces.Language) ([]*interfaces.Tool, error) {
	if len(selectedLanguages) == 0 {
		return nil, fmt.Errorf("no languages selected")
	}

	// Filter managers based on selected languages
	filteredManagers := filterManagersByLanguages(availableManagers, selectedLanguages)
	if len(filteredManagers) == 0 {
		return nil, nil // No relevant managers to show
	}

	// Create a selector for managers
	selector := components.NewBaseSelector("Language Managers")
	
	// Set up the selector with filtered managers
	selector.SetItems(
		utils.ConvertToInterfaceSlice(filteredManagers),
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

	// Run the selector
	p := tea.NewProgram(selector)
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("manager selection failed: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*components.BaseSelector); ok {
		if !selectorModel.Finished() {
			return nil, fmt.Errorf("selection cancelled")
		}

		// Convert selected items to managers
		selected := selectorModel.GetSelected()
		managers := make([]*interfaces.Tool, 0, len(selected))
		for _, item := range selected {
			if manager, ok := item.(*interfaces.Tool); ok {
				managers = append(managers, manager)
			}
		}
		return managers, nil
	}

	return nil, fmt.Errorf("failed to get manager selector model")
}

// Helper function to filter managers based on selected languages
func filterManagersByLanguages(managers []*interfaces.Tool, languages []*interfaces.Language) []*interfaces.Tool {
	filtered := make([]*interfaces.Tool, 0)
	languageNames := make(map[string]bool)
	
	// Create a map of selected language names
	for _, lang := range languages {
		languageNames[lang.Name] = true
	}
	
	// Filter managers that support any of the selected languages
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