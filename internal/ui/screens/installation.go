package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/pipeline"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Task State ---

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusRetrying
	StatusRollingBack
	StatusDone
	StatusFailed
	StatusRollbackFailed
)

type TaskState struct {
	ID          string
	Description string
	Status      TaskStatus
	Progress    float64 // 0.0 to 1.0 for progress bar
	Error       error
	StartTime   time.Time
	EndTime     time.Time
}

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

	spinner spinner.Model // Spinner for active tasks
	
	// State for display
	tasks      []*TaskState       
	taskMap    map[string]*TaskState 
	progresses map[string]*progress.Model // Store pointers to progress models
	activeTaskCount int
	logMessages []string // Simple log for now
	// TODO: Add more structured state later (e.g., map[taskID]taskState for progress bars)
}

func NewInstallationScreen(progChan <-chan pipeline.ProgressEvent) *InstallationScreen {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = styles.InfoStyle // Use an accent color for the spinner

	return &InstallationScreen{
		title:       "Installation Progress",
		progressChan: progChan,
		taskMap:      make(map[string]*TaskState),
		tasks:        make([]*TaskState, 0),
		progresses: make(map[string]*progress.Model), // Initialize map for pointers
		spinner:    sp,
	}
}

// --- Bubble Tea Interface --- 

func (s *InstallationScreen) Init() tea.Cmd {
	// Start listening and start the spinner ticking
	return tea.Batch(s.listenForProgress(), s.spinner.Tick)
}

func (s *InstallationScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		// Update progress bar widths
		for _, p := range s.progresses { // Iterate through pointers
			// Assuming progress bar width is related to overall width minus some padding
			newWidth := s.width - 10 // Example width calculation
			if newWidth < 10 { newWidth = 10 }
			p.Width = newWidth // Modify width directly on the pointer
		}
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Only allow exit via keypress if installation is actually finished
			if s.finished {
				return s, tea.Quit
			}
			// TODO: Implement cancellation signal to pipeline?
			// For now, don't quit if not finished.
			return s, nil 
		}

	// Handle spinner tick if installation is ongoing
	case spinner.TickMsg:
		var cmd tea.Cmd
		if !s.finished && s.activeTaskCount > 0 {
			s.spinner, cmd = s.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		return s, tea.Batch(cmds...)

	// Handle progress bar animation
	case progress.FrameMsg:
		var updatedCmds []tea.Cmd
		for _, p := range s.progresses { // p is *progress.Model, Use _ for id
			if p.Percent() < 1.0 {
				_, cmd := p.Update(msg) // Update returns progress.Model, use _ for newModel
				// We need to update the model pointed to, but Update returns a value.
				// This pattern with bubbles is often tricky. Let's assume Update modifies in place
				// or we handle it differently. Reverting to simpler update for now.
				// If Update doesn't modify in place, we'd need to create a new model and update the pointer:
				// if m, ok := newModel.(progress.Model); ok { *p = m } 
				// For now, just append the command if one is returned.
				if cmd != nil {
					updatedCmds = append(updatedCmds, cmd)
				}
			}
		}
		return s, tea.Batch(updatedCmds...)

	// Handle messages from the progress channel listener
	case progressMsg:
		var cmdsToBatch []tea.Cmd
		
		switch event := msg.event.(type) {
		case pipeline.TaskStart:
			// Add new task to state
			newTask := &TaskState{
				ID:          event.TaskID,
				Description: event.Description,
				Status:      StatusRunning,
				StartTime:   time.Now(),
				Progress:    -1, // Indeterminate initially
			}
			s.tasks = append(s.tasks, newTask)
			s.taskMap[event.TaskID] = newTask
			s.activeTaskCount++
			// Potentially create a progress bar if needed later

		case pipeline.TaskProgress:
			if task, ok := s.taskMap[event.TaskID]; ok {
				task.Progress = event.Percent / 100.0 // Convert percentage to 0.0-1.0
				if task.Progress < 0 { task.Progress = 0 } // Clamp if indeterminate was sent
				if task.Progress > 1 { task.Progress = 1 } 

				// Initialize or update progress bar
				p, pOk := s.progresses[event.TaskID] // p is now *progress.Model or nil
				if !pOk {
					progWidth := s.width - 10 
					if progWidth < 10 { progWidth = 10 }
					newProgress := progress.New(progress.WithDefaultGradient())
					newProgress.Width = progWidth
					s.progresses[event.TaskID] = &newProgress // Store pointer
					p = &newProgress // Use the new pointer
				}
				// Send command to update the progress bar animation
				// Call SetPercent directly on the pointer
				cmdsToBatch = append(cmdsToBatch, p.SetPercent(task.Progress))
			}

		case pipeline.TaskLog:
			// Simple log for now - append to a shared log or task-specific?
			// Append to general log for now
			s.logMessages = append(s.logMessages, fmt.Sprintf("LOG [%s]: %s", event.TaskID, event.Line))

		case pipeline.TaskEnd:
			if task, ok := s.taskMap[event.TaskID]; ok {
				task.EndTime = time.Now()
				task.Error = event.Error
				if event.Success {
					task.Status = StatusDone
					task.Progress = 1.0 // Ensure progress bar is full on success
					if p, pOk := s.progresses[event.TaskID]; pOk {
						cmdsToBatch = append(cmdsToBatch, p.SetPercent(1.0))
					}
				} else {
					// Distinguish between normal fail and rollback fail?
					if strings.HasSuffix(task.ID, "-rollback") {
						task.Status = StatusRollbackFailed
					} else {
						task.Status = StatusFailed
					}
				}
				s.activeTaskCount--
				if s.activeTaskCount < 0 { s.activeTaskCount = 0 }
			}

		case pipeline.PipelineComplete:
			s.finished = true
			s.success = event.OverallSuccess
			s.finalError = event.FinalError
			s.activeTaskCount = 0 // Ensure counter is zero
			// Stop listening implicitly as channel will close
			return s, nil // Wait for user to press Enter/q to Quit
		}

		// After processing a progress message, continue listening if not complete
		if !s.finished {
		    cmdsToBatch = append(cmdsToBatch, s.listenForProgress()) 
		}
		return s, tea.Batch(cmdsToBatch...)

	case errorMsg:
		s.finished = true
		s.finalError = msg.err
		s.logMessages = append(s.logMessages, styles.ErrorStyle.Render(fmt.Sprintf("Error listening for progress: %v", msg.err)))
		return s, tea.Quit // Quit on listener error
	}

	// Also handle spinner ticks if no other message consumed it
	var spinnerCmd tea.Cmd
	s.spinner, spinnerCmd = s.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return s, tea.Batch(cmds...)
}

