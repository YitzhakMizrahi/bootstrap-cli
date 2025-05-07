package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// LanguageScreen uses the BaseSelector component for language selection.
type LanguageScreen struct {
	selector *components.BaseSelector
	finished bool
	title    string
	width    int
	height   int
}

// NewLanguageScreen creates a new LanguageScreen.
func NewLanguageScreen(title string, languages []*interfaces.Language, preselected []*interfaces.Language) *LanguageScreen {
	selector := components.NewBaseSelector(title, false)
	
	items := make([]interface{}, len(languages))
	for i, l := range languages { items[i] = l }
	selectedItems := make([]interface{}, len(preselected))
	for i, l := range preselected { selectedItems[i] = l }

	selector.SetItems(items, 
		func(item interface{}) string { if l, ok := item.(*interfaces.Language); ok { return l.Name }; return "" }, 
		func(item interface{}) string { if l, ok := item.(*interfaces.Language); ok { return l.Description }; return "" },
	)
	if len(selectedItems) > 0 {
		selector.SetSelectedDataItems(selectedItems)
	}

	s := &LanguageScreen{
		selector: selector,
		finished: false,
		title:    title,
	}
	return s
}

func (s *LanguageScreen) Init() tea.Cmd { 
    if s.selector != nil { return s.selector.Init() } 
    return nil
}

func (s *LanguageScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s *LanguageScreen) View() string {
    // ... (View logic remains the same) ...
	if s.selector == nil { return styles.ErrorStyle.Render("Error: Language selector not initialized.") }
	return s.selector.View()
}

func (s *LanguageScreen) Finished() bool { return s.finished }

// GetSelected specifically returns selected languages.
func (s *LanguageScreen) GetSelected() []*interfaces.Language {
    // ... (GetSelected logic remains the same) ...
	if s.selector != nil && s.selector.Finished() {
		items := s.selector.GetSelected() 
		langs := make([]*interfaces.Language, 0, len(items))
		for _, item := range items {
			if lang, ok := item.(*interfaces.Language); ok { langs = append(langs, lang) }
		}
		return langs
	}
	return nil
} 