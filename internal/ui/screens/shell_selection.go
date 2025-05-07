package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// ShellSelectionScreen uses the BaseSelector component for shell selection.
type ShellSelectionScreen struct {
	selector      *components.BaseSelector
	finished      bool
	title         string
	width         int
	height        int
    currentShell  string // Store the path of the system's current default shell
}

// NewShellSelectionScreen creates a new ShellSelectionScreen.
func NewShellSelectionScreen(title string, availableShells []*interfaces.Shell, currentSystemShell string, preselectedShellName string) *ShellSelectionScreen {
	selector := components.NewBaseSelector(title, true)
	
	items := make([]interface{}, len(availableShells))
	for i, s := range availableShells { items[i] = s }
	
	var selectedItemsInitial []interface{}
	// Preselect based on name if provided
	for _, s := range availableShells {
		if s.Name == preselectedShellName {
			selectedItemsInitial = append(selectedItemsInitial, s)
			break
		}
	}

	selector.SetItems(items, 
		func(item interface{}) string { // Title function
			if s, ok := item.(*interfaces.Shell); ok {
                // Indicate if it's the current system default
                if s.Path == currentSystemShell {
                    return s.Name + " (current default)"
                }
				return s.Name 
			}
			return ""
		}, 
		func(item interface{}) string { // Description function
			if s, ok := item.(*interfaces.Shell); ok { return s.Description }
			return ""
		},
	)
	if len(selectedItemsInitial) > 0 {
		selector.SetSelectedDataItems(selectedItemsInitial)
	}

	s := &ShellSelectionScreen{
		selector:     selector,
		finished:     false,
		title:        title,
        currentShell: currentSystemShell,
	}
	return s
}

func (s *ShellSelectionScreen) Init() tea.Cmd { 
    if s.selector != nil { return s.selector.Init() } 
    return nil
}

func (s *ShellSelectionScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s *ShellSelectionScreen) View() string {
	if s.selector == nil { return styles.ErrorStyle.Render("Error: Shell selector not initialized.") }
	return s.selector.View()
}

func (s *ShellSelectionScreen) Finished() bool { return s.finished }

// GetSelected returns the selected *interfaces.Shell object (first selected item).
func (s *ShellSelectionScreen) GetSelected() *interfaces.Shell {
	if s.selector != nil && s.selector.Finished() {
		items := s.selector.GetSelected() // Returns []interface{}
		if len(items) > 0 {
			if shell, ok := items[0].(*interfaces.Shell); ok { return shell }
		}
	}
	return nil // Default or indicate no selection
} 