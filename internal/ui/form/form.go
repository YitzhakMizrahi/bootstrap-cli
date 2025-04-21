package form

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InputType defines the type of input
type InputType int

const (
	// TextInput is a simple text input
	TextInput InputType = iota
	// NumberInput is a numeric input
	NumberInput
	// PasswordInput is a password input (hidden)
	PasswordInput
	// EmailInput is an email input with validation
	EmailInput
	// URLInput is a URL input with validation
	URLInput
)

// FormField represents a field in a form
type FormField struct {
	Name        string
	Label       string
	Type        InputType
	Default     string
	Required    bool
	Validate    func(string) error
	Description string
	Value       string
	Error       string
}

// Form represents a collection of form fields
type Form struct {
	Fields []*FormField
	Width  int
}

// NewForm creates a new form
func NewForm() *Form {
	return &Form{
		Fields: []*FormField{},
		Width:  80,
	}
}

// AddField adds a field to the form
func (f *Form) AddField(field *FormField) {
	f.Fields = append(f.Fields, field)
}

// SetWidth sets the width of the form
func (f *Form) SetWidth(width int) {
	f.Width = width
}

// Display displays the form and collects input
func (f *Form) Display() error {
	reader := bufio.NewReader(os.Stdin)

	for _, field := range f.Fields {
		// Display field label and description
		fmt.Printf("\n%s\n", field.Label)
		if field.Description != "" {
			fmt.Printf("%s\n", field.Description)
		}

		// Display default value if provided
		defaultStr := ""
		if field.Default != "" {
			defaultStr = fmt.Sprintf(" (%s)", field.Default)
		}

		// Display required indicator
		requiredStr := ""
		if field.Required {
			requiredStr = " *"
		}

		// Prompt for input
		fmt.Printf("%s%s%s: ", field.Name, requiredStr, defaultStr)

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		// Trim whitespace
		input = strings.TrimSpace(input)

		// Use default if input is empty
		if input == "" && field.Default != "" {
			input = field.Default
		}

		// Validate required field
		if field.Required && input == "" {
			field.Error = "This field is required"
			fmt.Printf("Error: %s\n", field.Error)
			continue
		}

		// Validate input type
		switch field.Type {
		case NumberInput:
			if input != "" {
				_, err := strconv.ParseFloat(input, 64)
				if err != nil {
					field.Error = "Please enter a valid number"
					fmt.Printf("Error: %s\n", field.Error)
					continue
				}
			}
		case EmailInput:
			if input != "" && !isValidEmail(input) {
				field.Error = "Please enter a valid email address"
				fmt.Printf("Error: %s\n", field.Error)
				continue
			}
		case URLInput:
			if input != "" && !isValidURL(input) {
				field.Error = "Please enter a valid URL"
				fmt.Printf("Error: %s\n", field.Error)
				continue
			}
		}

		// Run custom validation if provided
		if field.Validate != nil {
			if err := field.Validate(input); err != nil {
				field.Error = err.Error()
				fmt.Printf("Error: %s\n", field.Error)
				continue
			}
		}

		// Store the value
		field.Value = input
		field.Error = ""
	}

	return nil
}

// GetValues returns a map of field names to values
func (f *Form) GetValues() map[string]string {
	values := make(map[string]string)
	for _, field := range f.Fields {
		values[field.Name] = field.Value
	}
	return values
}

// HasErrors checks if any fields have errors
func (f *Form) HasErrors() bool {
	for _, field := range f.Fields {
		if field.Error != "" {
			return true
		}
	}
	return false
}

// isValidEmail checks if a string is a valid email address
func isValidEmail(email string) bool {
	// Simple email validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

// isValidURL checks if a string is a valid URL
func isValidURL(url string) bool {
	// Simple URL validation
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	if len(url) < 8 {
		return false
	}
	return true
} 