package display

import (
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	f := NewFormatter()
	if f == nil {
		t.Error("NewFormatter() returned nil")
	}
	if f.Color != ColorReset {
		t.Errorf("Expected default color to be %s, got %s", ColorReset, f.Color)
	}
	if f.Style != "" {
		t.Errorf("Expected default style to be empty, got %s", f.Style)
	}
	if f.Width != 80 {
		t.Errorf("Expected default width to be 80, got %d", f.Width)
	}
	if f.Alignment != "left" {
		t.Errorf("Expected default alignment to be 'left', got %s", f.Alignment)
	}
}

func TestFormatterWithColor(t *testing.T) {
	f := NewFormatter().WithColor(ColorRed)
	if f.Color != ColorRed {
		t.Errorf("Expected color to be %s, got %s", ColorRed, f.Color)
	}
}

func TestFormatterWithStyle(t *testing.T) {
	f := NewFormatter().WithStyle(StyleBold)
	if f.Style != StyleBold {
		t.Errorf("Expected style to be %s, got %s", StyleBold, f.Style)
	}
}

func TestFormatterWithWidth(t *testing.T) {
	f := NewFormatter().WithWidth(100)
	if f.Width != 100 {
		t.Errorf("Expected width to be 100, got %d", f.Width)
	}
}

func TestFormatterWithAlignment(t *testing.T) {
	f := NewFormatter().WithAlignment("center")
	if f.Alignment != "center" {
		t.Errorf("Expected alignment to be 'center', got %s", f.Alignment)
	}
}

func TestFormatterFormat(t *testing.T) {
	f := NewFormatter().WithColor(ColorRed).WithStyle(StyleBold)
	text := "Hello, World!"
	result := f.Format(text)
	
	if !strings.Contains(result, text) {
		t.Errorf("Expected formatted text to contain '%s', got '%s'", text, result)
	}
	if !strings.Contains(result, ColorRed) {
		t.Errorf("Expected formatted text to contain color code %s", ColorRed)
	}
	if !strings.Contains(result, StyleBold) {
		t.Errorf("Expected formatted text to contain style code %s", StyleBold)
	}
	if !strings.Contains(result, ColorReset) {
		t.Errorf("Expected formatted text to contain reset code %s", ColorReset)
	}
}

func TestFormatterBox(t *testing.T) {
	f := NewFormatter()
	text := "Hello, World!"
	result := f.Box(text)
	
	if !strings.Contains(result, BoxTopLeft) {
		t.Errorf("Expected box to contain top-left corner %s", BoxTopLeft)
	}
	if !strings.Contains(result, BoxTopRight) {
		t.Errorf("Expected box to contain top-right corner %s", BoxTopRight)
	}
	if !strings.Contains(result, BoxBottomLeft) {
		t.Errorf("Expected box to contain bottom-left corner %s", BoxBottomLeft)
	}
	if !strings.Contains(result, BoxBottomRight) {
		t.Errorf("Expected box to contain bottom-right corner %s", BoxBottomRight)
	}
	if !strings.Contains(result, text) {
		t.Errorf("Expected box to contain text '%s'", text)
	}
}

func TestFormatterHeader(t *testing.T) {
	f := NewFormatter()
	text := "Test Header"
	
	// Test level 1 header
	result := f.Header(text, 1)
	if !strings.Contains(result, strings.ToUpper(text)) {
		t.Errorf("Expected level 1 header to contain uppercase text '%s'", strings.ToUpper(text))
	}
	if !strings.Contains(result, strings.Repeat("=", len(text))) {
		t.Errorf("Expected level 1 header to contain underline")
	}
	
	// Test level 2 header
	result = f.Header(text, 2)
	if !strings.Contains(result, text) {
		t.Errorf("Expected level 2 header to contain text '%s'", text)
	}
	if !strings.Contains(result, strings.Repeat("-", len(text))) {
		t.Errorf("Expected level 2 header to contain underline")
	}
	
	// Test default header
	result = f.Header(text, 3)
	if !strings.Contains(result, text) {
		t.Errorf("Expected default header to contain text '%s'", text)
	}
}

func TestFormatterList(t *testing.T) {
	f := NewFormatter()
	items := []string{"Item 1", "Item 2", "Item 3"}
	
	// Test with default bullet
	result := f.List(items, "")
	for _, item := range items {
		if !strings.Contains(result, item) {
			t.Errorf("Expected list to contain item '%s'", item)
		}
	}
	
	// Test with custom bullet
	customBullet := "*"
	result = f.List(items, customBullet)
	for _, item := range items {
		if !strings.Contains(result, customBullet+" "+item) {
			t.Errorf("Expected list to contain '%s %s'", customBullet, item)
		}
	}
}

func TestFormatterTable(t *testing.T) {
	f := NewFormatter()
	headers := []string{"Header 1", "Header 2"}
	rows := [][]string{
		{"Cell 1", "Cell 2"},
		{"Cell 3", "Cell 4"},
	}
	
	result := f.Table(headers, rows)
	
	// Check headers
	for _, header := range headers {
		if !strings.Contains(result, header) {
			t.Errorf("Expected table to contain header '%s'", header)
		}
	}
	
	// Check cells
	for _, row := range rows {
		for _, cell := range row {
			if !strings.Contains(result, cell) {
				t.Errorf("Expected table to contain cell '%s'", cell)
			}
		}
	}
	
	// Test empty table
	emptyResult := f.Table([]string{}, [][]string{})
	if emptyResult != "" {
		t.Error("Expected empty table to return empty string")
	}
}

