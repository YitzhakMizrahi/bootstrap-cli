package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PluginRepository represents a repository of shell plugins
type PluginRepository struct {
	Name        string
	URL         string
	Description string
	Plugins     []*Plugin
}

// PluginSearchResult represents the result of a plugin search
type PluginSearchResult struct {
	Name        string
	Path        string
	Version     string
	Description string
	Repository  string
}

// RepositoryManager manages plugin repositories
type RepositoryManager struct {
	Repositories []*PluginRepository
	CacheDir     string
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager() (*RepositoryManager, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".bootstrap-cli", "plugin-cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	return &RepositoryManager{
		Repositories: make([]*PluginRepository, 0),
		CacheDir:     cacheDir,
	}, nil
}

// AddRepository adds a new plugin repository
func (rm *RepositoryManager) AddRepository(name, url, description string) error {
	// Check if repository already exists
	for _, repo := range rm.Repositories {
		if repo.Name == name || repo.URL == url {
			return fmt.Errorf("repository %s already exists", name)
		}
	}

	repo := &PluginRepository{
		Name:        name,
		URL:         url,
		Description: description,
		Plugins:     make([]*Plugin, 0),
	}

	// Fetch repository metadata
	if err := rm.updateRepository(repo); err != nil {
		return err
	}

	rm.Repositories = append(rm.Repositories, repo)
	return nil
}

// RemoveRepository removes a plugin repository
func (rm *RepositoryManager) RemoveRepository(name string) error {
	for i, repo := range rm.Repositories {
		if repo.Name == name {
			rm.Repositories = append(rm.Repositories[:i], rm.Repositories[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("repository %s not found", name)
}

// ListRepositories returns a list of all repositories
func (rm *RepositoryManager) ListRepositories() []*PluginRepository {
	return rm.Repositories
}

// UpdateRepositories updates all repositories
func (rm *RepositoryManager) UpdateRepositories() error {
	for _, repo := range rm.Repositories {
		if err := rm.updateRepository(repo); err != nil {
			return fmt.Errorf("failed to update repository %s: %v", repo.Name, err)
		}
	}
	return nil
}

// SearchPlugins searches for plugins across all repositories
func (rm *RepositoryManager) SearchPlugins(query string) ([]*PluginSearchResult, error) {
	results := make([]*PluginSearchResult, 0)

	for _, repo := range rm.Repositories {
		for _, plugin := range repo.Plugins {
			if matchesSearch(plugin, query) {
				results = append(results, &PluginSearchResult{
					Name:        plugin.Name,
					Path:        plugin.Path,
					Version:     plugin.Version,
					Description: plugin.Description,
					Repository:  repo.Name,
				})
			}
		}
	}

	return results, nil
}

// InstallPlugin installs a plugin from a repository
func (rm *RepositoryManager) InstallPlugin(repoName, pluginName string, shell *Shell) error {
	// Find the repository
	var repo *PluginRepository
	for _, r := range rm.Repositories {
		if r.Name == repoName {
			repo = r
			break
		}
	}
	if repo == nil {
		return fmt.Errorf("repository %s not found", repoName)
	}

	// Find the plugin
	var plugin *Plugin
	for _, p := range repo.Plugins {
		if p.Name == pluginName {
			plugin = p
			break
		}
	}
	if plugin == nil {
		return fmt.Errorf("plugin %s not found in repository %s", pluginName, repoName)
	}

	// Download and install the plugin
	return rm.downloadAndInstallPlugin(plugin, shell)
}

// updateRepository updates a repository's plugin list
func (rm *RepositoryManager) updateRepository(repo *PluginRepository) error {
	// In a real implementation, this would fetch the repository metadata
	// from the repository URL and update the plugin list
	// For now, we'll simulate it with more realistic data
	
	// Clear existing plugins
	repo.Plugins = []*Plugin{}
	
	// Add some example plugins
	plugins := []*Plugin{
		{
			Name:        "zsh-autosuggestions",
			Path:        filepath.Join(rm.CacheDir, "zsh-autosuggestions", "zsh-autosuggestions.zsh"),
			Version:     "0.7.0",
			Description: "Fish-like autosuggestions for zsh",
			Enabled:     false,
			Config:      make(map[string]string),
			Dependencies: []string{},
		},
		{
			Name:        "zsh-syntax-highlighting",
			Path:        filepath.Join(rm.CacheDir, "zsh-syntax-highlighting", "zsh-syntax-highlighting.zsh"),
			Version:     "0.8.0",
			Description: "Fish shell like syntax highlighting for zsh",
			Enabled:     false,
			Config:      make(map[string]string),
			Dependencies: []string{},
		},
		{
			Name:        "zsh-completions",
			Path:        filepath.Join(rm.CacheDir, "zsh-completions", "zsh-completions.plugin.zsh"),
			Version:     "0.34.0",
			Description: "Additional completion definitions for zsh",
			Enabled:     false,
			Config:      make(map[string]string),
			Dependencies: []string{},
		},
		{
			Name:        "git",
			Path:        filepath.Join(rm.CacheDir, "git", "git.plugin.zsh"),
			Version:     "2.3.0",
			Description: "Git aliases and functions for zsh",
			Enabled:     false,
			Config:      make(map[string]string),
			Dependencies: []string{},
		},
		{
			Name:        "docker",
			Path:        filepath.Join(rm.CacheDir, "docker", "docker.plugin.zsh"),
			Version:     "1.0.0",
			Description: "Docker aliases and functions for zsh",
			Enabled:     false,
			Config:      make(map[string]string),
			Dependencies: []string{},
		},
	}
	
	// Add plugins to the repository
	repo.Plugins = plugins
	
	return nil
}

// downloadAndInstallPlugin downloads and installs a plugin
func (rm *RepositoryManager) downloadAndInstallPlugin(plugin *Plugin, shell *Shell) error {
	// In a real implementation, this would:
	// 1. Download the plugin from the repository
	// 2. Verify its integrity
	// 3. Install it to the appropriate location
	// 4. Add it to the shell's plugin list
	
	// For now, we'll simulate it with a more realistic implementation
	
	// Create the plugin directory if it doesn't exist
	pluginDir := filepath.Dir(plugin.Path)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}
	
	// Create a dummy plugin file
	pluginContent := fmt.Sprintf(`# %s plugin v%s
# %s

# This is a dummy plugin file for testing purposes.
# In a real implementation, this would be downloaded from the repository.

echo "Plugin %s loaded"
`, plugin.Name, plugin.Version, plugin.Description, plugin.Name)
	
	if err := os.WriteFile(plugin.Path, []byte(pluginContent), 0644); err != nil {
		return fmt.Errorf("failed to create plugin file: %w", err)
	}
	
	// Add the plugin to the shell
	return shell.AddPluginWithMetadata(
		plugin.Name,
		plugin.Path,
		plugin.Version,
		plugin.Description,
		plugin.Dependencies,
	)
}

// matchesSearch checks if a plugin matches a search query
func matchesSearch(plugin *Plugin, query string) bool {
	// Simple case-insensitive substring matching
	// In a real implementation, this could be more sophisticated
	// (e.g., fuzzy matching, tag-based search, etc.)
	return containsIgnoreCase(plugin.Name, query) ||
		containsIgnoreCase(plugin.Description, query)
}

// containsIgnoreCase checks if a string contains another string, ignoring case
func containsIgnoreCase(s, substr string) bool {
	// Simple implementation - in a real implementation,
	// this could use proper Unicode case folding
	return len(substr) == 0 || len(s) >= len(substr) &&
		strings.ToLower(s[0:len(substr)]) == strings.ToLower(substr)
} 