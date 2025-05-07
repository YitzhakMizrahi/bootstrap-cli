// Package app provides the main application model for the bootstrap-cli UI.
package app

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/detector"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
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
		if m.currentScreen == WelcomeScreen && m.welcomeScreen != nil {
			m.screenReady = true
		}
		if m.currentScreen == SystemScreen && m.systemScreen != nil {
			m.screenReady = true
		}
		m.stepIndicator.SetCurrentStep(int(m.currentScreen))
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
		// Add other screens as needed
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "shift+tab":
			if m.currentScreen > 0 {
				m.currentScreen--
				m.stepIndicator.SetCurrentStep(int(m.currentScreen))
				return m, nil
			}
			return m, nil
		case "right", "tab":
			if int(m.currentScreen) < 6 {
				m.currentScreen++
				m.stepIndicator.SetCurrentStep(int(m.currentScreen))
				m.initScreenIfNeeded()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd

	switch m.currentScreen {
	case WelcomeScreen:
		if m.welcomeScreen == nil {
			m.welcomeScreen = screens.ShowWelcomeScreen()
		}
		var newModel tea.Model
		newModel, cmd = m.welcomeScreen.Update(msg)
		m.welcomeScreen = newModel.(*screens.WelcomeScreen)
		if m.welcomeScreen.Finished() {
			return m, m.nextScreen()
		}
		return m, cmd
	case SystemScreen:
		if m.systemScreen == nil {
			sysInfo, err := system.Detect()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.systemScreen = screens.ShowSystemInfo(sysInfo)
		}
		var newModel tea.Model
		newModel, cmd = m.systemScreen.Update(msg)
		m.systemScreen = newModel.(*screens.SystemScreen)
		if m.systemScreen.Finished() {
			return m, m.nextScreen()
		}
		return m, cmd
	case ToolScreen:
		if m.toolScreen == nil {
			tools, err := m.config.LoadTools()
			if err != nil {
				m.err = err
				return m, nil
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
			return m, m.toolScreen.Init()
		}
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
		return m, cmd
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
			return m, m.fontScreen.Init()
		}
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
		return m, cmd
	case LanguageScreen:
		if m.languageScreen == nil {
			managers, err := m.config.LoadLanguageManagers()
			if err != nil {
				m.err = err
				return m, nil
			}
			languages, err := m.config.LoadLanguages()
			if err != nil {
				m.err = err
				return m, nil
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
			return m, m.languageScreen.Init()
		}
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
		return m, cmd
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