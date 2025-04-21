package dotfiles

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNewDotfilesManager(t *testing.T) {
	manager := NewDotfilesManager("https://github.com/user/dotfiles", "/tmp/dotfiles", "/tmp/dotfiles_backup")
	
	if manager.RepoURL != "https://github.com/user/dotfiles" {
		t.Errorf("Expected RepoURL to be 'https://github.com/user/dotfiles', got '%s'", manager.RepoURL)
	}
	
	if manager.LocalPath != "/tmp/dotfiles" {
		t.Errorf("Expected LocalPath to be '/tmp/dotfiles', got '%s'", manager.LocalPath)
	}
	
	if manager.BackupPath != "/tmp/dotfiles_backup" {
		t.Errorf("Expected BackupPath to be '/tmp/dotfiles_backup', got '%s'", manager.BackupPath)
	}
	
	if manager.IsCloned != false {
		t.Errorf("Expected IsCloned to be false, got %v", manager.IsCloned)
	}
}

func TestBackupAndRestoreDotfiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dotfiles_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test file
	testFile := filepath.Join(tempDir, ".testrc")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Create a backup directory
	backupDir := filepath.Join(tempDir, "backup")
	err = os.MkdirAll(backupDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backup directory: %v", err)
	}
	
	// Create a manager
	manager := NewDotfilesManager("", "", backupDir)
	
	// Backup the file
	err = manager.BackupDotfiles([]string{".testrc"})
	if err != nil {
		t.Fatalf("Failed to backup dotfiles: %v", err)
	}
	
	// Check if the file was backed up
	backupFile := filepath.Join(backupDir, ".testrc")
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Fatalf("Backup file does not exist: %v", err)
	}
	
	// Remove the original file
	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("Failed to remove original file: %v", err)
	}
	
	// Restore the file
	err = manager.RestoreDotfiles([]string{".testrc"})
	if err != nil {
		t.Fatalf("Failed to restore dotfiles: %v", err)
	}
	
	// Check if the file was restored
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Restored file does not exist: %v", err)
	}
	
	// Check the content of the restored file
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}
	
	if string(content) != "test content" {
		t.Errorf("Expected content to be 'test content', got '%s'", string(content))
	}
}

func TestListDotfiles(t *testing.T) {
	// Skip this test if git is not available
	if _, err := os.Stat("/usr/bin/git"); os.IsNotExist(err) {
		t.Skip("Git is not available, skipping test")
	}
	
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dotfiles_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Initialize a git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}
	
	// Create some test files
	files := []string{
		".zshrc",
		".bashrc",
		".gitconfig",
		"config/starship.toml",
	}
	
	for _, file := range files {
		filePath := filepath.Join(tempDir, file)
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		
		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}
	
	// Add files to git
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tempDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add files to git: %v", err)
	}
	
	// Commit files
	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to commit files: %v", err)
	}
	
	// Create a manager
	manager := NewDotfilesManager("", tempDir, "")
	manager.IsCloned = true
	
	// List dotfiles
	dotfiles, err := manager.ListDotfiles()
	if err != nil {
		t.Fatalf("Failed to list dotfiles: %v", err)
	}
	
	// Check if all files are listed
	if len(dotfiles) != len(files) {
		t.Errorf("Expected %d files, got %d", len(files), len(dotfiles))
	}
	
	// Check if each file is in the list
	for _, file := range files {
		found := false
		for _, dotfile := range dotfiles {
			if dotfile == file {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("File '%s' not found in dotfiles list", file)
		}
	}
} 