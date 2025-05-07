// Adapted from tools.go
package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// FontScreen uses the BaseSelector component for font selection.
type FontScreen struct {
	selector *components.BaseSelector
	finished bool
	title    string
	width    int
	height   int
}

// NewFontScreen creates a new FontScreen.
func NewFontScreen(title string, fonts []*interfaces.Font, preselected []*interfaces.Font) *FontScreen {
	selector := components.NewBaseSelector(title)
	
	// Convert fonts and preselected to []interface{} for BaseSelector
	items := make([]interface{}, len(fonts))
	for i, f := range fonts { items[i] = f }
	selectedItems := make([]interface{}, len(preselected))
	for i, f := range preselected { selectedItems[i] = f }

	selector.SetItems(items, 
		func(item interface{}) string { if f, ok := item.(*interfaces.Font); ok { return f.Name }; return "" }, 
		func(item interface{}) string { if f, ok := item.(*interfaces.Font); ok { return f.Description }; return "" },
	)
	if len(selectedItems) > 0 {
		selector.SetSelectedDataItems(selectedItems)
	}

	s := &FontScreen{
		selector: selector,
		finished: false,
		title:    title,
	}
	return s
}

func (s *FontScreen) Init() tea.Cmd { 
    if s.selector != nil {
        return s.selector.Init() 
    }
    return nil
}

func (s *FontScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		s.width = msg.Width
		s.height = msg.Height
		if s.selector != nil {
			newSelModel, newSelCmd := s.selector.Update(msg)
			if sel, ok := newSelModel.(*components.BaseSelector); ok { s.selector = sel }
			cmds = append(cmds, newSelCmd)
		}
		return s, tea.Batch(cmds...)
	default: 
		if s.selector != nil {
			newSelModel, newSelCmd := s.selector.Update(msg)
			if sel, ok := newSelModel.(*components.BaseSelector); ok {
				s.selector = sel
				if s.selector.Finished() {
					s.finished = true
				}
			}
			cmds = append(cmds, newSelCmd)
		}
	}
	return s, tea.Batch(cmds...)
}

func (s *FontScreen) View() string {
	if s.selector == nil { return styles.ErrorStyle.Render("Error: Font selector not initialized.") }
	return s.selector.View()
}

func (s *FontScreen) Finished() bool { return s.finished }

func (s *FontScreen) GetSelected() []*interfaces.Font {
	if s.selector != nil && s.selector.Finished() {
		items := s.selector.GetSelected() 
		fonts := make([]*interfaces.Font, 0, len(items))
		for _, item := range items {
			if font, ok := item.(*interfaces.Font); ok { fonts = append(fonts, font) }
		}
		return fonts
	}
	return nil
} 