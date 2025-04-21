package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const logo = `
██████╗  ██████╗  ██████╗ ████████╗███████╗████████╗██████╗  █████╗ ██████╗ 
██╔══██╗██╔═══██╗██╔═══██╗╚══██╔══╝██╔════╝╚══██╔══╝██╔══██╗██╔══██╗██╔══██╗
██████╔╝██║   ██║██║   ██║   ██║   ███████╗   ██║   ██████╔╝███████║██████╔╝
██╔══██╗██║   ██║██║   ██║   ██║   ╚════██║   ██║   ██╔══██╗██╔══██║██╔═══╝ 
██████╔╝╚██████╔╝╚██████╔╝   ██║   ███████║   ██║   ██║  ██║██║  ██║██║     
╚═════╝  ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     
`

var (
	// Style definitions
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	logoStyle = lipgloss.NewStyle().
		Foreground(highlight).
		Bold(true)

	infoStyle = lipgloss.NewStyle().
		Foreground(subtle).
		Italic(true)

	versionStyle = lipgloss.NewStyle().
		Foreground(special).
		Bold(true)
)

// DisplayWelcome shows the welcome screen with animations
func DisplayWelcome(version string) {
	// Clear screen
	fmt.Print("\033[H\033[2J")

	// Display logo with typing animation
	lines := strings.Split(logo, "\n")
	for _, line := range lines {
		fmt.Println(logoStyle.Render(line))
		time.Sleep(50 * time.Millisecond)
	}

	// Display version and info
	fmt.Println()
	fmt.Println(versionStyle.Render("Version: " + version))
	fmt.Println(infoStyle.Render("Your development environment, your way."))
	fmt.Println()
	
	// Display loading message
	fmt.Println(infoStyle.Render("Loading your personalized development experience..."))
	time.Sleep(500 * time.Millisecond)
}

// DisplayError shows an error message with style
func DisplayError(msg string) {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)
	
	fmt.Println(errorStyle.Render("Error: " + msg))
}

// DisplaySuccess shows a success message with style
func DisplaySuccess(msg string) {
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)
	
	fmt.Println(successStyle.Render("Success: " + msg))
} 