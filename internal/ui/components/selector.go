package components

import (
	"fmt"
	"io"

	// "github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles" // Import our styles
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// Import lipgloss for direct use if needed
)

// SelectorItem represents an item in the selection list
type SelectorItem struct {
	title       string
	description string
	item        interface{} // The actual data item
	selected    bool
}

// FilterValue implements list.Item interface
func (i SelectorItem) FilterValue() string { return i.title }

// Title returns the title for the list item, now handling single-select focus and multi-select state.
func (i SelectorItem) Title(isSingleSelect bool, isFocused bool) string { 
	var prefix string
    var titleStr string 
    var fullTitle string

    baseTitleStyle := styles.NormalTextStyle 
    if isFocused {
        // For focused items, both prefix (if applicable) and text might change style
        baseTitleStyle = styles.SelectedTextStyle 
    }

	if isSingleSelect {
		if isFocused {
			prefix = styles.SelectedTextStyle.Render("(o) ")
		} else {
			prefix = styles.NormalTextStyle.Render("( ) ")
		}
        titleStr = baseTitleStyle.Render(i.title) // Title text uses focused/normal style
        fullTitle = prefix + titleStr
	} else { // Multi-select mode
		currentCheckboxStyle := styles.NormalTextStyle // Style for [ ] or [x]
        actualTitleStyle := styles.NormalTextStyle   // Style for the title text next to checkbox

        if i.selected {
            currentCheckboxStyle = styles.SelectedTextStyle // [x] takes selected style
        }

        if isFocused {
            // If the row is focused, both checkbox and title text might take focused style
            currentCheckboxStyle = styles.SelectedTextStyle 
            actualTitleStyle = styles.SelectedTextStyle
        }
        
        // Construct prefix based on selection and focus
        if i.selected {
            prefix = currentCheckboxStyle.Render("[x] ")
	} else {
            prefix = currentCheckboxStyle.Render("[ ] ")
        }
        titleStr = actualTitleStyle.Render(i.title)
        fullTitle = prefix + titleStr
    }
	return fullTitle
}

// Description returns the description for the list item
func (i SelectorItem) Description() string {
	return i.description
}

// --- Custom Delegate ---
// itemDelegate wraps the default delegate to customize rendering
type itemDelegate struct {
    list.DefaultDelegate
    singleSelectMode bool
}

func newItemDelegate(singleSelect bool) itemDelegate {
    d := list.NewDefaultDelegate()

    // Apply Nord styles
    d.Styles.SelectedTitle = styles.SelectedTextStyle.Copy().UnsetPadding().Foreground(styles.ColorAccent)
    d.Styles.SelectedDesc = styles.UnselectedTextStyle.Copy().UnsetPadding() // Keep desc dimmer even if selected title is bright
    d.Styles.NormalTitle = styles.NormalTextStyle.Copy().UnsetPadding() 
    d.Styles.NormalDesc = styles.UnselectedTextStyle.Copy().UnsetPadding() 

    d.SetHeight(2) // Allow space for description
    d.SetSpacing(1)

    return itemDelegate{
        DefaultDelegate: d,
        singleSelectMode: singleSelect,
    }
}

