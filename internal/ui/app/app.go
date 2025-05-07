// Package app provides the main application model for the bootstrap-cli UI.
package app

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	base_iface "github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/pipeline"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/shell"
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
	WelcomeScreen Screen = iota // 0
	ShellSelectionScreen        // 1
	EssentialToolScreen         // 2
	ModernToolScreen            // 3
	FontScreen                  // 4
	LanguageScreen              // 5
	DotfilesScreen              // 6
	InstallationScreen          // 7 // New screen for progress (Moved before Finish)
	FinishScreen                // 8
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
	selectedTools     []*pipeline.Tool
	selectedFonts     []*interfaces.Font
	selectedLanguages []*interfaces.Language
	systemInfo        *system.Info // Store detected system info
	selectedShell     *interfaces.Shell // Changed type from string
	shellManager      interfaces.ShellManager // Added ShellManager
}

// New creates a new application model
func New(config *config.Loader) *Model {
	rand.Seed(time.Now().UnixNano())
	
	// Adjusted step names for indicator
	stepNames := []string{
		"Shell", 
		"Essential Tools", 
		"Modern Tools",    
		"Fonts",
		"Languages",
		"Dotfiles",
		"Installation", // Added Installation step
		"Finish",
	}
	stepIndicatorModel := components.NewModel(stepNames)

	welcomeModel := screens.NewWelcomeScreen()

	// Initialize ShellManager (assuming NewManager exists in shell package)
	// We might need to pass a logger or other dependencies to NewManager if required.
	shellMgr, err := shell.NewManager() // Placeholder, might need args
	if err != nil {
		// This is a critical failure during app initialization.
		// We should probably panic or return an error from New if this happens.
		// For now, let's log and proceed with a nil manager, which will cause issues later
		// but allows the app to start for dev purposes. Proper error handling needed.
		fmt.Fprintf(os.Stderr, "Error initializing ShellManager: %v\n", err)
		shellMgr = nil // Or handle error more gracefully
	}

	m := &Model{
		instanceID:    rand.Int63(),
		currentScreen:  WelcomeScreen, 
		activeModel:   welcomeModel,
		config:        config,
		shellManager:  shellMgr, // Assign initialized shell manager
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

// installCompleteMsg is sent when the background installation goroutine finishes
// It might contain an error if the installer logic itself failed (not just pipeline steps)
type installCompleteMsg struct {
	err error
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
		// Show indicator starting from Shell selection screen
		if m.currentScreen >= ShellSelectionScreen && m.currentScreen < FinishScreen {
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

	// Handle the message indicating background installation finished
	case installCompleteMsg:
		// This message currently signifies the goroutine running the installer finished.
		// The actual success/failure is handled within the InstallationScreen via pipeline events.
		// We might use this msg later for final cleanup or error display if the installer itself crashes.
		if msg.err != nil {
			m.err = fmt.Errorf("installer runner failed: %w", msg.err)
			// Potentially transition to an error screen or update InstallationScreen state
		}
		// No transition needed here usually, InstallationScreen handles its exit.
		return m, nil
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
			if screen.Finished() { cmds = append(cmds, m.transitionTo(ShellSelectionScreen)) }
		case *screens.ShellSelectionScreen: 
			if screen.Finished() { 
				m.selectedShell = screen.GetSelected() 
				cmds = append(cmds, m.transitionTo(EssentialToolScreen))
			}
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
			if screen.Finished() { 
				cmds = append(cmds, m.transitionTo(InstallationScreen))
			}
		
		// Installation screen handles its own lifecycle, no check needed here
		case *screens.InstallationScreen:
			// No transition logic needed here as InstallationScreen quits itself
			break
		
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
	// Adjust index mapping for new steps
	// Shell=0, EssentialTools=1, ModernTools=2, Fonts=3, Languages=4, Dotfiles=5, Installation=6, Finish=7
	if targetScreen >= ShellSelectionScreen && targetScreen <= InstallationScreen { 
		visualStepIndex = int(targetScreen) - 1 // Shell(1)->0, Essential(2)->1, ..., Installation(7)->6
	} 
	// If target is FinishScreen or WelcomeScreen, visualStepIndex remains -1 (indicator hidden)
	m.stepIndicator.SetCurrentStep(visualStepIndex)

	// --- Create New Screen --- 
	var newScreen tea.Model
	var initCmd tea.Cmd
	switch targetScreen {
	case WelcomeScreen: newScreen = screens.NewWelcomeScreen()
	case ShellSelectionScreen: 
		// 1. Load all defined shells from configuration
		definedShells, err := m.config.LoadShells()
		if err != nil { 
			m.err = fmt.Errorf("failed to load shell definitions: %w", err)
			newScreen = screens.NewWelcomeScreen() 
			break 
		}

		// 2. Get shells available on the system
		var systemShellsInfo []*interfaces.ShellInfo
		if m.shellManager != nil {
			systemShellsInfo, err = m.shellManager.ListAvailable()
			if err != nil {
				m.err = fmt.Errorf("failed to list available system shells: %w", err)
				// Proceeding with empty systemShellsInfo, or could go to WelcomeScreen
				systemShellsInfo = []*interfaces.ShellInfo{} 
			}
		} else {
			m.err = fmt.Errorf("shell manager not initialized")
			newScreen = screens.NewWelcomeScreen()
			break
		}

		// 3. Filter defined shells to only those available on the system
		availableDisplayShells := make([]*interfaces.Shell, 0)
		systemShellMap := make(map[string]bool)
		for _, sysShell := range systemShellsInfo {
			systemShellMap[sysShell.Type] = true // Assuming Type is like "bash", "zsh"
			systemShellMap[sysShell.Current] = true // sysShell.Current might be the name or path
		}

		for _, defShell := range definedShells {
			if systemShellMap[defShell.Name] { // Assuming defShell.Name matches sysShell.Type or sysShell.Current
				availableDisplayShells = append(availableDisplayShells, defShell)
			}
		}

		// 4. Get the current shell's path or name for pre-selection
		currentShellIdentifier := ""
		if m.shellManager != nil {
			currentShellInfo, err := m.shellManager.DetectCurrent()
			if err != nil {
				m.err = fmt.Errorf("failed to detect current shell: %w", err)
				// Proceeding without a pre-selected current shell
			} else if currentShellInfo != nil {
				currentShellIdentifier = currentShellInfo.Current // Or currentShellInfo.Path, depending on what NewShellSelectionScreen expects
			}
		} else {
			// shellManager not initialized, error already set above
		}
		
		preselectedName := ""
		if m.selectedShell != nil {
			preselectedName = m.selectedShell.Name // Get name for preselection if already chosen once
		}

		newScreen = screens.NewShellSelectionScreen(
			"Please select your primary shell:",
			availableDisplayShells, // Pass the filtered list of installable/configurable shells
			currentShellIdentifier,   // Pass the detected current shell (name or path)
			preselectedName,
		)
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
	case InstallationScreen:
		fmt.Println("Transitioning to Installation Screen...") // Use fmt for now

		// --- Prepare for Installation --- 
		// 1. Gather all selections (ensure this is done before this transition)
		selectedPipelineTools := m.SelectedTools()
		// TODO: Gather fonts, languages, shell, dotfiles
		selectedToolNames := make([]string, len(selectedPipelineTools))
		for i, tool := range selectedPipelineTools {
			selectedToolNames[i] = tool.Name
		}

		// 2. Load full pipeline tool definitions (PLACEHOLDER - Needs real implementation)
		selectedPipelineTools = make([]*pipeline.Tool, 0)
		fmt.Println("TODO: Placeholder: Using empty tool list for pipeline installation start.")
		// TODO: Implement loading based on selectedToolNames

		// 3. Prepare Platform and PackageManager for Installer
		sysInfo, err := system.Detect()
		if err != nil {
			m.err = fmt.Errorf("failed to detect system info for install: %w", err)
			newScreen = screens.NewWelcomeScreen()
			break
		}
		pkgManagerFactory := factory.NewPackageManagerFactory()
		pkgManagerImpl, err := pkgManagerFactory.GetPackageManager() // base_iface.PackageManager
		if err != nil {
			m.err = fmt.Errorf("failed to detect package manager for install: %w", err)
			newScreen = screens.NewWelcomeScreen()
			break
		}
		
		// Adapt the PackageManager
		var pipelinePackageManager pipeline.PackageManager = &packageManagerAdapter{impl: pkgManagerImpl}
		fmt.Println("TODO: Verify and complete PackageManager adapter implementation for pipeline.")

		pipelinePlatform := &pipeline.Platform{
			OS:             sysInfo.OS,
			Arch:           sysInfo.Arch,
			PackageManager: sysInfo.PackageType, // TODO: Use GetName from adapter?
			Shell:          sysInfo.Shell,
		}

		// 4. Create Installer
		installer, err := pipeline.NewInstaller(pipelinePlatform, pipelinePackageManager)
		if err != nil {
			m.err = fmt.Errorf("failed to create installer: %w", err)
			newScreen = screens.NewWelcomeScreen() 
			break
		}

		// 5. Create the Installation Screen, passing the READ end of the progress channel
		newScreen = screens.NewInstallationScreen(installer.ProgressChan)

		// 6. Create command to run the installation in the background
		installCmd := func() tea.Msg {
			fmt.Println("Starting background installation process...")
			err := installer.InstallSelections(selectedPipelineTools)
			fmt.Println("Background installation process finished.")
			return installCompleteMsg{err: err} 
		}
		
		initCmd = tea.Batch(newScreen.Init(), installCmd)

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
func filterToolsByCategory(tools []*pipeline.Tool, category string) []*pipeline.Tool {
	filtered := make([]*pipeline.Tool, 0)
	for _, tool := range tools {
		if string(tool.Category) == category {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// View method - Removing debug prints
func (m *Model) View() string {
	if !m.screenReady {
		return styles.AppStyle.Render(styles.SubtitleStyle.Render("Initializing...")) 
	}

	var finalView strings.Builder // Use builder for efficiency

	// --- Calculate Available Space (Inside App Margin) ---
	appHorizontalMargin := styles.AppStyle.GetHorizontalFrameSize()
	appVerticalMargin := styles.AppStyle.GetVerticalFrameSize()
	availableWidth := m.width - appHorizontalMargin
	availableHeight := m.height - appVerticalMargin
	if availableWidth < 0 { availableWidth = 0 }
	if availableHeight < 0 { availableHeight = 0 }

	// --- Determine Indicator Height & Render ---
	indicatorView := ""
	indicatorHeight := 0
	// Show indicator starting from Shell selection screen
	if m.currentScreen >= ShellSelectionScreen && m.currentScreen < FinishScreen { 
		m.stepIndicator.SetWidth(availableWidth) 
		indicatorView = m.stepIndicator.View()
		indicatorHeight = lipgloss.Height(indicatorView) 
        if indicatorHeight > 0 { 
            finalView.WriteString(indicatorView)
            finalView.WriteString("\n") 
            indicatorHeight++ 
        } 
	}

	// --- Determine Footer Height ---
	footerStr := ""
	if m.err != nil { footerStr = styles.ErrorStyle.Render(m.err.Error()) } else { footerStr = styles.HelpStyle.Render("q/Ctrl+c: quit") } 
	footerHeight := lipgloss.Height(footerStr) + 2 

	// --- Calculate Height for Active View ---
	activeViewHeight := availableHeight - indicatorHeight - footerHeight
    if activeViewHeight < 1 { activeViewHeight = 1 }

	// --- Get Raw Active View --- 
	activeViewRaw := ""
	if m.activeModel != nil {
		activeViewRaw = m.activeModel.View()
	} else {
		activeViewRaw = styles.ErrorStyle.Render("Error: No active screen model.")
	}
    
	// --- Constrain and Render Active View --- 
    activeViewStyled := lipgloss.NewStyle().
		Width(availableWidth).
		Height(activeViewHeight).
		Render(activeViewRaw)

    finalView.WriteString(activeViewStyled) 

    // --- Render Footer --- 
    finalView.WriteString("\n\n") 
	finalView.WriteString(footerStr)

	// Apply outer margin last
	return styles.AppStyle.Render(finalView.String())
}

// SelectedTools returns the selected tools
func (m *Model) SelectedTools() []*pipeline.Tool {
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

// GetSelectedShell returns the selected shell
func (m *Model) GetSelectedShell() *interfaces.Shell {
	return m.selectedShell
}

// Placeholder adapter - NEEDS REAL IMPLEMENTATION and matching interfaces defined
type packageManagerAdapter struct {
	impl base_iface.PackageManager // The implementation from internal/packages
}

func (a *packageManagerAdapter) Install(pkg string) error { return a.impl.Install(pkg) }
func (a *packageManagerAdapter) Uninstall(pkg string) error { return a.impl.Uninstall(pkg) }
func (a *packageManagerAdapter) IsInstalled(pkg string) (bool, error) {
	return a.impl.IsInstalled(pkg)
}
func (a *packageManagerAdapter) Update() error { return a.impl.Update() }
func (a *packageManagerAdapter) SetupSpecialPackage(pkg string) error { 
	return a.impl.SetupSpecialPackage(pkg) 
}
func (a *packageManagerAdapter) IsPackageAvailable(pkg string) bool { 
	return a.impl.IsPackageAvailable(pkg) 
} 