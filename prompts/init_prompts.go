package prompts

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/platform"
	"github.com/YitzhakMizrahi/bootstrap-cli/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Step represents a stage in the setup wizard
type Step int

const (
	StepShell Step = iota
	StepPluginManager
	StepPrompt
	StepCLITools
	StepLanguages
	StepPackageManagers
	StepEditors
	StepDotfilesPath
	StepOptions
	StepConfirm
	StepDone
)

// Option represents a selectable option
type Option struct {
	Title       string
	Description string
	Value       string
	Selected    bool // For multi-select
}

// wizardModel represents the state of the setup wizard
type wizardModel struct {
	platformInfo platform.Info
	steps        []Step
	currentStep  int
	config       types.UserConfig
	
	// Options for each step
	shellOptions      []Option
	pluginOptions     []Option
	promptOptions     []Option
	cliToolsOptions   []Option
	languagesOptions  []Option
	editorsOptions    []Option
	
	// Cursors for navigation
	shellCursor      int
	pluginCursor     int
	promptCursor     int
	
	// Text input for dotfiles path
	pathInput      textinput.Model
	
	// Options (checkboxes)
	useRelative    bool
	backupExisting bool
	devMode        bool
	
	quitting       bool
	err            error
}

// Initialize the wizard model
func initializeModel() (wizardModel, error) {
	platformInfo, err := platform.Detect()
	if err != nil {
		return wizardModel{}, fmt.Errorf("failed to detect platform: %w", err)
	}
	
	m := wizardModel{
		platformInfo: platformInfo,
		steps: []Step{
			StepShell,
			StepPluginManager,
			StepPrompt,
			StepCLITools,
			StepLanguages,
			StepPackageManagers,
			StepEditors,
			StepDotfilesPath,
			StepOptions,
			StepConfirm,
		},
		currentStep: 0,
		config: types.UserConfig{
			PackageManagers: make(map[string]string),
		},
		useRelative:    true,
		backupExisting: true,
	}
	
	// Setup shell options
	m.shellOptions = []Option{
		{Title: "zsh", Description: "The Z shell - powerful and feature-rich", Value: "zsh"},
		{Title: "bash", Description: "Bourne Again SHell - widely available", Value: "bash"},
		{Title: "fish", Description: "Friendly Interactive SHell - user-friendly defaults", Value: "fish"},
	}
	
	// Initialize prompt options
	m.promptOptions = []Option{
		{Title: "starship", Description: "Cross-shell customizable prompt", Value: "starship"},
		{Title: "powerlevel10k", Description: "Fast and customizable zsh theme", Value: "powerlevel10k"},
		{Title: "pure", Description: "Pretty, minimal and fast zsh prompt", Value: "pure"},
		{Title: "oh-my-posh", Description: "Cross-platform, customizable prompt", Value: "oh-my-posh"},
		{Title: "none", Description: "No custom prompt", Value: "none"},
	}
	
	// Initialize CLI tools options
	m.cliToolsOptions = []Option{
		{Title: "lsd", Description: "The next gen ls command", Value: "lsd", Selected: false},
		{Title: "bat", Description: "A cat clone with wings", Value: "bat", Selected: false},
		{Title: "fzf", Description: "Command-line fuzzy finder", Value: "fzf", Selected: false},
		{Title: "ripgrep", Description: "Fast search tool (rg)", Value: "ripgrep", Selected: false},
		{Title: "fd", Description: "Simple, fast alternative to find", Value: "fd", Selected: false},
		{Title: "jq", Description: "Lightweight JSON processor", Value: "jq", Selected: false},
		{Title: "htop", Description: "Interactive process viewer", Value: "htop", Selected: false},
		{Title: "lazygit", Description: "Simple terminal UI for git", Value: "lazygit", Selected: false},
		{Title: "tmux", Description: "Terminal multiplexer", Value: "tmux", Selected: false},
		{Title: "neofetch", Description: "System info written in bash", Value: "neofetch", Selected: false},
	}
	
	// Initialize languages options
	m.languagesOptions = []Option{
		{Title: "node", Description: "Node.js via nvm", Value: "node", Selected: false},
		{Title: "python", Description: "Python via pyenv", Value: "python", Selected: false},
		{Title: "go", Description: "Go programming language", Value: "go", Selected: false},
		{Title: "rust", Description: "Rust programming language via rustup", Value: "rust", Selected: false},
	}
	
	// Initialize editors options
	m.editorsOptions = []Option{
		{Title: "neovim", Description: "Hyperextensible Vim-based editor", Value: "neovim", Selected: false},
		{Title: "lazyvim", Description: "Neovim with LazyVim configuration", Value: "lazyvim", Selected: false},
		{Title: "astronvim", Description: "Neovim with AstroNvim configuration", Value: "astronvim", Selected: false},
		{Title: "vscode", Description: "Visual Studio Code", Value: "vscode", Selected: false},
	}
	
	// Initialize dotfiles path input
	m.pathInput = textinput.New()
	m.pathInput.Placeholder = "~/dotfiles"
	m.pathInput.Focus()
	
	return m, nil
}

