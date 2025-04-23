package ui

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/tools"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)

	categoryStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		MarginTop(1)

	warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		MarginTop(1)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	docStyle = lipgloss.NewStyle().
		Padding(1, 2, 1, 2)
)

// ToolItem represents a tool in the selection list
type ToolItem struct {
	tool     *install.Tool
	selected bool
}

// FilterValue implements list.Item interface
func (t ToolItem) FilterValue() string { return t.tool.Name }

// Title returns the title for the list item
func (t ToolItem) Title() string {
	title := t.tool.Name
	if t.selected {
		title = "✓ " + title
	}
	return title
}

// Description returns the description for the list item
func (t ToolItem) Description() string {
	return t.tool.Description
}

// ToolSelector is the main model for tool selection
type ToolSelector struct {
	categories []tools.ToolCategory
	lists      []list.Model
	activeList int
	selected   map[string]bool
	quitting   bool
	width     int
	height    int
	done      bool
}

// NewToolSelector creates a new tool selector
func NewToolSelector(toolList []string) *ToolSelector {
	items := make([]list.Item, len(toolList))
	for i, tool := range toolList {
		items[i] = ToolItem{
			tool: &install.Tool{
				Name:        tool,
				PackageName: tool,
			},
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#7D56F4"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#626262"))

	l := list.New(items, delegate, 0, 0)
	l.Title = "Available Tools"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = categoryStyle

	return &ToolSelector{
		lists:      []list.Model{l},
		selected:   make(map[string]bool),
		width:     80,
		height:    20,
	}
}

// Init implements tea.Model
func (m *ToolSelector) Init() tea.Cmd {
	return nil
}

// moveToNextCategory moves to the next category or finishes if on last category
func (m *ToolSelector) moveToNextCategory() tea.Cmd {
	if m.activeList < len(m.lists)-1 {
		m.activeList++
		return nil
	}
	// If we're on the last category, mark as done
	m.done = true
	return tea.Quit
}

// Update implements tea.Model
func (m *ToolSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		for i := range m.lists {
			m.lists[i].SetSize(msg.Width-4, msg.Height-8)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			// Switch between categories
			m.activeList = (m.activeList + 1) % len(m.lists)
			return m, nil
		case "shift+tab":
			// Switch between categories (reverse)
			m.activeList--
			if m.activeList < 0 {
				m.activeList = len(m.lists) - 1
			}
			return m, nil
		case " ":
			// Toggle selection of current item
			item := m.lists[m.activeList].SelectedItem().(ToolItem)
			m.selected[item.tool.Name] = !m.selected[item.tool.Name]

			// Update the item in the list
			items := m.lists[m.activeList].Items()
			for i, listItem := range items {
				if listItem.(ToolItem).tool.Name == item.tool.Name {
					items[i] = ToolItem{
						tool:     item.tool,
						selected: m.selected[item.tool.Name],
					}
					break
				}
			}
			m.lists[m.activeList].SetItems(items)
			return m, nil
		case "enter":
			// Move to next category or finish if on last category
			return m, m.moveToNextCategory()
		}
	}

	// Handle list updates
	var cmd tea.Cmd
	m.lists[m.activeList], cmd = m.lists[m.activeList].Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *ToolSelector) View() string {
	if m.quitting {
		return ""
	}

	// Create category navigation
	var categories []string
	for i, cat := range m.categories {
		if i == m.activeList {
			categories = append(categories, selectedStyle.Render("● "+cat.Name))
		} else {
			categories = append(categories, "○ "+cat.Name)
		}
	}
	nav := strings.Join(categories, " | ")

	// Create help text
	help := helpStyle.Render("↑/↓: navigate • space: select/deselect • enter: confirm category • tab: switch category • q: quit")
	warning := warningStyle.Render("Note: Installation requires sudo privileges. Run with: sudo bootstrap-cli tools install")

	// Combine all elements
	content := fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n%s",
		titleStyle.Render("Tool Selection"),
		nav,
		m.lists[m.activeList].View(),
		help,
		warning,
	)

	// Apply final styling
	return docStyle.Width(m.width - 4).Render(content)
}

// Finished returns true if the selector exited normally (not by quitting)
func (m *ToolSelector) Finished() bool {
	return m.done && !m.quitting
}

// GetSelectedTools returns the list of selected tool names
func (m *ToolSelector) GetSelectedTools() []string {
	var selected []string
	for name, isSelected := range m.selected {
		if isSelected {
			selected = append(selected, name)
		}
	}
	return selected
} 