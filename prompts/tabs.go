package prompts

import (
	"github.com/charmbracelet/lipgloss"
)

// TabState represents the state of the tabs
type TabState struct {
	Tabs       []string
	ActiveTab  int
	TabContent map[int]string
}

// NewTabState creates a new tab state
func NewTabState(tabs []string) TabState {
	return TabState{
		Tabs:       tabs,
		ActiveTab:  0,
		TabContent: make(map[int]string),
	}
}

// Next moves to the next tab
func (t *TabState) Next() {
	t.ActiveTab = (t.ActiveTab + 1) % len(t.Tabs)
}

// Prev moves to the previous tab
func (t *TabState) Prev() {
	t.ActiveTab = (t.ActiveTab - 1 + len(t.Tabs)) % len(t.Tabs)
}

// SetContent sets the content for a tab
func (t *TabState) SetContent(index int, content string) {
	t.TabContent[index] = content
}

// GetContent gets the content for the active tab
func (t *TabState) GetContent() string {
	return t.TabContent[t.ActiveTab]
}

// View renders the tabs
func (t *TabState) View(width int) string {
	if len(t.Tabs) == 0 {
		return ""
	}

	// Define styles for tabs
	tabBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "",
		BottomRight: "",
	}

	inactive := lipgloss.NewStyle().
		Border(tabBorder, false, false, false, false).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Foreground(lipgloss.Color("#BDBDBD")).
		Padding(0, 2).
		Render

	active := lipgloss.NewStyle().
		Border(tabBorder, false, false, false, false).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Bold(true).
		Padding(0, 2).
		Render

	// Render tabs
	var renderedTabs []string
	for i, tab := range t.Tabs {
		if i == t.ActiveTab {
			renderedTabs = append(renderedTabs, active(tab))
		} else {
			renderedTabs = append(renderedTabs, inactive(tab))
		}
	}

	// Join tabs
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Create a style for the content
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(width - 4) // Adjust width to account for padding/borders

	// Get the content for the active tab
	content := t.GetContent()

	// Join tabs and content
	return lipgloss.JoinVertical(lipgloss.Left, row, contentStyle.Render(content))
} 