// Package app provides the main application model for the bootstrap-cli UI.
package app

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/screens"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen represents a UI screen enum
type Screen int

const (
	WelcomeScreen Screen = iota // Includes System Info now
	EssentialToolScreen 
	ModernToolScreen    
	FontScreen
	LanguageScreen // Renamed back
	DotfilesScreen
	FinishScreen
)

// Model represents the main application model and aggregates all UI state.
type Model struct {
	instanceID    int64
	currentScreen Screen // Logical state
	activeModel   tea.Model // The currently displayed screen model
	config        *config.Loader
	width         int
	height        int
	stepIndicator components.Model
	err           error
	screenReady   bool // Flag to prevent rendering before first WindowSizeMsg

	// Stored selections - populated when selection screens finish
	selectedTools     []*interfaces.Tool
	selectedFonts     []*interfaces.Font
	selectedLanguages []*interfaces.Language
	systemInfo        *system.Info // Store detected system info
}

// New creates a new application model
func New(config *config.Loader) *Model {
	rand.Seed(time.Now().UnixNano())
	
	// Adjusted step names for indicator
	stepNames := []string{
		"Essential Tools", 
		"Modern Tools",    
		"Fonts",
		"Languages", // Simpler name again
		"Dotfiles",
		"Finish",
	}
	stepIndicatorModel := components.NewModel(stepNames)

	welcomeModel := screens.NewWelcomeScreen()

	m := &Model{
		instanceID:    rand.Int63(),
		currentScreen:  WelcomeScreen, 
		activeModel:   welcomeModel,
		config:        config,
		stepIndicator: stepIndicatorModel,
	}
	return m
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	// Initialize the first screen model and batch commands
	// Also trigger system detection
	initCmd := m.activeModel.Init()
	return tea.Batch(
		initCmd,
		detectSystem(), // Run system detection in the background
		tea.EnterAltScreen,
		tea.HideCursor,
	)
}

// --- Messages --- 

// detectSystemMsg is sent when system detection is complete
type detectSystemMsg struct {
	info *system.Info
	err  error
}

// detectSystem performs system detection asynchronously
func detectSystem() tea.Cmd {
	return func() tea.Msg {
		sysInfo, err := system.Detect()
		return detectSystemMsg{info: sysInfo, err: err}
	}
}

// --- Update Logic --- 

// Update implements tea.Model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit // Global quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.screenReady = true

		appHorizontalMargin := styles.AppStyle.GetHorizontalFrameSize()
		appVerticalMargin := styles.AppStyle.GetVerticalFrameSize()
		availableWidth := m.width - appHorizontalMargin
		availableHeight := m.height - appVerticalMargin
		if availableWidth < 0 { availableWidth = 0 }
		if availableHeight < 0 { availableHeight = 0 }

		// --- Recalculate static element heights ---
		indicatorView := ""
		indicatorHeight := 0
		if m.currentScreen >= EssentialToolScreen && m.currentScreen < FinishScreen {
			m.stepIndicator.SetWidth(availableWidth)
			indicatorView = m.stepIndicator.View()
			indicatorHeight = lipgloss.Height(indicatorView) + 1
		}

		// --- Calculate footer height (more accurately) ---
		footerStr := styles.HelpStyle.Render("q/Ctrl+c: quit") // Use the actual footer string
		if m.err != nil {
			footerStr = styles.ErrorStyle.Render(m.err.Error())
		}
		footerHeight := lipgloss.Height(footerStr) + 2 // +2 for newlines added in app.View

		// --- Calculate final dimensions for the child screen ---
		childWidth := availableWidth
		childHeight := availableHeight - indicatorHeight - footerHeight
		if childHeight < 1 { childHeight = 1 } // Keep absolute minimum

		// Forward calculated size to active model
		if m.activeModel != nil {
			var updatedModel tea.Model
			childMsg := tea.WindowSizeMsg{Width: childWidth, Height: childHeight}
			updatedModel, cmd = m.activeModel.Update(childMsg)
			m.activeModel = updatedModel
			cmds = append(cmds, cmd)
		}

	case detectSystemMsg:
		if msg.err != nil {
			m.err = fmt.Errorf("system detection failed: %w", msg.err)
		} else {
			m.systemInfo = msg.info
			m.err = nil
			if welcomeScreen, ok := m.activeModel.(*screens.WelcomeScreen); ok && welcomeScreen != nil {
				welcomeScreen.SetInfo(m.systemInfo)
			}
		}
		return m, nil // Consume this message
	}

	// --- Screen-specific Update Delegation & Transition Logic ---
	if m.activeModel != nil {
		var updatedModel tea.Model
		// Pass the *original* msg to the active model for non-WindowSize events
		if _, ok := msg.(tea.WindowSizeMsg); !ok { 
			updatedModel, cmd = m.activeModel.Update(msg)
			m.activeModel = updatedModel
			cmds = append(cmds, cmd)
		}

		// --- Screen Transition Logic ---
		switch screen := m.activeModel.(type) {
		case *screens.WelcomeScreen:
			if screen.Finished() { cmds = append(cmds, m.transitionTo(EssentialToolScreen)) }
		case *screens.EssentialToolScreen: 
			if screen.Finished() { 
				newTools := screen.GetSelected()
				existingModern := filterToolsByCategory(m.selectedTools, "modern")
				m.selectedTools = append(existingModern, newTools...)
				cmds = append(cmds, m.transitionTo(ModernToolScreen))
			}
		case *screens.ModernToolScreen:
			if screen.Finished() {
				newTools := screen.GetSelected()
				existingEssential := filterToolsByCategory(m.selectedTools, "essential")
				m.selectedTools = append(existingEssential, newTools...)
				cmds = append(cmds, m.transitionTo(FontScreen))
			}
		case *screens.FontScreen: 
			if screen.Finished() { m.selectedFonts = screen.GetSelected(); cmds = append(cmds, m.transitionTo(LanguageScreen)) }
		case *screens.LanguageScreen:
			if screen.Finished() { 
				m.selectedLanguages = screen.GetSelected() 
				cmds = append(cmds, m.transitionTo(DotfilesScreen)) 
			}
		case *screens.DotfilesScreen:
			if screen.Finished() { cmds = append(cmds, m.transitionTo(FinishScreen)) }
		case *screens.FinishScreen:
			; // Finish screen quits itself
		}
	}

	return m, tea.Batch(cmds...)
}

