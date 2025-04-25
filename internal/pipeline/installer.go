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
}

// NewInstaller creates a new installer instance
func NewInstaller(platform *Platform, pkgManager PackageManager) (*Installer, error) {
	context := NewInstallationContext(platform, pkgManager)
	pipeline := NewInstallationPipeline()
	
	return &Installer{
		Context:  context,
		Pipeline: pipeline,
		Logger:   log.New(log.Writer(), "[Installer] ", log.LstdFlags),
	}, nil
}

// Install installs a tool using the pipeline-based approach
func (i *Installer) Install(tool *Tool) error {
	i.Logger.Printf("Starting installation of %s", tool.Name)
	
	// Create installation steps
	steps := tool.GenerateInstallationSteps(i.Context.Platform, i.Context)
	
	// Add steps to pipeline
	for _, step := range steps {
		i.Pipeline.AddStep(step)
	}
	
	// Execute pipeline
	if err := i.Pipeline.Execute(); err != nil {
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

// InstallMultiple installs multiple tools in parallel
func (i *Installer) InstallMultiple(tools []*Tool) error {
	i.Logger.Printf("Starting installation of %d tools", len(tools))
	
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