func TestFormatterStatusMessages(t *testing.T) {
	f := NewFormatter()
	message := "Test message"
	
	// Test success message
	result := f.Success(message)
	if !strings.Contains(result, "✓") {
		t.Error("Expected success message to contain checkmark")
	}
	if !strings.Contains(result, message) {
		t.Errorf("Expected success message to contain '%s'", message)
	}
	if !strings.Contains(result, ColorGreen) {
		t.Error("Expected success message to be green")
	}
	
	// Test error message
	result = f.Error(message)
	if !strings.Contains(result, "✗") {
		t.Error("Expected error message to contain x-mark")
	}
	if !strings.Contains(result, message) {
		t.Errorf("Expected error message to contain '%s'", message)
	}
	if !strings.Contains(result, ColorRed) {
		t.Error("Expected error message to be red")
	}
	
	// Test warning message
	result = f.Warning(message)
	if !strings.Contains(result, "⚠") {
		t.Error("Expected warning message to contain warning symbol")
	}
	if !strings.Contains(result, message) {
		t.Errorf("Expected warning message to contain '%s'", message)
	}
	if !strings.Contains(result, ColorYellow) {
		t.Error("Expected warning message to be yellow")
	}
	
	// Test info message
	result = f.Info(message)
	if !strings.Contains(result, "ℹ") {
		t.Error("Expected info message to contain info symbol")
	}
	if !strings.Contains(result, message) {
		t.Errorf("Expected info message to contain '%s'", message)
	}
	if !strings.Contains(result, ColorBlue) {
		t.Error("Expected info message to be blue")
	}
}

func TestFormatterProgressBar(t *testing.T) {
	f := NewFormatter()
	
	// Test normal progress
	result := f.ProgressBar(50, 20)
	if !strings.Contains(result, "50%") {
		t.Error("Expected progress bar to show 50%")
	}
	if !strings.Contains(result, ProgressFilled) {
		t.Error("Expected progress bar to contain filled blocks")
	}
	if !strings.Contains(result, ProgressEmpty) {
		t.Error("Expected progress bar to contain empty blocks")
	}
	
	// Test minimum progress
	result = f.ProgressBar(-10, 20)
	if !strings.Contains(result, "0%") {
		t.Error("Expected negative progress to show 0%")
	}
	
	// Test maximum progress
	result = f.ProgressBar(150, 20)
	if !strings.Contains(result, "100%") {
		t.Error("Expected progress over 100% to show 100%")
	}
	
	// Test zero width
	result = f.ProgressBar(50, 0)
	if !strings.Contains(result, "50%") {
		t.Error("Expected progress bar with zero width to show percentage")
	}
}

func TestFormatterIndent(t *testing.T) {
	f := NewFormatter()
	text := "Line 1\nLine 2\nLine 3"
	
	// Test with indent level 1
	result := f.Indent(text, 1)
	expectedIndent := "  "
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, expectedIndent) {
			t.Errorf("Expected line to be indented with '%s', got '%s'", expectedIndent, line)
		}
	}
	
	// Test with zero indent level
	result = f.Indent(text, 0)
	if result != text {
		t.Error("Expected zero indent to return original text")
	}
	
	// Test with negative indent level
	result = f.Indent(text, -1)
	if result != text {
		t.Error("Expected negative indent to return original text")
	}
}

func TestFormatterCodeBlock(t *testing.T) {
	f := NewFormatter()
	code := "func main() {\n    fmt.Println('Hello, World!')\n}"
	language := "go"
	
	result := f.CodeBlock(code, language)
	
	if !strings.Contains(result, "```"+language) {
		t.Errorf("Expected code block to contain language marker ```%s", language)
	}
	if !strings.Contains(result, code) {
		t.Error("Expected code block to contain the code")
	}
	if !strings.Contains(result, "```") {
		t.Error("Expected code block to contain closing markers")
	}
}

func TestFormatterCollapsible(t *testing.T) {
	f := NewFormatter()
	title := "Section Title"
	content := "Section Content"
	
	// Test expanded section
	result := f.Collapsible(title, content, true)
	if !strings.Contains(result, "▼") {
		t.Error("Expected expanded section to show down arrow")
	}
	if !strings.Contains(result, title) {
		t.Errorf("Expected collapsible to contain title '%s'", title)
	}
	if !strings.Contains(result, content) {
		t.Errorf("Expected expanded section to contain content '%s'", content)
	}
	
	// Test collapsed section
	result = f.Collapsible(title, content, false)
	if !strings.Contains(result, "▶") {
		t.Error("Expected collapsed section to show right arrow")
	}
	if !strings.Contains(result, title) {
		t.Errorf("Expected collapsible to contain title '%s'", title)
	}
	if strings.Contains(result, content) {
		t.Error("Expected collapsed section to not contain content")
	}
}