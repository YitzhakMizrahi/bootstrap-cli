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
	help := styles.HelpStyle.Render("Press Enter to continue • Press q to quit")

	// Combine all elements
	content := fmt.Sprintf("%s\n\n%s\n\n%s",
		welcome,
		description,
		help,
	)

	return content
}

// Finished returns true if the screen exited normally (not by quitting)
func (w *WelcomeScreen) Finished() bool {
	return w.done && !w.quitting
}

// ShowWelcomeScreen displays the welcome screen and returns true if user wants to continue
func ShowWelcomeScreen() bool {
	screen := NewWelcomeScreen()
	p := tea.NewProgram(screen)

	model, err := p.Run()
	if err != nil {
		fmt.Printf("Error running welcome screen: %v\n", err)
		return false
	}

	if welcomeScreen, ok := model.(*WelcomeScreen); ok {
		return welcomeScreen.Finished()
	}

	return false
} 