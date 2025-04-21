package tools

import (
	"fmt"

	pkgmanager "github.com/YitzhakMizrahi/bootstrap-cli/internal/core/package"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/progress"
)

// Tool represents a development tool
type Tool struct {
	Name        string
	Description string
	PackageName string
	InstallCmd  string
	IsInstalled bool
}

// Installer handles tool installation
type Installer struct {
	packageManager *pkgmanager.PackageManager
	tools          map[string]*Tool
}

// NewInstaller creates a new tool installer
func NewInstaller(pmType string) *Installer {
	return &Installer{
		packageManager: pkgmanager.New(pmType),
		tools:          make(map[string]*Tool),
	}
}

// RegisterTool registers a tool for installation
func (i *Installer) RegisterTool(tool *Tool) {
	i.tools[tool.Name] = tool
}

// InstallTool installs a single tool
func (i *Installer) InstallTool(name string) error {
	tool, exists := i.tools[name]
	if !exists {
		return fmt.Errorf("tool %s not found", name)
	}

	// Check if the tool is already installed
	installed, err := i.packageManager.IsInstalled(tool.PackageName)
	if err != nil {
		return fmt.Errorf("failed to check if tool is installed: %w", err)
	}

	if installed {
		tool.IsInstalled = true
		fmt.Printf("Tool %s is already installed\n", name)
		return nil
	}

	// Install the tool
	fmt.Printf("Installing tool: %s\n", name)
	if err := i.packageManager.Install(tool.PackageName); err != nil {
		return fmt.Errorf("failed to install tool %s: %w", name, err)
	}

	tool.IsInstalled = true
	fmt.Printf("Tool %s installed successfully\n", name)
	return nil
}

// InstallTools installs multiple tools
func (i *Installer) InstallTools(names []string) error {
	// Create a progress bar
	bar := progress.NewProgressBar(len(names))
	bar.Display()

	// Install each tool
	for idx, name := range names {
		if err := i.InstallTool(name); err != nil {
			return fmt.Errorf("failed to install tool %s: %w", name, err)
		}
		bar.Update(idx + 1)
	}

	bar.Finish()
	return nil
}

// InstallAllTools installs all registered tools
func (i *Installer) InstallAllTools() error {
	// Get all tool names
	names := make([]string, 0, len(i.tools))
	for name := range i.tools {
		names = append(names, name)
	}

	// Install all tools
	return i.InstallTools(names)
}

// ListTools returns a list of all registered tools
func (i *Installer) ListTools() []*Tool {
	tools := make([]*Tool, 0, len(i.tools))
	for _, tool := range i.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetTool returns a tool by name
func (i *Installer) GetTool(name string) (*Tool, error) {
	tool, exists := i.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}
	return tool, nil
}

// IsToolInstalled checks if a tool is installed
func (i *Installer) IsToolInstalled(name string) (bool, error) {
	tool, err := i.GetTool(name)
	if err != nil {
		return false, err
	}

	// Check if the tool is already marked as installed
	if tool.IsInstalled {
		return true, nil
	}

	// Check if the tool is installed on the system
	installed, err := i.packageManager.IsInstalled(tool.PackageName)
	if err != nil {
		return false, fmt.Errorf("failed to check if tool is installed: %w", err)
	}

	// Update the tool's installed status
	tool.IsInstalled = installed
	return installed, nil
} 