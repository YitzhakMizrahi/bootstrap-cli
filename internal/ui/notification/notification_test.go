package notification

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// TestSaveAndLoadNotifications tests saving and loading notifications
func TestSaveAndLoadNotifications(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "notification-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a notification manager with custom storage path
	manager := NewNotificationManager()
	manager.storagePath = filepath.Join(tempDir, "notifications.json")
	
	// Add some notifications
	manager.Info("Test info message", "Info Title")
	manager.Success("Test success message", "Success Title")
	manager.Warning("Test warning message", "Warning Title")
	
	// Create a new manager to test loading
	manager2 := NewNotificationManager()
	manager2.storagePath = filepath.Join(tempDir, "notifications.json")
	
	// Load notifications
	if err := manager2.LoadNotifications(); err != nil {
		t.Fatalf("Failed to load notifications: %v", err)
	}
	
	// Check that notifications were loaded correctly
	if len(manager2.notifications) != 3 {
		t.Errorf("Expected 3 notifications, got %d", len(manager2.notifications))
	}
	
	// Check notification types
	if manager2.notifications[0].Type != InfoNotification {
		t.Errorf("Expected InfoNotification, got %v", manager2.notifications[0].Type)
	}
	if manager2.notifications[1].Type != SuccessNotification {
		t.Errorf("Expected SuccessNotification, got %v", manager2.notifications[1].Type)
	}
	if manager2.notifications[2].Type != WarningNotification {
		t.Errorf("Expected WarningNotification, got %v", manager2.notifications[2].Type)
	}
	
	// Check notification messages
	if manager2.notifications[0].Message != "Test info message" {
		t.Errorf("Expected 'Test info message', got '%s'", manager2.notifications[0].Message)
	}
	if manager2.notifications[1].Message != "Test success message" {
		t.Errorf("Expected 'Test success message', got '%s'", manager2.notifications[1].Message)
	}
	if manager2.notifications[2].Message != "Test warning message" {
		t.Errorf("Expected 'Test warning message', got '%s'", manager2.notifications[2].Message)
	}
	
	// Check notification titles
	if manager2.notifications[0].Title != "Info Title" {
		t.Errorf("Expected 'Info Title', got '%s'", manager2.notifications[0].Title)
	}
	if manager2.notifications[1].Title != "Success Title" {
		t.Errorf("Expected 'Success Title', got '%s'", manager2.notifications[1].Title)
	}
	if manager2.notifications[2].Title != "Warning Title" {
		t.Errorf("Expected 'Warning Title', got '%s'", manager2.notifications[2].Title)
	}
}

// TestClearWithPersistence tests that clearing notifications also clears the file
func TestClearWithPersistence(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "notification-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a notification manager with custom storage path
	manager := NewNotificationManager()
	manager.storagePath = filepath.Join(tempDir, "notifications.json")
	
	// Add some notifications
	manager.Info("Test info message", "Info Title")
	manager.Success("Test success message", "Success Title")
	
	// Clear notifications
	manager.Clear()
	
	// Create a new manager to test loading
	manager2 := NewNotificationManager()
	manager2.storagePath = filepath.Join(tempDir, "notifications.json")
	
	// Load notifications
	if err := manager2.LoadNotifications(); err != nil {
		t.Fatalf("Failed to load notifications: %v", err)
	}
	
	// Check that notifications were cleared
	if len(manager2.notifications) != 0 {
		t.Errorf("Expected 0 notifications, got %d", len(manager2.notifications))
	}
}

// TestNotificationExpiration tests that notifications expire after their duration
func TestNotificationExpiration(t *testing.T) {
	// Create a notification manager
	manager := NewNotificationManager()
	
	// Add a notification with a short duration
	manager.InfoWithDuration("This will expire", "Expiring Notification", 100*time.Millisecond)
	
	// Check that the notification was added
	if len(manager.notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(manager.notifications))
	}
	
	// Wait for the notification to expire
	time.Sleep(200 * time.Millisecond)
	
	// Manually trigger cleanup
	manager.cleanupExpiredNotifications()
	
	// Check that the notification was removed
	if len(manager.notifications) != 0 {
		t.Errorf("Expected 0 notifications after expiration, got %d", len(manager.notifications))
	}
	
	// Stop the cleanup goroutine
	manager.StopCleanup()
}

