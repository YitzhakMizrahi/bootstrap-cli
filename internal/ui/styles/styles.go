// Package styles provides minimal ASCII styling for UI components in the bootstrap-cli application.
package styles

import "github.com/charmbracelet/lipgloss"

// Nord-inspired color theme
const (
	NordPolarNight1 = lipgloss.Color("#2E3440") // Darkest Background
	NordPolarNight2 = lipgloss.Color("#3B4252")
	NordPolarNight3 = lipgloss.Color("#434C5E") // Darker Background / Borders
	NordPolarNight4 = lipgloss.Color("#4C566A")

	NordSnowStorm1 = lipgloss.Color("#D8DEE9") // Slightly Dim Text
	NordSnowStorm2 = lipgloss.Color("#E5E9F0") // Normal Text
	NordSnowStorm3 = lipgloss.Color("#ECEFF4") // Brighter Text / Highlights

	NordFrostGreen  = lipgloss.Color("#8FBCBB") // Accent Color 1
	NordFrostBlue   = lipgloss.Color("#88C0D0") // Accent Color 2
	NordFrostPurple = lipgloss.Color("#B48EAD") // Accent Color 3
	NordFrostLightBlue = lipgloss.Color("#81A1C1") // Accent Color 4

	NordAuroraRed    = lipgloss.Color("#BF616A") // Error
	NordAuroraOrange = lipgloss.Color("#D08770")
	NordAuroraYellow = lipgloss.Color("#EBCB8B") // Warning
	NordAuroraGreen  = lipgloss.Color("#A3BE8C") // Success
	NordAuroraPurple = lipgloss.Color("#B48EAD")

	// Map to our style variables
	ColorBackground    = NordPolarNight1
	ColorSubtleBorder  = NordPolarNight3
	ColorNormalText    = NordSnowStorm2
	ColorDimText       = NordSnowStorm1
	ColorBrightText    = NordSnowStorm3
	ColorAccent        = NordFrostGreen // Using Frost Green as primary accent
	ColorAccentAlt     = NordFrostBlue    // Alt accent
	ColorSuccess       = NordAuroraGreen
	ColorWarning       = NordAuroraYellow
	ColorError         = NordAuroraRed
	ColorSpinner       = ColorAccent
	ColorProgressEmpty = NordPolarNight3
	ColorProgressFull  = ColorAccent
)

var (
	// General
	BaseStyle = lipgloss.NewStyle().Padding(0, 1)

	AppStyle = lipgloss.NewStyle().Margin(1, 2)

	// Text Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBrightText). // Use brighter text for titles
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorDimText).
			MarginBottom(1)

	NormalTextStyle = lipgloss.NewStyle().
			Foreground(ColorNormalText)

	SelectedTextStyle = lipgloss.NewStyle().
				Foreground(ColorAccent). // Use Frost Green for selected text
				Bold(true)

	UnselectedTextStyle = lipgloss.NewStyle().
				Foreground(ColorDimText) // Dimmer text for unselected

	HelpStyle = lipgloss.NewStyle().
			Foreground(NordPolarNight4). // Darker gray for help
			Italic(true)

	// Status Messages
	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	ErrorStyle   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	InfoStyle    = lipgloss.NewStyle().Foreground(ColorAccentAlt) // Use Frost Blue for info

	// Borders & Layout
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSubtleBorder)

	FocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorAccent)

	// List Styles
	ListTitleStyle = TitleStyle.Copy().Foreground(ColorAccentAlt) // Frost blue title for lists
	ListItemStyle  = NormalTextStyle.Copy().Padding(0, 0, 0, 2)  // Indent list items
	// Use styling on the delegate itself rather than a prefix string here
	// ListSelectedItemStyle = SelectedTextStyle.Copy().Padding(0, 0, 0, 0).SetString("➤ ")

	// Help / Key Map Style for Bubble Tea list
	KeyMapStyle = HelpStyle.Copy().Italic(false) // Reuse HelpStyle, maybe not italic

	// Specific Components
	StepIndicatorHeaderStyle = SubtitleStyle.Copy().MarginBottom(0)
	// Styles for individual steps - will be used in step_indicator.go View
	StepStyle = lipgloss.NewStyle().Padding(0, 1) // Base padding for each step
	StepCompletedStyle = StepStyle.Copy().Foreground(ColorDimText) // Dim completed steps
	StepCurrentStyle   = StepStyle.Copy().Foreground(ColorAccent).Bold(true).Background(ColorSubtleBorder) // Highlight current step
	StepPendingStyle   = StepStyle.Copy().Foreground(ColorNormalText)
	StepErrorStyle     = StepStyle.Copy().Foreground(ColorBrightText).Background(ColorError).Bold(true)

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().Foreground(ColorSpinner)

	// Progress Bar
	ProgressStyle = lipgloss.NewStyle().
			Foreground(ColorProgressFull).
			Background(ColorProgressEmpty)
)

// Helper functions from your original file (can be kept if useful)
// We might integrate these into lipgloss layouts directly.

// JoinVertical joins multiple strings vertically with a blank line
// func JoinVertical(parts ...string) string {
//  return strings.Join(parts, "\n\n")
// }

// JoinHorizontal joins multiple strings horizontally with two spaces
// func JoinHorizontal(parts ...string) string {
//  return strings.Join(parts, "  ")
// }

// ASCII helpers for headers and separators (can be recreated with lipgloss borders/padding)
// func AsciiHeader(title string) string {
//  return "=== " + title + " ==="
// }

// func AsciiSeparator() string {
//  return strings.Repeat("-", 40)
// } 