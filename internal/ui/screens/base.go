// Package screens provides the UI screens for the bootstrap-cli TUI
package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

// BaseScreen defines the interface that all screens must implement
type BaseScreen interface {
	tea.Model
	Done() bool
}

// BaseScreenModel provides common functionality for all screens
type BaseScreenModel struct {
	done bool
}

// Done returns whether the screen is finished and ready to move to the next screen
func (b *BaseScreenModel) Done() bool {
	return b.done
}

// SetDone marks the screen as done
func (b *BaseScreenModel) SetDone(done bool) {
	b.done = done
} 