// TestNotificationWithDuration tests creating notifications with duration
func TestNotificationWithDuration(t *testing.T) {
	// Create a notification manager
	manager := NewNotificationManager()
	
	// Add notifications with different durations
	manager.InfoWithDuration("Short duration", "Short", 1*time.Second)
	manager.SuccessWithDuration("Medium duration", "Medium", 5*time.Second)
	manager.WarningWithDuration("Long duration", "Long", 10*time.Second)
	manager.ErrorWithDuration("No expiration", "No Exp", 0)
	
	// Check that all notifications were added
	if len(manager.notifications) != 4 {
		t.Errorf("Expected 4 notifications, got %d", len(manager.notifications))
	}
	
	// Check that expiration times were set correctly
	if manager.notifications[0].ExpiresAt.IsZero() {
		t.Error("Expected non-zero expiration time for notification with duration")
	}
	
	if manager.notifications[1].ExpiresAt.IsZero() {
		t.Error("Expected non-zero expiration time for notification with duration")
	}
	
	if manager.notifications[2].ExpiresAt.IsZero() {
		t.Error("Expected non-zero expiration time for notification with duration")
	}
	
	if !manager.notifications[3].ExpiresAt.IsZero() {
		t.Error("Expected zero expiration time for notification without duration")
	}
	
	// Stop the cleanup goroutine
	manager.StopCleanup()
}

// TestNotificationPriority tests that notifications are sorted by priority
func TestNotificationPriority(t *testing.T) {
	// Create a notification manager
	manager := NewNotificationManager()
	
	// Add notifications with different priorities
	manager.AddNotificationWithPriority(InfoNotification, LowPriority, "Low priority message", "Low Priority")
	manager.AddNotificationWithPriority(InfoNotification, CriticalPriority, "Critical priority message", "Critical Priority")
	manager.AddNotificationWithPriority(InfoNotification, NormalPriority, "Normal priority message", "Normal Priority")
	manager.AddNotificationWithPriority(InfoNotification, HighPriority, "High priority message", "High Priority")
	
	// Check that notifications are sorted by priority (highest first)
	if len(manager.notifications) != 4 {
		t.Errorf("Expected 4 notifications, got %d", len(manager.notifications))
	}
	
	// Check that the first notification is the highest priority
	if manager.notifications[0].Priority != CriticalPriority {
		t.Errorf("Expected first notification to be CriticalPriority, got %v", manager.notifications[0].Priority)
	}
	
	// Check that the second notification is the second highest priority
	if manager.notifications[1].Priority != HighPriority {
		t.Errorf("Expected second notification to be HighPriority, got %v", manager.notifications[1].Priority)
	}
	
	// Check that the third notification is the third highest priority
	if manager.notifications[2].Priority != NormalPriority {
		t.Errorf("Expected third notification to be NormalPriority, got %v", manager.notifications[2].Priority)
	}
	
	// Check that the fourth notification is the lowest priority
	if manager.notifications[3].Priority != LowPriority {
		t.Errorf("Expected fourth notification to be LowPriority, got %v", manager.notifications[3].Priority)
	}
}

// TestDefaultPriority tests that notifications get default priorities based on their type
func TestDefaultPriority(t *testing.T) {
	// Create a notification manager
	manager := NewNotificationManager()
	
	// Add notifications with different types
	manager.Info("Info message", "Info")
	manager.Success("Success message", "Success")
	manager.Warning("Warning message", "Warning")
	manager.Error("Error message", "Error")
	
	// Check that notifications have the correct default priorities
	if len(manager.notifications) != 4 {
		t.Errorf("Expected 4 notifications, got %d", len(manager.notifications))
	}
	
	// Find notifications by type
	var infoNotification, successNotification, warningNotification, errorNotification *Notification
	
	for _, notification := range manager.notifications {
		switch notification.Type {
		case InfoNotification:
			infoNotification = notification
		case SuccessNotification:
			successNotification = notification
		case WarningNotification:
			warningNotification = notification
		case ErrorNotification:
			errorNotification = notification
		}
	}
	
	// Check that Info notifications have NormalPriority
	if infoNotification.Priority != NormalPriority {
		t.Errorf("Expected Info notification to have NormalPriority, got %v", infoNotification.Priority)
	}
	
	// Check that Success notifications have NormalPriority
	if successNotification.Priority != NormalPriority {
		t.Errorf("Expected Success notification to have NormalPriority, got %v", successNotification.Priority)
	}
	
	// Check that Warning notifications have HighPriority
	if warningNotification.Priority != HighPriority {
		t.Errorf("Expected Warning notification to have HighPriority, got %v", warningNotification.Priority)
	}
	
	// Check that Error notifications have CriticalPriority
	if errorNotification.Priority != CriticalPriority {
		t.Errorf("Expected Error notification to have CriticalPriority, got %v", errorNotification.Priority)
	}
}

