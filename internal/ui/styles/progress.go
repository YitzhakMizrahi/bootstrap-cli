// Package styles provides styling for UI components in the bootstrap-cli application.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	// StepIndicatorHeaderStyle is used for the step indicator header
	StepIndicatorHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	// StepCompletedStyle is used for completed steps
	StepCompletedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B"))

	// StepCurrentStyle is used for the current step
	StepCurrentStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFB86C"))

	// StepPendingStyle is used for pending steps
	StepPendingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))
) 