package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// LanguageScreen represents the language selection screen
type LanguageScreen struct {
	selector *components.BaseSelector
	selectedLanguages []*interfaces.Language
	selectedManagers  []*interfaces.Tool
	finished bool
}

// NewLanguageScreen creates a new language selection screen
func NewLanguageScreen(languages []*interfaces.Language, managers []*interfaces.Tool) *LanguageScreen {
	items := make([]interface{}, 0, len(languages)+len(managers))
	for _, m := range managers {
		items = append(items, m)
	}
	for _, l := range languages {
		items = append(items, l)
	}
	selector := components.NewBaseSelector("Select Languages and Version Managers")
	selector.SetItems(items,
		func(i interface{}) string {
			switch v := i.(type) {
			case *interfaces.Tool:
				return "Manager: " + v.Name
			case *interfaces.Language:
				return v.Name
			default:
				return "Unknown"
			}
		},
		func(i interface{}) string {
			switch v := i.(type) {
			case *interfaces.Tool:
				return v.Description
			case *interfaces.Language:
				return v.Description
			default:
				return ""
			}
		},
	)
	return &LanguageScreen{
		selector: selector,
		finished: false,
	}
}

// Init implements tea.Model
func (s *LanguageScreen) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (s *LanguageScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	newModel, cmd := s.selector.Update(msg)
	if selector, ok := newModel.(*components.BaseSelector); ok {
		s.selector = selector
		if s.selector.Finished() {
			selected := s.selector.GetSelected()
			languages := make([]*interfaces.Language, 0)
			managers := make([]*interfaces.Tool, 0)
			for _, item := range selected {
				if lang, ok := item.(*interfaces.Language); ok {
					languages = append(languages, lang)
				} else if tool, ok := item.(*interfaces.Tool); ok {
					managers = append(managers, tool)
				}
			}
			s.selectedLanguages = languages
			s.selectedManagers = managers
			s.finished = true
		}
	}
	return s, cmd
}

// View implements tea.Model
func (s *LanguageScreen) View() string {
	return styles.TitleStyle.Render("Select Languages and Version Managers") + "\n\n" + s.selector.View()
}

// Finished returns true if the screen was completed
func (s *LanguageScreen) Finished() bool {
	return s.finished
}

// GetSelectedLanguages returns the selected languages
func (s *LanguageScreen) GetSelectedLanguages() []*interfaces.Language {
	return s.selectedLanguages
}

// GetSelectedManagers returns the selected managers
func (s *LanguageScreen) GetSelectedManagers() []*interfaces.Tool {
	return s.selectedManagers
}

// ShowLanguageSelection shows the language selection screen
func (s *LanguageScreen) ShowLanguageSelection(availableLanguages []*interfaces.Language) *components.BaseSelector {
	// Create a selector for languages
	selector := components.NewBaseSelector("Language Runtimes")
	
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
	return selector
}

// ShowManagerSelection prompts the user to select language managers based on selected languages
func (s *LanguageScreen) ShowManagerSelection(availableManagers []*interfaces.Tool, selectedLanguages []*interfaces.Language) *components.BaseSelector {
	if len(selectedLanguages) == 0 {
		return nil
	}

	// Filter managers based on selected languages
	filteredManagers := filterManagersByLanguages(availableManagers, selectedLanguages)
	if len(filteredManagers) == 0 {
		return nil // No relevant managers to show
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
	return selector
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