package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectionItem represents an item in the selection list
type SelectionItem struct {
	title       string
	description string
	value       interface{}
}

// Title returns the title of the item
func (i SelectionItem) Title() string { return i.title }

// Description returns the description of the item
func (i SelectionItem) Description() string { return i.description }

// FilterValue returns the value to filter on
func (i SelectionItem) FilterValue() string { return i.title }

// SelectionComponent represents an interactive selection component
type SelectionComponent struct {
	list          list.Model
	selectedItems []SelectionItem
	multiSelect   bool
	title        string
	width        int
	height       int
}

// NewSelectionComponent creates a new selection component
func NewSelectionComponent(title string, items []SelectionItem, multiSelect bool) *SelectionComponent {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	return &SelectionComponent{
		list:        l,
		multiSelect: multiSelect,
		title:       title,
		width:       80,
		height:      20,
	}
}

// SetSize sets the width and height of the component
func (s *SelectionComponent) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.list.SetSize(width, height)
}

// GetSelected returns the selected items
func (s *SelectionComponent) GetSelected() []SelectionItem {
	if s.multiSelect {
		return s.selectedItems
	}
	if item, ok := s.list.SelectedItem().(SelectionItem); ok {
		return []SelectionItem{item}
	}
	return nil
}

// Update handles component updates
func (s *SelectionComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "enter":
			if !s.multiSelect {
				return s, tea.Quit
			}
			if item, ok := s.list.SelectedItem().(SelectionItem); ok {
				// Toggle selection
				found := false
				for i, selected := range s.selectedItems {
					if selected.title == item.title {
						s.selectedItems = append(s.selectedItems[:i], s.selectedItems[i+1:]...)
						found = true
						break
					}
				}
				if !found {
					s.selectedItems = append(s.selectedItems, item)
				}
			}
		case "space":
			if s.multiSelect {
				if item, ok := s.list.SelectedItem().(SelectionItem); ok {
					// Toggle selection
					found := false
					for i, selected := range s.selectedItems {
						if selected.title == item.title {
							s.selectedItems = append(s.selectedItems[:i], s.selectedItems[i+1:]...)
							found = true
							break
						}
					}
					if !found {
						s.selectedItems = append(s.selectedItems, item)
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		s.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

// View renders the component
func (s *SelectionComponent) View() string {
	if s.multiSelect {
		// Add selection indicators
		items := s.list.Items()
		for i, item := range items {
			if si, ok := item.(SelectionItem); ok {
				selected := false
				for _, selectedItem := range s.selectedItems {
					if selectedItem.title == si.title {
						selected = true
						break
					}
				}
				prefix := "[ ]"
				if selected {
					prefix = "[x]"
				}
				si.title = fmt.Sprintf("%s %s", prefix, strings.TrimPrefix(si.title, "[ ] "))
				items[i] = si
			}
		}
		s.list.SetItems(items)
	}

	return s.list.View()
}

// Init initializes the component
func (s *SelectionComponent) Init() tea.Cmd {
	return nil
} 