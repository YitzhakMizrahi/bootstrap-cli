package interfaces

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Tool represents a development tool that can be installed
type Tool struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	
	// Package management
	PackageName string  `yaml:"package_name"` // Package name for the default package manager
	PackageNames struct {
		APT    string `yaml:"apt"`
		Brew   string `yaml:"brew"`
		DNF    string `yaml:"dnf"`
		Pacman string `yaml:"pacman"`
	} `yaml:"package_names"`

	Version string `yaml:"version"`
	VerifyCommand string `yaml:"verify_command"`
	SystemDependencies []string `yaml:"system_dependencies"`
	Dependencies []string `yaml:"dependencies"`
	PostInstall []struct {
		Command string `yaml:"command"`
		Description string `yaml:"description"`
	} `yaml:"post_install"`

	ShellConfig struct {
		Aliases   map[string]string `yaml:"aliases"`
		Functions map[string]string `yaml:"functions"`
		Env       map[string]string `yaml:"env"`
	} `yaml:"shell_config"`

	Files []struct {
		Source      string `yaml:"source"`
		Destination string `yaml:"destination"`
		Type        string `yaml:"type"`
		Permissions int    `yaml:"permissions"`
		Content     string `yaml:"content"`
	} `yaml:"files"`
}

// ToolInstaller represents a tool installation service
type ToolInstaller interface {
	Install(tool *Tool) error
	Verify(tool *Tool) error
	IsInstalled(tool *Tool) bool
}

// NewInstaller creates a new tool installer
func NewInstaller(pm PackageManager) ToolInstaller {
	return &toolInstaller{
		pm: pm,
	}
}

type toolInstaller struct {
	pm PackageManager
}

func (i *toolInstaller) Install(tool *Tool) error {
	// Try package-manager specific package name first
	var packageName string
	switch i.pm.GetName() {
	case "apt":
		packageName = tool.PackageNames.APT
	case "brew":
		packageName = tool.PackageNames.Brew
	case "dnf":
		packageName = tool.PackageNames.DNF
	case "pacman":
		packageName = tool.PackageNames.Pacman
	}

	// Fall back to default package name if specific one not found
	if packageName == "" {
		packageName = tool.PackageName
	}

	return i.pm.Install(packageName)
}

func (i *toolInstaller) Verify(tool *Tool) error {
	if tool.VerifyCommand == "" {
		return nil
	}
	return runCommand(tool.VerifyCommand)
}

func (i *toolInstaller) IsInstalled(tool *Tool) bool {
	// Try package-manager specific package name first
	var packageName string
	switch i.pm.GetName() {
	case "apt":
		packageName = tool.PackageNames.APT
	case "brew":
		packageName = tool.PackageNames.Brew
	case "dnf":
		packageName = tool.PackageNames.DNF
	case "pacman":
		packageName = tool.PackageNames.Pacman
	}

	// Fall back to default package name if specific one not found
	if packageName == "" {
		packageName = tool.PackageName
	}

	return i.pm.IsInstalled(packageName)
}

// runCommand executes a shell command
func runCommand(cmd string) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	
	// Create a command with the parts
	command := exec.Command(parts[0], parts[1:]...)
	
	// Capture output
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	
	// Run the command
	err := command.Run()
	if err != nil {
		return fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}
	
	return nil
} 