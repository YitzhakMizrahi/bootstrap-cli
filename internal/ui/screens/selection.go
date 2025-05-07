package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectionScreen is a reusable selection UI component for lists of items.
type SelectionScreen struct {
	title        string
	selector     *components.BaseSelector
	itemRenderer func(interface{}) string
	descRenderer func(interface{}) string
	finished     bool
	above        []func() string
	below        []func() string
}

// NewSelectionScreen creates a new SelectionScreen for the given items and options.
func NewSelectionScreen(title string, items []interface{}, itemRenderer, descRenderer func(interface{}) string, opts ...SelectionScreenOption) *SelectionScreen {
	selector := components.NewBaseSelector("")
	selector.SetItems(items, itemRenderer, descRenderer)
	s := &SelectionScreen{
		title:        title,
		selector:     selector,
		itemRenderer: itemRenderer,
		descRenderer: descRenderer,
		finished:     false,
		above:        nil,
		below:        nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// SelectionScreenOption configures a SelectionScreen.
type SelectionScreenOption func(*SelectionScreen)

// WithAboveComponents adds a function to render components above the selection list.
func WithAboveComponents(components ...func() string) SelectionScreenOption {
	return func(s *SelectionScreen) {
		s.above = append(s.above, components...)
	}
}

// WithBelowComponents adds a function to render components below the selection list.
func WithBelowComponents(components ...func() string) SelectionScreenOption {
	return func(s *SelectionScreen) {
		s.below = append(s.below, components...)
	}
}

// WithSelectedItems pre-selects items in the SelectionScreen.
func WithSelectedItems(selectedItems []interface{}) SelectionScreenOption {
	return func(s *SelectionScreen) {
		s.selector.SetSelectedItems(selectedItems)
	}
}

// Init initializes the SelectionScreen.
func (s *SelectionScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages for the SelectionScreen.
func (s *SelectionScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	newModel, cmd := s.selector.Update(msg)
	if selector, ok := newModel.(*components.BaseSelector); ok {
		s.selector = selector
		if s.selector.Finished() {
			s.finished = true
		}
	}
	return s, cmd
}

// View renders the SelectionScreen UI.
func (s *SelectionScreen) View() string {
	parts := make([]string, 0, 4)
	for _, above := range s.above {
		parts = append(parts, above())
	}
	parts = append(parts, styles.TitleStyle.Render(s.title))
	parts = append(parts, s.selector.View())
	for _, below := range s.below {
		parts = append(parts, below())
	}
	return styles.JoinVertical(parts...)
}

// Finished returns true if the selection is complete.
func (s *SelectionScreen) Finished() bool {
	return s.finished
}

// GetSelected returns the selected items from the SelectionScreen.
func (s *SelectionScreen) GetSelected() map[int]interface{} {
	return s.selector.GetSelected()
}

// Selector returns the underlying selector component.
func (s *SelectionScreen) Selector() *components.BaseSelector {
	return s.selector
} 