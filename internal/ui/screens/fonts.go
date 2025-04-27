// Package screens provides UI screens for the bootstrap-cli application.
package screens

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// FontScreen handles font selection
type FontScreen struct {
	selector *components.BaseSelector
	loader   *config.Loader
}

// NewFontScreen creates a new font selection screen
func NewFontScreen(loader *config.Loader) *FontScreen {
	availableFonts, _ := loader.LoadFonts()
	selector := components.NewBaseSelector("Select Fonts")
	selector.SetItems(
		utils.ConvertToInterfaceSlice(availableFonts),
		func(item interface{}) string {
			if font, ok := item.(*interfaces.Font); ok {
				return font.Name
			}
			return ""
		},
		func(item interface{}) string {
			if font, ok := item.(*interfaces.Font); ok {
				return font.Description
			}
			return ""
		},
	)
	return &FontScreen{
		selector: selector,
		loader:   loader,
	}
}

// Init implements tea.Model
func (s *FontScreen) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (s *FontScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	newModel, cmd := s.selector.Update(msg)
	if selector, ok := newModel.(*components.BaseSelector); ok {
		s.selector = selector
	}
	return s, cmd
}

// View implements tea.Model
func (s *FontScreen) View() string {
	return styles.TitleStyle.Render("Select Fonts") + "\n\n" + s.selector.View()
}

// Finished returns true if the screen was completed
func (s *FontScreen) Finished() bool {
	return s.selector.Finished()
}

// GetSelected returns the selected fonts
func (s *FontScreen) GetSelected() []*interfaces.Font {
	selected := s.selector.GetSelected()
	fonts := make([]*interfaces.Font, 0, len(selected))
	for _, item := range selected {
		if font, ok := item.(*interfaces.Font); ok {
			fonts = append(fonts, font)
		}
	}
	return fonts
} 