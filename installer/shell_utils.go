package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// GetCurrentShell returns the current shell
func GetCurrentShell() (string, error) {
	// Get the current user's shell from SHELL env var
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return "", fmt.Errorf("could not determine current shell")
	}
	
	// Extract the shell name from the path
	parts := strings.Split(shellPath, "/")
	return parts[len(parts)-1], nil
}

// GetAvailableShells returns a list of available shells
func GetAvailableShells() ([]string, error) {
	// Try to read from /etc/shells
	data, err := os.ReadFile("/etc/shells")
	if err != nil {
		return nil, fmt.Errorf("could not read available shells: %w", err)
	}
	
	lines := strings.Split(string(data), "\n")
	var shells []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			// Extract the shell name from the path
			parts := strings.Split(line, "/")
			shells = append(shells, parts[len(parts)-1])
		}
	}
	
	return shells, nil
}

// SetDefaultShell sets the default shell for the current user
func SetDefaultShell(shell string) error {
	shellPath, err := exec.LookPath(shell)
	if err != nil {
		return fmt.Errorf("shell '%s' not found in PATH: %w", shell, err)
	}
	
	// Check if the shell is in /etc/shells
	shellsData, err := os.ReadFile("/etc/shells")
	if err != nil {
		return fmt.Errorf("could not read available shells: %w", err)
	}
	
	shellsContent := string(shellsData)
	if !strings.Contains(shellsContent, shellPath) {
		return fmt.Errorf("shell '%s' is not in /etc/shells", shellPath)
	}
	
	// Use chsh to change the shell
	cmd := exec.Command("chsh", "-s", shellPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// IsShellSupported checks if a shell is supported
func IsShellSupported(shell string) bool {
	_, err := exec.LookPath(shell)
	return err == nil
}

// RestartShell executes the new shell
func RestartShell(shell string) error {
	shellPath, err := exec.LookPath(shell)
	if err != nil {
		return fmt.Errorf("shell '%s' not found in PATH: %w", shell, err)
	}
	
	fmt.Printf("🔄 Restarting with shell: %s\n", shell)
	
	// Execute the new shell
	return syscall.Exec(shellPath, []string{shellPath}, os.Environ())
}

// AddToRcFile adds content to the shell's rc file
func AddToRcFile(shell, content string) error {
	var rcFile string
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}
	
	switch shell {
	case "bash":
		rcFile = fmt.Sprintf("%s/.bashrc", home)
	case "zsh":
		rcFile = fmt.Sprintf("%s/.zshrc", home)
	case "fish":
		// Create the directory if it doesn't exist
		configDir := fmt.Sprintf("%s/.config/fish", home)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("could not create fish config directory: %w", err)
		}
		rcFile = fmt.Sprintf("%s/config.fish", configDir)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	
	// Check if the file exists and create it if not
	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		if err := os.WriteFile(rcFile, []byte("# Shell configuration\n\n"), 0644); err != nil {
			return fmt.Errorf("could not create rc file: %w", err)
		}
	}
	
	// Read the current content
	data, err := os.ReadFile(rcFile)
	if err != nil {
		return fmt.Errorf("could not read rc file: %w", err)
	}
	
	// Check if the content already exists
	if strings.Contains(string(data), content) {
		// Content already exists, no need to add it again
		return nil
	}
	
	// Add the content to the file
	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open rc file: %w", err)
	}
	defer f.Close()
	
	// Add a newline before the content if the file doesn't end with one
	if len(data) > 0 && data[len(data)-1] != '\n' {
		content = "\n" + content
	}
	
	if _, err := f.WriteString(content + "\n"); err != nil {
		return fmt.Errorf("could not write to rc file: %w", err)
	}
	
	return nil
}

// SourceRcFile sources the shell's rc file
func SourceRcFile(shell string) error {
	var rcFile string
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}
	
	switch shell {
	case "bash":
		rcFile = fmt.Sprintf("%s/.bashrc", home)
	case "zsh":
		rcFile = fmt.Sprintf("%s/.zshrc", home)
	case "fish":
		rcFile = fmt.Sprintf("%s/.config/fish/config.fish", home)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	
	// Check if the file exists
	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		return fmt.Errorf("rc file does not exist: %s", rcFile)
	}
	
	// Source the file
	var cmd *exec.Cmd
	switch shell {
	case "bash", "zsh":
		cmd = exec.Command(shell, "-c", fmt.Sprintf("source %s", rcFile))
	case "fish":
		cmd = exec.Command(shell, "-c", fmt.Sprintf("source %s", rcFile))
	}
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
} 