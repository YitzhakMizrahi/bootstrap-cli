package screens

import (
	"fmt"
	"sort"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
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
	return &ToolScreen{
		tools:    tools,
		finished: false,
		canceled: false,
	}
}

// ShowToolSelection displays the tool selection screen
func (s *ToolScreen) ShowToolSelection() ([]*interfaces.Tool, error) {
	// Sort tools by category
	toolsByCategory := make(map[string][]*interfaces.Tool)
	for _, tool := range s.tools {
		category := tool.Category
		toolsByCategory[category] = append(toolsByCategory[category], tool)
	}

	// Get sorted categories
	categories := make([]string, 0, len(toolsByCategory))
	for category := range toolsByCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	// Create selector
	s.selector = components.NewBaseSelector("Select Development Tools")

	// Set up the selector with tools directly
	s.selector.SetItems(
		s.tools, // Pass tools directly, no need for conversion
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

	// Run the selector
	p := tea.NewProgram(s.selector)
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("tool selection failed: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*components.BaseSelector); ok {
		if !selectorModel.Finished() {
			s.canceled = true
			return nil, nil
		}

		// Convert selected items to tools
		selected := selectorModel.GetSelected()
		s.selected = make([]*interfaces.Tool, 0, len(selected))
		for _, item := range selected {
			if tool, ok := item.(*interfaces.Tool); ok {
				s.selected = append(s.selected, tool)
			}
		}

		s.finished = true
		return s.selected, nil
	}

	return nil, fmt.Errorf("failed to get tool selector model")
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