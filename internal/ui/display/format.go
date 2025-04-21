package display

import (
	"fmt"
	"strings"
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

// Style constants for text formatting
const (
	StyleBold      = "\033[1m"
	StyleDim       = "\033[2m"
	StyleItalic    = "\033[3m"
	StyleUnderline = "\033[4m"
	StyleBlink     = "\033[5m"
	StyleReverse   = "\033[7m"
	StyleHidden    = "\033[8m"
)

// Box drawing characters
const (
	BoxTopLeft     = "┌"
	BoxTopRight    = "┐"
	BoxBottomLeft  = "└"
	BoxBottomRight = "┘"
	BoxHorizontal  = "─"
	BoxVertical    = "│"
)

// Progress bar characters
const (
	ProgressEmpty  = "░"
	ProgressFilled = "█"
)

// Formatter provides methods for text formatting
type Formatter struct {
	Color     string
	Style     string
	Width     int
	Alignment string
}

// NewFormatter creates a new formatter with default settings
func NewFormatter() *Formatter {
	return &Formatter{
		Color:     ColorReset,
		Style:     "",
		Width:     80,
		Alignment: "left",
	}
}

// WithColor sets the color for the formatter
func (f *Formatter) WithColor(color string) *Formatter {
	f.Color = color
	return f
}

// WithStyle sets the style for the formatter
func (f *Formatter) WithStyle(style string) *Formatter {
	f.Style = style
	return f
}

// WithWidth sets the width for the formatter
func (f *Formatter) WithWidth(width int) *Formatter {
	f.Width = width
	return f
}

// WithAlignment sets the alignment for the formatter
func (f *Formatter) WithAlignment(alignment string) *Formatter {
	f.Alignment = strings.ToLower(alignment)
	return f
}

// Format formats the text with the current settings
func (f *Formatter) Format(text string) string {
	// Apply color and style
	formatted := f.Color + f.Style + text + ColorReset

	// Apply alignment
	switch f.Alignment {
	case "center":
		padding := (f.Width - len(text)) / 2
		if padding > 0 {
			formatted = strings.Repeat(" ", padding) + formatted
		}
	case "right":
		padding := f.Width - len(text)
		if padding > 0 {
			formatted = strings.Repeat(" ", padding) + formatted
		}
	}

	return formatted
}

// Box creates a box around the text
func (f *Formatter) Box(text string) string {
	lines := strings.Split(text, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Add padding
	maxWidth += 2

	// Create top border
	result := BoxTopLeft + strings.Repeat(BoxHorizontal, maxWidth) + BoxTopRight + "\n"

	// Add each line
	for _, line := range lines {
		padding := strings.Repeat(" ", maxWidth-len(line))
		result += BoxVertical + " " + line + padding + BoxVertical + "\n"
	}

	// Add bottom border
	result += BoxBottomLeft + strings.Repeat(BoxHorizontal, maxWidth) + BoxBottomRight

	return f.Format(result)
}

// Header creates a header with optional underline
func (f *Formatter) Header(text string, level int) string {
	var result string
	switch level {
	case 1:
		result = f.WithStyle(StyleBold).Format(strings.ToUpper(text))
		result += "\n" + strings.Repeat("=", len(text))
	case 2:
		result = f.WithStyle(StyleBold).Format(text)
		result += "\n" + strings.Repeat("-", len(text))
	default:
		result = f.WithStyle(StyleBold).Format(text)
	}
	return result
}

// List creates a formatted list
func (f *Formatter) List(items []string, bullet string) string {
	if bullet == "" {
		bullet = "•"
	}

	var result strings.Builder
	for _, item := range items {
		result.WriteString(fmt.Sprintf("%s %s\n", bullet, item))
	}
	return f.Format(result.String())
}

// Table creates a formatted table
func (f *Formatter) Table(headers []string, rows [][]string) string {
	if len(headers) == 0 || len(rows) == 0 {
		return ""
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Create the table
	var result strings.Builder

	// Add headers
	for i, header := range headers {
		result.WriteString("| ")
		result.WriteString(f.WithStyle(StyleBold).Format(
			header + strings.Repeat(" ", colWidths[i]-len(header))))
		result.WriteString(" ")
	}
	result.WriteString("|\n")

	// Add separator
	for i, width := range colWidths {
		result.WriteString("|")
		result.WriteString(strings.Repeat("-", width+2))
		if i == len(colWidths)-1 {
			result.WriteString("|")
		}
	}
	result.WriteString("\n")

	// Add rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				result.WriteString("| ")
				result.WriteString(cell + strings.Repeat(" ", colWidths[i]-len(cell)))
				result.WriteString(" ")
			}
		}
		result.WriteString("|\n")
	}

	return f.Format(result.String())
}

// Success formats a success message
func (f *Formatter) Success(message string) string {
	return f.WithColor(ColorGreen).WithStyle(StyleBold).Format("✓ " + message)
}

// Error formats an error message
func (f *Formatter) Error(message string) string {
	return f.WithColor(ColorRed).WithStyle(StyleBold).Format("✗ " + message)
}

// Warning formats a warning message
func (f *Formatter) Warning(message string) string {
	return f.WithColor(ColorYellow).WithStyle(StyleBold).Format("⚠ " + message)
}

// Info formats an info message
func (f *Formatter) Info(message string) string {
	return f.WithColor(ColorBlue).WithStyle(StyleBold).Format("ℹ " + message)
}

// ProgressBar creates a progress bar with percentage
func (f *Formatter) ProgressBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	if width <= 0 {
		width = 20
	}

	filledWidth := (percent * width) / 100
	emptyWidth := width - filledWidth

	bar := strings.Repeat(ProgressFilled, filledWidth) + strings.Repeat(ProgressEmpty, emptyWidth)
	percentage := fmt.Sprintf("%d%%", percent)

	return f.Format(fmt.Sprintf("[%s] %s", bar, percentage))
}

// Indent creates an indented text block
func (f *Formatter) Indent(text string, indentLevel int) string {
	if indentLevel <= 0 {
		return text
	}

	indent := strings.Repeat("  ", indentLevel)
	lines := strings.Split(text, "\n")
	
	var result strings.Builder
	for i, line := range lines {
		result.WriteString(indent + line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}
	
	return f.Format(result.String())
}

// CodeBlock creates a code block with optional syntax highlighting
func (f *Formatter) CodeBlock(code string, language string) string {
	// For now, we'll just use a simple box with the language name
	// In a real implementation, we would use a syntax highlighter library
	header := fmt.Sprintf("```%s", language)
	footer := "```"
	
	// Create a box with the code
	boxedCode := f.Box(code)
	
	return f.Format(header + "\n" + boxedCode + "\n" + footer)
}

// Collapsible creates a collapsible section
func (f *Formatter) Collapsible(title string, content string, expanded bool) string {
	var result strings.Builder
	
	// Add the title with an expand/collapse indicator
	indicator := "▼"
	if !expanded {
		indicator = "▶"
		content = "" // Hide content if collapsed
	}
	
	result.WriteString(f.WithStyle(StyleBold).Format(indicator + " " + title) + "\n")
	
	// Add the content if expanded
	if expanded && content != "" {
		result.WriteString(f.Indent(content, 1) + "\n")
	}
	
	return f.Format(result.String())
} 