// Update handles user input
func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.currentStep > 0 {
				m.currentStep--
				// Reset cursor position when going back
				m.promptCursor = 0
				return m, nil
			}
			return m, tea.Quit
		}
		
		switch m.steps[m.currentStep] {
		case StepShell:
			switch msg.String() {
			case "down", "j":
				m.shellCursor++
				if m.shellCursor >= len(m.shellOptions) {
					m.shellCursor = 0
				}
			case "up", "k":
				m.shellCursor--
				if m.shellCursor < 0 {
					m.shellCursor = len(m.shellOptions) - 1
				}
			case "enter":
				m.config.Shell = m.shellOptions[m.shellCursor].Value
				// Filter plugin options based on shell
				m.filterPluginOptions()
				m.currentStep++
				// Reset cursor for next step
				m.promptCursor = 0
			}
		
		case StepPluginManager:
			switch msg.String() {
			case "down", "j":
				m.pluginCursor++
				if m.pluginCursor >= len(m.pluginOptions) {
					m.pluginCursor = 0
				}
			case "up", "k":
				m.pluginCursor--
				if m.pluginCursor < 0 {
					m.pluginCursor = len(m.pluginOptions) - 1
				}
			case "enter":
				m.config.PluginManager = m.pluginOptions[m.pluginCursor].Value
				m.currentStep++
				// Reset cursor for next step
				m.promptCursor = 0
			}
			
		case StepPrompt:
			switch msg.String() {
			case "down", "j":
				m.promptCursor++
				if m.promptCursor >= len(m.promptOptions) {
					m.promptCursor = 0
				}
			case "up", "k":
				m.promptCursor--
				if m.promptCursor < 0 {
					m.promptCursor = len(m.promptOptions) - 1
				}
			case "enter":
				m.config.Prompt = m.promptOptions[m.promptCursor].Value
				m.currentStep++
				// Reset cursor for next step
				m.promptCursor = 0
			}
			
		case StepCLITools:
			switch msg.String() {
			case "down", "j":
				// Move cursor down without affecting selection
				m.promptCursor++
				if m.promptCursor >= len(m.cliToolsOptions) {
					m.promptCursor = 0
				}
			case "up", "k":
				// Move cursor up without affecting selection
				m.promptCursor--
				if m.promptCursor < 0 {
					m.promptCursor = len(m.cliToolsOptions) - 1
				}
			case " ":
				// Toggle selection of current item only
				m.cliToolsOptions[m.promptCursor].Selected = !m.cliToolsOptions[m.promptCursor].Selected
			case "enter":
				// Collect selected CLI tools
				var tools []string
				for _, opt := range m.cliToolsOptions {
					if opt.Selected {
						tools = append(tools, opt.Value)
					}
				}
				m.config.CLITools = tools
				m.currentStep++
				// Reset cursor for next step
				m.promptCursor = 0
			}
			
		case StepLanguages:
			// Similar to CLI Tools but with language options
			switch msg.String() {
			case "down", "j":
				// Move cursor down
				m.promptCursor++
				if m.promptCursor >= len(m.languagesOptions) {
					m.promptCursor = 0
				}
			case "up", "k":
				// Move cursor up
				m.promptCursor--
				if m.promptCursor < 0 {
					m.promptCursor = len(m.languagesOptions) - 1
				}
			case " ":
				// Toggle selection
				m.languagesOptions[m.promptCursor].Selected = !m.languagesOptions[m.promptCursor].Selected
			case "enter":
				// Collect selected languages
				var langs []string
				for _, opt := range m.languagesOptions {
					if opt.Selected {
						langs = append(langs, opt.Value)
					}
				}
				m.config.Languages = langs
				m.currentStep++
				// Reset cursor for next step
				m.promptCursor = 0
			}
			
		case StepPackageManagers:
			// Handle package managers
			if len(m.config.Languages) == 0 {
				m.currentStep++
				return m, nil
			}
			
			// Assign default package managers
			for _, lang := range m.config.Languages {
				switch lang {
				case "node":
					m.config.PackageManagers["node"] = "pnpm"
				case "python":
					m.config.PackageManagers["python"] = "pip"
				case "go":
					// No package manager for Go
				case "rust":
					// Cargo is the default for Rust
					m.config.PackageManagers["rust"] = "cargo"
				}
			}
			
			if msg.String() == "enter" {
				m.currentStep++
			}
			
		case StepEditors:
			// Similar to Languages
			switch msg.String() {
			case "down", "j":
				// Move cursor down
				m.promptCursor++
				if m.promptCursor >= len(m.editorsOptions) {
					m.promptCursor = 0
				}
			case "up", "k":
				// Move cursor up
				m.promptCursor--
				if m.promptCursor < 0 {
					m.promptCursor = len(m.editorsOptions) - 1
				}
			case " ":
				// Toggle selection
				m.editorsOptions[m.promptCursor].Selected = !m.editorsOptions[m.promptCursor].Selected
			case "enter":
				// Collect selected editors
				var editors []string
				for _, opt := range m.editorsOptions {
					if opt.Selected {
						editors = append(editors, opt.Value)
					}
				}
				m.config.Editors = editors
				m.currentStep++
			}
			
		case StepDotfilesPath:
			// Handle text input for dotfiles path
			var cmd tea.Cmd
			m.pathInput, cmd = m.pathInput.Update(msg)
			
			if msg.Type == tea.KeyEnter {
				path := strings.TrimSpace(m.pathInput.Value())
				if path == "" {
					path = "~/dotfiles" // Default
				}
				m.config.DotfilesPath = path
				m.currentStep++
			}
			
			return m, cmd
			
		case StepOptions:
			// Handle checkbox options
			switch msg.String() {
			case "r":
				m.useRelative = !m.useRelative
			case "b":
				m.backupExisting = !m.backupExisting
			case "d":
				m.devMode = !m.devMode
			case "enter":
				m.config.UseRelativeLinks = m.useRelative
				m.config.BackupExisting = m.backupExisting
				m.config.DevMode = m.devMode
				m.currentStep++
			}
			
		case StepConfirm:
			// Handle confirmation
			if msg.Type == tea.KeyEnter || msg.String() == "y" {
				m.currentStep++
				return m, tea.Quit
			}
			if msg.String() == "n" {
				// Go back to the beginning
				m.currentStep = 0
			}
		
		case StepDone:
			return m, tea.Quit
		}
	}
	
	return m, nil
}

