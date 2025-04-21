package notification

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// NotificationType defines the type of notification
type NotificationType int

const (
	// InfoNotification is an informational notification
	InfoNotification NotificationType = iota
	// SuccessNotification is a success notification
	SuccessNotification
	// WarningNotification is a warning notification
	WarningNotification
	// ErrorNotification is an error notification
	ErrorNotification
)

// NotificationPriority defines the priority of a notification
type NotificationPriority int

const (
	// LowPriority is a low priority notification
	LowPriority NotificationPriority = iota
	// NormalPriority is a normal priority notification
	NormalPriority
	// HighPriority is a high priority notification
	HighPriority
	// CriticalPriority is a critical priority notification
	CriticalPriority
)

// Notification represents a notification message
type Notification struct {
	Type      NotificationType
	Priority  NotificationPriority
	Message   string
	Title     string
	Duration  time.Duration
	Timestamp time.Time
	ExpiresAt time.Time
	Category  string // Category of the notification
}

// NotificationManager manages notifications
type NotificationManager struct {
	notifications    []*Notification
	maxNotifications int
	width           int
	storagePath     string
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager() *NotificationManager {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	
	// Create storage path in user's home directory
	storagePath := filepath.Join(homeDir, ".bootstrap-cli", "notifications.json")
	
	manager := &NotificationManager{
		notifications:    []*Notification{},
		maxNotifications: 5,
		width:           80,
		storagePath:     storagePath,
		cleanupInterval: 1 * time.Minute,
		stopCleanup:     make(chan struct{}),
	}
	
	// Start cleanup goroutine
	go manager.startCleanup()
	
	return manager
}

// SetMaxNotifications sets the maximum number of notifications to display
func (m *NotificationManager) SetMaxNotifications(max int) {
	m.maxNotifications = max
}

// SetWidth sets the width of the notifications
func (m *NotificationManager) SetWidth(width int) {
	m.width = width
}

// AddNotification adds a notification to the manager
func (m *NotificationManager) AddNotification(notification *Notification) {
	// Set timestamp if not set
	if notification.Timestamp.IsZero() {
		notification.Timestamp = time.Now()
	}
	
	// Set expiration time if duration is set
	if notification.Duration > 0 {
		notification.ExpiresAt = notification.Timestamp.Add(notification.Duration)
	}
	
	// Set default priority if not set
	if notification.Priority == 0 {
		// Set priority based on notification type
		switch notification.Type {
		case InfoNotification:
			notification.Priority = NormalPriority
		case SuccessNotification:
			notification.Priority = NormalPriority
		case WarningNotification:
			notification.Priority = HighPriority
		case ErrorNotification:
			notification.Priority = CriticalPriority
		default:
			notification.Priority = NormalPriority
		}
	}
	
	// Set default category if not set
	if notification.Category == "" {
		notification.Category = "General"
	}
	
	// Add notification to the list
	m.notifications = append(m.notifications, notification)
	
	// Sort notifications by priority (highest first)
	m.sortNotifications()
	
	// Trim notifications if exceeding max
	if len(m.notifications) > m.maxNotifications {
		m.notifications = m.notifications[len(m.notifications)-m.maxNotifications:]
	}
	
	// Save notifications to file
	m.saveNotifications()
	
	// Display the notification
	m.displayNotification(notification)
}

// Display displays all notifications
func (m *NotificationManager) Display() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")
	
	// Group notifications by category
	notificationsByCategory := m.groupNotificationsByCategory()
	
	// Display notifications by category
	for category, notifications := range notificationsByCategory {
		// Display category header
		fmt.Printf("\033[1m%s\033[0m\n", category)
		
		// Display each notification in the category
		for _, notification := range notifications {
			m.displayNotification(notification)
		}
		
		// Add a separator between categories
		fmt.Println()
	}
}

// displayNotification displays a single notification
func (m *NotificationManager) displayNotification(notification *Notification) {
	// Get the color and icon based on the notification type
	color, icon := m.getStyleForType(notification.Type)
	
	// Create the notification box
	box := m.createNotificationBox(notification, color, icon)
	
	// Display the notification
	fmt.Println(box)
}

// getStyleForType returns the color and icon for a notification type
func (m *NotificationManager) getStyleForType(notificationType NotificationType) (string, string) {
	switch notificationType {
	case InfoNotification:
		return "\033[34m", "ℹ" // Blue
	case SuccessNotification:
		return "\033[32m", "✓" // Green
	case WarningNotification:
		return "\033[33m", "⚠" // Yellow
	case ErrorNotification:
		return "\033[31m", "✗" // Red
	default:
		return "\033[0m", "?"
	}
}

