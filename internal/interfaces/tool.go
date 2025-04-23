package interfaces

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Tool represents a development tool that can be installed
type Tool struct {
	Name        string
	Description string
	Category    string
	Tags        []string
	
	// Package management
	PackageName string  // Package name for the default package manager
	PackageNames struct {
		APT    string `yaml:"apt"`
		Brew   string `yaml:"brew"`
		DNF    string `yaml:"dnf"`
		Pacman string `yaml:"pacman"`
	}

	Version string
	VerifyCommand string
	PostInstall []struct {
		Command string
		Description string
	}

	ShellConfig struct {
		Aliases   map[string]string
		Functions map[string]string
		Env       map[string]string
	}

	Files []struct {
		Source      string
		Destination string
		Type        string
		Permissions int
		Content     string
	}
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