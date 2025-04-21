package prompts

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// mockStdin creates a mock stdin for testing
func mockStdin(input string) (*os.File, func()) {
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	
	// Write the mock input
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()
	
	// Return cleanup function
	return oldStdin, func() {
		os.Stdin = oldStdin
	}
}

func TestNewPrompt(t *testing.T) {
	message := "Enter your name"
	prompt := NewPrompt(TextPrompt, message)

	if prompt.Type != TextPrompt {
		t.Errorf("Expected Type to be TextPrompt, got %v", prompt.Type)
	}
	if prompt.Message != message {
		t.Errorf("Expected Message to be '%s', got '%s'", message, prompt.Message)
	}
	if prompt.Color != "\033[34m" {
		t.Errorf("Expected Color to be '\\033[34m', got '%s'", prompt.Color)
	}
}

func TestTextPrompt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		default_ string
		required bool
		validate func(string) error
		want     string
		wantErr  bool
	}{
		{
			name:  "basic input",
			input: "test\n",
			want:  "test",
		},
		{
			name:     "empty input with default",
			input:    "\n",
			default_: "default",
			want:     "default",
		},
		{
			name:     "required input - empty then valid",
			input:    "\ntest\n",
			required: true,
			want:     "test",
		},
		{
			name:     "validation - invalid then valid",
			input:    "short\nvalid input\n",
			validate: func(s string) error {
				if len(s) < 10 {
					return errors.New("input must be at least 10 characters")
				}
				return nil
			},
			want: "valid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock stdin
			oldStdin, cleanup := mockStdin(tt.input)
			defer cleanup()

			// Create prompt
			prompt := NewPrompt(TextPrompt, "Test prompt")
			if tt.default_ != "" {
				prompt.WithDefault(tt.default_)
			}
			if tt.required {
				prompt.WithRequired(true)
			}
			if tt.validate != nil {
				prompt.WithValidation(tt.validate)
			}

			// Get response
			got, err := prompt.Ask()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ask() = %v, want %v", got, tt.want)
			}

			// Restore stdin
			os.Stdin = oldStdin
		})
	}
}

func TestConfirmPrompt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		default_ string
		want     string
	}{
		{
			name:  "yes input",
			input: "y\n",
			want:  "yes",
		},
		{
			name:  "no input",
			input: "n\n",
			want:  "no",
		},
		{
			name:     "empty input with default yes",
			input:    "\n",
			default_: "y",
			want:     "yes",
		},
		{
			name:     "empty input with default no",
			input:    "\n",
			default_: "n",
			want:     "no",
		},
		{
			name:  "invalid then valid input",
			input: "invalid\ny\n",
			want:  "yes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock stdin
			oldStdin, cleanup := mockStdin(tt.input)
			defer cleanup()

			// Create prompt
			prompt := NewPrompt(ConfirmPrompt, "Test prompt")
			if tt.default_ != "" {
				prompt.WithDefault(tt.default_)
			}

			// Get response
			got, err := prompt.Ask()
			if err != nil {
				t.Errorf("Ask() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Ask() = %v, want %v", got, tt.want)
			}

			// Restore stdin
			os.Stdin = oldStdin
		})
	}
}

func TestSelectPrompt(t *testing.T) {
	options := []string{"Option 1", "Option 2", "Option 3"}
	tests := []struct {
		name     string
		input    string
		default_ string
		want     string
		wantErr  bool
	}{
		{
			name:  "valid selection",
			input: "2\n",
			want:  "Option 2",
		},
		{
			name:     "empty input with default",
			input:    "\n",
			default_: "1",
			want:     "Option 1",
		},
		{
			name:  "invalid then valid input",
			input: "invalid\n4\n2\n",
			want:  "Option 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock stdin
			oldStdin, cleanup := mockStdin(tt.input)
			defer cleanup()

			// Create prompt
			prompt := NewPrompt(SelectPrompt, "Test prompt").
				WithOptions(options)
			if tt.default_ != "" {
				prompt.WithDefault(tt.default_)
			}

			// Get response
			got, err := prompt.Ask()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ask() = %v, want %v", got, tt.want)
			}

			// Restore stdin
			os.Stdin = oldStdin
		})
	}
}

func TestPromptWithDescription(t *testing.T) {
	prompt := NewPrompt(TextPrompt, "Test prompt")
	description := "This is a test description"
	
	prompt.WithDescription(description)
	
	if prompt.Description != description {
		t.Errorf("Expected Description to be '%s', got '%s'", description, prompt.Description)
	}
}

func TestPromptWithColor(t *testing.T) {
	prompt := NewPrompt(TextPrompt, "Test prompt")
	color := "\033[32m" // Green
	
	prompt.WithColor(color)
	
	if prompt.Color != color {
		t.Errorf("Expected Color to be '%s', got '%s'", color, prompt.Color)
	}
}

func TestSelectPromptNoOptions(t *testing.T) {
	prompt := NewPrompt(SelectPrompt, "Test prompt")
	
	_, err := prompt.Ask()
	if err == nil {
		t.Error("Expected error for select prompt with no options")
	}
	if !strings.Contains(err.Error(), "no options provided") {
		t.Errorf("Expected error message to contain 'no options provided', got '%s'", err.Error())
	}
} 