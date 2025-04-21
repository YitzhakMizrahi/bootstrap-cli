package progress

import (
	"fmt"
	"strings"
	"time"
)

// Color constants for ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Style represents the visual style of the progress bar
type Style struct {
	Prefix    string
	Suffix    string
	FillChar  string
	EmptyChar string
	Color     string
}

// DefaultStyle returns the default progress bar style
func DefaultStyle() Style {
	return Style{
		Prefix:    "[",
		Suffix:    "]",
		FillChar:  "=",
		EmptyChar: " ",
		Color:     ColorBlue,
	}
}

// ProgressBar represents a progress bar
type ProgressBar struct {
	Total       int
	Current     int
	Width       int
	Style       Style
	ShowPercent bool
	ShowTime    bool
	startTime   time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		Total:       total,
		Current:     0,
		Width:       30,
		Style:       DefaultStyle(),
		ShowPercent: true,
		ShowTime:    true,
		startTime:   time.Now(),
	}
}

// SetStyle sets the style of the progress bar
func (p *ProgressBar) SetStyle(style Style) {
	p.Style = style
}

// Update updates the progress bar
func (p *ProgressBar) Update(current int) {
	p.Current = current
	p.Display()
}

// formatDuration formats a duration in a human-readable format
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	
	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// Display displays the progress bar
func (p *ProgressBar) Display() {
	// Calculate the percentage
	percent := float64(p.Current) / float64(p.Total) * 100

	// Calculate the number of filled characters
	filled := int(float64(p.Width) * float64(p.Current) / float64(p.Total))
	empty := p.Width - filled

	// Create the progress bar
	bar := p.Style.Prefix +
		p.Style.Color +
		strings.Repeat(p.Style.FillChar, filled) +
		strings.Repeat(p.Style.EmptyChar, empty) +
		ColorReset +
		p.Style.Suffix

	// Add percentage if enabled
	percentStr := ""
	if p.ShowPercent {
		percentStr = fmt.Sprintf(" %3.0f%%", percent)
	}

	// Add elapsed time if enabled
	timeStr := ""
	if p.ShowTime {
		elapsed := time.Since(p.startTime)
		timeStr = fmt.Sprintf(" [%s]", formatDuration(elapsed))
	}

	// Print the progress bar
	fmt.Printf("\r%s%s%s", bar, percentStr, timeStr)
}

// Finish finishes the progress bar
func (p *ProgressBar) Finish() {
	p.Current = p.Total
	p.Display()
	fmt.Println()
}

// Spinner represents a spinner
type Spinner struct {
	Frames []string
	Index  int
	Prefix string
	Suffix string
	Color  string
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		Index:  0,
		Prefix: "",
		Suffix: "",
		Color:  ColorBlue,
	}
}

// SetColor sets the color of the spinner
func (s *Spinner) SetColor(color string) {
	s.Color = color
}

// Update updates the spinner
func (s *Spinner) Update(message string) {
	// Get the current frame
	frame := s.Frames[s.Index]

	// Update the index
	s.Index = (s.Index + 1) % len(s.Frames)

	// Print the spinner with color
	fmt.Printf("\r%s%s%s%s%s %s", s.Prefix, s.Color, frame, ColorReset, s.Suffix, message)
}

// Finish finishes the spinner
func (s *Spinner) Finish(message string) {
	fmt.Printf("\r%s%s✓%s %s\n", s.Prefix, ColorGreen, ColorReset, message)
} 