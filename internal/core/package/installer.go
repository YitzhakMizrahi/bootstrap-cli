package pkgmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Tool represents a tool that can be installed
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PackageName string // The name of the package in the package manager
	Version     string `json:"version"`
	Category    string
	URL         string
	Installed   bool   `json:"installed"`
	InstallTime time.Time
}

// Installer handles tool installation
type Installer struct {
	packageManager PackageManager
	tools          map[string]*Tool
	installDir     string
	stateFile      string
}

// NewInstaller creates a new installer
func NewInstaller(packageManager PackageManager, installDir string, stateFile string) *Installer {
	return &Installer{
		packageManager: packageManager,
		tools:          make(map[string]*Tool),
		installDir:     installDir,
		stateFile:      stateFile,
	}
}

// RegisterTool registers a tool with the installer
func (i *Installer) RegisterTool(tool *Tool) {
	i.tools[tool.Name] = tool
}

// InstallTool installs a single tool
func (i *Installer) InstallTool(toolName string) error {
	tool, exists := i.tools[toolName]
	if !exists {
		return fmt.Errorf("tool %s not found", toolName)
	}

	// Check if the tool is already installed
	if tool.Installed {
		return fmt.Errorf("tool %s is already installed", toolName)
	}

	// Install the tool using the package manager
	if err := i.packageManager.Install(tool.PackageName); err != nil {
		return fmt.Errorf("failed to install tool %s: %w", toolName, err)
	}

	// Update tool status
	tool.Installed = true
	tool.InstallTime = time.Now()

	return i.SaveState()
}

// InstallTools installs multiple tools
func (i *Installer) InstallTools(toolNames []string) error {
	for _, toolName := range toolNames {
		if err := i.InstallTool(toolName); err != nil {
			return fmt.Errorf("failed to install tool %s: %w", toolName, err)
		}
	}
	return nil
}

// UninstallTool uninstalls a single tool
func (i *Installer) UninstallTool(toolName string) error {
	tool, exists := i.tools[toolName]
	if !exists {
		return fmt.Errorf("tool %s not found", toolName)
	}

	// Check if the tool is installed
	if !tool.Installed {
		return fmt.Errorf("tool %s is not installed", toolName)
	}

	// Uninstall the tool using the package manager
	if err := i.packageManager.Uninstall(tool.PackageName); err != nil {
		return fmt.Errorf("failed to uninstall tool %s: %w", toolName, err)
	}

	// Update tool status
	tool.Installed = false
	tool.InstallTime = time.Time{}

	return i.SaveState()
}

// IsToolInstalled checks if a tool is installed
func (i *Installer) IsToolInstalled(toolName string) bool {
	tool, exists := i.tools[toolName]
	if !exists {
		return false
	}
	return tool.Installed
}

// GetTool returns a tool by name
func (i *Installer) GetTool(toolName string) (*Tool, error) {
	tool, exists := i.tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", toolName)
	}
	return tool, nil
}

// ListTools returns a list of all registered tools
func (i *Installer) ListTools() []*Tool {
	tools := make([]*Tool, 0, len(i.tools))
	for _, tool := range i.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ListInstalledTools returns a list of installed tools
func (i *Installer) ListInstalledTools() []*Tool {
	tools := make([]*Tool, 0, len(i.tools))
	for _, tool := range i.tools {
		if tool.Installed {
			tools = append(tools, tool)
		}
	}
	return tools
}

// SaveState saves the current state to disk
func (i *Installer) SaveState() error {
	// Create the install directory if it doesn't exist
	if err := os.MkdirAll(i.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Create a file to store the tool state
	stateFile := filepath.Join(i.installDir, i.stateFile)
	
	// Marshal tools to JSON
	data, err := json.MarshalIndent(i.tools, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tools: %w", err)
	}

	// Write to file
	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}
	
	return nil
}

// LoadState loads the state from disk
func (i *Installer) LoadState() error {
	// Check if the install directory exists
	if _, err := os.Stat(i.installDir); os.IsNotExist(err) {
		return fmt.Errorf("install directory does not exist: %s", i.installDir)
	}

	// Check if the tool state file exists
	stateFile := filepath.Join(i.installDir, i.stateFile)
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return fmt.Errorf("tool state file does not exist: %s", stateFile)
	}
	
	// Read from file
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state file: %w", err)
	}

	if err := json.Unmarshal(data, &i.tools); err != nil {
		return fmt.Errorf("failed to unmarshal tools: %w", err)
	}
	
	return nil
} 