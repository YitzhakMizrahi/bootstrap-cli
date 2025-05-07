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

	// Create the bar
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	// Format with percentage
	return fmt.Sprintf("%s %3.0f%%", styles.ProgressStyle.Render(bar), percent*100)
} 