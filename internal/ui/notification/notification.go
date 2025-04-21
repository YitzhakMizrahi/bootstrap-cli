package notification

import (
	"fmt"
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

// Notification represents a notification message
type Notification struct {
	Type      NotificationType
	Message   string
	Title     string
	Duration  time.Duration
	Timestamp time.Time
}

// NotificationManager manages notifications
type NotificationManager struct {
	notifications []*Notification
	maxNotifications int
	width          int
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		notifications:    []*Notification{},
		maxNotifications: 5,
		width:           80,
	}
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
	
	// Add notification to the list
	m.notifications = append(m.notifications, notification)
	
	// Trim notifications if exceeding max
	if len(m.notifications) > m.maxNotifications {
		m.notifications = m.notifications[len(m.notifications)-m.maxNotifications:]
	}
	
	// Display the notification
	m.displayNotification(notification)
}

// Display displays all notifications
func (m *NotificationManager) Display() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")
	
	// Display each notification
	for _, notification := range m.notifications {
		m.displayNotification(notification)
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
}

// Info displays an info notification
func (m *NotificationManager) Info(message string, title string) {
	m.AddNotification(&Notification{
		Type:    InfoNotification,
		Message: message,
		Title:   title,
	})
}

// Success displays a success notification
func (m *NotificationManager) Success(message string, title string) {
	m.AddNotification(&Notification{
		Type:    SuccessNotification,
		Message: message,
		Title:   title,
	})
}

// Warning displays a warning notification
func (m *NotificationManager) Warning(message string, title string) {
	m.AddNotification(&Notification{
		Type:    WarningNotification,
		Message: message,
		Title:   title,
	})
}

// Error displays an error notification
func (m *NotificationManager) Error(message string, title string) {
	m.AddNotification(&Notification{
		Type:    ErrorNotification,
		Message: message,
		Title:   title,
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