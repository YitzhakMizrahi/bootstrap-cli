package prompts

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CountdownTimer represents a countdown timer
type CountdownTimer struct {
	Duration time.Duration
	StartTime time.Time
	Remaining time.Duration
	OnComplete func() // Function to call when countdown completes
	Style lipgloss.Style
	Message string
}

// NewCountdownTimer creates a new countdown timer
func NewCountdownTimer(seconds int, message string, onComplete func()) CountdownTimer {
	return CountdownTimer{
		Duration: time.Duration(seconds) * time.Second,
		Remaining: time.Duration(seconds) * time.Second,
		OnComplete: onComplete,
		Style: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")),
		Message: message,
	}
}

// Start initializes the timer and returns a tick command
func (t *CountdownTimer) Start() tea.Cmd {
	t.StartTime = time.Now()
	return t.Tick()
}

// Tick is a command that advances the timer
func (t *CountdownTimer) Tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Update handles tick messages
func (t *CountdownTimer) Update(msg tea.Msg) (CountdownTimer, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		elapsed := time.Since(t.StartTime)
		t.Remaining = t.Duration - elapsed
		
		if t.Remaining <= 0 {
			// Timer completed
			return *t, t.Complete()
		}
		
		// Continue ticking
		return *t, t.Tick()
	}
	
	return *t, nil
}

// Complete returns a command that executes the OnComplete function
func (t *CountdownTimer) Complete() tea.Cmd {
	return func() tea.Msg {
		if t.OnComplete != nil {
			t.OnComplete()
		}
		return completeMsg{}
	}
}

// View renders the timer
func (t *CountdownTimer) View() string {
	seconds := int(t.Remaining.Seconds())
	
	// Format the remaining time
	var timeText string
	if seconds > 0 {
		timeText = fmt.Sprintf("%d", seconds)
	} else {
		timeText = "0"
	}
	
	styledTime := t.Style.Render(timeText)
	
	return fmt.Sprintf("%s %s...", t.Message, styledTime)
}

// Custom message types
type tickMsg struct{}
type completeMsg struct{}

// TimerModel is a Bubbletea model for a standalone timer
type TimerModel struct {
	Timer CountdownTimer
	Width int
	Height int
}

// NewTimerModel creates a new timer model
func NewTimerModel(seconds int, message string, onComplete func()) TimerModel {
	return TimerModel{
		Timer: NewCountdownTimer(seconds, message, onComplete),
		Width: 80,
		Height: 24,
	}
}

func (m TimerModel) Init() tea.Cmd {
	return m.Timer.Start()
}

func (m TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
		
	case tea.KeyMsg:
		// Handle key presses
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		
	case completeMsg:
		// Timer completed, quit the program
		return m, tea.Quit
	}
	
	var cmd tea.Cmd
	m.Timer, cmd = m.Timer.Update(msg)
	return m, cmd
}

func (m TimerModel) View() string {
	// Create a centered, styled view
	mainStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(m.Width / 2).
		Align(lipgloss.Center)
	
	content := m.Timer.View()
	
	// Add key instructions
	footer := "\nPress q to skip, Ctrl+C to cancel"
	
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		mainStyle.Render(content + footer),
	)
}

// RunTimer runs a standalone timer
func RunTimer(seconds int, message string, onComplete func()) error {
	model := NewTimerModel(seconds, message, onComplete)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
} 