// Package components provides UI components for the bootstrap-cli application.
package components

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	StatusPending   = "pending"
	StatusCurrent   = "current"
	StatusCompleted = "completed"
	StatusError     = "error" // Added for more comprehensive status
)

// Step represents a single step in a sequence.
type Step struct {
	Name   string
	Status string // e.g., StatusPending, StatusCurrent, StatusCompleted, StatusError
}

// Model is the Bubble Tea model for the step indicator.
type Model struct {
	steps      []Step
	currentIdx int // Index of the current step
	width      int // Store available width for layout calculations
	Title      string // Optional title for the step group, e.g., "Installation Progress"
}

// NewModel creates a new step indicator model with the given step names.
// Initially, all steps are pending and the first step is current.
func NewModel(stepNames []string) Model {
	steps := make([]Step, len(stepNames))
	for i, name := range stepNames {
		steps[i] = Step{Name: name, Status: StatusPending}
	}

	m := Model{
		steps:      steps,
		currentIdx: -1, // No step is current initially, call SetCurrentStep to activate
	}
	if len(steps) > 0 {
		m.SetCurrentStep(0) // Make the first step current by default if steps exist
	}
	return m
}

// SetWidth sets the width for the component, used for layout.
func (m *Model) SetWidth(width int) {
	m.width = width
}

// Init does nothing for this simple component.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetCurrentStepMsg is a message to update the current step.
type SetCurrentStepMsg struct {
	Index int
}

// SetStepStatusMsg is a message to update a specific step's status.
type SetStepStatusMsg struct {
	Index  int
	Status string
	Name   string // Optional: if you want to update the name too
}

// Update handles messages for the step indicator.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width // Store width for responsive rendering
		return m, nil
	case SetCurrentStepMsg:
		m.SetCurrentStep(msg.Index)
		return m, nil
	case SetStepStatusMsg:
		if msg.Index >= 0 && msg.Index < len(m.steps) {
			m.steps[msg.Index].Status = msg.Status
			if msg.Name != "" { // Allow updating name too
				m.steps[msg.Index].Name = msg.Name
			}
		}
		return m, nil
	}
	return m, nil
}

// View renders the step indicator.
func (m Model) View() string {
	if len(m.steps) == 0 { return "" }

	// separator := styles.HelpStyle.Copy().Faint(true).SetString(" ").String() // Unused variable removed
    
    // Define styles for the step blocks
    baseStepStyle := lipgloss.NewStyle().
        Padding(0, 1). // Horizontal padding within each block
        MarginRight(1) // Space between blocks

	completedStyle := baseStepStyle.Copy().
		Foreground(styles.NordAuroraGreen). // Green text for completed
        Faint(true) // Make completed steps less prominent

	currentStyle := baseStepStyle.Copy().
        Foreground(styles.ColorBrightText). // Bright text for current
        Background(styles.ColorAccent). // Accent background
		Bold(true)

	pendingStyle := baseStepStyle.Copy().
        Foreground(styles.NordPolarNight4) // Dark gray for pending

	errorStyle := baseStepStyle.Copy().
        Foreground(styles.ColorBrightText).
        Background(styles.ColorError).
		Bold(true)

	var stepViews []string
	// Use blank identifier for index 'i' as it's not used
	for _, step := range m.steps { 
		var styledStep string
        name := step.Name // Use only the name

		switch step.Status {
		case StatusCompleted:
			styledStep = completedStyle.Render("✓ " + name) // Checkmark prefix
		case StatusCurrent:
			styledStep = currentStyle.Render(name) // Current step stands out
		case StatusError:
			styledStep = errorStyle.Render("✘ " + name) // Error prefix
		case StatusPending:
			fallthrough
		default:
			styledStep = pendingStyle.Render(name) // Pending are faint
		}
		stepViews = append(stepViews, styledStep)
	}

    // Join horizontally. Lipgloss handles joining styled blocks.
	joinedSteps := lipgloss.JoinHorizontal(lipgloss.Top, stepViews...)

	return joinedSteps // Return joined steps
}

// SetCurrentStep updates the current step and recalculates statuses.
func (m *Model) SetCurrentStep(idx int) {
	if idx < 0 || idx >= len(m.steps) {
		// Consider logging an error or handling this case more gracefully
		m.currentIdx = -1 // or idx to clamp, or just return
		// Clear all statuses or set to pending if idx is invalid
		for i := range m.steps {
			m.steps[i].Status = StatusPending
		}
		return
	}
	m.currentIdx = idx
	for i := range m.steps {
		if i < idx {
			m.steps[i].Status = StatusCompleted
		} else if i == idx {
			m.steps[i].Status = StatusCurrent
		} else {
			m.steps[i].Status = StatusPending
		}
	}
}

// SetSteps allows replacing the entire list of steps.
func (m *Model) SetSteps(stepNames []string) {
	newSteps := make([]Step, len(stepNames))
	for i, name := range stepNames {
		newSteps[i] = Step{Name: name, Status: StatusPending}
	}
	m.steps = newSteps
	if len(m.steps) > 0 {
		m.SetCurrentStep(0) // Default to first step as current
	} else {
		m.currentIdx = -1
	}
}

// GetSteps returns the current steps.
func (m *Model) GetSteps() []Step {
	return m.steps
}

// CurrentStepIndex returns the index of the current step.
func (m *Model) CurrentStepIndex() int {
	return m.currentIdx
} 