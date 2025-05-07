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
	steps      []Step
	currentIdx int
}

// NewStepIndicator creates a new step indicator component
func NewStepIndicator(steps []Step) *StepIndicator {
	return &StepIndicator{
		steps:      steps,
		currentIdx: 0,
	}
}

// SetCurrentStep sets the current step index
func (s *StepIndicator) SetCurrentStep(idx int) {
	if idx >= 0 && idx < len(s.steps) {
		s.currentIdx = idx
		for i := range s.steps {
			if i < idx {
				s.steps[i].Status = "completed"
			} else if i == idx {
				s.steps[i].Status = "current"
			} else {
				s.steps[i].Status = "pending"
			}
		}
	}
}

// View renders the step indicator
func (s *StepIndicator) View() string {
	var parts []string

	for i, step := range s.steps {
		var styledStep string
		switch step.Status {
		case "completed":
			styledStep = styles.StepCompletedStyle.Render("✔ " + step.Name)
		case "current":
			// Make current step more visually prominent
			styledStep = styles.StepCurrentStyle.Copy().Underline(true).Bold(true).Render("➤ " + step.Name)
		default:
			styledStep = styles.StepPendingStyle.Render("● " + step.Name)
		}
		if i == s.currentIdx {
			styledStep = styles.StepCurrentStyle.Copy().Underline(true).Bold(true).Render("➤ " + step.Name)
		}
		parts = append(parts, styledStep)
	}

	indicator := strings.Join(parts, "  ")
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