package screens

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/detector"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// SystemScreen handles the system information display
type SystemScreen struct {
	info     *system.Info
	quitting bool
	done     bool
	width    int
	height   int
}

// NewSystemScreen creates a new system info screen
func NewSystemScreen(info *system.Info) *SystemScreen {
	return &SystemScreen{
		info: info,
	}
}

// Init implements tea.Model
func (s *SystemScreen) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (s *SystemScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			s.quitting = true
			return s, tea.Quit
		case "enter":
			s.done = true
			return s, tea.Quit
		}
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	}

	return s, nil
}

// View implements tea.Model
func (s *SystemScreen) View() string {
	if s.quitting {
		return ""
	}

	// Create system info display
	title := styles.TitleStyle.Render("System Information")
	
	// Get package manager info
	pmInfo := "Not detected"
	if pmType, err := detector.DetectPackageManager(); err == nil {
		pmInfo = string(pmType)
	}

	// Create info sections
	systemInfo := styles.InfoStyle.Render(fmt.Sprintf("OS: %s %s", s.info.Distro, s.info.Version))
	archInfo := styles.InfoStyle.Render(fmt.Sprintf("Architecture: %s", s.info.Arch))
	pmInfoStyled := styles.InfoStyle.Render(fmt.Sprintf("Package Manager: %s", pmInfo))

	// Combine all elements
	content := fmt.Sprintf("%s\n\n%s\n%s\n%s",
		title,
		systemInfo,
		archInfo,
		pmInfoStyled,
	)

	return content
}

// Finished returns true if the screen exited normally (not by quitting)
func (s *SystemScreen) Finished() bool {
	return s.done && !s.quitting
}

// ShowSystemInfo returns the SystemScreen model to be managed by the main app
func ShowSystemInfo(info *system.Info) *SystemScreen {
	screen := NewSystemScreen(info)
	return screen
} 