package installer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/platform"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tool represents a CLI tool that can be installed
type Tool struct {
	Name           string
	Description    string
	BrewPackage    string
	AptPackage     string
	DnfPackage     string
	PacmanPackage  string
	ZypperPackage  string
	ChocoPackage   string
	SnapPackage    string
	InstallScript  func() error // Custom installation function
}

// ToolPackages maps tool names to their installation info
var ToolPackages = map[string]Tool{
	"lsd": {
		Name:           "lsd",
		Description:    "The next gen ls command",
		BrewPackage:    "lsd",
		AptPackage:     "lsd",
		DnfPackage:     "lsd",
		PacmanPackage:  "lsd",
		ZypperPackage:  "lsd", 
		ChocoPackage:   "lsd",
		InstallScript:  installLsdFromBinary,
	},
	"bat": {
		Name:         "bat",
		Description:  "A cat clone with wings",
		BrewPackage:  "bat",
		AptPackage:   "bat",
		DnfPackage:   "bat",
		PacmanPackage: "bat",
		ChocoPackage: "bat",
	},
	"fzf": {
		Name:         "fzf",
		Description:  "Command-line fuzzy finder",
		BrewPackage:  "fzf",
		AptPackage:   "fzf",
		DnfPackage:   "fzf",
		PacmanPackage: "fzf",
		ChocoPackage: "fzf",
	},
	"ripgrep": {
		Name:         "ripgrep",
		Description:  "Fast search tool (rg)",
		BrewPackage:  "ripgrep",
		AptPackage:   "ripgrep",
		DnfPackage:   "ripgrep",
		PacmanPackage: "ripgrep",
		ChocoPackage: "ripgrep",
	},
	"fd": {
		Name:         "fd",
		Description:  "Simple, fast alternative to find",
		BrewPackage:  "fd",
		AptPackage:   "fd-find",
		DnfPackage:   "fd-find",
		PacmanPackage: "fd",
		ChocoPackage: "fd",
	},
	"jq": {
		Name:         "jq",
		Description:  "Lightweight JSON processor",
		BrewPackage:  "jq",
		AptPackage:   "jq",
		DnfPackage:   "jq",
		PacmanPackage: "jq",
		ChocoPackage: "jq",
	},
	"htop": {
		Name:         "htop",
		Description:  "Interactive process viewer",
		BrewPackage:  "htop",
		AptPackage:   "htop",
		DnfPackage:   "htop",
		PacmanPackage: "htop",
		ChocoPackage: "htop",
	},
	"lazygit": {
		Name:         "lazygit",
		Description:  "Simple terminal UI for git",
		BrewPackage:  "lazygit",
		AptPackage:   "lazygit",
		DnfPackage:   "lazygit",
		PacmanPackage: "lazygit",
		ChocoPackage: "lazygit",
	},
	"tmux": {
		Name:         "tmux",
		Description:  "Terminal multiplexer",
		BrewPackage:  "tmux",
		AptPackage:   "tmux",
		DnfPackage:   "tmux",
		PacmanPackage: "tmux",
		ChocoPackage: "tmux",
	},
	"neofetch": {
		Name:         "neofetch",
		Description:  "System info written in bash",
		BrewPackage:  "neofetch",
		AptPackage:   "neofetch",
		DnfPackage:   "neofetch",
		PacmanPackage: "neofetch",
		ChocoPackage: "neofetch",
	},
}

// Language represents a programming language that can be installed
type Language struct {
	Name         string
	Description  string
	InstallCmd   func() error
	Managers     []string // Available package managers for this language
}

// LanguagePackages maps language names to their installation info
var LanguagePackages = map[string]Language{
	"node": {
		Name:        "Node.js",
		Description: "JavaScript runtime",
		InstallCmd:  installNode,
		Managers:    []string{"npm", "yarn", "pnpm"},
	},
	"python": {
		Name:        "Python",
		Description: "Python programming language",
		InstallCmd:  installPython,
		Managers:    []string{"pip", "pipenv", "poetry"},
	},
	"go": {
		Name:        "Go",
		Description: "Go programming language",
		InstallCmd:  installGo,
		Managers:    []string{},
	},
	"rust": {
		Name:        "Rust",
		Description: "Rust programming language",
		InstallCmd:  installRust,
		Managers:    []string{"cargo"},
	},
}

// ToolInstallModel represents a model for tool installation UI
type ToolInstallModel struct {
	title        string
	tools        []string
	current      int
	width        int
	height       int
	spinner      spinner.Model
	progress     progress.Model
	done         bool
	results      map[string]string // tool -> status
	errors       map[string]error
	packageMgr   platform.PackageManager
	activeInstall bool
	installChan  chan toolResultMsg
}

// NewToolInstallModel creates a new model for tool installation UI
func NewToolInstallModel(title string, tools []string, packageMgr platform.PackageManager) ToolInstallModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return ToolInstallModel{
		title:        title,
		tools:        tools,
		spinner:      s,
		progress:     p,
		results:      make(map[string]string),
		errors:       make(map[string]error),
		packageMgr:   packageMgr,
		installChan:  make(chan toolResultMsg),
	}
}

// Init initializes the model
func (m ToolInstallModel) Init() tea.Cmd {
	if len(m.tools) == 0 {
		return tea.Quit
	}
	
	return tea.Batch(
		m.spinner.Tick,
		m.waitForToolInstall(),
		m.startToolInstalls(),
	)
}

