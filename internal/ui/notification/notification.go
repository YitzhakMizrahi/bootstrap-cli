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
	Type          NotificationType
	Priority      NotificationPriority
	Message       string
	Title         string
	Duration      time.Duration
	Timestamp     time.Time
	ExpiresAt     time.Time
	Category      string // Category of the notification
	ParentCategory string // Parent category for nested categories
	Actions       []Action // Actions associated with this notification
}

// Action represents an action that can be performed on a notification
type Action struct {
	Label    string                 // Label to display for the action
	Callback func() error           // Function to call when the action is triggered
	Style    ActionStyle            // Style for the action button
	Data     map[string]interface{} // Additional data for the action
}

// ActionStyle defines the style for an action button
type ActionStyle struct {
	Color     string // Text color
	BgColor   string // Background color
	Icon      string // Icon to display
	Bold      bool   // Whether to display the text in bold
	Underline bool   // Whether to underline the text
}

// DefaultActionStyle returns the default style for action buttons
func DefaultActionStyle() ActionStyle {
	return ActionStyle{
		Color:     "\033[37m", // White
		BgColor:   "\033[44m", // Blue background
		Icon:      "🔘",
		Bold:      true,
		Underline: false,
	}
}

// PrimaryActionStyle returns the style for primary action buttons
func PrimaryActionStyle() ActionStyle {
	return ActionStyle{
		Color:     "\033[37m", // White
		BgColor:   "\033[42m", // Green background
		Icon:      "✅",
		Bold:      true,
		Underline: false,
	}
}

// SecondaryActionStyle returns the style for secondary action buttons
func SecondaryActionStyle() ActionStyle {
	return ActionStyle{
		Color:     "\033[37m", // White
		BgColor:   "\033[43m", // Yellow background
		Icon:      "ℹ️",
		Bold:      false,
		Underline: false,
	}
}

// DangerActionStyle returns the style for danger action buttons
func DangerActionStyle() ActionStyle {
	return ActionStyle{
		Color:     "\033[37m", // White
		BgColor:   "\033[41m", // Red background
		Icon:      "⚠️",
		Bold:      true,
		Underline: false,
	}
}

// CategoryStyle defines the style for a category
type CategoryStyle struct {
	Color string
	Icon  string
}

