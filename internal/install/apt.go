package install

import (
	"fmt"
	"os/exec"
)

type AptPackageManager struct {
	// Add any necessary fields here
}

func (a *AptPackageManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("dpkg", "-l", pkg)
	return cmd.Run() == nil
}

func (a *AptPackageManager) Update() error {
	cmd := exec.Command("sudo", "apt-get", "update")
	return cmd.Run()
}

func (a *AptPackageManager) SetupSpecialPackage(pkg string) error {
	switch pkg {
	case "lsd":
		// Add the PPA repository for lsd
		addRepoCmd := exec.Command("sudo", "add-apt-repository", "ppa:aslatter/ppa", "-y")
		if err := addRepoCmd.Run(); err != nil {
			return fmt.Errorf("failed to add lsd PPA: %w", err)
		}
		
		// Update package list
		if err := a.Update(); err != nil {
			return fmt.Errorf("failed to update after adding lsd PPA: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported special package: %s", pkg)
	}
} 