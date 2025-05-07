package screens

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// WelcomeScreen handles the welcome screen display
type WelcomeScreen struct {
	quitting bool
	done     bool
	width    int
	height   int
}

// NewWelcomeScreen creates a new welcome screen
func NewWelcomeScreen() *WelcomeScreen {
	return &WelcomeScreen{}
}

// Init implements tea.Model
func (w *WelcomeScreen) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (w *WelcomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			w.quitting = true
			return w, tea.Quit
		case "enter":
			w.done = true
			return w, tea.Quit
		}
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
	}

	return w, nil
}

// View implements tea.Model
func (w *WelcomeScreen) View() string {
	if w.quitting {
		return ""
	}

	// Create welcome message
	welcome := styles.TitleStyle.Render("✨ Welcome to Bootstrap CLI ✨")
	description := styles.InfoStyle.Render("Setup your development environment with ease")

	// Combine all elements
	content := fmt.Sprintf("%s\n\n%s",
		welcome,
		description,
	)

	return content
}

// Finished returns true if the screen exited normally (not by quitting)
func (w *WelcomeScreen) Finished() bool {
	return w.done && !w.quitting
}

// ShowWelcomeScreen returns the WelcomeScreen model to be managed by the main app
func ShowWelcomeScreen() *WelcomeScreen {
	screen := NewWelcomeScreen()
	return screen
} 