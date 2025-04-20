package installer

import (
	"fmt"
	"os"
	"os/exec"
)

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// appendToFile adds content to the end of a file
func appendToFile(filepath, content string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	
	return nil
}

// Package manager installation helpers
func brewInstall(pkg string) error {
	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// aptInstall installs a package using apt
func aptInstall(pkg string) error {
    // Use environment variables to ensure non-interactive mode
    env := os.Environ()
    env = append(env, "DEBIAN_FRONTEND=noninteractive")
    
    // First update
    cmd := exec.Command("sudo", "apt-get", "update", "-q")
    cmd.Env = env
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run() // Ignore error and continue
    
    // Then install with -y and non-interactive options
    cmd = exec.Command("sudo", "apt-get", "install", "-y", "-q", pkg)
    cmd.Env = env
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func dnfInstall(pkg string) error {
	cmd := exec.Command("sudo", "dnf", "install", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pacmanInstall(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func zypperInstall(pkg string) error {
	cmd := exec.Command("sudo", "zypper", "install", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func chocoInstall(pkg string) error {
	cmd := exec.Command("choco", "install", pkg, "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetSudoSession obtains a sudo session up front to avoid prompting during installation
func GetSudoSession() error {
    // Check if we already have sudo privileges
    testCmd := exec.Command("sudo", "-n", "true")
    if err := testCmd.Run(); err == nil {
        // We already have sudo privileges
        return nil
    }

    fmt.Println("🔐 Requesting sudo privileges for installation...")
    
    // Request sudo password up front
    cmd := exec.Command("sudo", "true")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    return cmd.Run()
}

// PreventServicePrompts configures apt to prevent service restart prompts
func PreventServicePrompts() error {
    // Only applicable on Debian/Ubuntu systems
    if _, err := os.Stat("/etc/apt/apt.conf.d"); os.IsNotExist(err) {
        return nil // Not a Debian/Ubuntu system
    }
    
    // Create the config if it doesn't exist
    configFile := "/etc/apt/apt.conf.d/90bootstrap-no-restart"
    configContent := `// Configured by bootstrap-cli to prevent service restart prompts
DPkg::Options {
    "--force-confdef";
    "--force-confold";
}
Dpkg::NoTriggers "true";
Dpkg::Options::="--force-confdef";
Dpkg::Options::="--force-confold";
`
    
    // Create a temporary file with the content
    tempFile, err := os.CreateTemp("", "bootstrap-apt-conf")
    if err != nil {
        return fmt.Errorf("failed to create temp file: %w", err)
    }
    defer os.Remove(tempFile.Name())
    
    if _, err := tempFile.WriteString(configContent); err != nil {
        return fmt.Errorf("failed to write to temp file: %w", err)
    }
    
    if err := tempFile.Close(); err != nil {
        return fmt.Errorf("failed to close temp file: %w", err)
    }
    
    // Use sudo to move the file to the correct location
    cmd := exec.Command("sudo", "cp", tempFile.Name(), configFile)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to copy config file: %w", err)
    }
    
    return nil
}

// CleanupServicePromptConfig removes the temporary apt configuration
func CleanupServicePromptConfig() error {
    configFile := "/etc/apt/apt.conf.d/90bootstrap-no-restart"
    
    // Check if file exists
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        return nil // File doesn't exist, nothing to do
    }
    
    // Remove the file
    cmd := exec.Command("sudo", "rm", configFile)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to remove config file: %w", err)
    }
    
    return nil
}