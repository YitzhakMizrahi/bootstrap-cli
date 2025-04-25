// Package components provides UI components for the bootstrap-cli application.
package components

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	"github.com/manifoldco/promptui"
)

// BasicPrompt represents a simple yes/no or selection prompt
type BasicPrompt struct {
	label string
	items []string
}

// NewBasicPrompt creates a new basic prompt
func NewBasicPrompt(label string, items []string) *BasicPrompt {
	return &BasicPrompt{
		label: label,
		items: items,
	}
}

// Run executes the prompt and returns the selected item
func (p *BasicPrompt) Run() (string, error) {
	prompt := promptui.Select{
		Label: styles.InfoStyle.Render(p.label),
		Items: p.items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | cyan }}",
			Active:   "➤ {{ . | cyan }}",
			Inactive: "  {{ . | white }}",
			Selected: "{{ . | green }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	return result, nil
}

// RunYesNo executes a yes/no prompt and returns a boolean
func (p *BasicPrompt) RunYesNo() (bool, error) {
	result, err := p.Run()
	if err != nil {
		return false, err
	}
	return result == "Yes", nil
}

// RunWithInput executes a prompt that requires text input
func (p *BasicPrompt) RunWithInput() (string, error) {
	prompt := promptui.Prompt{
		Label: styles.InfoStyle.Render(p.label),
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	return result, nil
} 