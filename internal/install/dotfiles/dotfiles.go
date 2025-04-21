package dotfiles

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DotfilesManager handles dotfiles operations
type DotfilesManager struct {
	RepoURL     string
	LocalPath   string
	BackupPath  string
	IsCloned    bool
}

// NewDotfilesManager creates a new DotfilesManager
func NewDotfilesManager(repoURL, localPath, backupPath string) *DotfilesManager {
	return &DotfilesManager{
		RepoURL:    repoURL,
		LocalPath:  localPath,
		BackupPath: backupPath,
		IsCloned:   false,
	}
}

// CloneRepository clones the dotfiles repository
func (d *DotfilesManager) CloneRepository() error {
	// Check if the repository is already cloned
	if d.IsCloned {
		return fmt.Errorf("repository already cloned at %s", d.LocalPath)
	}

	// Create the local directory if it doesn't exist
	if err := os.MkdirAll(d.LocalPath, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", d.RepoURL, d.LocalPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w\nOutput: %s", err, string(output))
	}

	d.IsCloned = true
	return nil
}

// BackupDotfiles creates a backup of the current dotfiles
func (d *DotfilesManager) BackupDotfiles(files []string) error {
	// Create the backup directory if it doesn't exist
	if err := os.MkdirAll(d.BackupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Backup each file
	for _, file := range files {
		sourcePath := filepath.Join(homeDir, file)
		destPath := filepath.Join(d.BackupPath, file)

		// Create the destination directory if it doesn't exist
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}

		// Check if the source file exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			// Skip if the file doesn't exist
			continue
		}

		// Copy the file
		cmd := exec.Command("cp", "-r", sourcePath, destPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to backup file %s: %w\nOutput: %s", file, err, string(output))
		}
	}

	return nil
}

// RestoreDotfiles restores dotfiles from the backup
func (d *DotfilesManager) RestoreDotfiles(files []string) error {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Restore each file
	for _, file := range files {
		sourcePath := filepath.Join(d.BackupPath, file)
		destPath := filepath.Join(homeDir, file)

		// Check if the source file exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			// Skip if the file doesn't exist
			continue
		}

		// Create the destination directory if it doesn't exist
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}

		// Copy the file
		cmd := exec.Command("cp", "-r", sourcePath, destPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to restore file %s: %w\nOutput: %s", file, err, string(output))
		}
	}

	return nil
}

// SyncDotfiles synchronizes dotfiles with the repository
func (d *DotfilesManager) SyncDotfiles(files []string) error {
	// Check if the repository is cloned
	if !d.IsCloned {
		return fmt.Errorf("repository not cloned, clone it first")
	}

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Sync each file
	for _, file := range files {
		sourcePath := filepath.Join(homeDir, file)
		destPath := filepath.Join(d.LocalPath, file)

		// Check if the source file exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			// Skip if the file doesn't exist
			continue
		}

		// Create the destination directory if it doesn't exist
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}

		// Copy the file
		cmd := exec.Command("cp", "-r", sourcePath, destPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to sync file %s: %w\nOutput: %s", file, err, string(output))
		}
	}

	// Commit and push the changes
	cmd := exec.Command("git", "-C", d.LocalPath, "add", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add files to git: %w\nOutput: %s", err, string(output))
	}

	cmd = exec.Command("git", "-C", d.LocalPath, "commit", "-m", "Update dotfiles")
	output, err = cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "nothing to commit") {
		return fmt.Errorf("failed to commit changes: %w\nOutput: %s", err, string(output))
	}

	cmd = exec.Command("git", "-C", d.LocalPath, "push")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push changes: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ListDotfiles lists the dotfiles in the repository
func (d *DotfilesManager) ListDotfiles() ([]string, error) {
	// Check if the repository is cloned
	if !d.IsCloned {
		return nil, fmt.Errorf("repository not cloned, clone it first")
	}

	// List the files in the repository
	cmd := exec.Command("find", d.LocalPath, "-type", "f", "-not", "-path", "*/\\.*")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w\nOutput: %s", err, string(output))
	}

	// Parse the output
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	// Convert absolute paths to relative paths
	relativeFiles := make([]string, 0, len(files))
	for _, file := range files {
		if file == "" {
			continue
		}
		
		relPath, err := filepath.Rel(d.LocalPath, file)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path: %w", err)
		}
		
		relativeFiles = append(relativeFiles, relPath)
	}

	return relativeFiles, nil
} 