// startToolInstalls launches all tool installs in parallel
func (m ToolInstallModel) startToolInstalls() tea.Cmd {
	return func() tea.Msg {
		m.activeInstall = true
		
		// Start a goroutine to install all tools
		go func() {
			for i, toolName := range m.tools {
				tool, exists := ToolPackages[toolName]
				
				if !exists {
					m.installChan <- toolResultMsg{
						tool:   toolName,
						index:  i,
						status: "error",
						err:    fmt.Errorf("unknown tool"),
					}
					continue
				}
				
				// Check if already installed
				if isCommandAvailable(toolName) {
					// For apps that don't install as commands with the same name
					alreadyInstalled := true
					switch toolName {
					case "fd":
						// Check for fd/fdfind
						if !isCommandAvailable("fd") && !isCommandAvailable("fdfind") {
							alreadyInstalled = false
						}
					case "ripgrep":
						// Check for rg
						if !isCommandAvailable("rg") {
							alreadyInstalled = false
						}
					}
					
					if alreadyInstalled {
						m.installChan <- toolResultMsg{
							tool:   toolName,
							index:  i,
							status: "skipped",
						}
						continue
					}
				}
				
				// Get the package name based on package manager
				var packageName string
				var installErr error
				
				switch m.packageMgr {
				case platform.Homebrew:
					packageName = tool.BrewPackage
					if packageName != "" {
						cmd := exec.Command("brew", "install", packageName)
						cmd.Stdout = nil
						cmd.Stderr = nil
						installErr = cmd.Run()
					} else {
						installErr = fmt.Errorf("no brew package available")
					}
				case platform.Apt:
					packageName = tool.AptPackage
					if packageName != "" {
						cmd := exec.Command("sudo", "apt-get", "install", "-y", packageName)
						cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
						cmd.Stdout = nil
						cmd.Stderr = nil
						installErr = cmd.Run()
					} else {
						installErr = fmt.Errorf("no apt package available")
					}
				case platform.Dnf:
					packageName = tool.DnfPackage
					if packageName != "" {
						cmd := exec.Command("sudo", "dnf", "install", "-y", packageName)
						cmd.Stdout = nil
						cmd.Stderr = nil
						installErr = cmd.Run()
					} else {
						installErr = fmt.Errorf("no dnf package available")
					}
				case platform.Pacman:
					packageName = tool.PacmanPackage
					if packageName != "" {
						cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", packageName)
						cmd.Stdout = nil
						cmd.Stderr = nil
						installErr = cmd.Run()
					} else {
						installErr = fmt.Errorf("no pacman package available")
					}
				case platform.Zypper:
					packageName = tool.ZypperPackage
					if packageName != "" {
						cmd := exec.Command("sudo", "zypper", "install", "-y", packageName)
						cmd.Stdout = nil
						cmd.Stderr = nil
						installErr = cmd.Run()
					} else {
						installErr = fmt.Errorf("no zypper package available")
					}
				case platform.Chocolatey:
					packageName = tool.ChocoPackage
					if packageName != "" {
						cmd := exec.Command("choco", "install", packageName, "-y")
						cmd.Stdout = nil
						cmd.Stderr = nil
						installErr = cmd.Run()
					} else {
						installErr = fmt.Errorf("no chocolatey package available")
					}
				default:
					installErr = fmt.Errorf("unsupported package manager: %s", m.packageMgr)
				}
				
				// Try alternative installation if package manager failed
				if installErr != nil && tool.InstallScript != nil {
					// Try alternative installation method
					installScriptErr := tool.InstallScript()
					if installScriptErr == nil {
						// Alternative installation succeeded
						installErr = nil
					}
				}
				
				if installErr != nil {
					m.installChan <- toolResultMsg{
						tool:   toolName,
						index:  i,
						status: "error",
						err:    installErr,
					}
				} else {
					m.installChan <- toolResultMsg{
						tool:   toolName,
						index:  i,
						status: "success",
					}
				}
				
				// Add a small delay for UI
				time.Sleep(300 * time.Millisecond)
			}
			
			// Close the channel when all installs are complete
			close(m.installChan)
		}()
		
		return nil
	}
}

// waitForToolInstall waits for tool installations to complete
func (m ToolInstallModel) waitForToolInstall() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.installChan
		if !ok {
			return toolResultMsg{
				status: "done",
			}
		}
		return msg
	}
}

// Update updates the model
func (m ToolInstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case toolResultMsg:
		if msg.status == "done" {
			m.done = true
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Tick(time.Millisecond*750, func(time.Time) tea.Msg { return nil }),
				tea.Quit,
			)
		}
		
		// Record result
		m.results[msg.tool] = msg.status
		if msg.err != nil {
			m.errors[msg.tool] = msg.err
		}

		// Count completed tools
		completed := len(m.results)
		
		// Update progress based on completed tools
		progCmd := m.progress.SetPercent(float64(completed) / float64(len(m.tools)))
		
		// Check if we've completed all tools
		if completed >= len(m.tools) {
			m.done = true
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Tick(time.Millisecond*750, func(time.Time) tea.Msg { return nil }),
				tea.Quit,
			)
		}

		// Keep waiting for more results
		return m, tea.Batch(
			progCmd,
			m.waitForToolInstall(),
		)
	}
	
	return m, nil
}

// View renders the current view of the model
func (m ToolInstallModel) View() string {
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
		s.WriteString(fmt.Sprintf("✅ Installed %d/%d tools\n\n", successCount, len(m.tools)))
		
		// Show status for each tool
		successMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#73F59F")).SetString("✓")
		errorMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).SetString("✗")
		skippedMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAE19")).SetString("•")
		
		for _, toolName := range m.tools {
			displayName := toolName
			if tool, ok := ToolPackages[toolName]; ok {
				displayName = tool.Name
				if tool.Description != "" {
					displayName = fmt.Sprintf("%s (%s)", tool.Name, tool.Description)
				}
			}
			
			marker := "  "
			if status, ok := m.results[toolName]; ok {
				switch status {
				case "success":
					marker = successMark.String() + " "
				case "error":
					marker = errorMark.String() + " "
				case "skipped":
					marker = skippedMark.String() + " "
				}
			}
			
			s.WriteString(marker + displayName)
			
			// Add error message if there was one
			if err, hasError := m.errors[toolName]; hasError {
				errorMsg := fmt.Sprintf(" - Error: %v", err)
				s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).Render(errorMsg))
			}
			
			s.WriteString("\n")
		}
		
		return doneStyle.Render(s.String())
	}

	// Get count of completed tools
	completed := len(m.results)
	
	// Print previously installed tools with their statuses
	var s strings.Builder
	
	// Add title
	s.WriteString(titleStyle.Render(m.title) + "\n\n")
	
	// Add progress bar
	progressText := fmt.Sprintf("Installing %d/%d tools", completed, len(m.tools))
	s.WriteString(progressText + "\n")
	s.WriteString(m.progress.View() + "\n\n")
	
	// Add tool list with status indicators
	for i, toolName := range m.tools {
		marker := "  "
		tool, exists := ToolPackages[toolName]
		
		displayName := toolName
		if exists {
			displayName = tool.Name
		}
		
		if status, ok := m.results[toolName]; ok {
			switch status {
			case "success":
				marker = successMark.String() + " "
			case "error":
				marker = errorMark.String() + " "
			case "skipped":
				marker = skippedMark.String() + " "
			}
		} else if completed < len(m.tools) && i == completed {
			// Current tool being installed
			marker = m.spinner.View() + " "
			displayName = currentPkgStyle.Render(displayName)
		}
		
		s.WriteString(marker + displayName + "\n")
	}
	
	// Add help text at the bottom
	s.WriteString("\nPress Ctrl+C to cancel\n")
	
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Render(s.String())
}

