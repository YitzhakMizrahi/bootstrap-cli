package components

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectorItem represents an item in the selection list
type SelectorItem struct {
	title       string
	description string
	item        interface{}
	selected    bool
}

// FilterValue implements list.Item interface
func (i SelectorItem) FilterValue() string { return i.title }

// Title returns the title for the list item
func (i SelectorItem) Title() string {
	checkbox := "[ ]"
	if i.selected {
		checkbox = "[✓]"
		return styles.SelectedStyle.Render(fmt.Sprintf("%s %s", checkbox, i.title))
	}
	return styles.UnselectedStyle.Render(fmt.Sprintf("%s %s", checkbox, i.title))
}

// Description returns the description for the list item
func (i SelectorItem) Description() string { return i.description }

// BaseSelector is the main model for selection
type BaseSelector struct {
	list      list.Model
	selected  map[int]interface{}
	quitting  bool
	width     int
	height    int
	done      bool
	title     string
	selectable func(interface{}) bool
}

// NewBaseSelector creates a new base selector
func NewBaseSelector(title string) *BaseSelector {
	// Create a custom delegate
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(2) // Give more space for items
	delegate.SetSpacing(1) // Add spacing between items
	delegate.UpdateFunc = func(_ tea.Msg, _ *list.Model) tea.Cmd {
		return nil
	}

	// Create and configure the list
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(false)
	l.Title = styles.TitleStyle.Render(title)

	return &BaseSelector{
		list:     l,
		selected: make(map[int]interface{}),
		title:    title,
	}
}

// Init implements tea.Model
func (s *BaseSelector) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and updates the selector state
func (s *BaseSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			s.quitting = true
			s.done = false
			return s, tea.Quit
		case "enter":
			// Only handle Enter if we have selections
			if len(s.selected) > 0 {
				s.done = true
				s.quitting = false
				return s, nil
			}
		case " ":
			// Get current item and index
			idx := s.list.Index()
			items := s.list.Items()
			if idx >= len(items) {
				return s, nil
			}

			// Get the current item and cast it
			currentItem, ok := items[idx].(SelectorItem)
			if !ok {
				return s, nil
			}

			// Only allow selection if item is selectable
			if s.selectable == nil || s.selectable(currentItem.item) {
				// Toggle selection state
				if _, exists := s.selected[idx]; exists {
					delete(s.selected, idx)
					currentItem.selected = false
				} else {
					s.selected[idx] = currentItem.item
					currentItem.selected = true
				}

				// Update just this item in the list
				items[idx] = currentItem
				s.list.SetItems(items)
			}
		}
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetWidth(msg.Width)
		s.list.SetHeight(msg.Height)
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

// View implements tea.Model
func (s *BaseSelector) View() string {
	if s.quitting {
		return ""
	}

	// Create help text
	help := styles.HelpStyle.Render("↑/↓: navigate • space: select/deselect • enter: confirm • q: quit")

	// Combine all elements with proper spacing
	content := fmt.Sprintf("%s\n\n%s\n\n%s",
		s.list.Title,
		s.list.View(),
		help,
	)

	return content
}

// Finished returns true if the selector was completed normally (not quit)
func (s *BaseSelector) Finished() bool {
	return s.done && !s.quitting
}

// GetSelected returns the map of selected items
func (s *BaseSelector) GetSelected() map[int]interface{} {
	if !s.Finished() {
		return nil
	}
	return s.selected
}

// SetItems sets the items to be displayed in the selector
func (s *BaseSelector) SetItems(items interface{}, titleFn func(interface{}) string, descFn func(interface{}) string) {
	// Convert items to a slice of list.Item
	var listItems []list.Item

	// Handle different types of input
	switch v := items.(type) {
	case []interface{}:
		listItems = make([]list.Item, len(v))
		for i, item := range v {
			listItems[i] = SelectorItem{
				title:       titleFn(item),
				description: descFn(item),
				item:        item,
				selected:    false,
			}
		}
	case []*interfaces.Tool:
		listItems = make([]list.Item, len(v))
		for i, item := range v {
			listItems[i] = SelectorItem{
				title:       titleFn(item),
				description: descFn(item),
				item:        item,
				selected:    false,
			}
		}
	case []*interfaces.Font:
		listItems = make([]list.Item, len(v))
		for i, item := range v {
			listItems[i] = SelectorItem{
				title:       titleFn(item),
				description: descFn(item),
				item:        item,
				selected:    false,
			}
		}
	case []*interfaces.Language:
		listItems = make([]list.Item, len(v))
		for i, item := range v {
			listItems[i] = SelectorItem{
				title:       titleFn(item),
				description: descFn(item),
				item:        item,
				selected:    false,
			}
		}
	default:
		fmt.Printf("Warning: Unhandled item type in SetItems: %T\n", items)
		return
	}

	// Only clear selections if this is a new set of items
	if len(s.list.Items()) == 0 {
		s.selected = make(map[int]interface{})
	}
	
	s.list.SetItems(listItems)
}

// SetSize sets the width and height of the selector
func (s *BaseSelector) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.list.SetWidth(width)
	s.list.SetHeight(height)
} 