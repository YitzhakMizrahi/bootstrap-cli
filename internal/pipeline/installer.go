package pipeline

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Installer manages the installation of tools using a pipeline-based approach
type Installer struct {
	Context  *InstallationContext
	Pipeline *InstallationPipeline
	Logger   *log.Logger
	// Add a field to hold the read-end of the channel for the UI
	ProgressChan <-chan ProgressEvent
	progressChanWriter chan<- ProgressEvent // Internal write-end for the pipeline
}

// NewInstaller creates a new installer instance
func NewInstaller(platform *Platform, pkgManager PackageManager) (*Installer, error) {
	context := NewInstallationContext(platform, pkgManager)
	
	// Create a buffered channel for progress events
	// Buffer size 100 is arbitrary, adjust as needed.
	progChan := make(chan ProgressEvent, 100)

	pipeline := NewInstallationPipeline(progChan) // Pass write-end to pipeline
	
	return &Installer{
		Context:  context,
		Pipeline: pipeline,
		Logger:   log.New(log.Writer(), "[Installer] ", log.LstdFlags),
		ProgressChan: progChan, // Expose read-end
		progressChanWriter: progChan, // Keep write-end internally
	}, nil
}

// Install installs a tool using the pipeline-based approach
func (i *Installer) Install(tool *Tool) error {
	i.Logger.Printf("Starting installation of %s", tool.Name)
	
	// Create installation steps (including dependency resolution for single tool install)
	steps := tool.GenerateInstallationSteps(i.Context.Platform, i.Context, false)
	
	// Use a new pipeline for single install, passing the installer's channel writer
	p := NewInstallationPipeline(i.progressChanWriter)
	p.Logger = i.Logger
	p.State = NewInstallationState() // Give it its own state tracker for single install?
	
	// Add steps to the new pipeline
	for _, step := range steps {
		p.AddStep(step)
	}
	
	// Execute the pipeline
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
	
	i.Logger.Printf("Successfully installed %s", tool.Name)
	return nil
}

// InstallMultipleUnsafe_DEPRECATED installs multiple tools in parallel *without* resolving inter-tool dependencies correctly.
// This is generally unsafe if tools depend on each other. Use InstallSelections instead.
// TODO: Remove this method once InstallSelections is fully integrated.
func (i *Installer) InstallMultipleUnsafe_DEPRECATED(tools []*Tool) error {
	i.Logger.Printf("[DEPRECATED] Starting unsafe parallel installation of %d tools", len(tools))
	
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
	
	i.Logger.Printf("Successfully installed all tools")
	return nil
}

// InstallSelections installs a collection of tools, respecting inter-tool dependencies.
func (i *Installer) InstallSelections(selectedTools []*Tool) error {
	if len(selectedTools) == 0 {
		i.Logger.Printf("No tools selected for installation.")
		return nil
	}
	i.Logger.Printf("Starting dependency-aware installation for %d tools...", len(selectedTools))

	// 1. Build Combined Dependency Graph
	// Use the context's graph, ensure it's clean or create a new one for this run?
	// Let's assume the context provided to the installer is fresh for this operation.
	// If not, context might need a Reset() or similar method.
	// Clear existing graph/installed state for this run
	i.Context.dependencyGraph = NewDependencyGraph()
	i.Context.installedTools = make(map[string]bool)
	toolMap := make(map[string]*Tool)
	for _, tool := range selectedTools {
		i.Context.AddTool(tool) // AddTool populates the dependency graph
		toolMap[tool.Name] = tool
		i.Logger.Printf("Added tool %s to graph with dependencies: %v", tool.Name, tool.Dependencies)
	}

	// 2. Calculate Overall Installation Order
	installOrder, err := i.Context.dependencyGraph.GetInstallOrder()
	if err != nil {
		return fmt.Errorf("failed to calculate installation order: %w", err)
	}
	i.Logger.Printf("Calculated installation order: %v", installOrder)

	// 3. Create Single Pipeline & Generate Ordered Steps
	// Use the installer's main pipeline instance. Ensure it's correctly initialized with the channel.
	i.Pipeline = NewInstallationPipeline(i.progressChanWriter) // Use the installer's write channel
	i.Pipeline.Logger = i.Logger 
	i.Pipeline.State = i.Context.State 

	addedSteps := make(map[string]bool) // Prevent adding steps for the same tool multiple times if it appears in order due to complex deps

	for _, toolName := range installOrder {
		if _, alreadyAdded := addedSteps[toolName]; alreadyAdded {
			continue
		}

		toolToInstall, exists := toolMap[toolName]
		if !exists {
			// This means a dependency was identified in the graph but wasn't in the initial selection list.
			// TODO: Handle this - either load the missing dependency tool config or error out.
			// For now, we assume all necessary tools (including dependencies) are in the initial selectedTools list, 
			// which might be incorrect. A robust implementation needs to load dependencies on demand.
			i.Logger.Printf("Warning: Tool %s found in install order but not in initial selection. Skipping its steps.", toolName)
			continue
		}

		i.Logger.Printf("Generating installation steps for: %s", toolName)
		// Generate steps for this specific tool, *excluding* the dependency resolution step itself.
		steps := toolToInstall.GenerateInstallationSteps(i.Context.Platform, i.Context, true)
		for _, step := range steps {
			i.Pipeline.AddStep(step)
			i.Logger.Printf("  Added step: %s", step.Name)
		}
		addedSteps[toolName] = true
	}

	// 4. Execute the single, ordered pipeline
	i.Logger.Printf("Executing combined installation pipeline with %d steps...", len(i.Pipeline.Steps))
	if err := i.Pipeline.Execute(); err != nil {
		// Execute already logs errors and attempts rollback, and sends events
		return fmt.Errorf("installation pipeline failed: %w", err)
	}

	// 5. Final Environment Setup / Verification? 
	// The original Install method did SetupEnvironment and VerifyInstallation *after* the pipeline.
	// Should we do a global environment setup (e.g., applying shell changes) once at the end?
	// Or rely on individual tool steps? Let's assume individual steps handle their own verification/setup for now.
	// A final shell config apply might be needed.
	i.Logger.Printf("Installation pipeline completed successfully.")
	return nil
}

// Uninstall removes a tool and its dependencies
func (i *Installer) Uninstall(tool *Tool) error {
	i.Logger.Printf("Starting uninstallation of %s", tool.Name)
	
	// Create uninstallation steps
	steps := []InstallationStep{
		{
			Name: "Remove package",
			Action: func() error {
				strategy := tool.GetInstallStrategy(i.Context.Platform)
				if pkgName, ok := strategy.PackageNames[i.Context.Platform.PackageManager]; ok {
					return i.Context.PackageManager.Uninstall(pkgName)
				}
				return i.Context.PackageManager.Uninstall(tool.Name)
			},
		},
		{
			Name: "Clean up configuration",
			Action: func() error {
				// TODO: Implement configuration cleanup
				return nil
			},
		},
	}
	
	// Add steps to pipeline
	for _, step := range steps {
		i.Pipeline.AddStep(step)
	}
	
	// Execute pipeline
	if err := i.Pipeline.Execute(); err != nil {
		return fmt.Errorf("uninstallation failed: %w", err)
	}
	
	i.Logger.Printf("Successfully uninstalled %s", tool.Name)
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