// toolResultMsg represents a message with tool installation result
type toolResultMsg struct {
	tool   string
	index  int
	status string
	err    error
}

// InstallToolsWithUI installs tools with a bubbletea UI
func InstallToolsWithUI(title string, tools []string, packageManager platform.PackageManager) error {
	if len(tools) == 0 {
		fmt.Println("No tools to install.")
		return nil
	}
	
	model := NewToolInstallModel(title, tools, packageManager)
	p := tea.NewProgram(model)
	
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Error running tool installation UI: %v\n", err)
	}
	
	// Print errors after completion
	if len(model.errors) > 0 {
		fmt.Println("\nThe following tools encountered errors:")
		for tool, err := range model.errors {
			fmt.Printf("  - %s: %v\n", tool, err)
		}
	}
	
	return err
}

// InstallCLITools installs the selected CLI tools
func InstallCLITools(tools []string) error {
	platformInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}
	
	primaryPM, err := platform.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}
	
	fmt.Printf("🛠️ Installing CLI tools using %s...\n", primaryPM)
	
	// If we should use the UI-based installer
	if usePackageManagerUI() {
		switch primaryPM {
		case platform.Homebrew:
			// Get the actual package names for all tools
			var packages []string
			for _, toolName := range tools {
				tool, exists := ToolPackages[toolName]
				if !exists {
					fmt.Printf("⚠️  Unknown tool: %s, skipping...\n", toolName)
					continue
				}
				packages = append(packages, tool.BrewPackage)
			}
			return InstallPackagesWithUI("CLI Tools", packages, "brew")
			
		case platform.Apt:
			// Get the actual package names for all tools
			var packages []string
			for _, toolName := range tools {
				tool, exists := ToolPackages[toolName]
				if !exists {
					fmt.Printf("⚠️  Unknown tool: %s, skipping...\n", toolName)
					continue
				}
				packages = append(packages, tool.AptPackage)
			}
			return InstallPackagesWithUI("CLI Tools", packages, "apt")
			
		default:
			// Fall back to regular installation for other package managers
			fmt.Printf("⚠️ UI-based installation not supported for %s, using standard installer\n", primaryPM)
		}
	}
	
	// Traditional (non-UI) installation process
	for _, toolName := range tools {
		tool, exists := ToolPackages[toolName]
		if !exists {
			fmt.Printf("⚠️  Unknown tool: %s, skipping...\n", toolName)
			continue
		}
		
		fmt.Printf("📦 Installing %s...\n", tool.Name)
		
		// Check if already installed
		if isCommandAvailable(toolName) {
			fmt.Printf("✅ %s is already installed, skipping\n", tool.Name)
			continue
		}
		
		// Try to install with the appropriate package manager
		var installErr error
		switch primaryPM {
		case platform.Homebrew:
			installErr = brewInstall(tool.BrewPackage)
		case platform.Apt:
			installErr = aptInstall(tool.AptPackage)
		case platform.Dnf:
			installErr = dnfInstall(tool.DnfPackage)
		case platform.Pacman:
			installErr = pacmanInstall(tool.PacmanPackage)
		case platform.Zypper:
			installErr = zypperInstall(tool.ZypperPackage)
		case platform.Chocolatey:
			installErr = chocoInstall(tool.ChocoPackage)
		default:
			installErr = fmt.Errorf("unsupported package manager: %s", primaryPM)
		}
		
		if installErr != nil {
			fmt.Printf("⚠️  Failed to install %s: %v\n", tool.Name, installErr)
			// Try alternative installation methods if available
			if tool.InstallScript != nil {
				fmt.Printf("🔄 Trying alternative installation method for %s...\n", tool.Name)
				if err := tool.InstallScript(); err != nil {
					fmt.Printf("❌ Alternative installation failed: %v\n", err)
				} else {
					fmt.Printf("✅ Installed %s using alternative method\n", tool.Name)
				}
			}
		} else {
			fmt.Printf("✅ Successfully installed %s\n", tool.Name)
		}
	}
	
	return nil
}

// usePackageManagerUI determines if we should use the bubbletea UI for package installation
func usePackageManagerUI() bool {
	// Check for environment variable to disable UI
	if os.Getenv("BOOTSTRAP_NO_UI") == "1" {
		return false
	}
	
	// Check if terminal supports UI (has columns/rows information)
	if _, err := os.Stdout.Stat(); err != nil {
		return false
	}
	
	// Check if terminal is interactive
	if !isTerminal() {
		return false
	}
	
	return true
}

// isTerminal checks if standard output is connected to a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	
	// On Unix systems, check if it's a character device
	if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	
	return false
}

// InstallLanguages installs the selected programming languages
func InstallLanguages(languages []string, packageManagers map[string]string) error {
	if len(languages) == 0 {
		return nil
	}
	
	fmt.Printf("🧪 Installing programming languages...\n")
	
	// If we should use the UI-based installer
	if usePackageManagerUI() {
		return InstallLanguagesWithUI("Programming Languages", languages, packageManagers)
	}
	
	// Traditional (non-UI) installation process
	for _, lang := range languages {
		language, exists := LanguagePackages[lang]
		if !exists {
			fmt.Printf("⚠️  Unknown language: %s, skipping...\n", lang)
			continue
		}
		
		fmt.Printf("🧪 Installing %s...\n", language.Name)
		
		if err := language.InstallCmd(); err != nil {
			fmt.Printf("❌ Failed to install %s: %v\n", language.Name, err)
		} else {
			fmt.Printf("✅ Successfully installed %s\n", language.Name)
			
			// Install selected package manager if applicable
			if mgr, ok := packageManagers[lang]; ok && mgr != "" {
				fmt.Printf("📦 Setting up %s package manager: %s\n", language.Name, mgr)
				if err := installPackageManager(lang, mgr); err != nil {
					fmt.Printf("❌ Failed to set up package manager %s: %v\n", mgr, err)
				}
			}
		}
	}
	
	return nil
}

