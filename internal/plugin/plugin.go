package plugin

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"sync"
)

// Plugin represents a loaded plugin
type Plugin struct {
	Name        string
	Path        string
	Description string
	Version     string
	Handle      *plugin.Plugin
}

// Manager handles plugin loading and management
type Manager struct {
	plugins     map[string]*Plugin
	pluginDir   string
	maxPlugins  int
	autoReload  bool
	mu          sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager(pluginDir string, maxPlugins int, autoReload bool) *Manager {
	return &Manager{
		plugins:    make(map[string]*Plugin),
		pluginDir:  pluginDir,
		maxPlugins: maxPlugins,
		autoReload: autoReload,
	}
}

// LoadPlugin loads a plugin from the given path
func (m *Manager) LoadPlugin(path string) (*Plugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if we've reached the plugin limit
	if len(m.plugins) >= m.maxPlugins {
		return nil, fmt.Errorf("maximum number of plugins (%d) reached", m.maxPlugins)
	}

	// Load the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %v", err)
	}

	// Get plugin metadata
	nameSym, err := p.Lookup("Name")
	if err != nil {
		return nil, fmt.Errorf("plugin missing Name symbol: %v", err)
	}
	name, ok := nameSym.(*string)
	if !ok {
		return nil, fmt.Errorf("plugin Name symbol has wrong type")
	}

	descSym, err := p.Lookup("Description")
	if err != nil {
		return nil, fmt.Errorf("plugin missing Description symbol: %v", err)
	}
	desc, ok := descSym.(*string)
	if !ok {
		return nil, fmt.Errorf("plugin Description symbol has wrong type")
	}

	versionSym, err := p.Lookup("Version")
	if err != nil {
		return nil, fmt.Errorf("plugin missing Version symbol: %v", err)
	}
	version, ok := versionSym.(*string)
	if !ok {
		return nil, fmt.Errorf("plugin Version symbol has wrong type")
	}

	// Create plugin instance
	plugin := &Plugin{
		Name:        *name,
		Path:        path,
		Description: *desc,
		Version:     *version,
		Handle:      p,
	}

	// Store plugin
	m.plugins[plugin.Name] = plugin

	return plugin, nil
}

// UnloadPlugin unloads a plugin by name
func (m *Manager) UnloadPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	delete(m.plugins, name)
	return nil
}

// GetPlugin returns a loaded plugin by name
func (m *Manager) GetPlugin(name string) (*Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// ListPlugins returns a list of all loaded plugins
func (m *Manager) ListPlugins() []*Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// LoadPlugins loads all plugins from the plugin directory
func (m *Manager) LoadPlugins() error {
	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// Find all .so files in the plugin directory
	files, err := filepath.Glob(filepath.Join(m.pluginDir, "*.so"))
	if err != nil {
		return fmt.Errorf("failed to scan plugin directory: %v", err)
	}

	// Load each plugin
	for _, file := range files {
		plugin, err := m.LoadPlugin(file)
		if err != nil {
			return fmt.Errorf("failed to load plugin %s: %v", file, err)
		}
		log.Printf("Loaded plugin: %s (version %s)", plugin.Name, plugin.Version)
	}

	return nil
}

// ReloadPlugins reloads all plugins
func (m *Manager) ReloadPlugins() error {
	if !m.autoReload {
		return fmt.Errorf("auto-reload is disabled")
	}

	// Unload all plugins
	m.mu.Lock()
	for name := range m.plugins {
		delete(m.plugins, name)
	}
	m.mu.Unlock()

	// Reload all plugins
	return m.LoadPlugins()
} 