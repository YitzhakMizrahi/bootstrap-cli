package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// Installer manages the installation of tools using a pipeline-based approach
type Installer struct {
	Context  *InstallationContext
	Pipeline *InstallationPipeline
	Logger   interfaces.Logger
	// Add a field to hold the read-end of the channel for the UI
	ProgressChan <-chan ProgressEvent
	progressChanWriter chan<- ProgressEvent // Internal write-end for the pipeline
}

// NewInstaller creates a new installer instance
func NewInstaller(platform *Platform, pkgManager PackageManager) (*Installer, error) {
	// Create a buffered channel for progress events
	progChan := make(chan ProgressEvent, 100)

	// Create context first, passing the channel
	context := NewInstallationContext(platform, pkgManager, progChan)

	// Pipeline creation is handled within InstallSelections/Install now
	// pipeline := NewInstallationPipeline(context) // Remove pipeline creation here
	
	return &Installer{
		Context:  context,
		// Pipeline: pipeline, // Remove field storage if pipeline is per-execution
		Logger:   context.Logger.(interfaces.Logger), // Use interface type directly
		ProgressChan: progChan, // Expose read-end
		progressChanWriter: progChan, // Keep write-end internally
	}, nil
}

// Install installs a tool using the pipeline-based approach
func (i *Installer) Install(tool *Tool) error {
	i.Logger.Info("Starting installation of %s", tool.Name)
	
	// Generate steps (including dependency resolution)
	steps := tool.GenerateInstallationSteps(i.Context.Platform, i.Context, false) 
	
	// Create and execute a pipeline specifically for this single tool install
	p := NewInstallationPipeline(i.Context) // Pass the shared context
	for _, step := range steps {
		p.AddStep(step)
	}
	if err := p.Execute(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Setup environment
	if err := i.Context.SetupEnvironment(tool); err != nil {
		return fmt.Errorf("environment setup failed: %w", err)
	}
	
	// Verify installation
	if err := i.Context.VerifyInstallation(tool); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	
	i.Logger.Info("Successfully installed %s", tool.Name)
	return nil
}

// InstallMultipleUnsafe_DEPRECATED installs multiple tools in parallel *without* resolving inter-tool dependencies correctly.
// This is generally unsafe if tools depend on each other. Use InstallSelections instead.
// TODO: Remove this method once InstallSelections is fully integrated.
func (i *Installer) InstallMultipleUnsafe_DEPRECATED(tools []*Tool) error {
	i.Logger.Info("[DEPRECATED] Starting unsafe parallel installation of %d tools", len(tools))
	
	// Create a channel to collect errors
	errChan := make(chan error, len(tools))
	
	// Install each tool in a goroutine
	for _, tool := range tools {
		go func(t *Tool) {
			errChan <- i.Install(t)
		}(tool)
	}
	
	// Collect results
	var errors []error
	for range tools {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}
	
	// Return combined error if any installations failed
	if len(errors) > 0 {
		return fmt.Errorf("some installations failed: %v", errors)
	}
	
	i.Logger.Info("Successfully installed all tools")
	return nil
}

// InstallSelections installs a collection of selected items, respecting dependencies.
func (i *Installer) InstallSelections(
	selectedTools []*Tool, 
	manageDotfiles bool, dotfilesRepoURL string,
	selectedFonts []*interfaces.Font,
	selectedLanguages []*interfaces.Language,
	selectedShell *interfaces.Shell,
) error { 
	if len(selectedTools) == 0 && !manageDotfiles && len(selectedFonts) == 0 && len(selectedLanguages) == 0 && selectedShell == nil {
		i.Logger.Info("No items selected for installation.")
		return nil
	}
	i.Logger.Info("Starting dependency-aware installation...")

	// 1. Build Combined Dependency Graph for Tools
	// TODO: Include dependencies from fonts, languages, dotfiles (e.g., git)
	i.Context.dependencyGraph = NewDependencyGraph()
	i.Context.installedTools = make(map[string]bool)
	toolMap := make(map[string]*Tool)
	for _, tool := range selectedTools {
		i.Context.AddTool(tool) 
		toolMap[tool.Name] = tool
		i.Logger.Info("Added tool %s to graph with dependencies: %v", tool.Name, tool.Dependencies)
	}
	// TODO: Add implicit dependencies (like git for dotfiles, curl/unzip for fonts?)

	// 2. Calculate Overall Installation Order (currently only based on tools)
	installOrder, err := i.Context.dependencyGraph.GetInstallOrder()
	if err != nil {
		return fmt.Errorf("failed to calculate installation order: %w", err)
	}
	i.Logger.Info("Calculated installation order: %v", installOrder)

	// 3. Create Single Pipeline & Generate Ordered Steps
	// Create the pipeline using the installer's context (which has the channel)
	pipeline := NewInstallationPipeline(i.Context)
	// No need to set Logger/State again as NewInstallationPipeline does it from context
	i.Pipeline = pipeline // Store the pipeline instance for this run? Or just execute?

	addedSteps := make(map[string]bool) 

	// Add Tool Steps in Order
	for _, toolName := range installOrder {
        if _, alreadyAdded := addedSteps[toolName]; alreadyAdded {
			continue
		}
		toolToInstall, exists := toolMap[toolName]
		if !exists {
            // TODO: Handle loading missing dependency tool definitions
            i.Logger.Info("Warning: Tool %s found in install order but not in initial selection. Skipping its steps.", toolName)
            continue
        }
		i.Logger.Info("Generating installation steps for: %s", toolName)
		steps := toolToInstall.GenerateInstallationSteps(i.Context.Platform, i.Context, true) // skip dependency step
		for _, step := range steps {
			i.Pipeline.AddStep(step)
			i.Logger.Info("  Added step: %s", step.Name)
		}
		addedSteps[toolName] = true
	}

	// Add Font Steps
	if len(selectedFonts) > 0 {
		i.Logger.Info("Adding steps for %d fonts...", len(selectedFonts))
		for _, font := range selectedFonts {
			i.Logger.Info("Generating steps for font: %s", font.Name)
			fontSteps := GenerateFontInstallSteps(font, i.Context.Platform) // Call font step generator
			for _, step := range fontSteps {
				i.Pipeline.AddStep(step)
				i.Logger.Info("  Added font step: %s", step.Name)
			}
		}
	}
	
	// Add Language Steps 
	if len(selectedLanguages) > 0 {
		i.Logger.Info("Adding steps for %d languages...", len(selectedLanguages))
		for _, lang := range selectedLanguages {
			i.Logger.Info("Generating steps for language: %s", lang.Name)
			// Pass context to generator as it might be needed for strategy decisions
			langSteps := GenerateLanguageInstallSteps(lang, i.Context) 
			for _, step := range langSteps {
				i.Pipeline.AddStep(step)
				i.Logger.Info("  Added language step: %s", step.Name)
			}
		}
	}

	// Add Dotfiles Steps (if selected)
	if manageDotfiles && dotfilesRepoURL != "" {
		i.Logger.Info("Adding dotfiles clone steps for repo: %s", dotfilesRepoURL)
		// TODO: Determine appropriate targetDir (e.g., ~/.dotfiles)
		homeDir, _ := os.UserHomeDir() // Handle potential error
		targetDir := filepath.Join(homeDir, ".dotfiles") // Example target
		dotfileSteps := GenerateDotfileCloneSteps(dotfilesRepoURL, targetDir)
		for _, step := range dotfileSteps {
			i.Pipeline.AddStep(step)
			i.Logger.Info("  Added dotfiles step: %s", step.Name)
		}
		// TODO: Add symlinking steps after clone
	}

	// Add Shell Configuration Steps (if selected)
	if selectedShell != nil {
		i.Logger.Info("Adding steps for shell configuration: %s", selectedShell.Name)
		shellSteps := GenerateShellConfigSteps(selectedShell, i.Context)
		for _, step := range shellSteps {
			i.Pipeline.AddStep(step)
			i.Logger.Info("  Added shell config step: %s", step.Name)
		}
	}

	// 4. Execute the single, ordered pipeline
	i.Logger.Info("Executing combined installation pipeline with %d steps...", len(i.Pipeline.Steps))
	if err := i.Pipeline.Execute(); err != nil {
		return fmt.Errorf("installation pipeline failed: %w", err)
	}

	// 5. Final Environment Setup ?
	i.Logger.Info("Installation pipeline completed successfully.")
	return nil
}

// Uninstall removes a tool and its dependencies
func (i *Installer) Uninstall(tool *Tool) error {
	i.Logger.Info("Starting uninstallation of %s", tool.Name)
	
	// Create uninstallation steps
	steps := []InstallationStep{
		{
			Name: "Remove package",
			Action: func(ctx *InstallationContext) error {
				strategy := tool.GetInstallStrategy(ctx.Platform)
				if pkgName, ok := strategy.PackageNames[ctx.Platform.PackageManager]; ok {
					return ctx.PackageManager.Uninstall(pkgName)
				}
				return ctx.PackageManager.Uninstall(tool.Name)
			},
		},
		{
			Name: "Clean up configuration",
			Action: func(ctx *InstallationContext) error {
				ctx.Logger.Info("TODO: Implement configuration cleanup for %s", tool.Name)
				return nil
			},
		},
	}
	
	// Create and execute a pipeline for uninstall
	p := NewInstallationPipeline(i.Context) // Pass the shared context
	for _, step := range steps {
		p.AddStep(step)
	}
	if err := p.Execute(); err != nil {
		return fmt.Errorf("uninstallation failed: %w", err)
	}
	
	i.Logger.Info("Successfully uninstalled %s", tool.Name)
	return nil
}

// GetStatus returns the current installation status of a tool
func (i *Installer) GetStatus(_ *Tool) string {
	return i.Context.State.Status
}

// GetProgress returns the current progress of the installation pipeline
func (i *Installer) GetProgress() string {
	return i.Pipeline.GetProgress()
}

// executeWithRetry executes a command with retries
func executeWithRetry(cmd *exec.Cmd, maxRetries int, delay time.Duration) error {
	var lastErr error
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(delay)
		}
		
		if err := cmd.Run(); err != nil {
			lastErr = err
			continue
		}
		
		return nil
	}
	
	return fmt.Errorf("command failed after %d attempts: %v", maxRetries, lastErr)
} 