// LanguageInstallModel represents a model for language installation UI
type LanguageInstallModel struct {
	title        string
	languages    []string
	current      int
	width        int
	height       int
	spinner      spinner.Model
	progress     progress.Model
	done         bool
	results      map[string]string // language -> status
	errors       map[string]error
	pkgManagers  map[string]string
	activeInstall bool
	installChan  chan languageResultMsg
}

// NewLanguageInstallModel creates a new model for language installation UI
func NewLanguageInstallModel(title string, languages []string, packageManagers map[string]string) LanguageInstallModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return LanguageInstallModel{
		title:        title,
		languages:    languages,
		spinner:      s,
		progress:     p,
		results:      make(map[string]string),
		errors:       make(map[string]error),
		pkgManagers:  packageManagers,
		installChan:  make(chan languageResultMsg),
	}
}

// Init initializes the model
func (m LanguageInstallModel) Init() tea.Cmd {
	if len(m.languages) == 0 {
		return tea.Quit
	}
	
	return tea.Batch(
		m.spinner.Tick,
		m.waitForLanguageInstall(),
		m.startLanguageInstalls(),
	)
}

// startLanguageInstalls launches all language installs in parallel
func (m LanguageInstallModel) startLanguageInstalls() tea.Cmd {
	return func() tea.Msg {
		m.activeInstall = true
		
		// Start a goroutine to install all languages
		go func() {
			for i, lang := range m.languages {
				langPkg, exists := LanguagePackages[lang]
				
				if !exists {
					m.installChan <- languageResultMsg{
						lang:   lang,
						index:  i,
						status: "error",
						err:    fmt.Errorf("unknown language"),
					}
					continue
				}
				
				output, err := installLanguage(lang, langPkg, m.pkgManagers)
				
				if err != nil {
					m.installChan <- languageResultMsg{
						lang:   lang,
						index:  i,
						status: "error",
						err:    err,
						output: output,
					}
					continue
				}
				
				m.installChan <- languageResultMsg{
					lang:   lang,
					index:  i,
					status: "success",
					output: output,
				}
			}
			
			// Close the channel when all installs are complete
			close(m.installChan)
		}()
		
		return nil
	}
}

// waitForLanguageInstall waits for language installations to complete
func (m LanguageInstallModel) waitForLanguageInstall() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.installChan
		if !ok {
			return languageResultMsg{
				status: "done",
			}
		}
		return msg
	}
}

// Update updates the model
func (m LanguageInstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case languageResultMsg:
		if msg.status == "done" {
			m.done = true
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Tick(time.Millisecond*750, func(time.Time) tea.Msg { return nil }),
				tea.Quit,
			)
		}
		
		// Record result
		m.results[msg.lang] = msg.status
		if msg.err != nil {
			m.errors[msg.lang] = msg.err
		}

		// Count completed languages
		completed := len(m.results)
		
		// Update progress based on completed languages
		progCmd := m.progress.SetPercent(float64(completed) / float64(len(m.languages)))
		
		// Check if we've completed all languages
		if completed >= len(m.languages) {
			m.done = true
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Tick(time.Millisecond*750, func(time.Time) tea.Msg { return nil }),
				tea.Quit,
			)
		}

		// Keep waiting for more results
		return m, tea.Batch(
			progCmd,
			m.waitForLanguageInstall(),
		)
	}
	
	return m, nil
}

// View renders the current view of the model
func (m LanguageInstallModel) View() string {
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
		s.WriteString(fmt.Sprintf("✅ Installed %d/%d languages\n\n", successCount, len(m.languages)))
		
		// Show status for each language
		successMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#73F59F")).SetString("✓")
		errorMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).SetString("✗")
		
		for _, lang := range m.languages {
			displayName := lang
			if langPkg, ok := LanguagePackages[lang]; ok {
				displayName = langPkg.Name
			}
			
			marker := "  "
			if status, ok := m.results[lang]; ok {
				switch status {
				case "success":
					marker = successMark.String() + " "
				case "error":
					marker = errorMark.String() + " "
				}
			}
			
			s.WriteString(marker + displayName)
			
			// Add package manager info if applicable
			if pm, ok := m.pkgManagers[lang]; ok && pm != "" {
				s.WriteString(fmt.Sprintf(" (with %s)", pm))
			}
			
			// Add error message if there was one
			if err, hasError := m.errors[lang]; hasError {
				errorMsg := fmt.Sprintf(" - Error: %v", err)
				s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).Render(errorMsg))
			}
			
			s.WriteString("\n")
		}
		
		return doneStyle.Render(s.String())
	}

	// Get count of completed languages
	completed := len(m.results)
	
	// Print previously installed languages with their statuses
	var s strings.Builder
	
	// Add title
	s.WriteString(titleStyle.Render(m.title) + "\n\n")
	
	// Add progress info without the progress bar yet
	progressText := fmt.Sprintf("Installing %d/%d languages", completed, len(m.languages))
	s.WriteString(progressText + "\n\n")
	
	// Add language list with status indicators
	successMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#73F59F")).SetString("✓")
	errorMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).SetString("✗")
	pendingMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAE19")).SetString("•")
	
	for _, lang := range m.languages {
		displayName := lang
		if langPkg, ok := LanguagePackages[lang]; ok {
			displayName = langPkg.Name
		}
		
		// Determine the status marker and style
		marker := pendingMark.String() + " "
		lineStyle := lipgloss.NewStyle()
		
		if status, ok := m.results[lang]; ok {
			// Language has a result, show checkmark or error
			switch status {
			case "success":
				marker = successMark.String() + " "
			case "error":
				marker = errorMark.String() + " "
				lineStyle = lineStyle.Foreground(lipgloss.Color("#FF4672"))
			}
		} else if m.activeInstall {
			// Language is currently installing, show spinner
			// Find the next language to install (one without a result)
			for _, l := range m.languages {
				if _, ok := m.results[l]; !ok {
					if l == lang {
						marker = m.spinner.View() + " "
						lineStyle = lineStyle.Foreground(lipgloss.Color("#FF78D4"))
						break
					}
					break
				}
			}
		}
		
		s.WriteString(marker + lineStyle.Render(displayName))
		
		// Add package manager info if applicable
		if pm, ok := m.pkgManagers[lang]; ok && pm != "" {
			s.WriteString(fmt.Sprintf(" (with %s)", pm))
		}
		
		// Add error message if there was one
		if err, hasError := m.errors[lang]; hasError {
			errorMsg := fmt.Sprintf(" - Error: %v", err)
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4672")).Render(errorMsg))
		}
		
		s.WriteString("\n")
	}
	
	// Add a spacer before the progress bar
	s.WriteString("\n")
	
	// Add progress bar after the list of languages
	s.WriteString(m.progress.View() + "\n")
	
	// Add helpful instructions
	s.WriteString("\n" + lipgloss.NewStyle().Faint(true).Render("Press Ctrl+C to cancel"))
	
	// Apply a consistent style to the whole output
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Render(s.String())
}

