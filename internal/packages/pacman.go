package packages

import (
	"os"
	"os/exec"
	"strings"
)

// PacmanManager implements PackageManager for Pacman-based systems
type PacmanManager struct{}

// Name returns the name of the package manager
func (p *PacmanManager) Name() string {
	return string(Pacman)
}

// IsAvailable checks if pacman is available on the system
func (p *PacmanManager) IsAvailable() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

// Install installs the given packages using pacman
func (p *PacmanManager) Install(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}

	args := append([]string{
		"-S",
		"--noconfirm",
	}, packages...)

	cmd := exec.Command("pacman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Update updates the package list
func (p *PacmanManager) Update() error {
	cmd := exec.Command("pacman", "-Sy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// IsInstalled checks if a package is installed
func (p *PacmanManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("pacman", "-Q", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), pkg)
} 