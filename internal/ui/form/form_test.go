package form

import (
	"os"
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

func TestNewForm(t *testing.T) {
	form := NewForm()
	if form == nil {
		t.Error("Expected non-nil Form")
	}
	if len(form.Fields) != 0 {
		t.Errorf("Expected empty fields, got %d fields", len(form.Fields))
	}
	if form.Width != 80 {
		t.Errorf("Expected width 80, got %d", form.Width)
	}
}

func TestAddField(t *testing.T) {
	form := NewForm()
	field := &FormField{
		Name:  "test",
		Label: "Test Field",
		Type:  TextInput,
	}
	
	form.AddField(field)
	
	if len(form.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d fields", len(form.Fields))
	}
	if form.Fields[0] != field {
		t.Error("Expected field to be added to form")
	}
}

func TestSetWidth(t *testing.T) {
	form := NewForm()
	form.SetWidth(100)
	
	if form.Width != 100 {
		t.Errorf("Expected width 100, got %d", form.Width)
	}
}

func TestDisplay(t *testing.T) {
	// Create a form with a text field
	form := NewForm()
	form.AddField(&FormField{
		Name:     "name",
		Label:    "Name",
		Type:     TextInput,
		Required: true,
	})
	
	// Mock stdin with valid input
	_, cleanup := mockStdin("John Doe\n")
	defer cleanup()
	
	// Display the form
	err := form.Display()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check the field value
	if form.Fields[0].Value != "John Doe" {
		t.Errorf("Expected value 'John Doe', got '%s'", form.Fields[0].Value)
	}
}

func TestDisplayWithDefault(t *testing.T) {
	// Create a form with a text field with default value
	form := NewForm()
	form.AddField(&FormField{
		Name:     "name",
		Label:    "Name",
		Type:     TextInput,
		Default:  "Default Name",
		Required: true,
	})
	
	// Mock stdin with empty input (should use default)
	_, cleanup := mockStdin("\n")
	defer cleanup()
	
	// Display the form
	err := form.Display()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check the field value
	if form.Fields[0].Value != "Default Name" {
		t.Errorf("Expected value 'Default Name', got '%s'", form.Fields[0].Value)
	}
}

func TestDisplayWithValidation(t *testing.T) {
	// Create a form with a number field
	form := NewForm()
	form.AddField(&FormField{
		Name:     "age",
		Label:    "Age",
		Type:     NumberInput,
		Required: true,
	})
	
	// Mock stdin with invalid input followed by valid input
	_, cleanup := mockStdin("abc\n25\n")
	defer cleanup()
	
	// Display the form
	err := form.Display()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check the field value
	if form.Fields[0].Value != "25" {
		t.Errorf("Expected value '25', got '%s'", form.Fields[0].Value)
	}
}

func TestGetValues(t *testing.T) {
	// Create a form with multiple fields
	form := NewForm()
	form.AddField(&FormField{
		Name:  "name",
		Label: "Name",
		Type:  TextInput,
		Value: "John Doe",
	})
	form.AddField(&FormField{
		Name:  "age",
		Label: "Age",
		Type:  NumberInput,
		Value: "30",
	})
	
	// Get the values
	values := form.GetValues()
	
	// Check the values
	if values["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", values["name"])
	}
	if values["age"] != "30" {
		t.Errorf("Expected age '30', got '%s'", values["age"])
	}
}

func TestHasErrors(t *testing.T) {
	// Create a form with no errors
	form := NewForm()
	form.AddField(&FormField{
		Name:  "name",
		Label: "Name",
		Type:  TextInput,
		Value: "John Doe",
	})
	
	// Check for errors
	if form.HasErrors() {
		t.Error("Expected no errors")
	}
	
	// Add a field with an error
	form.AddField(&FormField{
		Name:  "age",
		Label: "Age",
		Type:  NumberInput,
		Value: "30",
		Error: "Invalid age",
	})
	
	// Check for errors
	if !form.HasErrors() {
		t.Error("Expected errors")
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"test@example", false},
		{"test@", false},
		{"@example.com", false},
		{"test@example.com.", false},
		{"", false},
	}
	
	for _, test := range tests {
		result := isValidEmail(test.email)
		if result != test.expected {
			t.Errorf("isValidEmail(%s) = %v, expected %v", test.email, result, test.expected)
		}
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", false},
		{"example.com", false},
		{"http://", false},
		{"", false},
	}
	
	for _, test := range tests {
		result := isValidURL(test.url)
		if result != test.expected {
			t.Errorf("isValidURL(%s) = %v, expected %v", test.url, result, test.expected)
		}
	}
} 