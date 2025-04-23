package dotfiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	assert.NotNil(t, manager)
	assert.NotEmpty(t, manager.baseDir)
}

func TestInitialize(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "dotfiles-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create manager with temporary directory
	manager := &Manager{
		baseDir: tmpDir,
	}

	// Test initialization
	err = manager.Initialize()
	require.NoError(t, err)

	// Verify directory structure
	categories := []string{"shell", "editor", "git", "terminal"}
	for _, category := range categories {
		path := filepath.Join(tmpDir, category)
		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	}
}

func TestApplyDotfile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "dotfiles-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create manager with temporary directory
	manager := &Manager{
		baseDir: tmpDir,
	}

	// Create test dotfile
	dotfile := &interfaces.Dotfile{
		Category: "shell",
		Files: []interfaces.FileConfig{
			{
				Source:      "test.sh",
				Destination: ".test.sh",
				Type:        "content",
				Backup:      true,
			},
		},
	}

	// Create source file
	sourcePath := filepath.Join(tmpDir, "shell", "test.sh")
	err = os.MkdirAll(filepath.Dir(sourcePath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(sourcePath, []byte("test content"), 0644)
	require.NoError(t, err)

	// Test applying dotfile
	err = manager.ApplyDotfile(dotfile)
	require.NoError(t, err)

	// Verify file was created
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	destPath := filepath.Join(homeDir, ".test.sh")
	content, err := os.ReadFile(destPath)
	require.NoError(t, err)
	assert.Equal(t, "test content", string(content))

	// Clean up
	os.Remove(destPath)
}

func TestProcessFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "dotfiles-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create manager with temporary directory
	manager := &Manager{
		baseDir: tmpDir,
	}

	// Create test dotfile
	dotfile := &interfaces.Dotfile{
		Category: "shell",
	}

	tests := []struct {
		name     string
		file     interfaces.FileConfig
		setup    func() error
		verify   func() error
		wantErr  bool
	}{
		{
			name: "content file",
			file: interfaces.FileConfig{
				Source:      "test.sh",
				Destination: ".test.sh",
				Type:        "content",
				Backup:      true,
			},
			setup: func() error {
				sourcePath := filepath.Join(tmpDir, "shell", "test.sh")
				return os.WriteFile(sourcePath, []byte("test content"), 0644)
			},
			verify: func() error {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				destPath := filepath.Join(homeDir, ".test.sh")
				content, err := os.ReadFile(destPath)
				if err != nil {
					return err
				}
				if string(content) != "test content" {
					return assert.AnError
				}
				return os.Remove(destPath)
			},
			wantErr: false,
		},
		{
			name: "symlink",
			file: interfaces.FileConfig{
				Source:      "test.sh",
				Destination: ".test.sh",
				Type:        "symlink",
				Backup:      true,
			},
			setup: func() error {
				sourcePath := filepath.Join(tmpDir, "shell", "test.sh")
				return os.WriteFile(sourcePath, []byte("test content"), 0644)
			},
			verify: func() error {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				destPath := filepath.Join(homeDir, ".test.sh")
				info, err := os.Lstat(destPath)
				if err != nil {
					return err
				}
				if info.Mode()&os.ModeSymlink == 0 {
					return assert.AnError
				}
				return os.Remove(destPath)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				err := tt.setup()
				require.NoError(t, err)
			}

			err := manager.processFile(dotfile, tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.verify != nil {
				err := tt.verify()
				require.NoError(t, err)
			}
		})
	}
}

func TestBackupFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "dotfiles-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create manager with temporary directory
	manager := &Manager{
		baseDir: tmpDir,
	}

	// Create a test file to backup
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("original content"), 0644)
	require.NoError(t, err)

	// Test backup with default suffix
	err = manager.backupFile(testFile, ".bak")
	require.NoError(t, err)

	// Verify backup was created
	backupFile := testFile + ".bak"
	_, err = os.Stat(backupFile)
	require.NoError(t, err)

	// Verify original file no longer exists
	_, err = os.Stat(testFile)
	require.True(t, os.IsNotExist(err))

	// Test backup with custom suffix
	testFile2 := filepath.Join(tmpDir, "test2.txt")
	err = os.WriteFile(testFile2, []byte("original content"), 0644)
	require.NoError(t, err)

	err = manager.backupFile(testFile2, ".backup")
	require.NoError(t, err)

	// Verify backup was created with custom suffix
	backupFile2 := testFile2 + ".backup"
	_, err = os.Stat(backupFile2)
	require.NoError(t, err)

	// Test backup of non-existent file (should not error)
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
	err = manager.backupFile(nonExistentFile, ".bak")
	require.NoError(t, err)
}

func TestProcessNonExistentFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "dotfiles-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create manager with temporary directory
	manager := &Manager{
		baseDir: tmpDir,
	}

	// Create test dotfile with non-existent source
	dotfile := &interfaces.Dotfile{
		Category: "shell",
	}

	file := interfaces.FileConfig{
		Source:      "nonexistent.sh",
		Destination: ".test.sh",
		Type:        "content",
		Backup:      true,
	}

	// Test processing non-existent file
	err = manager.processFile(dotfile, file)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read source file")
} 