// createNotificationBox creates a box for a notification
func (m *NotificationManager) createNotificationBox(notification *Notification, color string, icon string) string {
	// Reset color
	reset := "\033[0m"
	
	// Create the title line
	title := notification.Title
	if title == "" {
		title = "Notification"
	}
	
	// Create the message lines
	messageLines := strings.Split(notification.Message, "\n")
	
	// Calculate the box width
	boxWidth := m.width
	
	// Create the top border
	topBorder := color + "┌" + strings.Repeat("─", boxWidth-2) + "┐" + reset + "\n"
	
	// Create the title line
	titleLine := color + "│ " + reset + color + icon + " " + title + reset
	titleLine += strings.Repeat(" ", boxWidth-len(title)-5)
	titleLine += color + "│" + reset + "\n"
	
	// Create the separator line
	separatorLine := color + "├" + strings.Repeat("─", boxWidth-2) + "┤" + reset + "\n"
	
	// Create the message lines
	messageBox := ""
	for _, line := range messageLines {
		// Truncate line if too long
		if len(line) > boxWidth-4 {
			line = line[:boxWidth-4]
		}
		
		// Pad line if too short
		paddedLine := line + strings.Repeat(" ", boxWidth-len(line)-4)
		
		messageBox += color + "│ " + reset + paddedLine + color + " │" + reset + "\n"
	}
	
	// Create the bottom border
	bottomBorder := color + "└" + strings.Repeat("─", boxWidth-2) + "┘" + reset
	
	// Combine all parts
	return topBorder + titleLine + separatorLine + messageBox + bottomBorder
}

// Clear clears all notifications
func (m *NotificationManager) Clear() {
	m.notifications = []*Notification{}
	m.saveNotifications()
}

// Info displays an info notification
func (m *NotificationManager) Info(message string, title string) {
	m.AddNotification(&Notification{
		Type:     InfoNotification,
		Priority: NormalPriority,
		Message:  message,
		Title:    title,
		Category: "General",
	})
}

// Success displays a success notification
func (m *NotificationManager) Success(message string, title string) {
	m.AddNotification(&Notification{
		Type:     SuccessNotification,
		Priority: NormalPriority,
		Message:  message,
		Title:    title,
		Category: "General",
	})
}

// Warning displays a warning notification
func (m *NotificationManager) Warning(message string, title string) {
	m.AddNotification(&Notification{
		Type:     WarningNotification,
		Priority: HighPriority,
		Message:  message,
		Title:    title,
		Category: "General",
	})
}

// Error displays an error notification
func (m *NotificationManager) Error(message string, title string) {
	m.AddNotification(&Notification{
		Type:     ErrorNotification,
		Priority: CriticalPriority,
		Message:  message,
		Title:    title,
		Category: "General",
	})
}

// Notify is a convenience function to display a notification
func Notify(notificationType NotificationType, message string, title string) {
	manager := NewNotificationManager()
	manager.AddNotification(&Notification{
		Type:    notificationType,
		Message: message,
		Title:   title,
	})
}

// saveNotifications saves notifications to a file
func (m *NotificationManager) saveNotifications() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(m.storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Marshal notifications to JSON
	data, err := json.MarshalIndent(m.notifications, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notifications: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(m.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notifications to file: %w", err)
	}
	
	return nil
}

// LoadNotifications loads notifications from a file
func (m *NotificationManager) LoadNotifications() error {
	// Check if file exists
	if _, err := os.Stat(m.storagePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to load
	}
	
	// Read file
	data, err := os.ReadFile(m.storagePath)
	if err != nil {
		return fmt.Errorf("failed to read notifications file: %w", err)
	}
	
	// Unmarshal JSON
	var notifications []*Notification
	if err := json.Unmarshal(data, &notifications); err != nil {
		return fmt.Errorf("failed to unmarshal notifications: %w", err)
	}
	
	// Set notifications
	m.notifications = notifications
	
	return nil
}

// startCleanup starts the cleanup goroutine
func (m *NotificationManager) startCleanup() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			m.cleanupExpiredNotifications()
		case <-m.stopCleanup:
			return
		}
	}
}

// cleanupExpiredNotifications removes expired notifications
func (m *NotificationManager) cleanupExpiredNotifications() {
	now := time.Now()
	expired := false
	
	// Create a new slice to hold non-expired notifications
	var newNotifications []*Notification
	
	// Check each notification
	for _, notification := range m.notifications {
		// Skip notifications without expiration
		if notification.ExpiresAt.IsZero() {
			newNotifications = append(newNotifications, notification)
			continue
		}
		
		// Check if notification has expired
		if now.After(notification.ExpiresAt) {
			expired = true
			continue
		}
		
		// Keep non-expired notification
		newNotifications = append(newNotifications, notification)
	}
	
	// Update notifications if any expired
	if expired {
		m.notifications = newNotifications
		m.saveNotifications()
	}
}

// StopCleanup stops the cleanup goroutine
func (m *NotificationManager) StopCleanup() {
	close(m.stopCleanup)
}