// NotificationManager manages notifications
type NotificationManager struct {
	notifications    []*Notification
	maxNotifications int
	width           int
	storagePath     string
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
	filterCategory  string // Category to filter by
	categoryStyles  map[string]CategoryStyle // Styles for categories
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
	
	// Initialize category styles with defaults
	categoryStyles := map[string]CategoryStyle{
		"General": {
			Color: "\033[37m", // White
			Icon:  "📋",
		},
		"System": {
			Color: "\033[36m", // Cyan
			Icon:  "🖥️",
		},
		"Installation": {
			Color: "\033[32m", // Green
			Icon:  "📦",
		},
		"Security": {
			Color: "\033[31m", // Red
			Icon:  "🔒",
		},
		"Update": {
			Color: "\033[33m", // Yellow
			Icon:  "🔄",
		},
		"Error": {
			Color: "\033[31m", // Red
			Icon:  "❌",
		},
	}
	
	manager := &NotificationManager{
		notifications:    []*Notification{},
		maxNotifications: 5,
		width:           80,
		storagePath:     storagePath,
		cleanupInterval: 1 * time.Minute,
		stopCleanup:     make(chan struct{}),
		filterCategory:  "", // No filter by default
		categoryStyles:  categoryStyles,
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

// SetFilterCategory sets the category to filter by
func (m *NotificationManager) SetFilterCategory(category string) {
	m.filterCategory = category
}

// ClearFilter clears the category filter
func (m *NotificationManager) ClearFilter() {
	m.filterCategory = ""
}

// GetFilteredNotifications returns notifications filtered by the current filter
func (m *NotificationManager) GetFilteredNotifications() []*Notification {
	if m.filterCategory == "" {
		return m.notifications
	}
	
	var filteredNotifications []*Notification
	for _, notification := range m.notifications {
		if notification.Category == m.filterCategory {
			filteredNotifications = append(filteredNotifications, notification)
		}
	}
	
	return filteredNotifications
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

// Display shows all notifications in a hierarchical view
func (nm *NotificationManager) Display() string {
	if len(nm.notifications) == 0 {
		return ""
	}

	var sb strings.Builder
	
	// Get all parent categories and their child categories
	nestedCategories := nm.getNestedCategories()
	
	// Get all top-level categories (those without a parent)
	topLevelCategories := make(map[string]bool)
	for _, notification := range nm.notifications {
		if notification.ParentCategory == "" {
			topLevelCategories[notification.Category] = true
		}
	}
	
	// Sort top-level categories alphabetically
	categories := make([]string, 0, len(topLevelCategories))
	for category := range topLevelCategories {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	
	// Display top-level categories and their notifications
	for _, category := range categories {
		// Get style for this category
		style := nm.GetCategoryStyle(category)
		
		// Display category header with style
		sb.WriteString(fmt.Sprintf("\n%s%s%s:\n", style.Color, style.Icon, category))
		
		// Get notifications for this category
		notifications := nm.GetNotificationsByCategory(category)
		
		// Sort notifications by priority and timestamp
		sort.Slice(notifications, func(i, j int) bool {
			if notifications[i].Priority != notifications[j].Priority {
				return notifications[i].Priority > notifications[j].Priority
			}
			return notifications[i].Timestamp.After(notifications[j].Timestamp)
		})
		
		// Display notifications for this category
		for _, notification := range notifications {
			sb.WriteString(notification.String())
			sb.WriteString("\n")
		}
		
		// Display child categories if any
		if childCategories, exists := nestedCategories[category]; exists {
			// Sort child categories alphabetically
			sort.Strings(childCategories)
			
			for _, childCategory := range childCategories {
				// Get style for child category
				childStyle := nm.GetCategoryStyle(childCategory)
				
				// Display child category header with style and indentation
				sb.WriteString(fmt.Sprintf("  %s%s%s:\n", childStyle.Color, childStyle.Icon, childCategory))
				
				// Get notifications for this child category
				childNotifications := nm.GetNotificationsByCategory(childCategory)
				
				// Filter to only those with this parent
				var filteredChildNotifications []*Notification
				for _, notification := range childNotifications {
					if notification.ParentCategory == category {
						filteredChildNotifications = append(filteredChildNotifications, notification)
					}
				}
				
				// Sort child notifications by priority and timestamp
				sort.Slice(filteredChildNotifications, func(i, j int) bool {
					if filteredChildNotifications[i].Priority != filteredChildNotifications[j].Priority {
						return filteredChildNotifications[i].Priority > filteredChildNotifications[j].Priority
					}
					return filteredChildNotifications[i].Timestamp.After(filteredChildNotifications[j].Timestamp)
				})
				
				// Display child notifications with indentation
				for _, notification := range filteredChildNotifications {
					sb.WriteString("    ")
					sb.WriteString(notification.String())
					sb.WriteString("\n")
				}
			}
		}
	}

	return sb.String()
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
	groups := make(map[string][]*Notification)
	
	// First, group by parent category if it exists
	for _, notification := range m.notifications {
		if notification.ParentCategory != "" {
			groups[notification.ParentCategory] = append(groups[notification.ParentCategory], notification)
		} else {
			groups[notification.Category] = append(groups[notification.Category], notification)
		}
	}
	
	return groups
}

// getNestedCategories returns a map of parent categories to their child categories
func (m *NotificationManager) getNestedCategories() map[string][]string {
	nestedCategories := make(map[string][]string)
	
	// Collect all parent-child relationships
	for _, notification := range m.notifications {
		if notification.ParentCategory != "" {
			// Check if this child category is already in the list
			childExists := false
			for _, child := range nestedCategories[notification.ParentCategory] {
				if child == notification.Category {
					childExists = true
					break
				}
			}
			
			// Add the child category if it's not already in the list
			if !childExists {
				nestedCategories[notification.ParentCategory] = append(nestedCategories[notification.ParentCategory], notification.Category)
			}
		}
	}
	
	return nestedCategories
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

// SetCategoryStyle sets the style for a category
func (m *NotificationManager) SetCategoryStyle(category string, color string, icon string) {
	m.categoryStyles[category] = CategoryStyle{
		Color: color,
		Icon:  icon,
	}
}

// GetCategoryStyle gets the style for a category
func (m *NotificationManager) GetCategoryStyle(category string) CategoryStyle {
	style, exists := m.categoryStyles[category]
	if !exists {
		// Return default style if category doesn't have a custom style
		return CategoryStyle{
			Color: "\033[37m", // White
			Icon:  "📋",
		}
	}
	return style
}

// String returns a string representation of the notification
func (n *Notification) String() string {
	var sb strings.Builder
	
	// Add priority indicator
	switch n.Priority {
	case CriticalPriority:
		sb.WriteString("🔴 ")
	case HighPriority:
		sb.WriteString("🟡 ")
	case NormalPriority:
		sb.WriteString("🟢 ")
	case LowPriority:
		sb.WriteString("⚪ ")
	}

	// Add type indicator
	switch n.Type {
	case InfoNotification:
		sb.WriteString("ℹ️  ")
	case SuccessNotification:
		sb.WriteString("✅ ")
	case WarningNotification:
		sb.WriteString("⚠️  ")
	case ErrorNotification:
		sb.WriteString("❌ ")
	}

	// Add message
	sb.WriteString(n.Message)

	// Add timestamp if available
	if !n.Timestamp.IsZero() {
		sb.WriteString(fmt.Sprintf(" (%s)", n.Timestamp.Format("15:04:05")))
	}
	
	// Add actions if available
	if len(n.Actions) > 0 {
		sb.WriteString("\n    Actions: ")
		for i, action := range n.Actions {
			if i > 0 {
				sb.WriteString(" | ")
			}
			
			// Apply action style
			style := action.Style
			if style.Color == "" {
				style = DefaultActionStyle()
			}
			
			// Add icon
			if style.Icon != "" {
				sb.WriteString(style.Icon)
				sb.WriteString(" ")
			}
			
			// Add label with styling
			sb.WriteString(style.BgColor)
			sb.WriteString(style.Color)
			if style.Bold {
				sb.WriteString("\033[1m")
			}
			if style.Underline {
				sb.WriteString("\033[4m")
			}
			sb.WriteString(action.Label)
			sb.WriteString("\033[0m") // Reset all styles
		}
	}

	return sb.String()
}

// AddNotificationWithNestedCategory adds a notification with a nested category structure
func (m *NotificationManager) AddNotificationWithNestedCategory(notificationType NotificationType, message string, title string, category string, parentCategory string) {
	m.AddNotification(&Notification{
		Type:          notificationType,
		Message:       message,
		Title:         title,
		Category:      category,
		ParentCategory: parentCategory,
	})
}

// AddNotificationWithNestedCategoryAndPriority adds a notification with a nested category and priority
func (m *NotificationManager) AddNotificationWithNestedCategoryAndPriority(notificationType NotificationType, priority NotificationPriority, message string, title string, category string, parentCategory string) {
	m.AddNotification(&Notification{
		Type:          notificationType,
		Priority:      priority,
		Message:       message,
		Title:         title,
		Category:      category,
		ParentCategory: parentCategory,
	})
}

// AddNotificationWithNestedCategoryAndDuration adds a notification with a nested category and duration
func (m *NotificationManager) AddNotificationWithNestedCategoryAndDuration(notificationType NotificationType, message string, title string, category string, parentCategory string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:          notificationType,
		Message:       message,
		Title:         title,
		Category:      category,
		ParentCategory: parentCategory,
		Duration:      duration,
	})
}

// AddNotificationWithNestedCategoryPriorityAndDuration adds a notification with a nested category, priority, and duration
func (m *NotificationManager) AddNotificationWithNestedCategoryPriorityAndDuration(notificationType NotificationType, priority NotificationPriority, message string, title string, category string, parentCategory string, duration time.Duration) {
	m.AddNotification(&Notification{
		Type:          notificationType,
		Priority:      priority,
		Message:       message,
		Title:         title,
		Category:      category,
		ParentCategory: parentCategory,
		Duration:      duration,
	})
}

// AddAction adds an action to a notification
func (n *Notification) AddAction(label string, callback func() error) {
	n.Actions = append(n.Actions, Action{
		Label:    label,
		Callback: callback,
		Style:    DefaultActionStyle(),
		Data:     make(map[string]interface{}),
	})
}

// AddActionWithStyle adds an action with a specific style to a notification
func (n *Notification) AddActionWithStyle(label string, callback func() error, style ActionStyle) {
	n.Actions = append(n.Actions, Action{
		Label:    label,
		Callback: callback,
		Style:    style,
		Data:     make(map[string]interface{}),
	})
}

// AddActionWithData adds an action with additional data to a notification
func (n *Notification) AddActionWithData(label string, callback func() error, data map[string]interface{}) {
	n.Actions = append(n.Actions, Action{
		Label:    label,
		Callback: callback,
		Style:    DefaultActionStyle(),
		Data:     data,
	})
}

// AddActionWithStyleAndData adds an action with a specific style and additional data to a notification
func (n *Notification) AddActionWithStyleAndData(label string, callback func() error, style ActionStyle, data map[string]interface{}) {
	n.Actions = append(n.Actions, Action{
		Label:    label,
		Callback: callback,
		Style:    style,
		Data:     data,
	})
}

// ExecuteAction executes the action at the specified index
func (n *Notification) ExecuteAction(index int) error {
	if index < 0 || index >= len(n.Actions) {
		return fmt.Errorf("action index out of range: %d", index)
	}
	
	return n.Actions[index].Callback()
} 