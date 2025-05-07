package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// ModernToolScreen uses the BaseSelector component for modern tool selection.
type ModernToolScreen struct {
	selector *components.BaseSelector
	finished bool
	title    string 
	width    int
	height   int
}

// NewModernToolScreen creates a new ModernToolScreen.
func NewModernToolScreen(title string, tools []*interfaces.Tool, preselected []*interfaces.Tool) *ModernToolScreen {
	selector := components.NewBaseSelector(title, false)
	
	// Convert tools and preselected to []interface{} for BaseSelector
	items := make([]interface{}, len(tools))
	for i, t := range tools { items[i] = t }
	selectedItems := make([]interface{}, len(preselected))
	for i, t := range preselected { selectedItems[i] = t }

	selector.SetItems(items, 
		func(item interface{}) string { if t, ok := item.(*interfaces.Tool); ok { return t.Name }; return "" }, 
		func(item interface{}) string { if t, ok := item.(*interfaces.Tool); ok { return t.Description }; return "" },
	)
	if len(selectedItems) > 0 {
		selector.SetSelectedDataItems(selectedItems)
	}

	s := &ModernToolScreen{
		selector: selector,
		finished: false,
		title:    title,
	}
	return s
}

func (s *ModernToolScreen) Init() tea.Cmd { 
    if s.selector != nil { return s.selector.Init() }
    return nil
}

func (s *ModernToolScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if s.selector.Finished() { s.finished = true }
			}
			cmds = append(cmds, newSelCmd)
		}
	}
	return s, tea.Batch(cmds...)
}

func (s *ModernToolScreen) View() string {
	if s.selector == nil { return styles.ErrorStyle.Render("Error: Modern Tool selector not initialized.") }
	return s.selector.View()
}

func (s *ModernToolScreen) Finished() bool { return s.finished }

func (s *ModernToolScreen) GetSelected() []*interfaces.Tool {
	if s.selector != nil && s.selector.Finished() {
		items := s.selector.GetSelected()
		tools := make([]*interfaces.Tool, 0, len(items))
		for _, item := range items {
			if tool, ok := item.(*interfaces.Tool); ok { tools = append(tools, tool) }
		}
		return tools
	}
	return nil
} 