// View renders the current UI state
func (m wizardModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to exit.", m.err)
	}
	
	if m.currentStep >= len(m.steps) {
		return "Configuration complete!\n"
	}
	
	// Add step information header
	stepName := "Unknown"
	switch m.steps[m.currentStep] {
	case StepShell:
		stepName = "Shell Selection"
	case StepPluginManager:
		stepName = "Plugin Manager Selection"
	case StepPrompt:
		stepName = "Prompt Selection"
	case StepCLITools:
		stepName = "CLI Tools Selection"
	case StepLanguages:
		stepName = "Languages Selection"
	case StepPackageManagers:
		stepName = "Package Managers Selection"
	case StepEditors:
		stepName = "Editors Selection"
	case StepDotfilesPath:
		stepName = "Dotfiles Path"
	case StepOptions:
		stepName = "Options"
	case StepConfirm:
		stepName = "Confirmation"
	}
	
	// Create progress indicator
	header := fmt.Sprintf("Step %d/%d: %s", m.currentStep+1, len(m.steps), stepName)
	
	var content strings.Builder
	
	// Main content based on current step
	switch m.steps[m.currentStep] {
	case StepShell:
		content.WriteString("\nSelect your shell:\n\n")
		
		// Display shell options with proper indentation
		for i, option := range m.shellOptions {
			cursor := " "
			if i == m.shellCursor {
				cursor = ">"
			}
			content.WriteString(fmt.Sprintf("%s %s - %s\n", cursor, option.Title, option.Description))
		}
		
	case StepPluginManager:
		content.WriteString("\nSelect plugin manager:\n\n")
		
		// Display plugin manager options with proper indentation
		for i, option := range m.pluginOptions {
			cursor := " "
			if i == m.pluginCursor {
				cursor = ">"
			}
			content.WriteString(fmt.Sprintf("%s %s - %s\n", cursor, option.Title, option.Description))
		}
		
	case StepPrompt:
		content.WriteString("\nSelect prompt theme:\n\n")
		
		// Display prompt options with proper indentation
		for i, option := range m.promptOptions {
			cursor := " "
			if i == m.promptCursor {
				cursor = ">"
			}
			content.WriteString(fmt.Sprintf("%s %s - %s\n", cursor, option.Title, option.Description))
		}
		
	case StepCLITools:
		content.WriteString("\nSelect CLI tools (space to toggle, enter to continue):\n\n")
		
		// Display CLI tools options with checkboxes and indentation
		for i, option := range m.cliToolsOptions {
			cursor := " "
			if i == m.promptCursor {
				cursor = ">"
			}
			
			checkbox := "[ ]"
			if option.Selected {
				checkbox = "[x]"
			}
			
			content.WriteString(fmt.Sprintf("%s %s %s - %s\n", cursor, checkbox, option.Title, option.Description))
		}
		
	case StepLanguages:
		content.WriteString("\nSelect languages to install (space to toggle, enter to continue):\n\n")
		
		// Display languages options with checkboxes and indentation
		for i, option := range m.languagesOptions {
			cursor := " "
			if i == m.promptCursor {
				cursor = ">"
			}
			
			checkbox := "[ ]"
			if option.Selected {
				checkbox = "[x]"
			}
			
			content.WriteString(fmt.Sprintf("%s %s %s - %s\n", cursor, checkbox, option.Title, option.Description))
		}
		
	case StepPackageManagers:
		// Handle the package managers step content
		if len(m.config.Languages) == 0 {
			content.WriteString("\nNo languages selected. Press Enter to continue.\n")
		} else {
			content.WriteString("\nDefault package managers will be installed:\n\n")
			for _, lang := range m.config.Languages {
				var pkgManager string
				switch lang {
				case "node":
					pkgManager = "pnpm (Node.js package manager)"
				case "python":
					pkgManager = "pip (Python package manager)"
				case "go":
					pkgManager = "No separate package manager needed for Go"
				case "rust":
					pkgManager = "cargo (Rust package manager)"
				default:
					pkgManager = "Unknown"
				}
				content.WriteString(fmt.Sprintf("• %s: %s\n", lang, pkgManager))
			}
			content.WriteString("\nPress Enter to continue.\n")
		}
		
	case StepEditors:
		content.WriteString("\nSelect editors (space to toggle, enter to continue):\n\n")
		
		// Display editors options with checkboxes and indentation
		for i, option := range m.editorsOptions {
			cursor := " "
			if i == m.promptCursor {
				cursor = ">"
			}
			
			checkbox := "[ ]"
			if option.Selected {
				checkbox = "[x]"
			}
			
			content.WriteString(fmt.Sprintf("%s %s %s - %s\n", cursor, checkbox, option.Title, option.Description))
		}
		
	case StepDotfilesPath:
		content.WriteString("\nEnter path to your dotfiles repository:\n\n")
		content.WriteString(m.pathInput.View())
		content.WriteString("\n\n(Press Enter to continue)")
		
	case StepOptions:
		content.WriteString("\nOptions:\n\n")
		useRelativeCheckbox := " "
		backupExistingCheckbox := " "
		devModeCheckbox := " "
		
		if m.useRelative {
			useRelativeCheckbox = "x"
		}
		if m.backupExisting {
			backupExistingCheckbox = "x"
		}
		if m.devMode {
			devModeCheckbox = "x"
		}
		
		content.WriteString(fmt.Sprintf("[%s] (r) Use relative symlinks\n", useRelativeCheckbox))
		content.WriteString(fmt.Sprintf("[%s] (b) Backup existing dotfiles\n", backupExistingCheckbox))
		content.WriteString(fmt.Sprintf("[%s] (d) Developer mode (verbose output)\n", devModeCheckbox))
		content.WriteString("\nPress the key in parentheses to toggle, Enter to continue")
		
	case StepConfirm:
		content.WriteString("\nConfiguration Summary:\n\n")
		content.WriteString(fmt.Sprintf("Shell: %s\n", m.config.Shell))
		content.WriteString(fmt.Sprintf("Plugin Manager: %s\n", m.config.PluginManager))
		content.WriteString(fmt.Sprintf("Prompt: %s\n", m.config.Prompt))
		content.WriteString(fmt.Sprintf("CLI Tools: %s\n", strings.Join(m.config.CLITools, ", ")))
		content.WriteString(fmt.Sprintf("Languages: %s\n", strings.Join(m.config.Languages, ", ")))
		content.WriteString(fmt.Sprintf("Editors: %s\n", strings.Join(m.config.Editors, ", ")))
		content.WriteString(fmt.Sprintf("Dotfiles Path: %s\n", m.config.DotfilesPath))
		content.WriteString(fmt.Sprintf("Use Relative Links: %v\n", m.config.UseRelativeLinks))
		content.WriteString(fmt.Sprintf("Backup Existing: %v\n", m.config.BackupExisting))
		content.WriteString(fmt.Sprintf("Dev Mode: %v\n\n", m.config.DevMode))
		content.WriteString("Is this correct? (y/n)")
		
	case StepDone:
		content.WriteString("\n🎉 Configuration complete!\n")
	}
	
	// Navigation help
	navHelp := "\nPress Esc to go back, Ctrl+C to quit"
	
	// Combine all elements
	return fmt.Sprintf("%s\n%s\n%s", header, content.String(), navHelp)
}

