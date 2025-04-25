package utils

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
)

// Section represents a section in the progress bar
type Section struct {
	Name   string
	Status string
}

// RenderSimpleProgressBar renders a simple progress bar with the given percentage
func RenderSimpleProgressBar(percent float64) string {
	// Calculate the filled and empty portions
	filled := int(percent * 20)
	empty := 20 - filled

	// Create the bar
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	// Format with percentage
	return fmt.Sprintf("%s %3.0f%%", styles.ProgressStyle.Render(bar), percent*100)
}

// RenderSectionProgressBar creates a text-based progress bar showing the status of installation sections
func RenderSectionProgressBar(sections []Section) string {
	var result strings.Builder
	
	// Calculate the maximum section name length
	maxNameLen := 0
	for _, section := range sections {
		if len(section.Name) > maxNameLen {
			maxNameLen = len(section.Name)
		}
	}
	
	// Create the progress bar
	result.WriteString("\nInstallation Progress:\n")
	for _, section := range sections {
		// Pad the name for alignment
		paddedName := section.Name + strings.Repeat(" ", maxNameLen-len(section.Name))
		
		// Create the status indicator
		var status string
		switch section.Status {
		case "completed":
			status = "✓"
		case "current":
			status = "→"
		case "pending":
			status = "·"
		case "failed":
			status = "✗"
		default:
			status = " "
		}
		
		// Format the line
		result.WriteString(fmt.Sprintf("  %s [%s] %s\n", paddedName, status, section.Status))
	}
	
	return result.String()
} 