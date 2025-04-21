package languages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GoManager handles Go installation and version management
type GoManager struct {
	InstallPath      string
	CurrentVersion   string
	AvailableVersions []string
}

// NewGoManager creates a new GoManager
func NewGoManager(installPath string) *GoManager {
	return &GoManager{
		InstallPath:      installPath,
		AvailableVersions: []string{},
	}
}

// InstallGo installs a specific version of Go
func (g *GoManager) InstallGo(version string) error {
	// Check if Go is already installed
	if g.isGoInstalled() {
		currentVersion, err := g.GetCurrentVersion()
		if err == nil && currentVersion == version {
			return fmt.Errorf("Go version %s is already installed", version)
		}
	}

	// Create the installation directory if it doesn't exist
	if err := os.MkdirAll(g.InstallPath, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Determine the download URL based on the version and OS
	osName := getOSName()
	arch := getArchitecture()
	downloadURL := fmt.Sprintf("https://golang.org/dl/go%s.%s-%s.tar.gz", version, osName, arch)

	// Download and extract Go
	cmd := exec.Command("curl", "-L", "-o", filepath.Join(g.InstallPath, "go.tar.gz"), downloadURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download Go %s: %w\nOutput: %s", version, err, string(output))
	}

	// Extract the archive
	cmd = exec.Command("tar", "-C", g.InstallPath, "-xzf", filepath.Join(g.InstallPath, "go.tar.gz"))
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to extract Go %s: %w\nOutput: %s", version, err, string(output))
	}

	// Clean up the downloaded archive
	if err := os.Remove(filepath.Join(g.InstallPath, "go.tar.gz")); err != nil {
		return fmt.Errorf("failed to remove downloaded archive: %w", err)
	}

	// Set up environment variables
	if err := g.setupEnvironment(); err != nil {
		return fmt.Errorf("failed to set up environment: %w", err)
	}

	g.CurrentVersion = version
	return nil
}

// InstallLatestStable installs the latest stable version of Go
func (g *GoManager) InstallLatestStable() error {
	// Get the latest stable version
	cmd := exec.Command("curl", "-s", "https://golang.org/dl/?mode=json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get latest Go version: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to find the latest stable version
	// This is a simplified approach; in a real implementation, you would parse the JSON
	// and find the latest stable version
	latestVersion := "1.21.0" // This should be dynamically determined

	// Install the latest stable version
	return g.InstallGo(latestVersion)
}

// ListAvailableVersions lists all available Go versions
func (g *GoManager) ListAvailableVersions() ([]string, error) {
	// Get available versions from golang.org
	cmd := exec.Command("curl", "-s", "https://golang.org/dl/?mode=json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list available Go versions: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract version numbers
	// This is a simplified approach; in a real implementation, you would parse the JSON
	versions := []string{"1.20.0", "1.20.1", "1.20.2", "1.21.0", "1.21.1"}

	g.AvailableVersions = versions
	return versions, nil
}

// ListInstalledVersions lists all installed Go versions
func (g *GoManager) ListInstalledVersions() ([]string, error) {
	// Check if the Go installation directory exists
	goDir := filepath.Join(g.InstallPath, "go")
	if _, err := os.Stat(goDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Get the current version
	currentVersion, err := g.GetCurrentVersion()
	if err != nil {
		return []string{}, err
	}

	return []string{currentVersion}, nil
}

// UseVersion sets the current Go version
func (g *GoManager) UseVersion(version string) error {
	// Check if the version is installed
	installedVersions, err := g.ListInstalledVersions()
	if err != nil {
		return fmt.Errorf("failed to list installed versions: %w", err)
	}

	found := false
	for _, v := range installedVersions {
		if v == version {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Go version %s is not installed", version)
	}

	// Set up environment variables for the specified version
	if err := g.setupEnvironment(); err != nil {
		return fmt.Errorf("failed to set up environment: %w", err)
	}

	g.CurrentVersion = version
	return nil
}

// GetCurrentVersion returns the current Go version
func (g *GoManager) GetCurrentVersion() (string, error) {
	// Check if Go is installed
	if !g.isGoInstalled() {
		return "", fmt.Errorf("Go is not installed")
	}

	// Get the current Go version
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current Go version: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract the version
	// Example output: "go version go1.21.0 linux/amd64"
	parts := strings.Split(string(output), " ")
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected output format: %s", string(output))
	}

	version := strings.TrimPrefix(parts[2], "go")
	g.CurrentVersion = version
	return version, nil
}

// isGoInstalled checks if Go is installed
func (g *GoManager) isGoInstalled() bool {
	// Check if the Go installation directory exists
	goDir := filepath.Join(g.InstallPath, "go")
	if _, err := os.Stat(goDir); os.IsNotExist(err) {
		return false
	}

	// Check if the Go executable exists
	goExec := filepath.Join(goDir, "bin", "go")
	if _, err := os.Stat(goExec); os.IsNotExist(err) {
		return false
	}

	return true
}

// setupEnvironment sets up the environment variables for Go
func (g *GoManager) setupEnvironment() error {
	// Get the shell configuration file
	shellConfig := getShellConfig()

	// Check if the Go environment variables are already set
	content, err := os.ReadFile(shellConfig)
	if err != nil {
		return fmt.Errorf("failed to read shell configuration file: %w", err)
	}

	// If the Go environment variables are not set, add them
	if !strings.Contains(string(content), "GOROOT") {
		goEnv := fmt.Sprintf(`
# Go environment variables
export GOROOT="%s/go"
export GOPATH="$HOME/go"
export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"
`, g.InstallPath)

		f, err := os.OpenFile(shellConfig, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open shell configuration file: %w", err)
		}
		defer f.Close()

		if _, err := f.WriteString(goEnv); err != nil {
			return fmt.Errorf("failed to write to shell configuration file: %w", err)
		}
	}

	return nil
}

// getOSName returns the OS name for the download URL
func getOSName() string {
	switch os := runtime.GOOS; os {
	case "linux":
		return "linux"
	case "darwin":
		return "darwin"
	case "windows":
		return "windows"
	default:
		return "linux"
	}
}

// getArchitecture returns the architecture for the download URL
func getArchitecture() string {
	switch arch := runtime.GOARCH; arch {
	case "amd64":
		return "amd64"
	case "arm64":
		return "arm64"
	default:
		return "amd64"
	}
}

// getShellConfig returns the path to the shell configuration file
func getShellConfig() string {
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