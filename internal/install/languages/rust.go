package languages

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RustManager handles Rust installation and version management
type RustManager struct {
	InstallPath      string
	CurrentVersion   string
	AvailableVersions []string
}

// NewRustManager creates a new RustManager
func NewRustManager(installPath string) *RustManager {
	return &RustManager{
		InstallPath:      installPath,
		AvailableVersions: []string{},
	}
}

// InstallRustup installs rustup, the Rust toolchain installer
func (r *RustManager) InstallRustup() error {
	// Check if rustup is already installed
	if r.isRustupInstalled() {
		return fmt.Errorf("rustup is already installed")
	}

	// Create the installation directory if it doesn't exist
	if err := os.MkdirAll(r.InstallPath, 0755); err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	// Download and run the rustup installer
	cmd := exec.Command("curl", "--proto", "=https", "--tlsv1.2", "-sSf", "https://sh.rustup.rs")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader("1\n") // Choose option 1 for default installation

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install rustup: %w", err)
	}

	// Set up environment variables
	if err := r.setupEnvironment(); err != nil {
		return fmt.Errorf("failed to set up environment: %w", err)
	}

	return nil
}

// InstallRust installs a specific version of Rust
func (r *RustManager) InstallRust(version string) error {
	// Check if rustup is installed
	if !r.isRustupInstalled() {
		if err := r.InstallRustup(); err != nil {
			return fmt.Errorf("failed to install rustup: %w", err)
		}
	}

	// Install the specified version
	cmd := exec.Command("rustup", "install", version)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Rust %s: %w\nOutput: %s", version, err, string(output))
	}

	// Set the specified version as the default
	cmd = exec.Command("rustup", "default", version)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set Rust %s as default: %w\nOutput: %s", version, err, string(output))
	}

	r.CurrentVersion = version
	return nil
}

// InstallLatestStable installs the latest stable version of Rust
func (r *RustManager) InstallLatestStable() error {
	// Check if rustup is installed
	if !r.isRustupInstalled() {
		if err := r.InstallRustup(); err != nil {
			return fmt.Errorf("failed to install rustup: %w", err)
		}
	}

	// Install the latest stable version
	cmd := exec.Command("rustup", "install", "stable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install latest stable Rust: %w\nOutput: %s", err, string(output))
	}

	// Set stable as the default
	cmd = exec.Command("rustup", "default", "stable")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set stable Rust as default: %w\nOutput: %s", err, string(output))
	}

	// Get the current version
	version, err := r.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current Rust version: %w", err)
	}

	r.CurrentVersion = version
	return nil
}

// ListAvailableVersions lists all available Rust versions
func (r *RustManager) ListAvailableVersions() ([]string, error) {
	// Check if rustup is installed
	if !r.isRustupInstalled() {
		return nil, fmt.Errorf("rustup is not installed")
	}

	// Get available versions
	cmd := exec.Command("rustup", "toolchain", "list-remote")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list available Rust versions: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract version numbers
	versions := strings.Split(strings.TrimSpace(string(output)), "\n")
	r.AvailableVersions = versions
	return versions, nil
}

// ListInstalledVersions lists all installed Rust versions
func (r *RustManager) ListInstalledVersions() ([]string, error) {
	// Check if rustup is installed
	if !r.isRustupInstalled() {
		return nil, fmt.Errorf("rustup is not installed")
	}

	// Get installed versions
	cmd := exec.Command("rustup", "toolchain", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed Rust versions: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract version numbers
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	versions := make([]string, 0, len(lines))
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 0 {
			versions = append(versions, parts[0])
		}
	}

	return versions, nil
}

// UseVersion sets the current Rust version
func (r *RustManager) UseVersion(version string) error {
	// Check if rustup is installed
	if !r.isRustupInstalled() {
		return fmt.Errorf("rustup is not installed")
	}

	// Set the specified version as the default
	cmd := exec.Command("rustup", "default", version)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set Rust %s as default: %w\nOutput: %s", version, err, string(output))
	}

	r.CurrentVersion = version
	return nil
}

// GetCurrentVersion returns the current Rust version
func (r *RustManager) GetCurrentVersion() (string, error) {
	// Check if rustup is installed
	if !r.isRustupInstalled() {
		return "", fmt.Errorf("rustup is not installed")
	}

	// Get the current Rust version
	cmd := exec.Command("rustc", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current Rust version: %w\nOutput: %s", err, string(output))
	}

	// Parse the output to extract the version
	// Example output: "rustc 1.71.0 (8cdc39929 2023-07-05)"
	parts := strings.Fields(string(output))
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected output format: %s", string(output))
	}

	version := parts[1]
	r.CurrentVersion = version
	return version, nil
}

// isRustupInstalled checks if rustup is installed
func (r *RustManager) isRustupInstalled() bool {
	cmd := exec.Command("rustup", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// setupEnvironment sets up the environment variables for Rust
func (r *RustManager) setupEnvironment() error {
	// Get the shell configuration file
	shellConfig := getShellConfig()

	// Check if the Rust environment variables are already set
	content, err := os.ReadFile(shellConfig)
	if err != nil {
		return fmt.Errorf("failed to read shell configuration file: %w", err)
	}

	// If the Rust environment variables are not set, add them
	if !strings.Contains(string(content), "RUSTUP_HOME") {
		rustEnv := `
# Rust environment variables
export RUSTUP_HOME="$HOME/.rustup"
export CARGO_HOME="$HOME/.cargo"
export PATH="$CARGO_HOME/bin:$PATH"
`

		f, err := os.OpenFile(shellConfig, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open shell configuration file: %w", err)
		}
		defer f.Close()

		if _, err := f.WriteString(rustEnv); err != nil {
			return fmt.Errorf("failed to write to shell configuration file: %w", err)
		}
	}

	return nil
} 