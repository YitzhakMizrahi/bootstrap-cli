// Package app provides the main application model for the bootstrap-cli UI.
package app

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/detector"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents a UI screen
type Screen int

const (
	// WelcomeScreen is the initial screen shown to users
	WelcomeScreen Screen = iota
	// SystemScreen allows users to configure system settings
	SystemScreen
	// ToolScreen allows users to select development tools
	ToolScreen
	// FontScreen allows users to select fonts
	FontScreen
	// LanguageScreen allows users to select programming languages
	LanguageScreen
	// DotfilesScreen allows users to select dotfiles
	DotfilesScreen
	// FinishScreen shows the completion status
	FinishScreen
)

// Model represents the main application model
type Model struct {
	currentScreen Screen
	config       *config.Loader
	width        int
	height       int
	
	// Selected items
	selectedTools     []*interfaces.Tool
	selectedFonts     []*interfaces.Font
	selectedLanguages []*interfaces.Language
	selectedManagers  []*interfaces.Tool

	// System information
	pmType      interfaces.PackageManagerType
	systemReady bool

	// Components
	stepIndicator *components.StepIndicator

	// Child screens
	welcomeScreen *screens.WelcomeScreen
	systemScreen  *screens.SystemScreen
	toolScreen    *screens.ToolScreen
	fontScreen    *screens.FontScreen
	languageScreen *screens.LanguageScreen

	// State
	err    error
	loaded bool
}

// New creates a new application model
func New(config *config.Loader) *Model {
	steps := []components.Step{
		{Name: "Welcome", Status: "current"},
		{Name: "System", Status: "pending"},
		{Name: "Tools", Status: "pending"},
		{Name: "Fonts", Status: "pending"},
		{Name: "Languages", Status: "pending"},
		{Name: "Dotfiles", Status: "pending"},
		{Name: "Finish", Status: "pending"},
	}

	return &Model{
		currentScreen:  WelcomeScreen,
		config:        config,
		stepIndicator: components.NewStepIndicator(steps),
	}
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	// Enter alternate screen and hide cursor
	return tea.Batch(
		tea.EnterAltScreen,
		tea.HideCursor,
	)
}

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
				m.toolScreen.Update(msg)
			}
		case FontScreen:
			if m.fontScreen != nil {
				m.fontScreen.Update(msg)
			}
		case LanguageScreen:
			if m.languageScreen != nil {
				m.languageScreen.Update(msg)
			}
		// Add other screens as needed
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
			m.toolScreen = screens.NewToolScreen(tools)
		}
		var newModel tea.Model
		newModel, cmd = m.toolScreen.Update(msg)
		m.toolScreen = newModel.(*screens.ToolScreen)
		if m.toolScreen.Finished() {
			m.selectedTools = m.toolScreen.GetSelected()
			return m, m.nextScreen()
		}
		return m, cmd
	case FontScreen:
		if m.fontScreen == nil {
			m.fontScreen = screens.NewFontScreen(m.config)
		}
		var newModel tea.Model
		newModel, cmd = m.fontScreen.Update(msg)
		m.fontScreen = newModel.(*screens.FontScreen)
		if m.fontScreen.Finished() {
			m.selectedFonts = m.fontScreen.GetSelected()
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
			m.languageScreen = screens.NewLanguageScreen(languages, managers)
		}
		var newModel tea.Model
		newModel, cmd = m.languageScreen.Update(msg)
		m.languageScreen = newModel.(*screens.LanguageScreen)
		if m.languageScreen.Finished() {
			m.selectedLanguages = m.languageScreen.GetSelectedLanguages()
			m.selectedManagers = m.languageScreen.GetSelectedManagers()
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

// View implements tea.Model
func (m *Model) View() string {
	output := "\x1b[2J\x1b[H"
	output += m.stepIndicator.View() + "\n\n"
	switch m.currentScreen {
	case WelcomeScreen:
		if m.welcomeScreen != nil {
			output += m.welcomeScreen.View()
		}
	case SystemScreen:
		if m.systemScreen != nil {
			output += m.systemScreen.View()
		}
	case ToolScreen:
		if m.toolScreen != nil {
			output += m.toolScreen.View()
		}
	case FontScreen:
		if m.fontScreen != nil {
			output += m.fontScreen.View()
		}
	case LanguageScreen:
		if m.languageScreen != nil {
			output += m.languageScreen.View()
		}
	// ... handle other screens ...
	}
	if m.err != nil {
		output += "\n\n" + styles.ErrorStyle.Render(m.err.Error())
	} else {
		output += "\n\n" + styles.HelpStyle.Render("↑/↓: navigate • space: select/deselect • enter: confirm • q: quit")
	}
	return output
}

// nextScreen advances to the next screen and configures the appropriate selector
func (m *Model) nextScreen() tea.Cmd {
	steps := m.stepIndicator.GetSteps()
	steps[int(m.currentScreen)].Status = "completed"
	m.currentScreen++
	if int(m.currentScreen) < len(steps) {
		steps[int(m.currentScreen)].Status = "current"
	}
	m.stepIndicator.SetSteps(steps)

	switch m.currentScreen {
	case SystemScreen:
		sysInfo, err := system.Detect()
		if err != nil {
			m.err = err
			return nil
		}
		m.systemScreen = screens.ShowSystemInfo(sysInfo)
		return nil
	case ToolScreen:
		tools, err := m.config.LoadTools()
		if err != nil {
			m.err = err
			return nil
		}
		m.toolScreen = screens.NewToolScreen(tools)
		return nil
	case FontScreen:
		m.fontScreen = screens.NewFontScreen(m.config)
		return nil
	case LanguageScreen:
		managers, err := m.config.LoadLanguageManagers()
		if err != nil {
			m.err = err
			return nil
		}
		languages, err := m.config.LoadLanguages()
		if err != nil {
			m.err = err
			return nil
		}
		m.languageScreen = screens.NewLanguageScreen(languages, managers)
		return nil
	case DotfilesScreen:
		// ... similar logic for dotfiles ...
		return nil
	case FinishScreen:
		return tea.Quit
	}
	return nil
}

// CurrentScreen returns the current screen
func (m *Model) CurrentScreen() Screen {
	return m.currentScreen
}

// SelectedTools returns the selected tools
func (m *Model) SelectedTools() []*interfaces.Tool {
	return m.selectedTools
}

// SelectedFonts returns the selected fonts
func (m *Model) SelectedFonts() []*interfaces.Font {
	return m.selectedFonts
}

// SelectedLanguages returns the selected languages
func (m *Model) SelectedLanguages() []*interfaces.Language {
	return m.selectedLanguages
}

// SelectedManagers returns the selected language managers
func (m *Model) SelectedManagers() []*interfaces.Tool {
	return m.selectedManagers
} 