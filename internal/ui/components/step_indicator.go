// Package components provides UI components for the bootstrap-cli application.
package components

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
)

// Step represents a step in the installation process
type Step struct {
	Name   string
	Status string // "pending", "current", "completed"
}

// StepIndicator is a component that renders a beautiful step indicator showing installation progress
type StepIndicator struct {
	steps []Step
}

// NewStepIndicator creates a new step indicator component
func NewStepIndicator(steps []Step) *StepIndicator {
	return &StepIndicator{
		steps: steps,
	}
}

// View renders the step indicator
func (s *StepIndicator) View() string {
	var parts []string

	for i, step := range s.steps {
		// Add step number
		stepNum := fmt.Sprintf("%d", i+1)
		
		// Style based on status
		var styledStep string
		switch step.Status {
		case "completed":
			styledStep = styles.StepCompletedStyle.Render("✓ " + step.Name)
		case "current":
			styledStep = styles.StepCurrentStyle.Render("➤ " + step.Name)
		default:
			styledStep = styles.StepPendingStyle.Render("○ " + step.Name)
		}

		parts = append(parts, fmt.Sprintf("%s. %s", stepNum, styledStep))
	}

	// Join with separator
	indicator := strings.Join(parts, " → ")

	// Add header
	header := styles.StepIndicatorHeaderStyle.Render("Installation Progress")
	
	return fmt.Sprintf("%s\n%s\n", header, indicator)
}

// GetSteps returns the current steps
func (s *StepIndicator) GetSteps() []Step {
	return s.steps
}

// SetSteps updates the steps
func (s *StepIndicator) SetSteps(steps []Step) {
	s.steps = steps
} 