package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Spinner styles
	spinnerStyle = lipgloss.NewStyle().Foreground(highlight)

	// Progress bar styles
	progressBarStyle = lipgloss.NewStyle().Foreground(special)
	progressFullStyle = lipgloss.NewStyle().Foreground(highlight)
	progressEmptyStyle = lipgloss.NewStyle().Foreground(subtle)

	// Input field styles
	inputStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(highlight).
		Padding(0, 1)

	// Selection styles
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(highlight).
		Bold(true).
		PaddingLeft(2)
	
	unselectedItemStyle = lipgloss.NewStyle().
		Foreground(subtle).
		PaddingLeft(2)
)

// SpinnerModel represents a loading spinner
type SpinnerModel struct {
	spinner  spinner.Model
	text     string
	quitting bool
}

// NewSpinner creates a new spinner with the given text
func NewSpinner(text string) *SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return &SpinnerModel{
		spinner: s,
		text:    text,
	}
}

func (m *SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *SpinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.text)
}

// ProgressModel represents a progress bar
type ProgressModel struct {
	progress progress.Model
	percent  float64
	text     string
}

// NewProgress creates a new progress bar
func NewProgress(text string) *ProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	return &ProgressModel{
		progress: p,
		text:     text,
	}
}

func (m *ProgressModel) Init() tea.Cmd {
	return nil
}

func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case float64:
		m.percent = msg
		return m, nil
	default:
		return m, nil
	}
}

func (m *ProgressModel) View() string {
	prog := m.progress.ViewAs(m.percent)
	return fmt.Sprintf("%s\n%s", m.text, progressBarStyle.Render(prog))
}

// SelectionModel represents a selection menu
type SelectionModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	done     bool
}

// NewSelection creates a new selection menu
func NewSelection(choices []string) *SelectionModel {
	return &SelectionModel{
		choices:  choices,
		selected: make(map[int]struct{}),
	}
}

func (m *SelectionModel) Init() tea.Cmd {
	return nil
}

func (m *SelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, exists := m.selected[m.cursor]
			if exists {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "q", "ctrl+c":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *SelectionModel) View() string {
	var s strings.Builder

	for i, choice := range m.choices {
		prefix := "  "
		if i == m.cursor {
			prefix = "→ "
		}
		if _, exists := m.selected[i]; exists {
			prefix = prefix[:len(prefix)-2] + "✓ "
		}

		style := unselectedItemStyle
		if i == m.cursor {
			style = selectedItemStyle
		}

		s.WriteString(style.Render(prefix + choice + "\n"))
	}

	return s.String()
}

// GetSelected returns the selected items
func (m *SelectionModel) GetSelected() []string {
	var selected []string
	for i := range m.selected {
		selected = append(selected, m.choices[i])
	}
	return selected
}

// RunSpinner runs a spinner until the done channel is closed
func RunSpinner(text string, done chan struct{}) {
	p := tea.NewProgram(NewSpinner(text))
	go func() {
		<-done
		p.Quit()
	}()
	p.Run()
}

// RunProgress runs a progress bar
func RunProgress(text string, progress chan float64) {
	p := tea.NewProgram(NewProgress(text))
	go func() {
		for percent := range progress {
			p.Send(percent)
			if percent >= 1.0 {
				p.Quit()
				return
			}
		}
	}()
	p.Run()
}

// RunSelection runs a selection menu and returns the selected items
func RunSelection(title string, choices []string) []string {
	m := NewSelection(choices)
	p := tea.NewProgram(m)
	p.Run()
	return m.GetSelected()
} 