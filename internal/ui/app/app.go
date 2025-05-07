// Package app provides the main application model for the bootstrap-cli UI.
package app

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/detector"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// detectSystemMsg is sent when system detection is complete
type detectSystemMsg struct {
	pmType interfaces.PackageManagerType
	err    error
}

// detectSystem performs system detection asynchronously
func detectSystem() tea.Cmd {
	return func() tea.Msg {
		pmType, err := detector.DetectPackageManager()
		return detectSystemMsg{pmType: pmType, err: err}
	}
}

// Update implements tea.Model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Forward WindowSizeMsg to the current screen so lists get correct height
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.screenReady = true
		m.stepIndicator.SetCurrentStep(int(m.currentScreen))
		
		// Update only the current screen
		switch m.currentScreen {
		case WelcomeScreen:
			if m.welcomeScreen != nil {
				m.welcomeScreen.Update(msg)
			}
		case SystemScreen:
			if m.systemScreen != nil {
				m.systemScreen.Update(msg)
			}
		case ToolScreen:
			if m.toolScreen != nil {
				m.toolScreen.Selector().SetSize(m.width, m.height-8)
				m.toolScreen.Update(msg)
			}
		case FontScreen:
			if m.fontScreen != nil {
				m.fontScreen.Selector().SetSize(m.width, m.height-8)
				m.fontScreen.Update(msg)
			}
		case LanguageScreen:
			if m.languageScreen != nil {
				m.languageScreen.Selector().SetSize(m.width, m.height-8)
				m.languageScreen.Update(msg)
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "shift+tab":
			if m.currentScreen > 0 {
				m.currentScreen--
				m.stepIndicator.SetCurrentStep(int(m.currentScreen))
				// Initialize the screen immediately
				m.initScreenIfNeeded()
				// Send an initial update to the screen
				if m.currentScreen == ToolScreen && m.toolScreen != nil {
					m.toolScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				} else if m.currentScreen == FontScreen && m.fontScreen != nil {
					m.fontScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				} else if m.currentScreen == LanguageScreen && m.languageScreen != nil {
					m.languageScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				}
				return m, nil
			}
			return m, nil
		case "right", "tab":
			if int(m.currentScreen) < 6 {
				m.currentScreen++
				m.stepIndicator.SetCurrentStep(int(m.currentScreen))
				// Initialize the screen immediately
				m.initScreenIfNeeded()
				// Send an initial update to the screen
				if m.currentScreen == ToolScreen && m.toolScreen != nil {
					m.toolScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				} else if m.currentScreen == FontScreen && m.fontScreen != nil {
					m.fontScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				} else if m.currentScreen == LanguageScreen && m.languageScreen != nil {
					m.languageScreen.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd

	// Initialize screen if needed before processing updates
	m.initScreenIfNeeded()

	switch m.currentScreen {
	case WelcomeScreen:
		if m.welcomeScreen != nil {
			var newModel tea.Model
			newModel, cmd = m.welcomeScreen.Update(msg)
			m.welcomeScreen = newModel.(*screens.WelcomeScreen)
			if m.welcomeScreen.Finished() {
				return m, m.nextScreen()
			}
		}
	case SystemScreen:
		if m.systemScreen != nil {
			var newModel tea.Model
			newModel, cmd = m.systemScreen.Update(msg)
			m.systemScreen = newModel.(*screens.SystemScreen)
			if m.systemScreen.Finished() {
				return m, m.nextScreen()
			}
		}
	case ToolScreen:
		if m.toolScreen != nil {
			var newModel tea.Model
			newModel, cmd = m.toolScreen.Update(msg)
			m.toolScreen = newModel.(*screens.SelectionScreen)
			if m.toolScreen.Finished() {
				selectedMap := m.toolScreen.GetSelected()
				selectedTools := make([]*interfaces.Tool, 0, len(selectedMap))
				for _, item := range selectedMap {
					if tool, ok := item.(*interfaces.Tool); ok {
						selectedTools = append(selectedTools, tool)
					}
				}
				m.selectedTools = selectedTools
				return m, m.nextScreen()
			}
		}
	case FontScreen:
		if m.fontScreen != nil {
			var newModel tea.Model
			newModel, cmd = m.fontScreen.Update(msg)
			m.fontScreen = newModel.(*screens.SelectionScreen)
			if m.fontScreen.Finished() {
				selectedMap := m.fontScreen.GetSelected()
				selectedFonts := make([]*interfaces.Font, 0, len(selectedMap))
				for _, item := range selectedMap {
					if font, ok := item.(*interfaces.Font); ok {
						selectedFonts = append(selectedFonts, font)
					}
				}
				m.selectedFonts = selectedFonts
				return m, m.nextScreen()
			}
		}
	case LanguageScreen:
		if m.languageScreen != nil {
			var newModel tea.Model
			newModel, cmd = m.languageScreen.Update(msg)
			m.languageScreen = newModel.(*screens.SelectionScreen)
			if m.languageScreen.Finished() {
				selectedMap := m.languageScreen.GetSelected()
				selectedLanguages := make([]*interfaces.Language, 0, len(selectedMap))
				selectedManagers := make([]*interfaces.Tool, 0, len(selectedMap))
				for _, item := range selectedMap {
					switch v := item.(type) {
					case *interfaces.Tool:
						selectedManagers = append(selectedManagers, v)
					case *interfaces.Language:
						selectedLanguages = append(selectedLanguages, v)
					}
				}
				m.selectedLanguages = selectedLanguages
				m.selectedManagers = selectedManagers
				return m, m.nextScreen()
			}
		}
	case DotfilesScreen:
		// ... similar logic for dotfiles ...
		return m, nil
	case FinishScreen:
		return m, tea.Quit
	}
	return m, cmd
}

// CurrentScreen returns the current screen
// func (m *Model) CurrentScreen() Screen {
// 	return m.currentScreen
// }

// SelectedTools returns the selected tools
// func (m *Model) SelectedTools() []*interfaces.Tool {
// 	return m.selectedTools
// }

// SelectedFonts returns the selected fonts
// func (m *Model) SelectedFonts() []*interfaces.Font {
// 	return m.selectedFonts
// }

// SelectedLanguages returns the selected languages
// func (m *Model) SelectedLanguages() []*interfaces.Language {
// 	return m.selectedLanguages
// }

// SelectedManagers returns the selected language managers
// func (m *Model) SelectedManagers() []*interfaces.Tool {
// 	return m.selectedManagers
// } 