package progress

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func captureOutput(f func()) string {
	// Redirect stdout to capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function that produces output
	f()

	// Restore stdout and get the output
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewProgressBar(t *testing.T) {
	bar := NewProgressBar(100)

	if bar.Total != 100 {
		t.Errorf("Expected Total to be 100, got %d", bar.Total)
	}
	if bar.Width != 30 {
		t.Errorf("Expected Width to be 30, got %d", bar.Width)
	}
	if bar.Style.Prefix != "[" {
		t.Errorf("Expected Style.Prefix to be '[', got %s", bar.Style.Prefix)
	}
	if bar.Style.Suffix != "]" {
		t.Errorf("Expected Style.Suffix to be ']', got %s", bar.Style.Suffix)
	}
	if bar.Style.FillChar != "=" {
		t.Errorf("Expected Style.FillChar to be '=', got %s", bar.Style.FillChar)
	}
	if bar.Style.EmptyChar != " " {
		t.Errorf("Expected Style.EmptyChar to be ' ', got %s", bar.Style.EmptyChar)
	}
	if bar.Style.Color != ColorBlue {
		t.Errorf("Expected Style.Color to be ColorBlue, got %s", bar.Style.Color)
	}
	if !bar.ShowPercent {
		t.Error("Expected ShowPercent to be true")
	}
	if !bar.ShowTime {
		t.Error("Expected ShowTime to be true")
	}
}

func TestProgressBarDisplay(t *testing.T) {
	bar := NewProgressBar(10)
	bar.Width = 10 // Smaller width for easier testing
	bar.ShowTime = false // Disable time display for consistent testing

	output := captureOutput(func() {
		bar.Update(5)
	})

	expected := "\r[" + ColorBlue + "=====" + strings.Repeat(" ", 5) + ColorReset + "] 50%"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
}

func TestProgressBarCustomStyle(t *testing.T) {
	bar := NewProgressBar(10)
	bar.Width = 10
	bar.ShowTime = false

	customStyle := Style{
		Prefix:    "(",
		Suffix:    ")",
		FillChar:  "#",
		EmptyChar: "-",
		Color:     ColorGreen,
	}
	bar.SetStyle(customStyle)

	output := captureOutput(func() {
		bar.Update(5)
	})

	expected := "\r(" + ColorGreen + "#####" + strings.Repeat("-", 5) + ColorReset + ") 50%"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
}

func TestProgressBarFinish(t *testing.T) {
	bar := NewProgressBar(10)
	bar.Width = 10
	bar.ShowTime = false

	output := captureOutput(func() {
		bar.Finish()
	})

	expected := "\r[" + ColorBlue + strings.Repeat("=", 10) + ColorReset + "] 100%\n"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{3 * time.Second, "3s"},
		{65 * time.Second, "1m05s"},
		{3665 * time.Second, "1h01m05s"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("Expected formatDuration(%v) to be '%s', got '%s'", test.duration, test.expected, result)
		}
	}
}

func TestProgressBarTimeDisplay(t *testing.T) {
	bar := NewProgressBar(10)
	bar.Width = 10

	// Sleep for a short duration to test time display
	time.Sleep(2 * time.Second)

	output := captureOutput(func() {
		bar.Update(5)
	})

	if !strings.Contains(output, "[2s]") {
		t.Errorf("Expected output to contain time display, got '%s'", output)
	}
}

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner()

	if len(spinner.Frames) == 0 {
		t.Error("Expected Frames to be non-empty")
	}
	if spinner.Index != 0 {
		t.Errorf("Expected Index to be 0, got %d", spinner.Index)
	}
	if spinner.Color != ColorBlue {
		t.Errorf("Expected Color to be ColorBlue, got %s", spinner.Color)
	}
}

func TestSpinnerUpdate(t *testing.T) {
	spinner := NewSpinner()
	message := "Loading..."

	output := captureOutput(func() {
		spinner.Update(message)
	})

	if !strings.Contains(output, message) {
		t.Errorf("Expected output to contain '%s', got '%s'", message, output)
	}
	if !strings.Contains(output, ColorBlue) {
		t.Errorf("Expected output to contain color code '%s'", ColorBlue)
	}
}

func TestSpinnerFinish(t *testing.T) {
	spinner := NewSpinner()
	message := "Done!"

	output := captureOutput(func() {
		spinner.Finish(message)
	})

	expected := ColorGreen + "✓" + ColorReset + " Done!\n"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
	}
}

func TestSpinnerSetColor(t *testing.T) {
	spinner := NewSpinner()
	spinner.SetColor(ColorPurple)

	if spinner.Color != ColorPurple {
		t.Errorf("Expected Color to be ColorPurple, got %s", spinner.Color)
	}

	output := captureOutput(func() {
		spinner.Update("Loading...")
	})

	if !strings.Contains(output, ColorPurple) {
		t.Errorf("Expected output to contain color code '%s'", ColorPurple)
	}
} 