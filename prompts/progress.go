package prompts

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressIndicator represents a step progress indicator
type ProgressIndicator struct {
	Current     int
	Total       int
	Width       int
	ActiveColor string
	InactiveColor string
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(total int) ProgressIndicator {
	return ProgressIndicator{
		Current:      0,
		Total:        total,
		Width:        40,
		ActiveColor:  "#7D56F4",
		InactiveColor: "#3C3C3C",
	}
}

// Next advances the progress indicator
func (p *ProgressIndicator) Next() {
	if p.Current < p.Total {
		p.Current++
	}
}

// Prev moves the progress indicator back
func (p *ProgressIndicator) Prev() {
	if p.Current > 0 {
		p.Current--
	}
}

// View renders the progress indicator
func (p *ProgressIndicator) View() string {
	if p.Total == 0 {
		return ""
	}
	
	percentage := float64(p.Current) / float64(p.Total)
	
	// Create the progress bar
	var bar strings.Builder
	
	// Write progress
	bar.WriteString("[")
	
	completedWidth := int(float64(p.Width) * percentage)
	for i := 0; i < p.Width; i++ {
		if i < completedWidth {
			bar.WriteString("█")
		} else {
			bar.WriteString("░")
		}
	}
	
	bar.WriteString("]")
	
	// Style the progress bar
	styledBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.ActiveColor)).
		Render(bar.String())
	
	// Add percentage and step counter
	percentText := fmt.Sprintf(" %d%%", int(percentage*100))
	stepsText := fmt.Sprintf(" Step %d of %d", p.Current, p.Total)
	
	// Combine percentage and step counter
	infoText := percentText + stepsText
	
	// Style the text
	styledInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BBBBBB")).
		Render(infoText)
	
	return styledBar + styledInfo
}

// SimpleView renders a simpler view with just circles
func (p *ProgressIndicator) SimpleView() string {
	if p.Total == 0 {
		return ""
	}
	
	var circles []string
	
	for i := 0; i < p.Total; i++ {
		var circle string
		if i < p.Current {
			// Completed step
			circle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(p.ActiveColor)).
				Render("●")
		} else if i == p.Current {
			// Current step
			circle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(p.ActiveColor)).
				Bold(true).
				Render("◉")
		} else {
			// Future step
			circle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(p.InactiveColor)).
				Render("○")
		}
		circles = append(circles, circle)
	}
	
	return strings.Join(circles, " ")
} 