// languageResultMsg represents a message with language installation result
type languageResultMsg struct {
	lang   string
	index  int
	status string
	err    error
	output string
}

// InstallLanguagesWithUI installs languages with a bubbletea UI
func InstallLanguagesWithUI(title string, languages []string, packageManagers map[string]string) error {
	if len(languages) == 0 {
		fmt.Println("No languages to install.")
		return nil
	}
	
	model := NewLanguageInstallModel(title, languages, packageManagers)
	p := tea.NewProgram(model)
	
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Error running language installation UI: %v\n", err)
	}
	
	// Print errors after completion
	if len(model.errors) > 0 {
		fmt.Println("\nThe following languages encountered errors:")
		for lang, err := range model.errors {
			fmt.Printf("  - %s: %v\n", lang, err)
		}
	}
	
	return err
}

// Language installations

func installNode() error {
	// Check if node is already installed
	if _, err := exec.LookPath("node"); err == nil {
		fmt.Println("✅ Node.js is already installed")
		return nil
	}

	// Ensure we have curl
	if _, err := exec.LookPath("curl"); err != nil {
		return fmt.Errorf("curl not found, required for NVM installation")
	}
	
	// Get HOME directory
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME environment variable not set")
	}
	
	// Check if nvm exists
	nvmPath := fmt.Sprintf("%s/.nvm", home)
	if _, err := os.Stat(nvmPath); os.IsNotExist(err) {
		fmt.Println("🔄 Installing NVM (Node Version Manager)...")
		
		// Create directories if they don't exist
		if err := os.MkdirAll(nvmPath, 0755); err != nil {
			return fmt.Errorf("failed to create NVM directory: %w", err)
		}
		
		// Install nvm
		installCmd := `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash`
		cmd := exec.Command("bash", "-c", installCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install nvm: %w", err)
		}
		
		// Update shell config to source nvm
		shellFiles := []string{
			fmt.Sprintf("%s/.bashrc", home),
			fmt.Sprintf("%s/.zshrc", home),
			fmt.Sprintf("%s/.profile", home),
		}
		
		nvmInit := `
# NVM setup
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
`
		
		// Add NVM init to shell config files
		for _, file := range shellFiles {
			if _, err := os.Stat(file); err == nil {
				// Read current content
				content, err := os.ReadFile(file)
				if err != nil {
					fmt.Printf("⚠️ Warning: Could not read %s: %v\n", file, err)
					continue
				}
				
				// Check if NVM is already configured
				if strings.Contains(string(content), "NVM_DIR") {
					continue
				}
				
				// Append NVM init
				if err := appendToFile(file, nvmInit); err != nil {
					fmt.Printf("⚠️ Warning: Could not update %s: %v\n", file, err)
				} else {
					fmt.Printf("✅ Updated %s with NVM configuration\n", file)
				}
			}
		}
		
		// Write a temporary script to source NVM
		tmpScript := fmt.Sprintf("%s/nvm_install_script.sh", nvmPath)
		scriptContent := `#!/bin/bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm install --lts
nvm use --lts
node --version
`
		if err := os.WriteFile(tmpScript, []byte(scriptContent), 0755); err != nil {
			return fmt.Errorf("failed to create temporary NVM script: %w", err)
		}
		defer os.Remove(tmpScript)
	}
	
	// Install latest LTS version of Node
	fmt.Println("🔄 Installing Node.js LTS version...")
	
	cmd := exec.Command("bash", "-c", 
		`export NVM_DIR="$HOME/.nvm" && 
		[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && 
		nvm install --lts && nvm use --lts && node --version`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Node.js: %w", err)
	}
	
	fmt.Println("✅ Node.js LTS installed successfully")
	
	// Add node binary directory to PATH for the current session
	nvmCurrentPath := fmt.Sprintf("%s/.nvm/versions/node/**/bin", home)
	paths, err := filepath.Glob(nvmCurrentPath)
	if err == nil && len(paths) > 0 {
		newPath := paths[0] + ":" + os.Getenv("PATH")
		os.Setenv("PATH", newPath)
		fmt.Printf("✅ Added Node.js to PATH: %s\n", paths[0])
	}
	
	return nil
}

