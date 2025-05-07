package app

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents a UI screen
// (moved from app.go)
type Screen int

const (
	// WelcomeScreen is the initial welcome screen.
	WelcomeScreen Screen = iota
	// SystemScreen displays system information.
	SystemScreen
	// ToolScreen allows selection of development tools.
	ToolScreen
	// FontScreen allows selection of fonts.
	FontScreen
	// LanguageScreen allows selection of languages and managers.
	LanguageScreen
	// DotfilesScreen allows selection of dotfiles.
	DotfilesScreen
	// FinishScreen is the final confirmation screen.
	FinishScreen
)

// Model represents the main application model and aggregates all UI state.
// (moved from app.go)
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
	welcomeScreen   *screens.WelcomeScreen
	systemScreen    *screens.SystemScreen
	toolScreen      *screens.SelectionScreen
	fontScreen      *screens.SelectionScreen
	languageScreen  *screens.SelectionScreen

	// State
	err         error
	loaded      bool
	screenReady bool
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