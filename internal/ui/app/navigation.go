package app

import (
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

	switch m.currentScreen {
	case SystemScreen:
		sysInfo, err := system.Detect()
		if err != nil {
			m.err = err
			return nil
		}
		m.systemScreen = screens.ShowSystemInfo(sysInfo)
		return nil
	// ... repeat for other screens as in app.go ...
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
	// ... repeat for other screens as in app.go ...
	}
} 