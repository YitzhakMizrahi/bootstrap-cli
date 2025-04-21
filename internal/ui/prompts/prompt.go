package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PromptType represents the type of prompt
type PromptType int

const (
	// TextPrompt is a simple text input prompt
	TextPrompt PromptType = iota
	// ConfirmPrompt is a yes/no confirmation prompt
	ConfirmPrompt
	// SelectPrompt is a selection from multiple options
	SelectPrompt
)

// Prompt represents a user prompt
type Prompt struct {
	Type        PromptType
	Message     string
	Default     string
	Options     []string
	Required    bool
	Validate    func(string) error
	Color       string
	Description string
}

// NewPrompt creates a new prompt
func NewPrompt(promptType PromptType, message string) *Prompt {
	return &Prompt{
		Type:    promptType,
		Message: message,
		Color:   "\033[34m", // Default to blue
	}
}

// Ask displays the prompt and returns the user's response
func (p *Prompt) Ask() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	switch p.Type {
	case TextPrompt:
		return p.askText(reader)
	case ConfirmPrompt:
		return p.askConfirm(reader)
	case SelectPrompt:
		return p.askSelect(reader)
	default:
		return "", fmt.Errorf("unsupported prompt type")
	}
}

// askText prompts for text input
func (p *Prompt) askText(reader *bufio.Reader) (string, error) {
	for {
		// Display prompt
		if p.Default != "" {
			fmt.Printf("%s%s%s (%s): ", p.Color, p.Message, "\033[0m", p.Default)
		} else {
			fmt.Printf("%s%s%s: ", p.Color, p.Message, "\033[0m")
		}

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// Trim whitespace
		input = strings.TrimSpace(input)

		// Use default if input is empty
		if input == "" && p.Default != "" {
			input = p.Default
		}

		// Check if input is required
		if p.Required && input == "" {
			fmt.Println("Input is required. Please try again.")
			continue
		}

		// Validate input if validator is provided
		if p.Validate != nil {
			if err := p.Validate(input); err != nil {
				fmt.Printf("Invalid input: %s\n", err)
				continue
			}
		}

		return input, nil
	}
}

// askConfirm prompts for yes/no confirmation
func (p *Prompt) askConfirm(reader *bufio.Reader) (string, error) {
	defaultValue := strings.ToLower(p.Default)
	isDefaultYes := defaultValue == "y" || defaultValue == "yes"

	for {
		// Display prompt
		if p.Default != "" {
			fmt.Printf("%s%s%s [%s/%s]: ", p.Color, p.Message, "\033[0m",
				map[bool]string{true: "Y/n", false: "y/N"}[isDefaultYes],
				p.Description)
		} else {
			fmt.Printf("%s%s%s [y/n]: ", p.Color, p.Message, "\033[0m")
		}

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// Trim whitespace
		input = strings.TrimSpace(strings.ToLower(input))

		// Use default if input is empty
		if input == "" && p.Default != "" {
			input = defaultValue
		}

		// Parse response
		switch input {
		case "y", "yes":
			return "yes", nil
		case "n", "no":
			return "no", nil
		default:
			fmt.Println("Please answer 'yes' or 'no'")
		}
	}
}

// askSelect prompts for selection from options
func (p *Prompt) askSelect(reader *bufio.Reader) (string, error) {
	if len(p.Options) == 0 {
		return "", fmt.Errorf("no options provided for selection")
	}

	for {
		// Display options
		fmt.Printf("%s%s%s\n", p.Color, p.Message, "\033[0m")
		for i, option := range p.Options {
			fmt.Printf("%d) %s\n", i+1, option)
		}

		// Display prompt
		if p.Default != "" {
			fmt.Printf("Enter number (%s): ", p.Default)
		} else {
			fmt.Print("Enter number: ")
		}

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// Trim whitespace
		input = strings.TrimSpace(input)

		// Use default if input is empty
		if input == "" && p.Default != "" {
			input = p.Default
		}

		// Parse selection
		if num, err := strconv.Atoi(input); err == nil {
			if num > 0 && num <= len(p.Options) {
				return p.Options[num-1], nil
			}
		}

		fmt.Printf("Please enter a number between 1 and %d\n", len(p.Options))
	}
}

// WithDefault sets the default value for the prompt
func (p *Prompt) WithDefault(value string) *Prompt {
	p.Default = value
	return p
}

// WithOptions sets the options for a select prompt
func (p *Prompt) WithOptions(options []string) *Prompt {
	p.Options = options
	return p
}

// WithRequired sets whether the prompt is required
func (p *Prompt) WithRequired(required bool) *Prompt {
	p.Required = required
	return p
}

// WithValidation sets a validation function for the prompt
func (p *Prompt) WithValidation(validate func(string) error) *Prompt {
	p.Validate = validate
	return p
}

// WithColor sets the color for the prompt
func (p *Prompt) WithColor(color string) *Prompt {
	p.Color = color
	return p
}

// WithDescription sets a description for the prompt
func (p *Prompt) WithDescription(description string) *Prompt {
	p.Description = description
	return p
} 