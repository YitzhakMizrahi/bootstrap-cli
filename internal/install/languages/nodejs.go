package languages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// NodeJSManager handles Node.js installation and version management
type NodeJSManager struct {
	InstallPath string
	CurrentVersion string
	AvailableVersions []string
}

// NewNodeJSManager creates a new NodeJSManager
func NewNodeJSManager(installPath string) *NodeJSManager {
	return &NodeJSManager{
		InstallPath: installPath,
		AvailableVersions: []string{},
	}
}

// InstallNVM installs the Node Version Manager
func (n *NodeJSManager) InstallNVM() error {
	// Check if NVM is already installed
	if n.isNVMInstalled() {
		return fmt.Errorf("NVM is already installed")
	}

	// Create the installation directory if it doesn't exist
	if err := os.MkdirAll(n.InstallPath, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Download and install NVM
	cmd := exec.Command("curl", "-o-", "https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download NVM: %w\nOutput: %s", err, string(output))
	}

	// Execute the installation script
	installScript := filepath.Join(n.InstallPath, "nvm-install.sh")
	if err := os.WriteFile(installScript, output, 0755); err != nil {
		return fmt.Errorf("failed to write installation script: %w", err)
	}

	cmd = exec.Command("bash", installScript)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install NVM: %w\nOutput: %s", err, string(output))
	}

	// Clean up the installation script
	if err := os.Remove(installScript); err != nil {
		return fmt.Errorf("failed to remove installation script: %w", err)
	}

	return nil
}

// InstallNodeJS installs a specific version of Node.js
func (n *NodeJSManager) InstallNodeJS(version string) error {
	// Check if NVM is installed
	if !n.isNVMInstalled() {
		return fmt.Errorf("NVM is not installed, install it first")
	}

	// Install the specified version of Node.js
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm install %s", n.getNVMInitScript(), version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Node.js %s: %w\nOutput: %s", version, err, string(output))
	}

	// Set the installed version as the current version
	cmd = exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm use %s", n.getNVMInitScript(), version))
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set Node.js %s as current version: %w\nOutput: %s", version, err, string(output))
	}

	n.CurrentVersion = version
	return nil
}

// InstallLatestLTS installs the latest LTS version of Node.js
func (n *NodeJSManager) InstallLatestLTS() error {
	// Check if NVM is installed
	if !n.isNVMInstalled() {
		return fmt.Errorf("NVM is not installed, install it first")
	}

	// Install the latest LTS version of Node.js
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm install --lts", n.getNVMInitScript()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install latest LTS version of Node.js: %w\nOutput: %s", err, string(output))
	}

	// Set the installed version as the current version
	cmd = exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm use --lts", n.getNVMInitScript()))
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set latest LTS version as current version: %w\nOutput: %s", err, string(output))
	}

	// Get the current version
	cmd = exec.Command("bash", "-c", fmt.Sprintf("source %s && node -v", n.getNVMInitScript()))
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get current Node.js version: %w\nOutput: %s", err, string(output))
	}

	n.CurrentVersion = strings.TrimSpace(string(output))
	return nil
}

// InstallLatestCurrent installs the latest current version of Node.js
func (n *NodeJSManager) InstallLatestCurrent() error {
	// Check if NVM is installed
	if !n.isNVMInstalled() {
		return fmt.Errorf("NVM is not installed, install it first")
	}

	// Install the latest current version of Node.js
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm install node", n.getNVMInitScript()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install latest current version of Node.js: %w\nOutput: %s", err, string(output))
	}

	// Set the installed version as the current version
	cmd = exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm use node", n.getNVMInitScript()))
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set latest current version as current version: %w\nOutput: %s", err, string(output))
	}

	// Get the current version
	cmd = exec.Command("bash", "-c", fmt.Sprintf("source %s && node -v", n.getNVMInitScript()))
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get current Node.js version: %w\nOutput: %s", err, string(output))
	}

	n.CurrentVersion = strings.TrimSpace(string(output))
	return nil
}

// ListAvailableVersions lists all available Node.js versions
func (n *NodeJSManager) ListAvailableVersions() ([]string, error) {
	// Check if NVM is installed
	if !n.isNVMInstalled() {
		return nil, fmt.Errorf("NVM is not installed, install it first")
	}

	// List all available Node.js versions
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm ls-remote", n.getNVMInitScript()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list available Node.js versions: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract version numbers
	versions := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "v") && !strings.Contains(line, "->") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				version := strings.TrimPrefix(parts[0], "v")
				versions = append(versions, version)
			}
		}
	}

	n.AvailableVersions = versions
	return versions, nil
}

// ListInstalledVersions lists all installed Node.js versions
func (n *NodeJSManager) ListInstalledVersions() ([]string, error) {
	// Check if NVM is installed
	if !n.isNVMInstalled() {
		return nil, fmt.Errorf("NVM is not installed, install it first")
	}

	// List all installed Node.js versions
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm ls", n.getNVMInitScript()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed Node.js versions: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract version numbers
	versions := []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "v") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				version := strings.TrimPrefix(parts[0], "v")
				versions = append(versions, version)
			}
		}
	}

	return versions, nil
}

// UseVersion sets the current Node.js version
func (n *NodeJSManager) UseVersion(version string) error {
	// Check if NVM is installed
	if !n.isNVMInstalled() {
		return fmt.Errorf("NVM is not installed, install it first")
	}

	// Set the current Node.js version
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && nvm use %s", n.getNVMInitScript(), version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set Node.js version to %s: %w\nOutput: %s", version, err, string(output))
	}

	n.CurrentVersion = version
	return nil
}

// GetCurrentVersion returns the current Node.js version
func (n *NodeJSManager) GetCurrentVersion() (string, error) {
	// Check if NVM is installed
	if !n.isNVMInstalled() {
		return "", fmt.Errorf("NVM is not installed, install it first")
	}

	// Get the current Node.js version
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && node -v", n.getNVMInitScript()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current Node.js version: %w\nOutput: %s", err, string(output))
	}

	version := strings.TrimSpace(string(output))
	n.CurrentVersion = version
	return version, nil
}

// isNVMInstalled checks if NVM is installed
func (n *NodeJSManager) isNVMInstalled() bool {
	// Check if the NVM directory exists
	nvmDir := filepath.Join(os.Getenv("HOME"), ".nvm")
	if _, err := os.Stat(nvmDir); os.IsNotExist(err) {
		return false
	}

	// Check if the NVM script exists
	nvmScript := filepath.Join(nvmDir, "nvm.sh")
	if _, err := os.Stat(nvmScript); os.IsNotExist(err) {
		return false
	}

	return true
}

// getNVMInitScript returns the path to the NVM initialization script
func (n *NodeJSManager) getNVMInitScript() string {
	return filepath.Join(os.Getenv("HOME"), ".nvm", "nvm.sh")
} 