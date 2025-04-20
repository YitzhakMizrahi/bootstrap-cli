package prompts

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// CustomSpinner wraps the bubbles spinner with styling
type CustomSpinner struct {
	Model  spinner.Model
	Text   string
	Style  lipgloss.Style
}

// NewSpinner creates a new spinner with default styling
func NewSpinner(text string) CustomSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4"))
	
	return CustomSpinner{
		Model: s,
		Text:  text,
		Style: style,
	}
}

// WithStyle sets a custom style for the spinner
func (s *CustomSpinner) WithStyle(style lipgloss.Style) *CustomSpinner {
	s.Style = style
	return s
}

// WithSpinnerType sets a custom spinner type
func (s *CustomSpinner) WithSpinnerType(spinnerType spinner.Spinner) *CustomSpinner {
	s.Model.Spinner = spinnerType
	return s
}

// View renders the spinner with text
func (s *CustomSpinner) View() string {
	spinnerView := s.Style.Render(s.Model.View())
	return spinnerView + " " + s.Text
} 