func installPython() error {
	// Check if Python is already installed
	if _, err := exec.LookPath("python3"); err == nil {
		fmt.Println("✅ Python is already installed")
		// Check the version
		cmd := exec.Command("python3", "--version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // ignore errors
		return nil
	}
	
	// Get HOME directory
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME environment variable not set")
	}
	
	// Check if pyenv exists
	pyenvPath := fmt.Sprintf("%s/.pyenv", home)
	if _, err := os.Stat(pyenvPath); os.IsNotExist(err) {
		fmt.Println("🔄 Installing pyenv...")
		
		// Create directories if they don't exist
		if err := os.MkdirAll(pyenvPath, 0755); err != nil {
			return fmt.Errorf("failed to create pyenv directory: %w", err)
		}
		
		// Install dependencies first
		fmt.Println("🔄 Installing pyenv dependencies...")
		platformInfo, err := platform.Detect()
		if err == nil && platformInfo.OS == platform.Linux {
			// Install dependencies based on distribution
			switch platformInfo.Distribution {
			case "ubuntu", "debian":
				depsCmd := exec.Command("sudo", "apt-get", "install", "-y", 
					"make", "build-essential", "libssl-dev", "zlib1g-dev", 
					"libbz2-dev", "libreadline-dev", "libsqlite3-dev", "wget", 
					"curl", "llvm", "libncurses5-dev", "libncursesw5-dev", 
					"xz-utils", "tk-dev", "libffi-dev", "liblzma-dev", "python-openssl")
				depsCmd.Stdout = os.Stdout
				depsCmd.Stderr = os.Stderr
				depsCmd.Run() // ignore errors, best effort
			}
		}
		
		// Install pyenv
		installCmd := `curl https://pyenv.run | bash`
		cmd := exec.Command("bash", "-c", installCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install pyenv: %w", err)
		}
		
		// Add pyenv to shell config
		shellFiles := []string{
			fmt.Sprintf("%s/.bashrc", home),
			fmt.Sprintf("%s/.zshrc", home),
			fmt.Sprintf("%s/.profile", home),
		}
		
		pyenvInit := `
# pyenv
export PYENV_ROOT="$HOME/.pyenv"
export PATH="$PYENV_ROOT/bin:$PATH"
if command -v pyenv 1>/dev/null 2>&1; then
  eval "$(pyenv init --path)"
  eval "$(pyenv init -)"
  eval "$(pyenv virtualenv-init -)"
fi
`
		
		// Add pyenv init to shell config files
		for _, file := range shellFiles {
			if _, err := os.Stat(file); err == nil {
				// Read current content
				content, err := os.ReadFile(file)
				if err != nil {
					fmt.Printf("⚠️ Warning: Could not read %s: %v\n", file, err)
					continue
				}
				
				// Check if pyenv is already configured
				if strings.Contains(string(content), "PYENV_ROOT") {
					continue
				}
				
				// Append pyenv init
				if err := appendToFile(file, pyenvInit); err != nil {
					fmt.Printf("⚠️ Warning: Could not update %s: %v\n", file, err)
				} else {
					fmt.Printf("✅ Updated %s with pyenv configuration\n", file)
				}
			}
		}
	}
	
	// Set up environment for pyenv
	os.Setenv("PYENV_ROOT", fmt.Sprintf("%s/.pyenv", home))
	path := fmt.Sprintf("%s/.pyenv/bin:%s", home, os.Getenv("PATH"))
	os.Setenv("PATH", path)
	
	// Install Python using pyenv
	fmt.Println("🔄 Installing Python 3.10 using pyenv...")
	cmd := exec.Command("bash", "-c", 
		fmt.Sprintf(`export PYENV_ROOT="%s/.pyenv" && 
		export PATH="$PYENV_ROOT/bin:$PATH" && 
		eval "$(~/.pyenv/bin/pyenv init -)" && 
		pyenv install 3.10.0 -s && 
		pyenv global 3.10.0 && 
		python --version`, home))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Python: %w", err)
	}
	
	fmt.Println("✅ Python installed successfully")
	return nil
}

func installGo() error {
	// Check if Go is already installed
	if _, err := exec.LookPath("go"); err == nil {
		fmt.Println("✅ Go is already installed")
		// Check the version
		cmd := exec.Command("go", "version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // ignore errors
		return nil
	}
	
	platformInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}
	
	primaryPM, err := platform.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}
	
	fmt.Println("🔄 Installing Go...")
	
	// Create a wrapper function to handle package manager output
	installWithPM := func(cmd *exec.Cmd) error {
		// Capture the output without letting it break our UI
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install Go: %v\n%s", err, stderr.String())
		}
		
		// Print important parts of the output without the noise
		fmt.Println("📦 Go installation completed successfully")
		return nil
	}
	
	// Install Go using package manager
	switch primaryPM {
	case platform.Homebrew:
		cmd := exec.Command("brew", "install", "go")
		return installWithPM(cmd)
		
	case platform.Apt:
		// Update package lists first (quietly)
		updateCmd := exec.Command("sudo", "apt-get", "update", "-qq")
		updateCmd.Run() // Ignore errors and continue with install
		
		// Install golang
		cmd := exec.Command("sudo", "apt-get", "install", "-y", "golang-go")
		return installWithPM(cmd)
		
	case platform.Dnf:
		cmd := exec.Command("sudo", "dnf", "install", "-y", "golang")
		return installWithPM(cmd)
		
	case platform.Pacman:
		cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "go")
		return installWithPM(cmd)
		
	case platform.Zypper:
		cmd := exec.Command("sudo", "zypper", "install", "-y", "go")
		return installWithPM(cmd)
		
	case platform.Chocolatey:
		cmd := exec.Command("choco", "install", "golang", "-y")
		return installWithPM(cmd)
		
	default:
		// Fallback to manual install for other platforms
		fmt.Println("⚠️ Package manager not supported for Go installation. Trying direct download...")
		return installGoFromSource()
	}
}