// TestPriorityWithDuration tests that notifications with priority and duration work correctly
func TestPriorityWithDuration(t *testing.T) {
	// Create a notification manager
	manager := NewNotificationManager()
	
	// Add notifications with different priorities and durations
	manager.AddNotificationWithPriorityAndDuration(InfoNotification, LowPriority, "Low priority message", "Low Priority", 5*time.Second)
	manager.AddNotificationWithPriorityAndDuration(InfoNotification, CriticalPriority, "Critical priority message", "Critical Priority", 10*time.Second)
	
	// Check that notifications are sorted by priority (highest first)
	if len(manager.notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(manager.notifications))
	}
	
	// Check that the first notification is the highest priority
	if manager.notifications[0].Priority != CriticalPriority {
		t.Errorf("Expected first notification to be CriticalPriority, got %v", manager.notifications[0].Priority)
	}
	
	// Check that the second notification is the lowest priority
	if manager.notifications[1].Priority != LowPriority {
		t.Errorf("Expected second notification to be LowPriority, got %v", manager.notifications[1].Priority)
	}
	
	// Check that the first notification has the correct duration
	if manager.notifications[0].Duration != 10*time.Second {
		t.Errorf("Expected first notification to have duration 10s, got %v", manager.notifications[0].Duration)
	}
	
	// Check that the second notification has the correct duration
	if manager.notifications[1].Duration != 5*time.Second {
		t.Errorf("Expected second notification to have duration 5s, got %v", manager.notifications[1].Duration)
	}
	
	// Check that the first notification has the correct expiration time
	expectedExpiration := manager.notifications[0].Timestamp.Add(10 * time.Second)
	if !manager.notifications[0].ExpiresAt.Equal(expectedExpiration) {
		t.Errorf("Expected first notification to expire at %v, got %v", expectedExpiration, manager.notifications[0].ExpiresAt)
	}
	
	// Check that the second notification has the correct expiration time
	expectedExpiration = manager.notifications[1].Timestamp.Add(5 * time.Second)
	if !manager.notifications[1].ExpiresAt.Equal(expectedExpiration) {
		t.Errorf("Expected second notification to expire at %v, got %v", expectedExpiration, manager.notifications[1].ExpiresAt)
	}
}

