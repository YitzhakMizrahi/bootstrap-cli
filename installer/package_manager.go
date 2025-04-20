package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles for the package manager UI
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	currentPkgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF78D4"))
	doneStyle       = lipgloss.NewStyle().Margin(1, 2)
	successMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("#73F59F")).SetString("✓")
	errorMark       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).SetString("✗")
	skippedMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAE19")).SetString("•")
)

// PackageStatus represents the status of a package installation
type PackageStatus string

const (
	StatusPending    PackageStatus = "pending"
	StatusInstalling PackageStatus = "installing"
	StatusSuccess    PackageStatus = "success"
	StatusFailed     PackageStatus = "failed"
	StatusSkipped    PackageStatus = "skipped"
)

// Package represents a package to be installed
type Package struct {
	Name        string
	Description string
	Status      PackageStatus
	Error       error
	Version     string
	StartTime   time.Time
	EndTime     time.Time
}

// SimplePackageModel represents a simple package installation model
type SimplePackageModel struct {
	title        string
	packages     []string
	current      int
	width        int
	height       int
	spinner      spinner.Model
	progress     progress.Model
	done         bool
	results      map[string]string // package name -> status ("success", "error", "skipped")
	errors       map[string]error
	packageMgr   string
}

// NewSimplePackageModel creates a new package model
func NewSimplePackageModel(title string, packages []string, packageMgr string) SimplePackageModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return SimplePackageModel{
		title:      title,
		packages:   packages,
		spinner:    s,
		progress:   p,
		results:    make(map[string]string),
		errors:     make(map[string]error),
		packageMgr: packageMgr,
	}
}

func (m SimplePackageModel) Init() tea.Cmd {
	if len(m.packages) == 0 {
		return tea.Quit
	}
	return tea.Batch(
		m.spinner.Tick,
		m.installCurrentPackage(),
	)
}

func (m SimplePackageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		if p, ok := progressModel.(progress.Model); ok {
			m.progress = p
		}
		return m, cmd

	case packageResultMsg:
		// Record result
		m.results[msg.pkg] = msg.status
		if msg.err != nil {
			m.errors[msg.pkg] = msg.err
		}

		// Update progress based on completed packages (not including current)
		progCmd := m.progress.SetPercent(float64(m.current) / float64(len(m.packages)))
		
		// Check if we've completed all packages
		m.current++
		if m.current >= len(m.packages) {
			m.done = true
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Tick(time.Millisecond*750, func(time.Time) tea.Msg { return nil }),
				tea.Quit,
			)
		}

		// Install the next package
		return m, tea.Batch(
			progCmd,
			m.installCurrentPackage(),
		)
	}
	
	return m, nil
}

func (m SimplePackageModel) View() string {
	if m.done {
		// Summary view after completion
		var s strings.Builder
		s.WriteString(titleStyle.Render(m.title) + "\n\n")
		
		// Count successes
		successCount := 0
		for _, status := range m.results {
			if status == "success" {
				successCount++
			}
		}
		
		// Show a summary of results
		s.WriteString(fmt.Sprintf("✅ Installed %d/%d tools\n\n", successCount, len(m.packages)))
		
		// Show status for each package with its status
		for _, pkg := range m.packages {
			marker := "  "
			
			if status, ok := m.results[pkg]; ok {
				switch status {
				case "success":
					marker = successMark.String() + " "
				case "error":
					marker = errorMark.String() + " "
				case "skipped":
					marker = skippedMark.String() + " "
				}
			}
			
			s.WriteString(marker + pkg)
			
			// Add error message if there was one
			if err, hasError := m.errors[pkg]; hasError {
				errorMsg := fmt.Sprintf(" - Error: %v", err)
				s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).Render(errorMsg))
			}
			
			s.WriteString("\n")
		}
		
		return doneStyle.Render(s.String())
	}

	// Get current package name
	pkgName := ""
	if m.current < len(m.packages) {
		pkgName = m.packages[m.current]
	}
	
	// Print previously installed packages with their statuses
	var s strings.Builder
	
	// Add title
	s.WriteString(titleStyle.Render(m.title) + "\n\n")
	
	// Add progress info
	progressText := fmt.Sprintf("Installing %d/%d tools", m.current, len(m.packages))
	s.WriteString(progressText + "\n")
	s.WriteString(m.progress.View() + "\n\n")
	
	// Add completed packages
	for i := 0; i < m.current; i++ {
		pkg := m.packages[i]
		marker := "  "
		
		if status, ok := m.results[pkg]; ok {
			switch status {
			case "success":
				marker = successMark.String() + " "
			case "error":
				marker = errorMark.String() + " "
			case "skipped":
				marker = skippedMark.String() + " "
			}
		}
		
		s.WriteString(marker + pkg + "\n")
	}
	
	// Add current package with spinner
	if m.current < len(m.packages) {
		current := m.spinner.View() + " " + currentPkgStyle.Render(pkgName)
		s.WriteString(current + "\n")
	}
	
	// Add pending packages
	for i := m.current + 1; i < len(m.packages); i++ {
		s.WriteString("  " + m.packages[i] + "\n")
	}
	
	// Add help text
	s.WriteString("\nPress Ctrl+C to cancel\n")
	
	// Apply a consistent style to the whole output
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Render(s.String())
}