// installGoFromSource installs Go directly from golang.org
func installGoFromSource() error {
	// Get HOME directory
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME environment variable not set")
	}
	
	// Create a temp directory
	tempDir, err := os.MkdirTemp("", "go-install")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Determine architecture and OS
	var arch, osName string
	switch runtime.GOOS {
	case "linux":
		osName = "linux"
	case "darwin":
		osName = "darwin"
	case "windows":
		osName = "windows"
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	
	switch runtime.GOARCH {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	case "386":
		arch = "386"
	default:
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
	
	// Download latest Go version (using a static link for now)
	version := "1.22.0" // Latest stable version at time of writing
	downloadURL := fmt.Sprintf("https://go.dev/dl/go%s.%s-%s.tar.gz", version, osName, arch)
	
	fmt.Printf("📥 Downloading Go %s from %s\n", version, downloadURL)
	
	// Download the archive
	archivePath := filepath.Join(tempDir, "go.tar.gz")
	downloadCmd := exec.Command("curl", "-L", downloadURL, "-o", archivePath)
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to download Go: %w", err)
	}
	
	// Extract Go to /usr/local (or similar)
	installPath := "/usr/local"
	if runtime.GOOS == "windows" {
		installPath = "C:\\Program Files"
	}
	
	fmt.Printf("📦 Extracting Go to %s\n", installPath)
	
	extractCmd := exec.Command("sudo", "tar", "-C", installPath, "-xzf", archivePath)
	extractCmd.Stdout = os.Stdout
	extractCmd.Stderr = os.Stderr
	if err := extractCmd.Run(); err != nil {
		return fmt.Errorf("failed to extract Go: %w", err)
	}
	
	// Add Go to shell config files
	shellFiles := []string{
		fmt.Sprintf("%s/.bashrc", home),
		fmt.Sprintf("%s/.zshrc", home),
		fmt.Sprintf("%s/.profile", home),
	}
	
	goInit := `
# Go setup
export GOPATH="$HOME/go"
export PATH="$PATH:/usr/local/go/bin:$GOPATH/bin"
`
	
	// Add Go init to shell config files
	for _, file := range shellFiles {
		if _, err := os.Stat(file); err == nil {
			// Read current content
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("⚠️ Warning: Could not read %s: %v\n", file, err)
				continue
			}
			
			// Check if Go is already configured
			if strings.Contains(string(content), "GOPATH") {
				continue
			}
			
			// Append Go init
			if err := appendToFile(file, goInit); err != nil {
				fmt.Printf("⚠️ Warning: Could not update %s: %v\n", file, err)
			} else {
				fmt.Printf("✅ Updated %s with Go configuration\n", file)
			}
		}
	}
	
	// Add Go to PATH for the current session
	goPath := "/usr/local/go/bin"
	if _, err := os.Stat(goPath); err == nil {
		newPath := goPath + ":" + os.Getenv("PATH")
		os.Setenv("PATH", newPath)
		fmt.Printf("✅ Added Go to PATH: %s\n", goPath)
	}
	
	fmt.Println("✅ Go installed successfully")
	return nil
}

func installRust() error {
	// Check if Rust is already installed
	if _, err := exec.LookPath("rustc"); err == nil {
		fmt.Println("✅ Rust is already installed")
		// Check the version
		cmd := exec.Command("rustc", "--version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // ignore errors
		return nil
	}
	
	// Get HOME directory
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME environment variable not set")
	}
	
	// Install Rust via rustup
	fmt.Println("🔄 Installing Rust via rustup...")
	
	// Create a temporary script for installation
	tmpScript := filepath.Join(os.TempDir(), "install_rust.sh")
	scriptContent := `#!/bin/bash
set -e
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --no-modify-path
source "$HOME/.cargo/env"
rustc --version
`
	if err := os.WriteFile(tmpScript, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to create temporary rustup script: %w", err)
	}
	defer os.Remove(tmpScript)
	
	// Run the installation script
	cmd := exec.Command("bash", tmpScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to install Rust: %w", err)
	}
	
	// Add Rust to shell config files
	shellFiles := []string{
		fmt.Sprintf("%s/.bashrc", home),
		fmt.Sprintf("%s/.zshrc", home),
		fmt.Sprintf("%s/.profile", home),
	}
	
	rustInit := `
# Rust setup
[ -f "$HOME/.cargo/env" ] && source "$HOME/.cargo/env"
`
	
	// Add Rust init to shell config files
	for _, file := range shellFiles {
		if _, err := os.Stat(file); err == nil {
			// Read current content
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("⚠️ Warning: Could not read %s: %v\n", file, err)
				continue
			}
			
			// Check if Rust is already configured
			if strings.Contains(string(content), ".cargo/env") {
				continue
			}
			
			// Append Rust init
			if err := appendToFile(file, rustInit); err != nil {
				fmt.Printf("⚠️ Warning: Could not update %s: %v\n", file, err)
			} else {
				fmt.Printf("✅ Updated %s with Rust configuration\n", file)
			}
		}
	}
	
	// Add cargo bin to PATH for the current session
	cargoPath := fmt.Sprintf("%s/.cargo/bin", home)
	if _, err := os.Stat(cargoPath); err == nil {
		newPath := cargoPath + ":" + os.Getenv("PATH")
		os.Setenv("PATH", newPath)
		fmt.Printf("✅ Added Rust to PATH: %s\n", cargoPath)
	}
	
	fmt.Println("✅ Rust installed successfully")
	return nil
}

// installLanguage captures the output for one language installation and returns a structured result
func installLanguage(lang string, langPkg Language, packageManagers map[string]string) (string, error) {
	// Install the language with completely suppressed output
	var stdout bytes.Buffer
	
	// Create a custom function to capture output while suppressing it completely
	installFn := func() error {
		// Save original stdout/stderr
		originalStdout := os.Stdout
		originalStderr := os.Stderr
		
		// Create pipes for capturing output
		r, w, _ := os.Pipe()
		
		// Redirect stdout/stderr to our pipe
		os.Stdout = w
		os.Stderr = w
		
		// Restore original stdout when done
		defer func() {
			os.Stdout = originalStdout
			os.Stderr = originalStderr
		}()
		
		// Run the actual install function
		err := langPkg.InstallCmd()
		
		// Close writer and collect output
		w.Close()
		io.Copy(&stdout, r)
		
		return err
	}
	
	err := installFn()
	
	if err != nil {
		return stdout.String(), err
	}
	
	// Install package manager if specified
	if mgr, ok := packageManagers[lang]; ok && mgr != "" {
		// Capture and suppress package manager output too
		pmInstallFn := func() error {
			// Save original stdout/stderr
			originalStdout := os.Stdout
			originalStderr := os.Stderr
			
			// Redirect to /dev/null or equivalent
			devNull, _ := os.Open(os.DevNull)
			os.Stdout = devNull
			os.Stderr = devNull
			
			// Restore original stdout when done
			defer func() {
				os.Stdout = originalStdout
				os.Stderr = originalStderr
				devNull.Close()
			}()
			
			// Install the package manager
			return installPackageManager(lang, mgr)
		}
		
		// Run the package manager installation with suppressed output
		if err := pmInstallFn(); err != nil {
			// Log the error but don't fail the language install
			fmt.Printf("⚠️ Failed to set up package manager %s: %v\n", mgr, err)
		}
	}
	
	return stdout.String(), nil
}

