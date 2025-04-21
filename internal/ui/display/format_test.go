package display

import (
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	f := NewFormatter()
	if f == nil {
		t.Error("Expected non-nil formatter")
	}
}

func TestFormatterWithColor(t *testing.T) {
	f := NewFormatter()
	colored := f.WithColor(ColorRed, "test")
	if !strings.Contains(colored, ColorRed) {
		t.Errorf("Expected color %s in result", ColorRed)
	}
}

func TestFormatterWithStyle(t *testing.T) {
	f := NewFormatter()
	styled := f.WithStyle(StyleBold, "test")
	if !strings.Contains(styled, StyleBold) {
		t.Errorf("Expected style %s in result", StyleBold)
	}
}

func TestFormatterWithWidth(t *testing.T) {
	f := NewFormatter()
	text := "test"
	width := 10
	padded := f.WithWidth(width, text)
	if len(padded) != width {
		t.Errorf("Expected width %d, got %d", width, len(padded))
	}
}

func TestFormatterWithAlignment(t *testing.T) {
	f := NewFormatter()
	text := "test"
	width := 10
	aligned := f.WithAlignment(AlignCenter, width, text)
	if len(aligned) != width {
		t.Errorf("Expected width %d, got %d", width, len(aligned))
	}
}

func TestFormatterFormat(t *testing.T) {
	f := NewFormatter()
	text := "test"
	formatted := f.Format(text)
	if formatted != text {
		t.Errorf("Expected %s, got %s", text, formatted)
	}
}

func TestFormatterBox(t *testing.T) {
	f := NewFormatter()
	text := "test"
	box := f.Box(text)
	if !strings.Contains(box, text) {
		t.Errorf("Expected box to contain text %s", text)
	}
}

func TestFormatterHeader(t *testing.T) {
	f := NewFormatter()
	text := "test"
	header := f.Header(text)
	if !strings.Contains(header, text) {
		t.Errorf("Expected header to contain text %s", text)
	}
}

func TestFormatterList(t *testing.T) {
	f := NewFormatter()
	items := []string{"item1", "item2"}
	list := f.List(items)
	for _, item := range items {
		if !strings.Contains(list, item) {
			t.Errorf("Expected list to contain item %s", item)
		}
	}
}

func TestFormatterTable(t *testing.T) {
	f := NewFormatter()
	headers := []string{"header1", "header2"}
	rows := [][]string{{"row1col1", "row1col2"}, {"row2col1", "row2col2"}}
	table := f.Table(headers, rows)
	for _, header := range headers {
		if !strings.Contains(table, header) {
			t.Errorf("Expected table to contain header %s", header)
		}
	}
	for _, row := range rows {
		for _, cell := range row {
			if !strings.Contains(table, cell) {
				t.Errorf("Expected table to contain cell %s", cell)
			}
		}
	}
}

func TestFormatterStatusMessages(t *testing.T) {
	f := NewFormatter()
	tests := []struct {
		name    string
		method  func(string) string
		message string
		color   string
		prefix  string
	}{
		{
			name:    "success message",
			method:  f.Success,
			message: "Operation completed",
			color:   ColorGreen,
			prefix:  "✓",
		},
		{
			name:    "error message",
			method:  f.Error,
			message: "Operation failed",
			color:   ColorRed,
			prefix:  "✗",
		},
		{
			name:    "warning message",
			method:  f.Warning,
			message: "Proceed with caution",
			color:   ColorYellow,
			prefix:  "⚠",
		},
		{
			name:    "info message",
			method:  f.Info,
			message: "System status",
			color:   ColorBlue,
			prefix:  "ℹ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method(tt.message)
			if !strings.Contains(result, tt.color) {
				t.Errorf("Expected color %s in result", tt.color)
			}
			if !strings.Contains(result, tt.prefix) {
				t.Errorf("Expected prefix %s in result", tt.prefix)
			}
			if !strings.Contains(result, tt.message) {
				t.Errorf("Expected message %s in result", tt.message)
			}
		})
	}
}