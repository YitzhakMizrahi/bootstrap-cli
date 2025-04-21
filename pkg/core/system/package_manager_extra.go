package system

import (
	"os/exec"
	"strings"
)

// DNF package manager implementation
type dnfManager struct {
	baseManager
}

func (m *dnfManager) Update() error {
	return m.runCmd("dnf", "check-update")
}

func (m *dnfManager) Upgrade() error {
	return m.runCmd("dnf", "upgrade", "-y")
}

func (m *dnfManager) Install(pkg string) error {
	return m.runCmd("dnf", "install", "-y", pkg)
}

func (m *dnfManager) InstallMany(pkgs []string) error {
	args := append([]string{"install", "-y"}, pkgs...)
	return m.runCmd("dnf", args...)
}

func (m *dnfManager) Remove(pkg string) error {
	return m.runCmd("dnf", "remove", "-y", pkg)
}

func (m *dnfManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("rpm", "-q", pkg)
	return cmd.Run() == nil
}

func (m *dnfManager) Search(query string) ([]Package, error) {
	cmd := exec.Command("dnf", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []Package
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			packages = append(packages, Package{Name: parts[0]})
		}
	}
	return packages, nil
}

// Pacman package manager implementation
type pacmanManager struct {
	baseManager
}

func (m *pacmanManager) Update() error {
	return m.runCmd("pacman", "-Sy")
}

func (m *pacmanManager) Upgrade() error {
	return m.runCmd("pacman", "-Syu", "--noconfirm")
}

func (m *pacmanManager) Install(pkg string) error {
	return m.runCmd("pacman", "-S", "--noconfirm", pkg)
}

func (m *pacmanManager) InstallMany(pkgs []string) error {
	args := append([]string{"-S", "--noconfirm"}, pkgs...)
	return m.runCmd("pacman", args...)
}

func (m *pacmanManager) Remove(pkg string) error {
	return m.runCmd("pacman", "-R", "--noconfirm", pkg)
}

func (m *pacmanManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("pacman", "-Q", pkg)
	return cmd.Run() == nil
}

func (m *pacmanManager) Search(query string) ([]Package, error) {
	cmd := exec.Command("pacman", "-Ss", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []Package
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, " ") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			name := strings.Split(parts[0], "/")[1] // Remove repo prefix
			packages = append(packages, Package{Name: name})
		}
	}
	return packages, nil
} 