// transitionTo is a helper to change the active screen model
func (m *Model) transitionTo(targetScreen Screen) tea.Cmd {
	// --- State Update ---
	m.currentScreen = targetScreen 
	visualStepIndex := -1
	// Adjust index mapping for removed steps
	// EssentialTools=0, ModernTools=1, Fonts=2, Languages=3, Dotfiles=4, Finish=5 (but finish isn't shown)
	if targetScreen >= EssentialToolScreen && targetScreen < FinishScreen { 
		visualStepIndex = int(targetScreen) - 1 
	}
	m.stepIndicator.SetCurrentStep(visualStepIndex)

	// --- Create New Screen --- 
	var newScreen tea.Model
	var initCmd tea.Cmd
	switch targetScreen {
	case WelcomeScreen: newScreen = screens.NewWelcomeScreen()
	case EssentialToolScreen: 
		tools, err := m.config.LoadTools()
		if err != nil { m.err = err; newScreen = screens.NewWelcomeScreen(); break }
		essentialTools := filterToolsByCategory(tools, "essential")
		preselectedEssential := filterToolsByCategory(m.selectedTools, "essential")
		newScreen = screens.NewEssentialToolScreen("", essentialTools, preselectedEssential)
	case ModernToolScreen:
		tools, err := m.config.LoadTools()
		if err != nil { m.err = err; newScreen = screens.NewWelcomeScreen(); break }
		modernTools := filterToolsByCategory(tools, "modern")
		preselectedModern := filterToolsByCategory(m.selectedTools, "modern")
		newScreen = screens.NewModernToolScreen("", modernTools, preselectedModern)
	case FontScreen:
		fonts, err := m.config.LoadFonts()
		if err != nil { m.err = err; newScreen = screens.NewWelcomeScreen(); break }
		newScreen = screens.NewFontScreen("", fonts, m.selectedFonts)
	case LanguageScreen:
		langs, errL := m.config.LoadLanguages()
		if errL != nil { m.err = fmt.Errorf("Lang load error: %v", errL); newScreen = screens.NewWelcomeScreen(); break }
		newScreen = screens.NewLanguageScreen("", langs, m.selectedLanguages)
	case DotfilesScreen: newScreen = screens.NewDotfilesScreen()
	case FinishScreen: newScreen = screens.NewFinishScreen()
	default: m.err = fmt.Errorf("invalid target screen: %d", targetScreen); newScreen = screens.NewWelcomeScreen()
	}

	m.activeModel = newScreen 

	// --- Initialize and Size New Screen --- 
	if m.activeModel != nil {
		// Recalculate available space 
		appHorizontalMargin := styles.AppStyle.GetHorizontalFrameSize()
		appVerticalMargin := styles.AppStyle.GetVerticalFrameSize()
		availableWidth := m.width - appHorizontalMargin
		availableHeight := m.height - appVerticalMargin
		if availableWidth < 0 { availableWidth = 0 }
		if availableHeight < 0 { availableHeight = 0 }

		indicatorHeight := 0
		// Use m.currentScreen (which is now targetScreen)
		if m.currentScreen >= EssentialToolScreen && m.currentScreen < FinishScreen {
			m.stepIndicator.SetWidth(availableWidth)
			indicatorHeight = lipgloss.Height(m.stepIndicator.View()) + 1
		}
		
		footerStr := styles.HelpStyle.Render("q/Ctrl+c: quit") // Base footer
		if m.err != nil { footerStr = styles.ErrorStyle.Render(m.err.Error()) }
		footerHeight := lipgloss.Height(footerStr) + 2

		childWidth := availableWidth
		childHeight := availableHeight - indicatorHeight - footerHeight
		if childHeight < 1 { childHeight = 1 } 

		// Send WindowSizeMsg 
		updatedModel, _ := m.activeModel.Update(tea.WindowSizeMsg{Width: childWidth, Height: childHeight})
		m.activeModel = updatedModel

		// Initialize model
		initCmd = m.activeModel.Init()
	}
	return initCmd
}

