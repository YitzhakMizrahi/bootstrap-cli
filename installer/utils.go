package installer

import (
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