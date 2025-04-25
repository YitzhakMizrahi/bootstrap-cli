package utils

import (
	"fmt"
	"strings"
)

// Section represents a section in the progress bar
type Section struct {
	Name   string
	Status string
}

// RenderProgressBar creates a text-based progress bar showing the status of installation sections
func RenderProgressBar(sections []struct {
	Name   string
	Status string
}, width int) string {
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