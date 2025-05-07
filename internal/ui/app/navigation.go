package app

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// nextScreen advances to the next screen and configures the appropriate selector
func (m *Model) nextScreen() tea.Cmd {
	steps := m.stepIndicator.GetSteps()
	if int(m.currentScreen) < len(steps) {
		steps[int(m.currentScreen)].Status = "completed"
	}
	m.currentScreen++
	if int(m.currentScreen) < len(steps) {
		steps[int(m.currentScreen)].Status = "current"
	}
	m.stepIndicator.SetSteps(steps)

	// Initialize the new screen immediately
	m.initScreenIfNeeded()

	// Send an initial update to the screen
	switch m.currentScreen {
	case SystemScreen:
		if m.systemScreen != nil {
			m.systemScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
	case ToolScreen:
		if m.toolScreen != nil {
			m.toolScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
	case FontScreen:
		if m.fontScreen != nil {
			m.fontScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
	case LanguageScreen:
		if m.languageScreen != nil {
			m.languageScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
	}

	return nil
}

// initScreenIfNeeded initializes the current screen if needed
func (m *Model) initScreenIfNeeded() {
	switch m.currentScreen {
	case WelcomeScreen:
		if m.welcomeScreen == nil {
			m.welcomeScreen = screens.ShowWelcomeScreen()
		}
	case SystemScreen:
		if m.systemScreen == nil {
			sysInfo, err := system.Detect()
			if err != nil {
				m.err = err
				return
			}
			m.systemScreen = screens.ShowSystemInfo(sysInfo)
		}
	case ToolScreen:
		if m.toolScreen == nil {
			tools, err := m.config.LoadTools()
			if err != nil {
				m.err = err
				return
			}
			toolItems := make([]interface{}, len(tools))
			for i, t := range tools {
				toolItems[i] = t
			}
			selected := make([]interface{}, len(m.selectedTools))
			for i, t := range m.selectedTools {
				selected[i] = t
			}
			m.toolScreen = screens.NewSelectionScreen(
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
				screens.WithSelectedItems(selected),
			)
			m.toolScreen.Init()
		}
	case FontScreen:
		if m.fontScreen == nil {
			availableFonts, _ := m.config.LoadFonts()
			fontItems := make([]interface{}, len(availableFonts))
			for i, f := range availableFonts {
				fontItems[i] = f
			}
			selected := make([]interface{}, len(m.selectedFonts))
			for i, f := range m.selectedFonts {
				selected[i] = f
			}
			m.fontScreen = screens.NewSelectionScreen(
				"Select Fonts",
				fontItems,
				func(item interface{}) string {
					if font, ok := item.(*interfaces.Font); ok {
						return font.Name
					}
					return ""
				},
				func(item interface{}) string {
					if font, ok := item.(*interfaces.Font); ok {
						return font.Description
					}
					return ""
				},
				screens.WithSelectedItems(selected),
			)
			m.fontScreen.Init()
		}
	case LanguageScreen:
		if m.languageScreen == nil {
			managers, err := m.config.LoadLanguageManagers()
			if err != nil {
				m.err = err
				return
			}
			languages, err := m.config.LoadLanguages()
			if err != nil {
				m.err = err
				return
			}
			languageItems := make([]interface{}, 0, len(languages)+len(managers))
			for _, mng := range managers {
				languageItems = append(languageItems, mng)
			}
			for _, lang := range languages {
				languageItems = append(languageItems, lang)
			}
			selected := make([]interface{}, 0, len(m.selectedLanguages)+len(m.selectedManagers))
			for _, l := range m.selectedLanguages {
				selected = append(selected, l)
			}
			for _, mng := range m.selectedManagers {
				selected = append(selected, mng)
			}
			m.languageScreen = screens.NewSelectionScreen(
				"Select Languages and Version Managers",
				languageItems,
				func(item interface{}) string {
					switch v := item.(type) {
					case *interfaces.Tool:
						return "Manager: " + v.Name
					case *interfaces.Language:
						return v.Name
					default:
						return "Unknown"
					}
				},
				func(item interface{}) string {
					switch v := item.(type) {
					case *interfaces.Tool:
						return v.Description
					case *interfaces.Language:
						return v.Description
					default:
						return ""
					}
				},
				screens.WithSelectedItems(selected),
			)
			m.languageScreen.Init()
		}
	}
} 