func (s *InstallationScreen) View() string {
	if s.width == 0 { // Avoid rendering before size is known
		return "Initializing..."
	}
	var content strings.Builder

	// Title
	content.WriteString(styles.TitleStyle.Render(s.title))
	content.WriteString("\n\n")

	// Display Tasks
	for _, task := range s.tasks {
		var line strings.Builder

		// Status Indicator
		switch task.Status {
		case StatusRunning, StatusRetrying, StatusRollingBack:
			line.WriteString(s.spinner.View() + " ")
		case StatusDone:
			line.WriteString(styles.SuccessStyle.Render("✓") + " ")
		case StatusFailed, StatusRollbackFailed:
			line.WriteString(styles.ErrorStyle.Render("✗") + " ")
		default: // Pending
			line.WriteString(styles.UnselectedTextStyle.Render("·") + " ") // Use UnselectedTextStyle
		}

		// Description
		desc := task.Description
		if task.Status == StatusRetrying {
			desc += " (Retrying...)"
		} else if task.Status == StatusRollingBack {
			desc += " (Rolling back...)"
		}
		line.WriteString(styles.NormalTextStyle.Render(desc))

		// Progress Bar (if applicable)
		if prog, ok := s.progresses[task.ID]; ok && task.Status != StatusDone && task.Status != StatusFailed && task.Status != StatusRollbackFailed {
			line.WriteString("\n  ") // Indent progress bar
			line.WriteString(prog.ViewAs(task.Progress))
		}

		// Error Message (if applicable)
		if task.Error != nil && (task.Status == StatusFailed || task.Status == StatusRollbackFailed) {
			line.WriteString("\n  ") // Indent error
			errorMsg := styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", task.Error))
			// Wrap error message if too long
			errorMsg = lipgloss.NewStyle().Width(s.width - 4).Render(errorMsg) // Adjust width as needed
			line.WriteString(errorMsg)
		}

		content.WriteString(line.String())
		content.WriteString("\n\n") // Add extra space between tasks
	}

	// Footer / Final Status
	footer := "\n"
	if s.finished {
		if s.success {
			footer += styles.SuccessStyle.Render("Installation Complete!")
		} else {
			footer += styles.ErrorStyle.Render(fmt.Sprintf("Installation Failed: %v", s.finalError))
		}
        footer += "\nPress Enter or q to exit."
	} else if s.activeTaskCount > 0 {
		footer += styles.HelpStyle.Render(fmt.Sprintf("Installation in progress (%d active)... (Press Ctrl+C to attempt cancel)", s.activeTaskCount))
	} else {
        footer += styles.HelpStyle.Render("Waiting for pipeline...") // Should not stay here long
    }

    // Combine content and footer, considering height limits
    availableHeight := s.height - lipgloss.Height(footer) - 2 // Account for footer and spacing
    if availableHeight < 0 { availableHeight = 0 }
    
    finalContent := lipgloss.NewStyle().MaxHeight(availableHeight).Render(content.String())
    
	// Use lipgloss.Place for centering, assuming full screen height/width is passed
	// Or just return with AppStyle if nested
	// return styles.AppStyle.Render(finalContent + footer)
	return lipgloss.Place(s.width, s.height, lipgloss.Left, lipgloss.Top, finalContent+footer)
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