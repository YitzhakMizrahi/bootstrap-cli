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
	Tags        []string `yaml:"tags,omitempty"`
	Languages   []string `yaml:"languages,omitempty"`  // List of supported languages
	
	// Package management
	PackageNames struct {
		APT    string `yaml:"apt"`
		Brew   string `yaml:"brew"`
		DNF    string `yaml:"dnf"`
		Pacman string `yaml:"pacman"`
	} `yaml:"package_names"`

	Version            string   `yaml:"version"`
	SystemDependencies []string `yaml:"system_dependencies,omitempty"`
	Dependencies       []struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Optional bool   `yaml:"optional,omitempty"`
	} `yaml:"dependencies,omitempty"`
	VerifyCommand string `yaml:"verify_command"`
	PostInstall   []struct {
		Command     string `yaml:"command"`
		Description string `yaml:"description"`
	} `yaml:"post_install,omitempty"`

	ShellConfig struct {
		Aliases   map[string]string `yaml:"aliases,omitempty"`
		Env       map[string]string `yaml:"env,omitempty"`
		Path      []string         `yaml:"path,omitempty"`
		Functions map[string]string `yaml:"functions,omitempty"`
	} `yaml:"shell_config,omitempty"`

	RequiresRestart bool   `yaml:"requires_restart,omitempty"`
	InstallPath     string `yaml:"install_path,omitempty"`
	ConfigFiles     []struct {
		Source      string `yaml:"source"`
		Destination string `yaml:"destination"`
		Template    bool   `yaml:"template,omitempty"`
		Mode        string `yaml:"mode,omitempty"`
	} `yaml:"config_files,omitempty"`
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

// SupportsLanguage checks if the tool supports a given language
func (t *Tool) SupportsLanguage(language string) bool {
	if t.Languages == nil {
		return false
	}
	for _, lang := range t.Languages {
		if lang == language {
			return true
		}
	}
	return false
} 