// TestNotificationCategories tests that notifications are properly categorized
func TestNotificationCategories(t *testing.T) {
	// Create a new notification manager
	manager := NewNotificationManager()
	
	// Add notifications with different categories
	manager.AddNotificationWithCategory(InfoNotification, "Info message", "Info title", "System")
	manager.AddNotificationWithCategory(SuccessNotification, "Success message", "Success title", "Installation")
	manager.AddNotificationWithCategory(WarningNotification, "Warning message", "Warning title", "System")
	manager.AddNotificationWithCategory(ErrorNotification, "Error message", "Error title", "Installation")
	
	// Get all categories
	categories := manager.GetCategories()
	
	// Check that we have the expected categories
	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}
	
	// Check that the categories are as expected
	expectedCategories := []string{"Installation", "System"}
	for i, category := range categories {
		if category != expectedCategories[i] {
			t.Errorf("Expected category %s, got %s", expectedCategories[i], category)
		}
	}
	
	// Get notifications by category
	systemNotifications := manager.GetNotificationsByCategory("System")
	installationNotifications := manager.GetNotificationsByCategory("Installation")
	
	// Check that we have the expected number of notifications in each category
	if len(systemNotifications) != 2 {
		t.Errorf("Expected 2 notifications in System category, got %d", len(systemNotifications))
	}
	
	if len(installationNotifications) != 2 {
		t.Errorf("Expected 2 notifications in Installation category, got %d", len(installationNotifications))
	}
	
	// Check that the notifications in each category have the expected types
	systemTypes := make(map[NotificationType]bool)
	for _, notification := range systemNotifications {
		systemTypes[notification.Type] = true
	}
	
	installationTypes := make(map[NotificationType]bool)
	for _, notification := range installationNotifications {
		installationTypes[notification.Type] = true
	}
	
	// Check System category
	if !systemTypes[InfoNotification] {
		t.Error("Expected InfoNotification in System category")
	}
	if !systemTypes[WarningNotification] {
		t.Error("Expected WarningNotification in System category")
	}
	
	// Check Installation category
	if !installationTypes[SuccessNotification] {
		t.Error("Expected SuccessNotification in Installation category")
	}
	if !installationTypes[ErrorNotification] {
		t.Error("Expected ErrorNotification in Installation category")
	}
}

// TestDefaultCategory tests that notifications get the default category if none is specified
func TestDefaultCategory(t *testing.T) {
	// Create a new notification manager
	manager := NewNotificationManager()
	
	// Add notifications without specifying a category
	manager.Info("Info message", "Info title")
	manager.Success("Success message", "Success title")
	manager.Warning("Warning message", "Warning title")
	manager.Error("Error message", "Error title")
	
	// Get all categories
	categories := manager.GetCategories()
	
	// Check that we have only the default category
	if len(categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(categories))
	}
	
	if categories[0] != "General" {
		t.Errorf("Expected category General, got %s", categories[0])
	}
	
	// Get notifications by category
	generalNotifications := manager.GetNotificationsByCategory("General")
	
	// Check that we have the expected number of notifications
	if len(generalNotifications) != 4 {
		t.Errorf("Expected 4 notifications in General category, got %d", len(generalNotifications))
	}
}

// TestNotificationGrouping tests that notifications are properly grouped by category
func TestNotificationGrouping(t *testing.T) {
	// Create a new notification manager
	manager := NewNotificationManager()
	
	// Add notifications with different categories
	manager.AddNotificationWithCategory(InfoNotification, "Info message", "Info title", "System")
	manager.AddNotificationWithCategory(SuccessNotification, "Success message", "Success title", "Installation")
	manager.AddNotificationWithCategory(WarningNotification, "Warning message", "Warning title", "System")
	manager.AddNotificationWithCategory(ErrorNotification, "Error message", "Error title", "Installation")
	
	// Get notifications grouped by category
	groupedNotifications := manager.groupNotificationsByCategory()
	
	// Check that we have the expected number of categories
	if len(groupedNotifications) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(groupedNotifications))
	}
	
	// Check that each category has the expected number of notifications
	if len(groupedNotifications["System"]) != 2 {
		t.Errorf("Expected 2 notifications in System category, got %d", len(groupedNotifications["System"]))
	}
	
	if len(groupedNotifications["Installation"]) != 2 {
		t.Errorf("Expected 2 notifications in Installation category, got %d", len(groupedNotifications["Installation"]))
	}
	
	// Check that the notifications in each category have the expected types
	systemTypes := make(map[NotificationType]bool)
	for _, notification := range groupedNotifications["System"] {
		systemTypes[notification.Type] = true
	}
	
	installationTypes := make(map[NotificationType]bool)
	for _, notification := range groupedNotifications["Installation"] {
		installationTypes[notification.Type] = true
	}
	
	// Check System category
	if !systemTypes[InfoNotification] {
		t.Error("Expected InfoNotification in System category")
	}
	if !systemTypes[WarningNotification] {
		t.Error("Expected WarningNotification in System category")
	}
	
	// Check Installation category
	if !installationTypes[SuccessNotification] {
		t.Error("Expected SuccessNotification in Installation category")
	}
	if !installationTypes[ErrorNotification] {
		t.Error("Expected ErrorNotification in Installation category")
	}
} 