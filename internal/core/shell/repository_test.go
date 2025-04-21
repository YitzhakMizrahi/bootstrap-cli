package shell

import (
	"os"
	"testing"
)

func TestRepositoryManagement(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "repository-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a repository manager
	manager, err := NewRepositoryManager()
	if err != nil {
		t.Fatalf("Failed to create repository manager: %v", err)
	}

	// Test adding a repository
	err = manager.AddRepository("test-repo", "https://example.com/repo", "Test repository")
	if err != nil {
		t.Errorf("Failed to add repository: %v", err)
	}

	// Test listing repositories
	repos := manager.ListRepositories()
	if len(repos) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(repos))
	}

	// Test removing a repository
	err = manager.RemoveRepository("test-repo")
	if err != nil {
		t.Errorf("Failed to remove repository: %v", err)
	}

	repos = manager.ListRepositories()
	if len(repos) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(repos))
	}
}

func TestPluginSearch(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "repository-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a repository manager
	manager, err := NewRepositoryManager()
	if err != nil {
		t.Fatalf("Failed to create repository manager: %v", err)
	}

	// Add a repository with a plugin
	err = manager.AddRepository("test-repo", "https://example.com/repo", "Test repository")
	if err != nil {
		t.Fatalf("Failed to add repository: %v", err)
	}

	// Test searching for plugins
	results, err := manager.SearchPlugins("example")
	if err != nil {
		t.Errorf("Failed to search plugins: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(results))
	}

	// Test searching for non-existent plugins
	results, err = manager.SearchPlugins("nonexistent")
	if err != nil {
		t.Errorf("Failed to search plugins: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 search results, got %d", len(results))
	}
}

func TestPluginInstallation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "repository-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a repository manager
	manager, err := NewRepositoryManager()
	if err != nil {
		t.Fatalf("Failed to create repository manager: %v", err)
	}

	// Add a repository with a plugin
	err = manager.AddRepository("test-repo", "https://example.com/repo", "Test repository")
	if err != nil {
		t.Fatalf("Failed to add repository: %v", err)
	}

	// Create a shell for testing
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Test installing a plugin
	err = manager.InstallPlugin("test-repo", "example-plugin", shell)
	if err != nil {
		t.Errorf("Failed to install plugin: %v", err)
	}

	// Verify the plugin was installed
	if !shell.HasPlugin("example-plugin") {
		t.Error("Plugin was not installed correctly")
	}
} 