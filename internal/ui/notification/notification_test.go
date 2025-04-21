package notification

import (
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput captures the output of a function
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

	var buf strings.Builder
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewNotificationManager(t *testing.T) {
	manager := NewNotificationManager()
	if manager == nil {
		t.Error("Expected non-nil NotificationManager")
	}
	if len(manager.notifications) != 0 {
		t.Errorf("Expected empty notifications, got %d notifications", len(manager.notifications))
	}
	if manager.maxNotifications != 5 {
		t.Errorf("Expected maxNotifications 5, got %d", manager.maxNotifications)
	}
	if manager.width != 80 {
		t.Errorf("Expected width 80, got %d", manager.width)
	}
}

func TestSetMaxNotifications(t *testing.T) {
	manager := NewNotificationManager()
	manager.SetMaxNotifications(10)
	if manager.maxNotifications != 10 {
		t.Errorf("Expected maxNotifications 10, got %d", manager.maxNotifications)
	}
}

func TestSetWidth(t *testing.T) {
	manager := NewNotificationManager()
	manager.SetWidth(100)
	if manager.width != 100 {
		t.Errorf("Expected width 100, got %d", manager.width)
	}
}

func TestAddNotification(t *testing.T) {
	manager := NewNotificationManager()
	notification := &Notification{
		Type:    InfoNotification,
		Message: "Test message",
		Title:   "Test title",
	}
	
	// Capture the output
	output := captureOutput(func() {
		manager.AddNotification(notification)
	})
	
	// Check that the notification was added
	if len(manager.notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(manager.notifications))
	}
	
	// Check that the output contains the notification
	if !strings.Contains(output, "Test title") {
		t.Errorf("Expected output to contain 'Test title', got '%s'", output)
	}
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected output to contain 'Test message', got '%s'", output)
	}
}

func TestDisplay(t *testing.T) {
	manager := NewNotificationManager()
	manager.AddNotification(&Notification{
		Type:    InfoNotification,
		Message: "Test message 1",
		Title:   "Test title 1",
	})
	manager.AddNotification(&Notification{
		Type:    SuccessNotification,
		Message: "Test message 2",
		Title:   "Test title 2",
	})
	
	// Capture the output
	output := captureOutput(func() {
		manager.Display()
	})
	
	// Check that the output contains both notifications
	if !strings.Contains(output, "Test title 1") {
		t.Errorf("Expected output to contain 'Test title 1', got '%s'", output)
	}
	if !strings.Contains(output, "Test message 1") {
		t.Errorf("Expected output to contain 'Test message 1', got '%s'", output)
	}
	if !strings.Contains(output, "Test title 2") {
		t.Errorf("Expected output to contain 'Test title 2', got '%s'", output)
	}
	if !strings.Contains(output, "Test message 2") {
		t.Errorf("Expected output to contain 'Test message 2', got '%s'", output)
	}
}

func TestGetStyleForType(t *testing.T) {
	manager := NewNotificationManager()
	
	tests := []struct {
		notificationType NotificationType
		expectedColor    string
		expectedIcon     string
	}{
		{InfoNotification, "\033[34m", "ℹ"},
		{SuccessNotification, "\033[32m", "✓"},
		{WarningNotification, "\033[33m", "⚠"},
		{ErrorNotification, "\033[31m", "✗"},
	}
	
	for _, test := range tests {
		color, icon := manager.getStyleForType(test.notificationType)
		if color != test.expectedColor {
			t.Errorf("getStyleForType(%v) color = %v, expected %v", test.notificationType, color, test.expectedColor)
		}
		if icon != test.expectedIcon {
			t.Errorf("getStyleForType(%v) icon = %v, expected %v", test.notificationType, icon, test.expectedIcon)
		}
	}
}

func TestCreateNotificationBox(t *testing.T) {
	manager := NewNotificationManager()
	notification := &Notification{
		Type:    InfoNotification,
		Message: "Test message",
		Title:   "Test title",
	}
	
	box := manager.createNotificationBox(notification, "\033[34m", "ℹ")
	
	// Check that the box contains the title and message
	if !strings.Contains(box, "Test title") {
		t.Errorf("Expected box to contain 'Test title', got '%s'", box)
	}
	if !strings.Contains(box, "Test message") {
		t.Errorf("Expected box to contain 'Test message', got '%s'", box)
	}
	
	// Check that the box has the correct structure
	if !strings.Contains(box, "┌") {
		t.Errorf("Expected box to contain top border, got '%s'", box)
	}
	if !strings.Contains(box, "┐") {
		t.Errorf("Expected box to contain top right corner, got '%s'", box)
	}
	if !strings.Contains(box, "└") {
		t.Errorf("Expected box to contain bottom left corner, got '%s'", box)
	}
	if !strings.Contains(box, "┘") {
		t.Errorf("Expected box to contain bottom right corner, got '%s'", box)
	}
}

func TestClear(t *testing.T) {
	manager := NewNotificationManager()
	manager.AddNotification(&Notification{
		Type:    InfoNotification,
		Message: "Test message",
		Title:   "Test title",
	})
	
	manager.Clear()
	
	if len(manager.notifications) != 0 {
		t.Errorf("Expected empty notifications, got %d notifications", len(manager.notifications))
	}
}

func TestNotificationMethods(t *testing.T) {
	manager := NewNotificationManager()
	
	// Test Info method
	output := captureOutput(func() {
		manager.Info("Info message", "Info title")
	})
	if !strings.Contains(output, "Info title") {
		t.Errorf("Expected output to contain 'Info title', got '%s'", output)
	}
	if !strings.Contains(output, "Info message") {
		t.Errorf("Expected output to contain 'Info message', got '%s'", output)
	}
	
	// Test Success method
	output = captureOutput(func() {
		manager.Success("Success message", "Success title")
	})
	if !strings.Contains(output, "Success title") {
		t.Errorf("Expected output to contain 'Success title', got '%s'", output)
	}
	if !strings.Contains(output, "Success message") {
		t.Errorf("Expected output to contain 'Success message', got '%s'", output)
	}
	
	// Test Warning method
	output = captureOutput(func() {
		manager.Warning("Warning message", "Warning title")
	})
	if !strings.Contains(output, "Warning title") {
		t.Errorf("Expected output to contain 'Warning title', got '%s'", output)
	}
	if !strings.Contains(output, "Warning message") {
		t.Errorf("Expected output to contain 'Warning message', got '%s'", output)
	}
	
	// Test Error method
	output = captureOutput(func() {
		manager.Error("Error message", "Error title")
	})
	if !strings.Contains(output, "Error title") {
		t.Errorf("Expected output to contain 'Error title', got '%s'", output)
	}
	if !strings.Contains(output, "Error message") {
		t.Errorf("Expected output to contain 'Error message', got '%s'", output)
	}
}

func TestNotify(t *testing.T) {
	// Test Notify function
	output := captureOutput(func() {
		Notify(InfoNotification, "Test message", "Test title")
	})
	if !strings.Contains(output, "Test title") {
		t.Errorf("Expected output to contain 'Test title', got '%s'", output)
	}
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected output to contain 'Test message', got '%s'", output)
	}
} 