// Package screens provides UI screens for the bootstrap-cli application.
package screens

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/components"
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
	return &FontScreen{
		selector: components.NewBaseSelector("Select Fonts"),
		loader:   loader,
	}
}

// ShowFontSelection prompts the user to select fonts
func (s *FontScreen) ShowFontSelection() ([]*interfaces.Font, error) {
	// Load available fonts
	availableFonts, err := s.loader.LoadFonts()
	if err != nil {
		return nil, fmt.Errorf("failed to load font configurations: %w", err)
	}

	// Set up the selector with fonts
	s.selector.SetItems(
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

	// Run the interactive UI
	p := tea.NewProgram(s.selector)
	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("UI error: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*components.BaseSelector); ok {
		if !selectorModel.Finished() {
			return nil, fmt.Errorf("selection cancelled")
		}
		
		// Convert selected items to fonts
		selected := selectorModel.GetSelected()
		fonts := make([]*interfaces.Font, 0, len(selected))
		for _, item := range selected {
			if font, ok := item.(*interfaces.Font); ok {
				fonts = append(fonts, font)
			}
		}
		return fonts, nil
	}

	return nil, fmt.Errorf("failed to get font selector model")
} 