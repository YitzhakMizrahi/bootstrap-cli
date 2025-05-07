package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DotfilesScreen is a placeholder screen.
type DotfilesScreen struct { done bool; quitting bool; width, height int }
func NewDotfilesScreen() *DotfilesScreen { return &DotfilesScreen{} }
func (s *DotfilesScreen) Init() tea.Cmd { return nil }
func (s *DotfilesScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": s.quitting = true; return s, tea.Quit
		case "enter": s.done = true; return s, nil // Signal finished
		}
	case tea.WindowSizeMsg: s.width = msg.Width; s.height = msg.Height
	}
	return s, nil
}
func (s *DotfilesScreen) View() string {
	title := styles.TitleStyle.Render("Dotfiles Configuration")
	body := styles.NormalTextStyle.Render("\nDotfiles setup placeholder.\n\nPress Enter to continue.")
	// Basic centering or layout
	return lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, title+"\n"+body)
}
func (s *DotfilesScreen) Finished() bool { return s.done && !s.quitting } 