// Helper function to filter tools by category
func filterToolsByCategory(tools []*interfaces.Tool, category string) []*interfaces.Tool {
	filtered := make([]*interfaces.Tool, 0)
	for _, tool := range tools {
		if tool.Category == category {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// View method - Using explicit height constraint + JoinVertical
func (m *Model) View() string {
	if !m.screenReady {
		return styles.AppStyle.Render(styles.SubtitleStyle.Render("Initializing...")) 
	}

	// --- Calculate Available Space ---
	appHorizontalMargin := styles.AppStyle.GetHorizontalFrameSize()
	appVerticalMargin := styles.AppStyle.GetVerticalFrameSize()
	availableWidth := m.width - appHorizontalMargin
	availableHeight := m.height - appVerticalMargin
	if availableWidth < 0 { availableWidth = 0 }
	if availableHeight < 0 { availableHeight = 0 }

	// --- Prepare Indicator --- 
	indicatorView := ""
	indicatorHeight := 0
	if m.currentScreen >= EssentialToolScreen && m.currentScreen < FinishScreen { 
		m.stepIndicator.SetWidth(availableWidth) // Set width first
		indicatorView = m.stepIndicator.View()
		indicatorHeight = lipgloss.Height(indicatorView) // Get actual height
        // Don't add newline height here, JoinVertical will handle spacing
	}

	// --- Prepare Footer --- 
	footerStr := ""
	if m.err != nil { footerStr = styles.ErrorStyle.Render(m.err.Error()) } else { footerStr = styles.HelpStyle.Render("q/Ctrl+c: quit") } 
	footerHeight := lipgloss.Height(footerStr)

    // --- Calculate Height for Active View ---
    spacingHeight := 0 
    if indicatorView != "" { spacingHeight++ } // Space after indicator
    spacingHeight += 2 // Space before footer
	activeViewHeight := availableHeight - indicatorHeight - footerHeight - spacingHeight
	if activeViewHeight < 1 { activeViewHeight = 1 } // Min height 1

	// --- Get Raw Active View ---
	activeViewRaw := ""
	if m.activeModel != nil {
		activeViewRaw = m.activeModel.View()
	} else {
		activeViewRaw = styles.ErrorStyle.Render("Error: No active screen model.")
	}
    
    // --- Create Constrained Active View Block ---
    activeViewBlock := lipgloss.NewStyle().
        Width(availableWidth).
        Height(activeViewHeight).
        MaxWidth(availableWidth).
        MaxHeight(activeViewHeight).
        Render(activeViewRaw)


    // --- Assemble Final View using Place or JoinVertical carefully ---
    var finalViewContent []string

    if indicatorView != "" {
        finalViewContent = append(finalViewContent, indicatorView)
        finalViewContent = append(finalViewContent, "") // Spacer
    }
    finalViewContent = append(finalViewContent, activeViewBlock)
    finalViewContent = append(finalViewContent, "") // Spacer
    finalViewContent = append(finalViewContent, footerStr)

    joinedView := lipgloss.JoinVertical(lipgloss.Top, finalViewContent...)

	// Apply outer AppStyle margin last
	return styles.AppStyle.Render(joinedView)
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