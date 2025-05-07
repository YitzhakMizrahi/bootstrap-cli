package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/pipeline"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// EssentialToolScreen uses the BaseSelector component for essential tool selection.
type EssentialToolScreen struct {
	selector *components.BaseSelector
	finished bool
	title    string // Keep title field even if not displayed by default
	width    int
	height   int
}

// NewEssentialToolScreen creates a new EssentialToolScreen.
func NewEssentialToolScreen(title string, tools []*pipeline.Tool, preselected []*pipeline.Tool) *EssentialToolScreen {
	selector := components.NewBaseSelector(title, false)
	
	// Convert tools and preselected to []interface{} for BaseSelector
	items := make([]interface{}, len(tools))
	for i, t := range tools { items[i] = t }
	selectedItems := make([]interface{}, len(preselected))
	for i, t := range preselected { selectedItems[i] = t }

	selector.SetItems(items, 
		func(item interface{}) string { 
			if t, ok := item.(*pipeline.Tool); ok { return t.Name }
			return ""
		}, 
		func(item interface{}) string { 
			if t, ok := item.(*pipeline.Tool); ok { return t.Description }
			return ""
		},
	)
	if len(selectedItems) > 0 {
		selector.SetSelectedDataItems(selectedItems)
	}

	s := &EssentialToolScreen{
		selector: selector,
		finished: false,
		title:    title,
	}
	return s
}

func (s *EssentialToolScreen) Init() tea.Cmd { 
    if s.selector != nil { return s.selector.Init() }
    return nil
}

func (s *EssentialToolScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s *EssentialToolScreen) View() string {
	if s.selector == nil { return styles.ErrorStyle.Render("Error: Essential Tool selector not initialized.") }
	return s.selector.View()
}

func (s *EssentialToolScreen) Finished() bool { return s.finished }

func (s *EssentialToolScreen) GetSelected() []*pipeline.Tool {
	if s.selector != nil && s.selector.Finished() {
		items := s.selector.GetSelected()
		tools := make([]*pipeline.Tool, 0, len(items))
		for _, item := range items {
			if tool, ok := item.(*pipeline.Tool); ok { tools = append(tools, tool) }
		}
		return tools
	}
	return nil
} 