// Package manager installations
func installPackageManager(language, manager string) error {
	// Get HOME directory
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME environment variable not set")
	}
	
	switch language {
	case "node":
		switch manager {
		case "npm":
			// npm comes with Node.js
			return nil
		case "yarn":
			cmd := exec.Command("bash", "-c", 
				fmt.Sprintf(`export NVM_DIR="%s/.nvm" && 
				[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && 
				npm install -g yarn`, home))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		case "pnpm":
			cmd := exec.Command("bash", "-c", 
				fmt.Sprintf(`export NVM_DIR="%s/.nvm" && 
				[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && 
				npm install -g pnpm`, home))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	case "python":
		switch manager {
		case "pip":
			// pip comes with Python
			return nil
		case "pipenv":
			cmd := exec.Command("bash", "-c", 
				fmt.Sprintf(`export PYENV_ROOT="%s/.pyenv" && 
				export PATH="$PYENV_ROOT/bin:$PATH" && 
				eval "$(~/.pyenv/bin/pyenv init -)" && 
				python -m pip install --user pipenv`, home))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		case "poetry":
			cmd := exec.Command("bash", "-c", 
				fmt.Sprintf(`export PYENV_ROOT="%s/.pyenv" && 
				export PATH="$PYENV_ROOT/bin:$PATH" && 
				eval "$(~/.pyenv/bin/pyenv init -)" && 
				python -m pip install --user poetry`, home))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}
	
	return fmt.Errorf("unknown package manager %s for language %s", manager, language)
}

// installLsdFromBinary installs lsd from GitHub releases when package manager doesn't have it
func installLsdFromBinary() error {
	// Determine OS and architecture
	osInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}
	
	// Only proceed if we're on Linux or macOS or Windows
	var archStr string
	switch runtime.GOARCH {
	case "amd64":
		archStr = "x86_64"
	case "arm64":
		archStr = "aarch64"
	default:
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
	
	tempDir, err := os.MkdirTemp("", "lsd-install")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	version := "1.0.0" // Latest stable version
	var downloadURL string
	var binaryName string
	
	// Build the download URL based on the detected OS
	switch osInfo.OS {
	case platform.Linux:
		downloadURL = fmt.Sprintf("https://github.com/lsd-rs/lsd/releases/download/v%s/lsd-v%s-%s-unknown-linux-gnu.tar.gz", version, version, archStr)
		binaryName = "lsd"
	case platform.MacOS:
		downloadURL = fmt.Sprintf("https://github.com/lsd-rs/lsd/releases/download/v%s/lsd-v%s-%s-apple-darwin.tar.gz", version, version, archStr)
		binaryName = "lsd"
	case platform.Windows:
		downloadURL = fmt.Sprintf("https://github.com/lsd-rs/lsd/releases/download/v%s/lsd-v%s-%s-pc-windows-msvc.zip", version, version, archStr)
		binaryName = "lsd.exe"
	default:
		return fmt.Errorf("unsupported OS: %s", osInfo.OS)
	}
	
	fmt.Printf("📥 Downloading lsd from %s\n", downloadURL)
	
	// Download the archive
	archivePath := filepath.Join(tempDir, "lsd-archive")
	var downloadCmd *exec.Cmd
	
	if osInfo.OS == platform.Windows {
		// Use PowerShell on Windows
		downloadCmd = exec.Command("powershell", "-Command", 
			fmt.Sprintf("Invoke-WebRequest -Uri '%s' -OutFile '%s'", downloadURL, archivePath))
	} else {
		// Use curl on Unix systems
		downloadCmd = exec.Command("curl", "-L", downloadURL, "-o", archivePath)
	}
	
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to download lsd: %w", err)
	}
	
	// Extract the archive
	extractDir := filepath.Join(tempDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	
	var extractCmd *exec.Cmd
	if strings.HasSuffix(downloadURL, ".tar.gz") {
		extractCmd = exec.Command("tar", "-xzf", archivePath, "-C", extractDir)
	} else if strings.HasSuffix(downloadURL, ".zip") {
		extractCmd = exec.Command("unzip", archivePath, "-d", extractDir)
	} else {
		return fmt.Errorf("unsupported archive format")
	}
	
	extractCmd.Stdout = os.Stdout
	extractCmd.Stderr = os.Stderr
	if err := extractCmd.Run(); err != nil {
		return fmt.Errorf("failed to extract lsd: %w", err)
	}
	
	// Find the binary
	var binaryPath string
	err = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == binaryName {
			binaryPath = path
			return filepath.SkipDir
		}
		return nil
	})
	
	if binaryPath == "" {
		return fmt.Errorf("could not find lsd binary in extracted files")
	}
	
	// Install to /usr/local/bin or equivalent
	var installPath string
	if osInfo.OS == platform.Windows {
		// On Windows, install to a directory in PATH
		installPath = "C:\\Program Files\\lsd"
		if err := os.MkdirAll(installPath, 0755); err != nil {
			return fmt.Errorf("failed to create installation directory: %w", err)
		}
		installPath = filepath.Join(installPath, binaryName)
	} else {
		// On Unix systems, install to /usr/local/bin
		installPath = "/usr/local/bin/lsd"
	}
	
	// Copy the binary to the install location
	fmt.Printf("📦 Installing lsd to %s\n", installPath)
	
	var installCmd *exec.Cmd
	if osInfo.OS == platform.Windows {
		installCmd = exec.Command("powershell", "-Command", 
			fmt.Sprintf("Copy-Item -Path '%s' -Destination '%s'", binaryPath, installPath))
	} else {
		installCmd = exec.Command("sudo", "cp", binaryPath, installPath)
	}
	
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install lsd: %w", err)
	}
	
	// Make it executable (Unix only)
	if osInfo.OS != platform.Windows {
		chmodCmd := exec.Command("sudo", "chmod", "+x", installPath)
		chmodCmd.Stdout = os.Stdout
		chmodCmd.Stderr = os.Stderr
		if err := chmodCmd.Run(); err != nil {
			return fmt.Errorf("failed to make lsd executable: %w", err)
		}
	}
	
	fmt.Println("✅ Successfully installed lsd")
	return nil
}