// Init initializes the model
func (m wizardModel) Init() tea.Cmd {
	return nil
}

// RunInitWizard starts the interactive setup wizard and returns the resulting config
func RunInitWizard() types.UserConfig {
	model, err := initializeModel()
	if err != nil {
		fmt.Printf("Error initializing wizard: %v\n", err)
		return types.UserConfig{}
	}
	
	// Use a larger terminal size to ensure all items are visible
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running wizard: %v\n", err)
		return types.UserConfig{}
	}
	
	// Return the final config
	m := finalModel.(wizardModel)
	return m.config
}

// filterPluginOptions updates plugin options based on shell selection
func (m *wizardModel) filterPluginOptions() {
	shell := m.config.Shell
	var options []Option
	
	// Add shell-specific plugin managers
	switch shell {
	case "zsh":
		options = append(options, 
			Option{Title: "zinit", Description: "Fast and feature-rich plugin manager for zsh", Value: "zinit"},
			Option{Title: "oh-my-zsh", Description: "Community-driven framework for zsh", Value: "oh-my-zsh"})
	case "bash":
		options = append(options, 
			Option{Title: "bash-it", Description: "Community bash framework", Value: "bash-it"},
			Option{Title: "oh-my-bash", Description: "Community-driven framework for bash", Value: "oh-my-bash"})
	case "fish":
		options = append(options, 
			Option{Title: "fisher", Description: "Plugin manager for fish", Value: "fisher"},
			Option{Title: "oh-my-fish", Description: "Community fish framework", Value: "oh-my-fish"})
	}
	
	// Always add "none" option
	options = append(options, Option{Title: "none", Description: "No plugin manager", Value: "none"})
	
	m.pluginOptions = options
	m.pluginCursor = 0
}