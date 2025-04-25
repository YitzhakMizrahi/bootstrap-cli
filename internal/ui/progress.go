package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	progressStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	completedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B"))

	currentStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFB86C"))

	pendingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))
)

// Section represents a section in the installation process
type Section struct {
	Name   string
	Status string // "pending", "current", "completed"
}

// RenderProgressBar renders a beautiful progress bar with sections
func RenderProgressBar(sections []Section, width int) string {
	var parts []string

	for i, section := range sections {
		// Add section number
		sectionNum := fmt.Sprintf("%d", i+1)
		
		// Style based on status
		var styledSection string
		switch section.Status {
		case "completed":
			styledSection = completedStyle.Render("✓ " + section.Name)
		case "current":
			styledSection = currentStyle.Render("➤ " + section.Name)
		default:
			styledSection = pendingStyle.Render("○ " + section.Name)
		}

		parts = append(parts, fmt.Sprintf("%s. %s", sectionNum, styledSection))
	}

	// Join with separator
	bar := strings.Join(parts, " → ")

	// Add header
	header := progressStyle.Render("Installation Progress")
	
	return fmt.Sprintf("%s\n%s\n", header, bar)
} 