// InfoWithDuration displays an info notification with a duration
func (m *NotificationManager) InfoWithDuration(message string, title string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:     InfoNotification,
		Priority: NormalPriority,
		Message:  message,
		Title:    title,
		Duration: duration,
	})
}

// SuccessWithDuration displays a success notification with a duration
func (m *NotificationManager) SuccessWithDuration(message string, title string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:     SuccessNotification,
		Priority: NormalPriority,
		Message:  message,
		Title:    title,
		Duration: duration,
	})
}

// WarningWithDuration displays a warning notification with a duration
func (m *NotificationManager) WarningWithDuration(message string, title string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:     WarningNotification,
		Priority: HighPriority,
		Message:  message,
		Title:    title,
		Duration: duration,
	})
}

// ErrorWithDuration displays an error notification with a duration
func (m *NotificationManager) ErrorWithDuration(message string, title string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:     ErrorNotification,
		Priority: CriticalPriority,
		Message:  message,
		Title:    title,
		Duration: duration,
	})
}

// sortNotifications sorts notifications by priority (highest first)
func (m *NotificationManager) sortNotifications() {
	// Sort notifications by priority (highest first)
	sort.Slice(m.notifications, func(i, j int) bool {
		// First sort by priority (highest first)
		if m.notifications[i].Priority != m.notifications[j].Priority {
			return m.notifications[i].Priority > m.notifications[j].Priority
		}
		
		// Then sort by timestamp (newest first)
		return m.notifications[i].Timestamp.After(m.notifications[j].Timestamp)
	})
}

// AddNotificationWithPriority adds a notification with a specific priority
func (m *NotificationManager) AddNotificationWithPriority(notificationType NotificationType, priority NotificationPriority, message string, title string) {
	m.AddNotification(&Notification{
		Type:     notificationType,
		Priority: priority,
		Message:  message,
		Title:    title,
	})
}

// AddNotificationWithPriorityAndDuration adds a notification with a specific priority and duration
func (m *NotificationManager) AddNotificationWithPriorityAndDuration(notificationType NotificationType, priority NotificationPriority, message string, title string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:     notificationType,
		Priority: priority,
		Message:  message,
		Title:    title,
		Duration: duration,
	})
}

// groupNotificationsByCategory groups notifications by their category
func (m *NotificationManager) groupNotificationsByCategory() map[string][]*Notification {
	// Create a map to store notifications by category
	notificationsByCategory := make(map[string][]*Notification)
	
	// Group notifications by category
	for _, notification := range m.notifications {
		notificationsByCategory[notification.Category] = append(notificationsByCategory[notification.Category], notification)
	}
	
	return notificationsByCategory
}

// GetNotificationsByCategory returns all notifications in a specific category
func (m *NotificationManager) GetNotificationsByCategory(category string) []*Notification {
	var filteredNotifications []*Notification
	
	for _, notification := range m.notifications {
		if notification.Category == category {
			filteredNotifications = append(filteredNotifications, notification)
		}
	}
	
	return filteredNotifications
}

// GetCategories returns all unique categories
func (m *NotificationManager) GetCategories() []string {
	categories := make(map[string]bool)
	
	for _, notification := range m.notifications {
		categories[notification.Category] = true
	}
	
	// Convert map keys to slice
	var result []string
	for category := range categories {
		result = append(result, category)
	}
	
	// Sort categories alphabetically
	sort.Strings(result)
	
	return result
}

// AddNotificationWithCategory adds a notification with a specific category
func (m *NotificationManager) AddNotificationWithCategory(notificationType NotificationType, message string, title string, category string) {
	m.AddNotification(&Notification{
		Type:     notificationType,
		Message:  message,
		Title:    title,
		Category: category,
	})
}

// AddNotificationWithCategoryAndPriority adds a notification with a specific category and priority
func (m *NotificationManager) AddNotificationWithCategoryAndPriority(notificationType NotificationType, priority NotificationPriority, message string, title string, category string) {
	m.AddNotification(&Notification{
		Type:     notificationType,
		Priority: priority,
		Message:  message,
		Title:    title,
		Category: category,
	})
}

// AddNotificationWithCategoryAndDuration adds a notification with a specific category and duration
func (m *NotificationManager) AddNotificationWithCategoryAndDuration(notificationType NotificationType, message string, title string, category string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:     notificationType,
		Message:  message,
		Title:    title,
		Category: category,
		Duration: duration,
	})
}

// AddNotificationWithCategoryPriorityAndDuration adds a notification with a specific category, priority, and duration
func (m *NotificationManager) AddNotificationWithCategoryPriorityAndDuration(notificationType NotificationType, priority NotificationPriority, message string, title string, category string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:     notificationType,
		Priority: priority,
		Message:  message,
		Title:    title,
		Category: category,
		Duration: duration,
	})
} 