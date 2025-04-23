package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// Manager handles dotfiles operations
type Manager struct {
	configLoader *config.ConfigLoader
	baseDir     string
}

// NewManager creates a new dotfiles manager
func NewManager() *Manager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}
	
	return &Manager{
		configLoader: config.NewConfigLoader("config"),
		baseDir:     filepath.Join(homeDir, ".dotfiles"),
	}
}

// Initialize sets up the dotfiles directory structure
func (m *Manager) Initialize() error {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(m.baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create dotfiles directory: %w", err)
	}

	// Create category subdirectories
	categories := []string{"shell", "editor", "git", "terminal"}
	for _, category := range categories {
		if err := os.MkdirAll(filepath.Join(m.baseDir, category), 0755); err != nil {
			return fmt.Errorf("failed to create category directory %s: %w", category, err)
		}
	}

	return nil
}

// CloneUserRepo clones a user's dotfiles repository
func (m *Manager) CloneUserRepo(repoURL string) error {
	// TODO: Implement git clone logic
	return nil
}

// ApplyDotfile applies a dotfile configuration
func (m *Manager) ApplyDotfile(dotfile *interfaces.Dotfile) error {
	// Create category directory
	categoryDir := filepath.Join(m.baseDir, dotfile.Category)
	if err := os.MkdirAll(categoryDir, 0755); err != nil {
		return fmt.Errorf("failed to create category directory: %w", err)
	}

	// Process each file in the configuration
	for _, file := range dotfile.Files {
		if err := m.processFile(dotfile, file); err != nil {
			return fmt.Errorf("failed to process file %s: %w", file.Source, err)
		}
	}

	// Apply shell configuration
	if err := m.applyShellConfig(dotfile); err != nil {
		return fmt.Errorf("failed to apply shell config: %w", err)
	}

	return nil
}

// processFile handles a single file configuration
func (m *Manager) processFile(dotfile *interfaces.Dotfile, file interfaces.FileConfig) error {
	// Determine source and destination paths
	sourcePath := file.Source
	if !strings.HasPrefix(sourcePath, "http") {
		sourcePath = filepath.Join(m.baseDir, dotfile.Category, file.Source)
	}

	destPath := file.Destination
	if !filepath.IsAbs(destPath) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		destPath = filepath.Join(homeDir, destPath)
	}

	// Create parent directories if needed
	if file.CreateParents {
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories: %w", err)
		}
	}

	// Handle different file types
	switch file.Type {
	case "content":
		return m.writeContentFile(sourcePath, destPath, file)
	case "symlink":
		return m.createSymlink(sourcePath, destPath, file)
	default:
		return fmt.Errorf("unsupported file type: %s", file.Type)
	}
}

// writeContentFile writes content to a file
func (m *Manager) writeContentFile(source, dest string, file interfaces.FileConfig) error {
	// Read source content
	content, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Backup existing file if needed
	if file.Backup {
		if err := m.backupFile(dest, file.BackupSuffix); err != nil {
			return fmt.Errorf("failed to backup file: %w", err)
		}
	}

	// Write content to destination
	if err := os.WriteFile(dest, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// createSymlink creates a symlink
func (m *Manager) createSymlink(source, dest string, file interfaces.FileConfig) error {
	// Backup existing file if needed
	if file.Backup {
		if err := m.backupFile(dest, file.BackupSuffix); err != nil {
			return fmt.Errorf("failed to backup file: %w", err)
		}
	}

	// Remove existing file/symlink
	if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing file: %w", err)
	}

	// Create symlink
	if err := os.Symlink(source, dest); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// backupFile creates a backup of an existing file
func (m *Manager) backupFile(path, suffix string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // No file to backup
	}

	backupPath := path + suffix
	return os.Rename(path, backupPath)
}

// applyShellConfig applies shell-specific configuration
func (m *Manager) applyShellConfig(dotfile *interfaces.Dotfile) error {
	// TODO: Implement shell configuration application
	return nil
} 