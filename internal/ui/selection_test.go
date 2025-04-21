package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewSelectionComponent(t *testing.T) {
	items := []SelectionItem{
		{title: "Item 1", description: "Description 1", value: 1},
		{title: "Item 2", description: "Description 2", value: 2},
	}

	// Test single select mode
	singleSelect := NewSelectionComponent("Test Single Select", items, false)
	if singleSelect == nil {
		t.Error("Expected non-nil SelectionComponent for single select")
	}
	if singleSelect.multiSelect {
		t.Error("Expected multiSelect to be false for single select")
	}
	if singleSelect.title != "Test Single Select" {
		t.Errorf("Expected title 'Test Single Select', got '%s'", singleSelect.title)
	}

	// Test multi select mode
	multiSelect := NewSelectionComponent("Test Multi Select", items, true)
	if multiSelect == nil {
		t.Error("Expected non-nil SelectionComponent for multi select")
	}
	if !multiSelect.multiSelect {
		t.Error("Expected multiSelect to be true for multi select")
	}
}

func TestSelectionComponentUpdate(t *testing.T) {
	items := []SelectionItem{
		{title: "Item 1", description: "Description 1", value: 1},
		{title: "Item 2", description: "Description 2", value: 2},
	}
	comp := NewSelectionComponent("Test Selection", items, true)

	// Test window size update
	model, _ := comp.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	if s, ok := model.(*SelectionComponent); ok {
		if s.width != 100 || s.height != 50 {
			t.Errorf("Expected width=100, height=50, got width=%d, height=%d", s.width, s.height)
		}
	} else {
		t.Error("Expected *SelectionComponent from Update")
	}

	// Test item selection in multi-select mode
	model, _ = comp.Update(tea.KeyMsg{Type: tea.KeySpace})
	if s, ok := model.(*SelectionComponent); ok {
		selected := s.GetSelected()
		if len(selected) != 1 {
			t.Errorf("Expected 1 selected item, got %d", len(selected))
		}
	}
}

func TestGetSelected(t *testing.T) {
	items := []SelectionItem{
		{title: "Item 1", description: "Description 1", value: 1},
		{title: "Item 2", description: "Description 2", value: 2},
	}

	// Test single select mode
	singleSelect := NewSelectionComponent("Test Single Select", items, false)
	selected := singleSelect.GetSelected()
	if len(selected) > 1 {
		t.Error("Single select mode should not return multiple items")
	}

	// Test multi select mode
	multiSelect := NewSelectionComponent("Test Multi Select", items, true)
	multiSelect.selectedItems = []SelectionItem{items[0], items[1]}
	selected = multiSelect.GetSelected()
	if len(selected) != 2 {
		t.Errorf("Expected 2 selected items in multi-select mode, got %d", len(selected))
	}
}

func TestView(t *testing.T) {
	items := []SelectionItem{
		{title: "Item 1", description: "Description 1", value: 1},
		{title: "Item 2", description: "Description 2", value: 2},
	}

	// Test single select mode view
	singleSelect := NewSelectionComponent("Test Single Select", items, false)
	view := singleSelect.View()
	if view == "" {
		t.Error("View should not return empty string")
	}

	// Test multi select mode view
	multiSelect := NewSelectionComponent("Test Multi Select", items, true)
	multiSelect.selectedItems = []SelectionItem{items[0]}
	view = multiSelect.View()
	if view == "" {
		t.Error("View should not return empty string")
	}
} 