// installCurrentPackage creates a command to install the current package
func (m SimplePackageModel) installCurrentPackage() tea.Cmd {
	return func() tea.Msg {
		if m.current >= len(m.packages) {
			return nil
		}
		
		pkg := m.packages[m.current]
		
		// Check if already installed
		if isCommandAvailable(pkg) {
			// For apps that don't install as commands with the same name
			switch pkg {
			case "fd-find":
				// Check for fdfind
				if !isCommandAvailable("fdfind") {
					break
				}
			case "ripgrep":
				// Check for rg
				if !isCommandAvailable("rg") {
					break
				}
			}
			
			// Wait a moment for UI to show this clearly
			time.Sleep(300 * time.Millisecond)
			return packageResultMsg{pkg: pkg, status: "skipped"}
		}
		
		// Install the package
		var err error
		switch m.packageMgr {
		case "brew":
			cmd := exec.Command("brew", "install", pkg)
			cmd.Stdout = nil
			cmd.Stderr = nil
			err = cmd.Run()
		case "apt":
			cmd := exec.Command("sudo", "apt-get", "install", "-y", pkg)
			cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
			cmd.Stdout = nil
			cmd.Stderr = nil
			err = cmd.Run()
		default:
			err = fmt.Errorf("unsupported package manager: %s", m.packageMgr)
		}
		
		status := "success"
		if err != nil {
			status = "error"
		}
		
		// Add a small delay for UI
		time.Sleep(500 * time.Millisecond)
		return packageResultMsg{pkg: pkg, status: status, err: err}
	}
}

// Package result message
type packageResultMsg struct {
	pkg    string
	status string
	err    error
}

// InstallPackagesWithUI installs packages with a bubbletea UI
func InstallPackagesWithUI(title string, packages []string, packageManager string) error {
	if len(packages) == 0 {
		fmt.Println("No packages to install.")
		return nil
	}
	
	model := NewSimplePackageModel(title, packages, packageManager)
	p := tea.NewProgram(model)
	
	_, err := p.Run()
	
	// Print errors after completion
	if len(model.errors) > 0 {
		fmt.Println("\nThe following packages encountered errors:")
		for pkg, err := range model.errors {
			fmt.Printf("  - %s: %v\n", pkg, err)
		}
	}
	
	// Count successful installations
	successCount := 0
	for _, status := range model.results {
		if status == "success" {
			successCount++
		}
	}
	
	fmt.Printf("\n✅ Successfully installed %d/%d packages\n", successCount, len(packages))
	
	// If no packages were successfully installed, return an error
	if successCount == 0 && len(packages) > 0 {
		return fmt.Errorf("failed to install any packages")
	}
	
	return err
}

// getVersionString formats a version string with proper prefix
func getVersionString(version string) string {
	if version == "" {
		return ""
	}
	return fmt.Sprintf(" (v%s)", version)
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
} 