package screens

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// ToolScreen represents the tool selection screen
type ToolScreen struct {
	tools     []*interfaces.Tool
	selector  *components.BaseSelector
	selected  []*interfaces.Tool
	finished  bool
	canceled  bool
}

// NewToolScreen creates a new tool selection screen
func NewToolScreen(tools []*interfaces.Tool) *ToolScreen {
	ts := &ToolScreen{
		tools:    tools,
		finished: false,
		canceled: false,
	}
	// Create selector
	ts.selector = components.NewBaseSelector("Select Development Tools")
	ts.selector.SetItems(
		tools,
		func(item interface{}) string {
			if tool, ok := item.(*interfaces.Tool); ok {
				if tool.Category != "" {
					return fmt.Sprintf("[%s] %s", tool.Category, tool.Name)
				}
				return tool.Name
			}
			return ""
		},
		func(item interface{}) string {
			if tool, ok := item.(*interfaces.Tool); ok {
				return tool.Description
			}
			return ""
		},
	)
	return ts
}

// Init implements tea.Model
func (s *ToolScreen) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (s *ToolScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	newModel, cmd := s.selector.Update(msg)
	if selector, ok := newModel.(*components.BaseSelector); ok {
		s.selector = selector
		if s.selector.Finished() {
			selected := s.selector.GetSelected()
			s.selected = make([]*interfaces.Tool, 0, len(selected))
			for _, item := range selected {
				if tool, ok := item.(*interfaces.Tool); ok {
					s.selected = append(s.selected, tool)
				}
			}
			s.finished = true
		}
	}
	return s, cmd
}

// View implements tea.Model
func (s *ToolScreen) View() string {
	return styles.TitleStyle.Render("Select Development Tools") + "\n\n" + s.selector.View()
}

// Finished returns true if the screen was completed
func (s *ToolScreen) Finished() bool {
	return s.finished
}

// GetSelected returns the selected tools
func (s *ToolScreen) GetSelected() []*interfaces.Tool {
	return s.selected
}

// IsFinished returns true if the screen was completed
func (s *ToolScreen) IsFinished() bool {
	return s.finished
}

// IsCanceled returns true if the screen was canceled
func (s *ToolScreen) IsCanceled() bool {
	return s.canceled
}

// GetSelectedTools returns the selected tools
func (s *ToolScreen) GetSelectedTools() []*interfaces.Tool {
	return s.selected
} 