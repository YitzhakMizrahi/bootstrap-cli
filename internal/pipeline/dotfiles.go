package pipeline

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	// TODO: Add import for os user home dir if needed for target path
	// TODO: Add import for logger if needed
)

// GenerateDotfileCloneSteps creates pipeline steps for cloning a dotfiles repository.
func GenerateDotfileCloneSteps(repoURL, targetDir string) []InstallationStep {
	steps := []InstallationStep{}

	// Step 1: Check for git dependency (optional, could be handled by main dependency graph)
	// steps = append(steps, InstallationStep{
	// 	Name:        "check-git-for-dotfiles",
	// 	Description: "Checking for git command",
	// 	Action: func() error {
	// 		if _, err := exec.LookPath("git"); err != nil {
	// 			return fmt.Errorf("git command not found, required for cloning dotfiles: %w", err)
	// 		}
	// 		return nil
	// 	},
	// })

	// Step 2: Clone the repository
	// Construct the full GitHub URL if only user/repo is provided
	fullRepoURL := repoURL
	if !strings.Contains(repoURL, "://") { // Basic check if it's not a full URL
		fullRepoURL = fmt.Sprintf("https://github.com/%s.git", repoURL)
	}

	cloneStep := InstallationStep{
		Name:        fmt.Sprintf("clone-dotfiles-%s", filepath.Base(repoURL)),
		Description: fmt.Sprintf("Cloning dotfiles from %s", fullRepoURL),
		Action: func() error {
			// TODO: Add logging via context logger
			fmt.Printf("Attempting to clone %s into %s\n", fullRepoURL, targetDir)
			cmd := exec.Command("git", "clone", "--depth=1", fullRepoURL, targetDir) // Shallow clone
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to clone dotfiles repo '%s': %w\nOutput:\n%s", fullRepoURL, err, string(output))
			}
			fmt.Printf("Successfully cloned dotfiles.\n")
			return nil
		},
		// TODO: Implement Rollback? (e.g., remove targetDir)
		Rollback: func() error {
			fmt.Printf("Attempting to roll back dotfiles clone by removing %s\n", targetDir)
			cmd := exec.Command("rm", "-rf", targetDir)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to remove dotfiles directory during rollback '%s': %w\nOutput:\n%s", targetDir, err, string(output))
			}
			return nil
		},
		Timeout:    5 * time.Minute, // Cloning can take time
		RetryCount: 1, // Usually no point retrying a failed clone immediately
	}
	steps = append(steps, cloneStep)

	// TODO: Add steps for symlinking configurations from the cloned repo
	// This would involve: 
	// 1. Determining the source files/dirs within targetDir.
	// 2. Determining the destination paths in $HOME.
	// 3. Using os.Symlink or similar.

	return steps
} 