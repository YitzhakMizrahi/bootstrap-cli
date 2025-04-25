// Package styles provides styling for UI components in the bootstrap-cli application.
package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")  // Purple
	successColor   = lipgloss.Color("#50FA7B")  // Green
	warningColor   = lipgloss.Color("#FFB86C")  // Orange
	inactiveColor  = lipgloss.Color("#6272A4")  // Gray
	textColor      = lipgloss.Color("#F8F8F2")  // Light gray
	highlightColor = lipgloss.Color("#FF79C6")  // Pink
	errorColor     = lipgloss.Color("#FF5555")  // Red

	// BaseStyle provides the base styling for all UI components
	BaseStyle = lipgloss.NewStyle().
		Foreground(textColor)

	// TitleStyle provides styling for screen titles and headers
	TitleStyle = BaseStyle.Copy().
		Bold(true).
		Foreground(highlightColor).
		Padding(0, 1).
		MarginBottom(1)

	// NormalStyle provides styling for regular text content
	NormalStyle = BaseStyle.Copy()

	// SelectedStyle provides styling for selected items in lists and menus
	SelectedStyle = BaseStyle.Copy().
		Foreground(successColor).
		Bold(true)

	// UnselectedStyle provides styling for unselected items in lists and menus
	UnselectedStyle = BaseStyle.Copy().
		Foreground(inactiveColor)

	// HelpStyle provides styling for help text and instructions
	HelpStyle = BaseStyle.Copy().
		Foreground(inactiveColor).
		Italic(true)

	// ErrorStyle provides styling for error messages
	ErrorStyle = BaseStyle.Copy().
		Foreground(errorColor).
		Bold(true)

	// HeaderStyle provides styling for section headers
	HeaderStyle = BaseStyle.Copy().
		Bold(true).
		Foreground(highlightColor).
		MarginBottom(1)

	// SuccessStyle provides styling for success messages
	SuccessStyle = BaseStyle.Copy().
		Foreground(successColor).
		Bold(true)

	// WarningStyle provides styling for warning messages
	WarningStyle = BaseStyle.Copy().
		Foreground(warningColor).
		Bold(true)

	// InfoStyle provides styling for informational messages
	InfoStyle = BaseStyle.Copy().
		Foreground(primaryColor).
		Bold(true)
)

// JoinVertical joins multiple strings vertically with proper spacing
func JoinVertical(parts ...string) string {
	return strings.Join(parts, "\n\n")
}

// JoinHorizontal joins multiple strings horizontally with proper spacing
func JoinHorizontal(parts ...string) string {
	return strings.Join(parts, "  ")
} 