// Render overrides the default render to use the updated Title method.
// It now directly prints the string from item.Title() and adds a left border for focused items.
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	if item, ok := listItem.(*SelectorItem); ok {
		desc := item.Description()
		// Calculate available width, considering potential border+padding
		// Border = 1, Padding = 1. Total width reduction = 2
		availableWidth := m.Width() 
		contentWidth := availableWidth - 2 
		if contentWidth < 0 {
			contentWidth = 0
		}

        isFocused := index == m.Index()

        // Get the fully pre-styled title string (with prefixes and internal styling)
        renderedTitle := item.Title(d.singleSelectMode, isFocused)
		var descStyle lipgloss.Style

		if isFocused {
            descStyle = d.Styles.SelectedDesc
		} else {
            descStyle = d.Styles.NormalDesc
		}
        
        renderedDesc := descStyle.Width(contentWidth).Render(desc)

        // Combine title and description vertically
        // Apply width limiting to the title string as well before joining
        // Note: Applying width directly to renderedTitle might break ANSI styles if not handled carefully.
        // For now, assume title fits or truncates visually. A better approach might involve lipgloss text wrapping.
        content := lipgloss.JoinVertical(lipgloss.Left, renderedTitle, renderedDesc)

        // Define styles for focused (border) and non-focused (padding)
        focusedStyle := lipgloss.NewStyle().
            Border(lipgloss.NormalBorder(), false, false, false, true).
            BorderForeground(styles.ColorAccent).
            PaddingLeft(1)
        
        normalStyle := lipgloss.NewStyle().
            PaddingLeft(2) // Match border(1) + padding(1) of focused style

        // Apply the appropriate style and print
        var finalOutput string
        if isFocused {
            finalOutput = focusedStyle.Render(content)
        } else {
            finalOutput = normalStyle.Render(content)
        }
        fmt.Fprint(w, finalOutput)

	} else {
		// Fallback to default delegate rendering for other item types (if any)
		d.DefaultDelegate.Render(w, m, index, listItem)
	}
}
// --- End Custom Delegate ---

// BaseSelector is the main model for selection
type BaseSelector struct {
	list           list.Model
	selectedItems  map[interface{}]struct{} // For multi-select
	currentItem    interface{} // For single-select result
	quitting       bool
	done           bool
	title          string
	singleSelectMode bool // New flag
}

// NewBaseSelector creates a new base selector
func NewBaseSelector(title string, singleSelect bool) *BaseSelector { // Added singleSelect param
	delegate := newItemDelegate(singleSelect) // Use custom delegate
	l := list.New([]list.Item{}, delegate, 0, 0)

    if title != "" {
        l.Title = title
        l.Styles.Title = styles.ListTitleStyle
        l.SetShowTitle(true)
    } else {
        l.SetShowTitle(false)
    }
    
	l.Styles.HelpStyle = styles.KeyMapStyle
	l.Styles.StatusBar = styles.KeyMapStyle.Copy()
	l.Styles.FilterPrompt = styles.InfoStyle.Copy()
	l.Styles.FilterCursor = styles.InfoStyle.Copy()
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)

	return &BaseSelector{
		list:           l,
		selectedItems:  make(map[interface{}]struct{}),
		title:          title,
		singleSelectMode: singleSelect, // Set the mode
	}
}

// Init implements tea.Model
func (s *BaseSelector) Init() tea.Cmd {
	// Don't enter alt screen here, let parent handle it
	return nil 
}

// Update handles keyboard input and updates the selector state
func (s *BaseSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width, msg.Height)
		return s, nil

	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering { break } // Let list handle filter keys

		switch msg.String() {
		case "ctrl+c", "q":
			s.quitting = true
			s.done = false
			return s, tea.Quit 
		case "enter":
            // Select current item and signal done
            if item, ok := s.list.SelectedItem().(*SelectorItem); ok {
                s.currentItem = item.item // Store the single selected item
            }
			s.done = true
			s.quitting = false
			return s, nil // Signal finished, parent handles transition
		case " ": 
            if s.singleSelectMode { break } // Ignore space in single select mode

            // Multi-select toggle logic
			currentItem, ok := s.list.SelectedItem().(*SelectorItem)
			if !ok { return s, nil }
			currentItem.selected = !currentItem.selected
			if currentItem.selected { s.selectedItems[currentItem.item] = struct{}{} } else { delete(s.selectedItems, currentItem.item) }
            // No SetItem needed, view update handles visual change
			return s, nil // Just update internal state
		}
	}

	// Pass unmatched messages down to the list for default handling (navigation, filtering)
	s.list, cmd = s.list.Update(msg)
	cmds = append(cmds, cmd)
	return s, tea.Batch(cmds...)
}

// View implements tea.Model
func (s *BaseSelector) View() string {
	if s.quitting { return styles.InfoStyle.Render("Selection cancelled.") }
	return s.list.View() // List handles rendering title, items, status, help
}

