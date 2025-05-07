package screens

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/pipeline"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Messages for internal screen updates ---

// progressMsg wraps a ProgressEvent coming from the pipeline channel
type progressMsg struct{
	event pipeline.ProgressEvent
}

// errorMsg indicates an error reading from the progress channel
type errorMsg struct{
	err error
}

// --- Model --- 

type InstallationScreen struct {
	title       string
	progressChan <-chan pipeline.ProgressEvent // Channel to receive events
	width       int
	height      int
	finished    bool
	finalError  error
	success     bool

	// State for display
	logMessages []string // Simple log for now
	// TODO: Add more structured state later (e.g., map[taskID]taskState for progress bars)
}

func NewInstallationScreen(progChan <-chan pipeline.ProgressEvent) *InstallationScreen {
	return &InstallationScreen{
		title:       "Installation Progress",
		progressChan: progChan,
		logMessages:  []string{"Starting installation process..."},
	}
}

// --- Bubble Tea Interface --- 

func (s *InstallationScreen) Init() tea.Cmd {
	// Start listening to the progress channel
	return s.listenForProgress()
}

func (s *InstallationScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// TODO: Implement cancellation signal to pipeline?
			s.finished = true // Mark as finished locally for now
			s.finalError = fmt.Errorf("installation cancelled by user")
			return s, tea.Quit // For now, just quit the UI
		}

	// Handle messages from the progress channel listener
	case progressMsg:
		// Format the event to string for logging
		var eventString string
		switch event := msg.event.(type) {
		case pipeline.TaskStart:
			eventString = event.String()
		case pipeline.TaskProgress:
			eventString = event.String()
		case pipeline.TaskLog:
			eventString = event.String()
		case pipeline.TaskEnd:
			eventString = event.String()
		case pipeline.PipelineComplete:
			eventString = event.String()
		default:
			eventString = fmt.Sprintf("Received unknown progress event type: %T", msg.event)
		}
		s.logMessages = append(s.logMessages, eventString) // Add string representation to log
		
		// Handle specific event types for state changes
		switch event := msg.event.(type) {
		case pipeline.PipelineComplete:
			s.finished = true
			s.success = event.OverallSuccess
			s.finalError = event.FinalError
			// No command needed, listener goroutine will exit because channel is closed
            // return s, tea.Quit // Quit the app when pipeline completes - Let user press Enter/q now
            return s, nil // Stop listening, wait for user exit
        // TODO: Add cases for TaskStart, TaskEnd, TaskProgress to update more structured state later
		}
		// After processing a progress message, continue listening if not complete
		if !s.finished {
		    return s, s.listenForProgress()
		} 
		return s, nil

	case errorMsg:
		s.finished = true
		s.finalError = msg.err
		s.logMessages = append(s.logMessages, styles.ErrorStyle.Render(fmt.Sprintf("Error listening for progress: %v", msg.err)))
		return s, tea.Quit // Quit on listener error
	}

	return s, nil
}

func (s *InstallationScreen) View() string {
	var content strings.Builder

	// Title
	content.WriteString(styles.TitleStyle.Render(s.title))
	content.WriteString("\n\n")

	// Display Log Messages (Simple view for now)
	maxLogLines := s.height - 6 // Adjust based on title/footer height
	if maxLogLines < 1 { maxLogLines = 1 }
	startIdx := 0
	if len(s.logMessages) > maxLogLines {
		startIdx = len(s.logMessages) - maxLogLines
	}
	for _, log := range s.logMessages[startIdx:] {
		content.WriteString(log)
		content.WriteString("\n")
	}

	// Footer / Final Status
	content.WriteString("\n")
	if s.finished {
		if s.success {
			content.WriteString(styles.SuccessStyle.Render("Installation Complete!"))
		} else {
			content.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Installation Failed: %v", s.finalError)))
		}
        content.WriteString("\nPress Enter or q to exit.") // Inform user how to exit now
	} else {
		content.WriteString(styles.HelpStyle.Render("Installation in progress... (Press Ctrl+C to attempt cancel)"))
	}

	// Use lipgloss.Place for centering, assuming full screen height/width is passed
	return lipgloss.Place(s.width, s.height, lipgloss.Left, lipgloss.Top, content.String())
	// Or apply AppStyle border if this screen is nested within app.Model
	// return styles.AppStyle.Render(content.String())
}

// listenForProgress returns a command that listens for the next message on the progress channel.
func (s *InstallationScreen) listenForProgress() tea.Cmd {
	return func() tea.Msg {
		event, ok := <-s.progressChan
		if !ok {
			// Channel closed, usually means pipeline finished successfully or failed critically
			// The PipelineComplete event should have already been sent before close.
			// We might not need a specific message here, but could return one if needed.
            // For safety, return an error indicating unexpected close if finished flag isn't set.
            if !s.finished {
                 return errorMsg{fmt.Errorf("progress channel closed unexpectedly")}
            }
            return nil // Normal closure
		}
		return progressMsg{event}
	}
} 