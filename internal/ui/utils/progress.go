package utils

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
)

// RenderSimpleProgressBar renders a simple progress bar with the given percentage
func RenderSimpleProgressBar(percent float64) string {
	// Calculate the filled and empty portions
	filled := int(percent * 20)
	empty := 20 - filled

	// Create the bar with modern Unicode characters
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	// Format with percentage and add a subtle border
	content := fmt.Sprintf("%s %3.0f%%", bar, percent*100)
	
	// Apply styling with padding and subtle border
	return styles.ProgressStyle.Render(content)
} 