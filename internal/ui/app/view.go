package app

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
)

// View implements tea.Model for Model
func (m *Model) View() string {
	if !m.screenReady {
		return ""
	}

	output := "\x1b[2J\x1b[H" // ANSI clear screen
	output += m.stepIndicator.View() + "\n\n"

	// Only render the current screen's content
	screenContent := ""
	switch m.currentScreen {
	case WelcomeScreen:
		if m.welcomeScreen != nil {
			screenContent = m.welcomeScreen.View()
		}
	case SystemScreen:
		if m.systemScreen != nil {
			screenContent = m.systemScreen.View()
		}
	case ToolScreen:
		if m.toolScreen != nil {
			screenContent = m.toolScreen.View()
		}
	case FontScreen:
		if m.fontScreen != nil {
			screenContent = m.fontScreen.View()
		}
	case LanguageScreen:
		if m.languageScreen != nil {
			screenContent = m.languageScreen.View()
		}
	// ... handle other screens ...
	}
	output += screenContent

	// Add error/help bar
	if m.err != nil {
		output += "\n\n" + styles.ErrorStyle.Render(m.err.Error())
	} else {
		output += "\n\n" + styles.HelpStyle.Render("↑/↓: navigate • space: select/deselect • enter: confirm • q: quit")
	}

	return output
} 