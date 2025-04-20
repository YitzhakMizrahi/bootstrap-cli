// cmd/install.go
package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/YitzhakMizrahi/bootstrap-cli/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/installer"
	"github.com/YitzhakMizrahi/bootstrap-cli/types"
)

// Installation stages
type InstallStage int

const (
	StageTools InstallStage = iota
	StageLanguages
	StageEditors
	StageDone
)

// installModel represents the state for the installation UI
type installModel struct {
	progress     progress.Model
	spinner      spinner.Model
	stage        InstallStage
	config       types.UserConfig
	stageCount   int
	currentItem  string
	totalItems   int
	itemsInstalled int
	err          error
	width        int
}

func newInstallModel() installModel {
	p := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	return installModel{
		progress:   p,
		spinner:    s,
		stage:      StageTools,
		stageCount: 3, // Tools, Languages, Editors
		width:      80,
	}
}

func (m installModel) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		func() tea.Msg {
			cfg, err := config.Load()
			if err != nil {
				return errMsg{err}
			}
			return configLoadedMsg{cfg}
		},
	)
}

// Custom messages
type errMsg struct{ err error }
func (e errMsg) Error() string { return e.err.Error() }

type configLoadedMsg struct{ config types.UserConfig }
type installItemMsg struct{ itemName string }
type stageCompleteMsg struct{ stage InstallStage }

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = msg.Width - 20
		if m.progress.Width < 20 {
			m.progress.Width = 20
		}
		return m, nil
		
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
		return m, nil
		
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
		
	case configLoadedMsg:
		m.config = msg.config
		
		// Determine total items for progress bar
		m.totalItems = len(m.config.CLITools) + len(m.config.Languages) + len(m.config.Editors)
		if m.totalItems == 0 {
			return m, tea.Quit
		}
		
		// Start installing CLI tools
		return m, runInstallStage(m)
		
	case installItemMsg:
		m.currentItem = msg.itemName
		m.itemsInstalled++
		cmd := m.progress.SetPercent(float64(m.itemsInstalled) / float64(m.totalItems))
		return m, cmd
		
	case stageCompleteMsg:
		// Move to next stage
		m.stage = m.stage + 1
		if m.stage >= StageDone {
			return m, tea.Quit
		}
		
		// Start next stage
		return m, runInstallStage(m)
		
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
		
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	
	return m, nil
}

func runInstallStage(m installModel) tea.Cmd {
	return func() tea.Msg {
		switch m.stage {
		case StageTools:
			if len(m.config.CLITools) > 0 {
				for _, tool := range m.config.CLITools {
					// Simulate installation here - this would actually call the installer
					installer.InstallCLITools([]string{tool})
					tea.NewProgram(&installModel{}).Send(installItemMsg{tool})
				}
			}
			return stageCompleteMsg{StageTools}
			
		case StageLanguages:
			if len(m.config.Languages) > 0 {
				for _, lang := range m.config.Languages {
					// Simulate installation here
					installer.InstallLanguages([]string{lang}, m.config.PackageManagers)
					tea.NewProgram(&installModel{}).Send(installItemMsg{lang})
				}
			}
			return stageCompleteMsg{StageLanguages}
			
		case StageEditors:
			if len(m.config.Editors) > 0 {
				for _, editor := range m.config.Editors {
					// Simulate installation here
					installer.InstallEditors([]string{editor})
					tea.NewProgram(&installModel{}).Send(installItemMsg{editor})
				}
			}
			return stageCompleteMsg{StageEditors}
		}
		
		return nil
	}
}

func (m installModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	
	// Build the view based on current stage
	var title string
	switch m.stage {
	case StageTools:
		title = "Installing CLI Tools"
	case StageLanguages:
		title = "Installing Languages"
	case StageEditors:
		title = "Installing Editors"
	case StageDone:
		return "✅ Installation complete!\n"
	}
	
	// Create a styled title
	styledTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render(title)
	
	// Progress indicator
	progressTitle := fmt.Sprintf("Progress: %d/%d", m.itemsInstalled, m.totalItems)
	progressDisplay := lipgloss.JoinVertical(
		lipgloss.Center,
		progressTitle,
		m.progress.View(),
	)
	
	// Current item being installed
	currentStatus := fmt.Sprintf("%s Installing: %s", m.spinner.View(), m.currentItem)
	
	// Overall status
	stageStatus := fmt.Sprintf("Stage: %d/%d", int(m.stage)+1, m.stageCount)
	
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		styledTitle,
		"",
		progressDisplay,
		currentStatus,
		stageStatus,
		"",
	)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install selected tools, languages, and editors",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("simple").Value.String() == "true" {
			// Run in simple mode (no UI)
			runSimpleInstall()
			return
		}
		
		// Run with the bubbletea UI
		model := newInstallModel()
		p := tea.NewProgram(model, tea.WithAltScreen())
		
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running install UI: %v\n", err)
			// Fallback to simple install
			runSimpleInstall()
		}
	},
}

// Simple install without the UI, as a fallback
func runSimpleInstall() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("❌ Failed to load config:", err)
		return
	}

	fmt.Println("📦 Installing tools...")
	installer.InstallCLITools(cfg.CLITools)

	fmt.Println("🧪 Installing languages...")
	installer.InstallLanguages(cfg.Languages, cfg.PackageManagers)

	fmt.Println("📝 Setting up editors...")
	installer.InstallEditors(cfg.Editors)

	fmt.Println("✅ Installation complete.")
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().Bool("simple", false, "Run in simple mode without UI")
}
