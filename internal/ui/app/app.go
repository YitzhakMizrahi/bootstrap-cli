// Package app provides the main application model for the bootstrap-cli UI.
package app

import (
	"fmt"
	"runtime"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/detector"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/utils"
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
	selector     *components.BaseSelector
	
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
	var cmd tea.Cmd

	// Update the current selector if it exists
	if m.selector != nil {
		var newModel tea.Model
		newModel, cmd = m.selector.Update(msg)
		if newSelector, ok := newModel.(*components.BaseSelector); ok {
			m.selector = newSelector
			// If selector is finished, process selections and move to next screen
			if m.selector.Finished() {
				switch m.currentScreen {
				case ToolScreen:
					selected := m.selector.GetSelected()
					tools := make([]*interfaces.Tool, 0, len(selected))
					for _, item := range selected {
						if tool, ok := item.(*interfaces.Tool); ok {
							tools = append(tools, tool)
						}
					}
					m.selectedTools = tools
				case FontScreen:
					selected := m.selector.GetSelected()
					fonts := make([]*interfaces.Font, 0, len(selected))
					for _, item := range selected {
						if font, ok := item.(*interfaces.Font); ok {
							fonts = append(fonts, font)
						}
					}
					m.selectedFonts = fonts
				case LanguageScreen:
					selected := m.selector.GetSelected()
					languages := make([]*interfaces.Language, 0)
					managers := make([]*interfaces.Tool, 0)
					for _, item := range selected {
						if lang, ok := item.(*interfaces.Language); ok {
							languages = append(languages, lang)
						} else if tool, ok := item.(*interfaces.Tool); ok {
							managers = append(managers, tool)
						}
					}
					m.selectedLanguages = languages
					m.selectedManagers = managers
				}
				return m, m.nextScreen()
			}
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case detectSystemMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.pmType = msg.pmType
		m.systemReady = true
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.selector != nil {
			m.selector.SetSize(msg.Width, msg.Height-6)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.currentScreen == WelcomeScreen {
				return m, tea.Batch(m.nextScreen(), detectSystem())
			}
			if m.currentScreen == SystemScreen && m.systemReady {
				return m, m.nextScreen()
			}
		}
	}

	return m, cmd
}

// View implements tea.Model
func (m *Model) View() string {
	// Clear screen and move cursor to top-left
	output := "\x1b[2J\x1b[H"
	
	// Header with step indicator
	output += m.stepIndicator.View() + "\n\n"

	// Main content based on current screen
	switch m.currentScreen {
	case WelcomeScreen:
		output += styles.TitleStyle.Render("Welcome to Bootstrap CLI") + "\n\n" +
			styles.BaseStyle.Render("Setup your development environment with ease") + "\n\n" +
			styles.HelpStyle.Render("Press Enter to begin setup...")
	case SystemScreen:
		output += styles.TitleStyle.Render("System Information") + "\n\n"
		if !m.systemReady {
			output += styles.BaseStyle.Render("Detecting system configuration...") + "\n"
		} else {
			output += styles.BaseStyle.Render("System detected:") + "\n\n" +
				styles.BaseStyle.Render("OS: " + runtime.GOOS) + "\n" +
				styles.BaseStyle.Render("Architecture: " + runtime.GOARCH) + "\n" +
				styles.BaseStyle.Render("Package Manager: " + string(m.pmType)) + "\n\n" +
				styles.HelpStyle.Render("Press Enter to continue...")
		}
	default:
		if m.selector != nil {
			output += m.selector.View()
		}
	}

	// Footer with error message or help text
	if m.err != nil {
		output += "\n\n" + styles.ErrorStyle.Render(m.err.Error())
	} else {
		output += "\n\n" + styles.HelpStyle.Render("↑/↓: navigate • space: select/deselect • enter: confirm • q: quit")
	}

	return output
}

// nextScreen advances to the next screen and configures the appropriate selector
func (m *Model) nextScreen() tea.Cmd {
	// Update step indicator
	steps := m.stepIndicator.GetSteps()
	steps[int(m.currentScreen)].Status = "completed"
	m.currentScreen++
	if int(m.currentScreen) < len(steps) {
		steps[int(m.currentScreen)].Status = "current"
	}
	m.stepIndicator.SetSteps(steps)
	
	switch m.currentScreen {
	case SystemScreen:
		// System detection screen - no selector needed
		return nil
	case ToolScreen:
		tools, err := m.config.LoadTools()
		if err != nil {
			m.err = err
			return nil
		}
		fmt.Printf("Loaded %d tools\n", len(tools))
		
		m.selector = components.NewBaseSelector("Select Development Tools")
		converted := utils.ConvertToInterfaceSlice(tools)
		fmt.Printf("Converted to %d interface items\n", len(converted))
		
		m.selector.SetItems(converted, 
			func(i interface{}) string {
				tool := i.(*interfaces.Tool)
				if tool.Category != "" {
					return fmt.Sprintf("[%s] %s", tool.Category, tool.Name)
				}
				return tool.Name
			},
			func(i interface{}) string { return i.(*interfaces.Tool).Description })
		m.selector.SetSize(m.width, m.height-6)
		return nil
	case FontScreen:
		fonts, err := m.config.LoadFonts()
		if err != nil {
			m.err = err
			return nil
		}
		m.selector = components.NewBaseSelector("Select Fonts")
		m.selector.SetItems(utils.ConvertToInterfaceSlice(fonts),
			func(i interface{}) string { return i.(*interfaces.Font).Name },
			func(i interface{}) string { return i.(*interfaces.Font).Description })
		m.selector.SetSize(m.width, m.height-6)
		return nil
	case LanguageScreen:
		// Load both language managers and languages
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

		// Combine managers and languages into a single list
		items := make([]interface{}, 0, len(managers)+len(languages))
		for _, m := range managers {
			items = append(items, m)
		}
		for _, l := range languages {
			items = append(items, l)
		}

		m.selector = components.NewBaseSelector("Select Languages and Version Managers")
		m.selector.SetItems(items,
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
			})
		m.selector.SetSize(m.width, m.height-6)
		return nil
	case DotfilesScreen:
		dotfiles, err := m.config.LoadDotfiles()
		if err != nil {
			m.err = err
			return nil
		}
		m.selector = components.NewBaseSelector("Select Dotfiles")
		m.selector.SetItems(utils.ConvertToInterfaceSlice(dotfiles),
			func(i interface{}) string { return i.(*interfaces.Dotfile).Name },
			func(i interface{}) string { return i.(*interfaces.Dotfile).Description })
		m.selector.SetSize(m.width, m.height-6)
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