func (s *BaseSelector) Finished() bool {
	return s.done && !s.quitting
}

// GetSelected returns selected items
func (s *BaseSelector) GetSelected() []interface{} {
	if !s.Finished() { return nil }
    if s.singleSelectMode {
        if s.currentItem != nil {
             return []interface{}{s.currentItem} // Return single item in a slice
        }
        return nil
    }
    // Multi-select logic
	var result []interface{}
	for item := range s.selectedItems { result = append(result, item) }
	return result
}

// SetItems prepares SelectorItem for the list from a slice of actual data items
func (s *BaseSelector) SetItems(items []interface{}, titleFn func(interface{}) string, descFn func(interface{}) string) {
	listItems := make([]list.Item, len(items))
	for i, dataItem := range items {
		_, isSelected := s.selectedItems[dataItem] // Preserve selection if item already exists
		listItems[i] = &SelectorItem{
			title:       titleFn(dataItem),
			description: descFn(dataItem),
			item:        dataItem,
			selected:    isSelected,
		}
	}
	s.list.SetItems(listItems)
}

// SetSize sets the width and height of the selector - usually called on tea.WindowSizeMsg
func (s *BaseSelector) SetSize(width, height int) {
	s.list.SetSize(width, height)
}

// SetSelectedDataItems allows pre-selecting items by providing the actual data items
func (s *BaseSelector) SetSelectedDataItems(dataItemsToSelect []interface{}) {
	s.selectedItems = make(map[interface{}]struct{}) // Clear previous selections
	for _, dataItem := range dataItemsToSelect {
		s.selectedItems[dataItem] = struct{}{}
	}
	// Update the 'selected' field on the list.Item wrappers (SelectorItem)
	currentListItems := s.list.Items()
	for i, listItem := range currentListItems {
		if si, ok := listItem.(*SelectorItem); ok {
			if _, ok := s.selectedItems[si.item]; ok {
				si.selected = true
			} else {
				si.selected = false
			}
			_ = s.list.SetItem(i, si) // Update item in list for visual consistency
		}
	}
}

// RunSelector is a helper to start this TUI component and get selected items.
// It expects items as []interface{} and functions to get title/description.
func RunSelector(title string, items []interface{}, titleFn func(interface{}) string, descFn func(interface{}) string, preselected []interface{}) ([]interface{}, error) {
	selector := NewBaseSelector(title, false)
	selector.SetItems(items, titleFn, descFn)
	if len(preselected) > 0 {
		selector.SetSelectedDataItems(preselected)
	}

	p := tea.NewProgram(selector)
	finalModel, err := p.StartReturningModel()
	if err != nil {
		return nil, fmt.Errorf("error running selector: %w", err)
	}

	castedModel, ok := finalModel.(*BaseSelector)
	if !ok {
		return nil, fmt.Errorf("could not cast final model to BaseSelector")
	}

	if castedModel.Finished() {
		return castedModel.GetSelected(), nil
	}
	return nil, nil // Selection cancelled or quit
}

// Specific SetItems for common types (Tool, Font, Language) - REMOVED
// These are convenience wrappers around the generic SetItems.

// func (s *BaseSelector) SetToolItems(tools []*interfaces.Tool, preselectedTools []*interfaces.Tool) {
// 	genericItems := make([]interface{}, len(tools))
// 	for i, t := range tools { genericItems[i] = t }
// 	
// 	genericPreselected := make([]interface{}, len(preselectedTools))
// 	for i, t := range preselectedTools { genericPreselected[i] = t }
// 
// 	s.SetItems(genericItems, 
// 		func(item interface{}) string { return item.(*interfaces.Tool).Name }, 
// 		func(item interface{}) string { return item.(*interfaces.Tool).Description },
// 	)
// 	if len(genericPreselected) > 0 {
// 		s.SetSelectedDataItems(genericPreselected)
// 	}
// }

// Add similar SetFontItems, SetLanguageItems etc. as needed

