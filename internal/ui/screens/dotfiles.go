package screens

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dotfileState int

const (
	dsAskManage dotfileState = iota
	dsAskURL
)

// DotfilesScreen prompts the user for dotfiles configuration.
type DotfilesScreen struct {
	state         dotfileState
	manageDotfiles bool
	repoURL       string
	textInput     textinput.Model
	quitting      bool
	finished      bool
	width         int
	height        int
	err           error // To display potential input errors
}

func NewDotfilesScreen() *DotfilesScreen {
	ti := textinput.New()
	ti.Placeholder = "github_username/repository_name"
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 50 // Initial width, will adjust on WindowSizeMsg

	return &DotfilesScreen{
		state:     dsAskManage,
		textInput: ti,
	}
}

func (s *DotfilesScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s *DotfilesScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		// Adjust text input width based on screen size if needed
		s.textInput.Width = min(s.width-4, 80) // Example adjustment

	case tea.KeyMsg:
		switch s.state {
		case dsAskManage:
			switch msg.String() {
			case "ctrl+c", "q":
				s.quitting = true
				return s, tea.Quit
			case "n", "N":
				s.manageDotfiles = false
				s.finished = true
				return s, nil // Signal finished
			case "y", "Y":
				s.manageDotfiles = true
				s.state = dsAskURL
				cmds = append(cmds, s.textInput.Focus()) // Focus input field
				return s, tea.Batch(cmds...)
			}
		case dsAskURL:
			switch msg.String() {
			case "ctrl+c", "q":
				s.quitting = true
				return s, tea.Quit
			case "enter":
				url := strings.TrimSpace(s.textInput.Value())
				if url == "" {
					s.err = fmt.Errorf("repository URL cannot be empty")
				} else {
					// Basic validation (contains /)
					if !strings.Contains(url, "/") {
						s.err = fmt.Errorf("invalid format, use user/repo")
					} else {
						s.repoURL = url
						s.finished = true
						s.err = nil
						return s, nil // Signal finished
					}
				}
			// If enter was pressed but input is invalid, fall through to update input
			fallthrough // Allow enter to potentially modify the input value itself
			default:
				// Handle text input updates
				s.textInput, cmd = s.textInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return s, tea.Batch(cmds...)
}

func (s *DotfilesScreen) View() string {
	var content strings.Builder

	title := styles.TitleStyle.Render("Dotfiles Configuration")
	content.WriteString(title)
	content.WriteString("\n\n")

	switch s.state {
	case dsAskManage:
		body := styles.NormalTextStyle.Render("Manage dotfiles by cloning a GitHub repository? (y/n)")
		content.WriteString(body)
	case dsAskURL:
		body := styles.NormalTextStyle.Render("Enter GitHub repository URL (e.g., username/repo):")
		content.WriteString(body)
		content.WriteString("\n\n")
		content.WriteString(s.textInput.View())
		if s.err != nil {
			content.WriteString("\n\n")
			content.WriteString(styles.ErrorStyle.Render(s.err.Error()))
		}
	}

	// Add help/footer
	content.WriteString("\n\n")
	content.WriteString(styles.HelpStyle.Render("Enter: Confirm"))

	// Use lipgloss.Place for centering
	return lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, content.String())
}

func (s *DotfilesScreen) Finished() bool { return s.finished && !s.quitting }

// GetSelection returns the user's choice regarding dotfiles.
func (s *DotfilesScreen) GetSelection() (manage bool, repoURL string) {
	if !s.Finished() {
		return false, ""
	}
	return s.manageDotfiles, s.repoURL
}

// Helper for min width calculation
func min(a, b int) int {
	if a < b { return a }
	return b
} 