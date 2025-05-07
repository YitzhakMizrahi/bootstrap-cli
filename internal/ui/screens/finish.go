package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FinishScreen is a placeholder screen.
type FinishScreen struct { done bool; quitting bool; width, height int }
func NewFinishScreen() *FinishScreen { return &FinishScreen{} }
func (s *FinishScreen) Init() tea.Cmd { return nil }
func (s *FinishScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter": // Enter or q quits here
			s.quitting = true
			return s, tea.Quit
		}
	case tea.WindowSizeMsg: s.width = msg.Width; s.height = msg.Height
	}
	return s, nil
}
func (s *FinishScreen) View() string {
	title := styles.TitleStyle.Render("Setup Complete!")
	body := styles.NormalTextStyle.Render("\nYour selections would be installed now.\n\nPress Enter or q to exit.")
	return lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, title+"\n"+body)
}
func (s *FinishScreen) Finished() bool { return s.done && !s.quitting } // Not strictly needed if Update always Quits 