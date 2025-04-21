package languages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PythonManager handles Python installation and version management
type PythonManager struct {
	InstallPath      string
	CurrentVersion   string
	AvailableVersions []string
}

// NewPythonManager creates a new PythonManager
func NewPythonManager(installPath string) *PythonManager {
	return &PythonManager{
		InstallPath:      installPath,
		AvailableVersions: []string{},
	}
}

// InstallPyenv installs the Python Version Manager (pyenv)
func (p *PythonManager) InstallPyenv() error {
	// Check if pyenv is already installed
	if p.isPyenvInstalled() {
		return fmt.Errorf("pyenv is already installed")
	}

	// Create the installation directory if it doesn't exist
	if err := os.MkdirAll(p.InstallPath, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Clone pyenv repository
	cmd := exec.Command("git", "clone", "https://github.com/pyenv/pyenv.git", filepath.Join(p.InstallPath, ".pyenv"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone pyenv repository: %w\nOutput: %s", err, string(output))
	}

	// Add pyenv to PATH in shell configuration files
	shellConfigFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".bashrc"),
		filepath.Join(os.Getenv("HOME"), ".zshrc"),
	}

	for _, configFile := range shellConfigFiles {
		if _, err := os.Stat(configFile); err == nil {
			// Check if pyenv is already in the config file
			content, err := os.ReadFile(configFile)
			if err != nil {
				continue
			}

			if !strings.Contains(string(content), "pyenv") {
				// Append pyenv configuration to the shell config file
				pyenvConfig := fmt.Sprintf(`
# pyenv configuration
export PYENV_ROOT="%s/.pyenv"
export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init --path)"
eval "$(pyenv init -)"
`, p.InstallPath)

				f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					continue
				}
				defer f.Close()

				if _, err := f.WriteString(pyenvConfig); err != nil {
					continue
				}
			}
		}
	}

	return nil
}

// InstallPython installs a specific version of Python
func (p *PythonManager) InstallPython(version string) error {
	// Check if pyenv is installed
	if !p.isPyenvInstalled() {
		return fmt.Errorf("pyenv is not installed, install it first")
	}

	// Install the specified version of Python
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv install %s", p.getShellConfig(), version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Python %s: %w\nOutput: %s", version, err, string(output))
	}

	// Set the installed version as the current version
	cmd = exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv global %s", p.getShellConfig(), version))
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set Python %s as current version: %w\nOutput: %s", version, err, string(output))
	}

	p.CurrentVersion = version
	return nil
}

// InstallLatestStable installs the latest stable version of Python
func (p *PythonManager) InstallLatestStable() error {
	// Check if pyenv is installed
	if !p.isPyenvInstalled() {
		return fmt.Errorf("pyenv is not installed, install it first")
	}

	// Get the latest stable version
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv install --list | grep '^  [0-9]\\.[0-9]\\.[0-9]' | sort -V | tail -n 1", p.getShellConfig()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get latest stable Python version: %w\nOutput: %s", err, string(output))
	}

	latestVersion := strings.TrimSpace(string(output))
	if latestVersion == "" {
		return fmt.Errorf("failed to determine latest stable Python version")
	}

	// Install the latest stable version
	return p.InstallPython(latestVersion)
}

// InstallSystemPython uses the system Python
func (p *PythonManager) InstallSystemPython() error {
	// Check if pyenv is installed
	if !p.isPyenvInstalled() {
		return fmt.Errorf("pyenv is not installed, install it first")
	}

	// Set the system Python as the current version
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv global system", p.getShellConfig()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set system Python as current version: %w\nOutput: %s", err, string(output))
	}

	p.CurrentVersion = "system"
	return nil
}

// ListAvailableVersions lists all available Python versions
func (p *PythonManager) ListAvailableVersions() ([]string, error) {
	// Check if pyenv is installed
	if !p.isPyenvInstalled() {
		return nil, fmt.Errorf("pyenv is not installed, install it first")
	}

	// List all available Python versions
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv install --list | grep '^  [0-9]\\.[0-9]\\.[0-9]'", p.getShellConfig()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list available Python versions: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract version numbers
	versions := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		version := strings.TrimSpace(line)
		if version != "" {
			versions = append(versions, version)
		}
	}

	p.AvailableVersions = versions
	return versions, nil
}

// ListInstalledVersions lists all installed Python versions
func (p *PythonManager) ListInstalledVersions() ([]string, error) {
	// Check if pyenv is installed
	if !p.isPyenvInstalled() {
		return nil, fmt.Errorf("pyenv is not installed, install it first")
	}

	// List all installed Python versions
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv versions --bare", p.getShellConfig()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed Python versions: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract version numbers
	versions := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		version := strings.TrimSpace(line)
		if version != "" && version != "system" {
			versions = append(versions, version)
		}
	}

	return versions, nil
}

// UseVersion sets the current Python version
func (p *PythonManager) UseVersion(version string) error {
	// Check if pyenv is installed
	if !p.isPyenvInstalled() {
		return fmt.Errorf("pyenv is not installed, install it first")
	}

	// Set the current Python version
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv global %s", p.getShellConfig(), version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set Python version to %s: %w\nOutput: %s", version, err, string(output))
	}

	p.CurrentVersion = version
	return nil
}

// GetCurrentVersion returns the current Python version
func (p *PythonManager) GetCurrentVersion() (string, error) {
	// Check if pyenv is installed
	if !p.isPyenvInstalled() {
		return "", fmt.Errorf("pyenv is not installed, install it first")
	}

	// Get the current Python version
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && pyenv version", p.getShellConfig()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current Python version: %w\nOutput: %s", err, string(output))
	}

	version := strings.TrimSpace(string(output))
	if strings.Contains(version, " ") {
		version = strings.Split(version, " ")[0]
	}

	p.CurrentVersion = version
	return version, nil
}

// isPyenvInstalled checks if pyenv is installed
func (p *PythonManager) isPyenvInstalled() bool {
	// Check if the pyenv directory exists
	pyenvDir := filepath.Join(p.InstallPath, ".pyenv")
	if _, err := os.Stat(pyenvDir); os.IsNotExist(err) {
		return false
	}

	// Check if the pyenv executable exists
	pyenvExec := filepath.Join(pyenvDir, "bin", "pyenv")
	if _, err := os.Stat(pyenvExec); os.IsNotExist(err) {
		return false
	}

	return true
}

// getShellConfig returns the path to the shell configuration file
func (p *PythonManager) getShellConfig() string {
	// Try to determine the current shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	// Determine the shell configuration file
	var configFile string
	if strings.Contains(shell, "zsh") {
		configFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
	} else {
		configFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
	}

	return configFile
} 