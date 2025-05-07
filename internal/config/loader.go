package config

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

//go:embed defaults/**
var defaultConfigs embed.FS

// Loader handles loading and parsing configuration files
type Loader struct {
	baseDir     string // User config directory
	defaultsDir string // Embedded defaults directory
	configFS    embed.FS
}

// NewLoader creates a new configuration loader
func NewLoader(baseDir string) *Loader {
	loader := &Loader{
		baseDir:     baseDir,
		defaultsDir: "defaults",
		configFS:    defaultConfigs,
	}
	
	return loader
}

// mergeConfigs merges user config into default config
func mergeConfigs[T any](defaultConfig, userConfig *T) *T {
	if userConfig == nil {
		return defaultConfig
	}
	if defaultConfig == nil {
		return userConfig
	}
	
	// Use reflection or yaml.Marshal/Unmarshal to merge structs
	defaultYAML, _ := yaml.Marshal(defaultConfig)
	userYAML, _ := yaml.Marshal(userConfig)
	
	merged := *defaultConfig // Create a copy of default
	_ = yaml.Unmarshal(defaultYAML, &merged)
	_ = yaml.Unmarshal(userYAML, &merged) // User config overrides defaults
	
	return &merged
}

// LoadTools loads all tool configurations
func (l *Loader) LoadTools() ([]*interfaces.Tool, error) {
	configs, err := l.loadConfigsFromDir("tools")
	if err != nil {
		return nil, err
	}
	tools, ok := configs.([]*interfaces.Tool)
	if !ok {
		return nil, fmt.Errorf("failed to convert configs to tools")
	}
	return tools, nil
}

// LoadFonts loads all font configurations
func (l *Loader) LoadFonts() ([]*interfaces.Font, error) {
	configs, err := l.loadConfigsFromDir("fonts")
	if err != nil {
		return nil, err
	}
	fonts, ok := configs.([]*interfaces.Font)
	if !ok {
		return nil, fmt.Errorf("failed to convert configs to fonts")
	}
	return fonts, nil
}

// LoadLanguages loads all language configurations
func (l *Loader) LoadLanguages() ([]*interfaces.Language, error) {
	configs, err := l.loadConfigsFromDir("languages")
	if err != nil {
		return nil, err
	}
	languages, ok := configs.([]*interfaces.Language)
	if !ok {
		return nil, fmt.Errorf("failed to convert configs to languages")
	}
	return languages, nil
}

// LoadDotfiles loads all dotfile configurations
func (l *Loader) LoadDotfiles() ([]*interfaces.Dotfile, error) {
	configs, err := l.loadConfigsFromDir("dotfiles")
	if err != nil {
		return nil, err
	}
	dotfiles, ok := configs.([]*interfaces.Dotfile)
	if !ok {
		return nil, fmt.Errorf("failed to convert configs to dotfiles")
	}
	return dotfiles, nil
}

// LoadShells loads all shell configurations
func (l *Loader) LoadShells() ([]*interfaces.Shell, error) {
	configs, err := l.loadConfigsFromDir("shells")
	if err != nil {
		return nil, err
	}
	shells, ok := configs.([]*interfaces.Shell)
	if !ok {
		return nil, fmt.Errorf("failed to convert configs to shells")
	}
	return shells, nil
}

// LoadLanguageManagers loads all language manager configurations
func (l *Loader) LoadLanguageManagers() ([]*interfaces.Tool, error) {
	dir := filepath.Join(l.defaultsDir, "language_managers")
	managers := make([]*interfaces.Tool, 0)

	entries, err := l.configFS.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading language managers directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "schema.yaml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := l.configFS.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading language manager file %s: %w", path, err)
		}

		var manager interfaces.Tool
		if err := yaml.Unmarshal(data, &manager); err != nil {
			return nil, fmt.Errorf("error parsing language manager %s: %w", path, err)
		}

		managers = append(managers, &manager)
	}

	return managers, nil
}

// loadConfigsFromDir loads all configurations from both default and user directories
func (l *Loader) loadConfigsFromDir(dir string) (interface{}, error) {
	var configs interface{}
	
	// Load defaults first
	defaultConfigs, err := l.loadDefaultConfigs(dir)
	if err != nil {
		return nil, fmt.Errorf("error loading default configs: %w", err)
	}
	
	// Load user configs
	userConfigs, err := l.loadUserConfigs(dir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading user configs: %w", err)
	}
	
	// Merge configs based on type
	switch dir {
	case "tools":
		defaultTools, ok := defaultConfigs.([]*interfaces.Tool)
		if !ok {
			return nil, fmt.Errorf("invalid default tools configuration type: expected []*interfaces.Tool, got %T", defaultConfigs)
		}
		var userTools []*interfaces.Tool
		if userConfigs != nil {
			userTools, ok = userConfigs.([]*interfaces.Tool)
			if !ok {
				return nil, fmt.Errorf("invalid user tools configuration type: expected []*interfaces.Tool, got %T", userConfigs)
			}
		}
		configs = l.mergeToolConfigs(defaultTools, userTools)
	case "fonts":
		defaultFonts, ok := defaultConfigs.([]*interfaces.Font)
		if !ok {
			return nil, fmt.Errorf("invalid default fonts configuration type: expected []*interfaces.Font, got %T", defaultConfigs)
		}
		var userFonts []*interfaces.Font
		if userConfigs != nil {
			userFonts, ok = userConfigs.([]*interfaces.Font)
			if !ok {
				return nil, fmt.Errorf("invalid user fonts configuration type: expected []*interfaces.Font, got %T", userConfigs)
			}
		}
		configs = l.mergeFontConfigs(defaultFonts, userFonts)
	case "languages":
		defaultLanguages, ok := defaultConfigs.([]*interfaces.Language)
		if !ok {
			return nil, fmt.Errorf("invalid default languages configuration type: expected []*interfaces.Language, got %T", defaultConfigs)
		}
		var userLanguages []*interfaces.Language
		if userConfigs != nil {
			userLanguages, ok = userConfigs.([]*interfaces.Language)
			if !ok {
				return nil, fmt.Errorf("invalid user languages configuration type: expected []*interfaces.Language, got %T", userConfigs)
			}
		}
		configs = l.mergeLanguageConfigs(defaultLanguages, userLanguages)
	case "dotfiles":
		defaultDotfiles, ok := defaultConfigs.([]*interfaces.Dotfile)
		if !ok {
			return nil, fmt.Errorf("invalid default dotfiles configuration type: expected []*interfaces.Dotfile, got %T", defaultConfigs)
		}
		var userDotfiles []*interfaces.Dotfile
		if userConfigs != nil {
			userDotfiles, ok = userConfigs.([]*interfaces.Dotfile)
			if !ok {
				return nil, fmt.Errorf("invalid user dotfiles configuration type: expected []*interfaces.Dotfile, got %T", userConfigs)
			}
		}
		configs = l.mergeDotfileConfigs(defaultDotfiles, userDotfiles)
	case "shells":
		defaultShells, ok := defaultConfigs.([]*interfaces.Shell)
		if !ok {
			return nil, fmt.Errorf("invalid default shells configuration type: expected []*interfaces.Shell, got %T", defaultConfigs)
		}
		var userShells []*interfaces.Shell
		if userConfigs != nil {
			userShells, ok = userConfigs.([]*interfaces.Shell)
			if !ok {
				return nil, fmt.Errorf("invalid user shells configuration type: expected []*interfaces.Shell, got %T", userConfigs)
			}
		}
		configs = l.mergeShellConfigs(defaultShells, userShells)
	case "language_managers":
		defaultManagers, ok := defaultConfigs.([]*interfaces.Tool)
		if !ok {
			return nil, fmt.Errorf("invalid default language managers configuration type: expected []*interfaces.Tool, got %T", defaultConfigs)
		}
		var userManagers []*interfaces.Tool
		if userConfigs != nil {
			userManagers, ok = userConfigs.([]*interfaces.Tool)
			if !ok {
				return nil, fmt.Errorf("invalid user language managers configuration type: expected []*interfaces.Tool, got %T", userConfigs)
			}
		}
		configs = l.mergeToolConfigs(defaultManagers, userManagers)
	default:
		return nil, fmt.Errorf("unknown configuration type: %s", dir)
	}
	
	return configs, nil
}

// loadDefaultConfigs loads configurations from embedded defaults
func (l *Loader) loadDefaultConfigs(dir string) (interface{}, error) {
	defaultDir := filepath.Join(l.defaultsDir, dir)
	
	var configs interface{}
	switch dir {
	case "tools":
		tools := make([]*interfaces.Tool, 0)
		
		// Function to load tools from a directory
		var loadToolsFromDir func(string) error
		loadToolsFromDir = func(dirPath string) error {
			entries, err := l.configFS.ReadDir(dirPath)
			if err != nil {
				return fmt.Errorf("error reading directory %s: %w", dirPath, err)
			}
			
			for _, entry := range entries {
				if entry.IsDir() {
					// Recursively load tools from subdirectory
					subdir := filepath.Join(dirPath, entry.Name())
					if err := loadToolsFromDir(subdir); err != nil {
						return err
					}
					continue
				}
				
				if !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "schema.yaml" {
					continue
				}
				
				path := filepath.Join(dirPath, entry.Name())
				
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", path, err)
				}
				
				var tool interfaces.Tool
				if err := yaml.Unmarshal(data, &tool); err != nil {
					return fmt.Errorf("error parsing tool %s: %w", path, err)
				}
				
				// Set category based on subdirectory if not already set
				if tool.Category == "" {
					rel, _ := filepath.Rel(defaultDir, dirPath)
					if rel != "." {
						tool.Category = rel
					}
				}
				
				tools = append(tools, &tool)
			}
			return nil
		}
		
		// Start loading from the root tools directory
		if err := loadToolsFromDir(defaultDir); err != nil {
			return nil, err
		}
		
		configs = tools
		
	case "fonts":
		fonts := make([]*interfaces.Font, 0)
		var loadFontsFromDir func(string) error
		loadFontsFromDir = func(dirPath string) error {
			entries, err := l.configFS.ReadDir(dirPath)
			if err != nil {
				return fmt.Errorf("error reading directory %s: %w", dirPath, err)
			}
			
			for _, entry := range entries {
				if entry.IsDir() {
					subdir := filepath.Join(dirPath, entry.Name())
					if err := loadFontsFromDir(subdir); err != nil {
						return err
					}
					continue
				}
				
				if !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "schema.yaml" {
					continue
				}
				
				path := filepath.Join(dirPath, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", path, err)
				}
				
				var font interfaces.Font
				if err := yaml.Unmarshal(data, &font); err != nil {
					return fmt.Errorf("error parsing font %s: %w", path, err)
				}
				fonts = append(fonts, &font)
			}
			return nil
		}
		
		if err := loadFontsFromDir(defaultDir); err != nil {
			return nil, err
		}
		configs = fonts
		
	case "languages":
		languages := make([]*interfaces.Language, 0)
		var loadLanguagesFromDir func(string) error
		loadLanguagesFromDir = func(dirPath string) error {
			entries, err := l.configFS.ReadDir(dirPath)
			if err != nil {
				return fmt.Errorf("error reading directory %s: %w", dirPath, err)
			}
			
			for _, entry := range entries {
				if entry.IsDir() {
					subdir := filepath.Join(dirPath, entry.Name())
					if err := loadLanguagesFromDir(subdir); err != nil {
						return err
					}
					continue
				}
				
				if !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "schema.yaml" {
					continue
				}
				
				path := filepath.Join(dirPath, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", path, err)
				}
				
				var language interfaces.Language
				if err := yaml.Unmarshal(data, &language); err != nil {
					return fmt.Errorf("error parsing language %s: %w", path, err)
				}
				languages = append(languages, &language)
			}
			return nil
		}
		
		if err := loadLanguagesFromDir(defaultDir); err != nil {
			return nil, err
		}
		configs = languages
		
	case "dotfiles":
		dotfiles := make([]*interfaces.Dotfile, 0)
		var loadDotfilesFromDir func(string) error
		loadDotfilesFromDir = func(dirPath string) error {
			entries, err := l.configFS.ReadDir(dirPath)
			if err != nil {
				return fmt.Errorf("error reading directory %s: %w", dirPath, err)
			}
			
			for _, entry := range entries {
				if entry.IsDir() {
					subdir := filepath.Join(dirPath, entry.Name())
					if err := loadDotfilesFromDir(subdir); err != nil {
						return err
					}
					continue
				}
				
				if !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "schema.yaml" {
					continue
				}
				
				path := filepath.Join(dirPath, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", path, err)
				}
				
				var dotfile interfaces.Dotfile
				if err := yaml.Unmarshal(data, &dotfile); err != nil {
					return fmt.Errorf("error parsing dotfile %s: %w", path, err)
				}
				dotfiles = append(dotfiles, &dotfile)
			}
			return nil
		}
		
		if err := loadDotfilesFromDir(defaultDir); err != nil {
			return nil, err
		}
		configs = dotfiles
		
	case "shells":
		shells := make([]*interfaces.Shell, 0)
		var loadShellsFromDir func(string) error
		loadShellsFromDir = func(dirPath string) error {
			entries, err := l.configFS.ReadDir(dirPath)
			if err != nil {
				return fmt.Errorf("error reading directory %s: %w", dirPath, err)
			}
			
			for _, entry := range entries {
				if entry.IsDir() {
					subdir := filepath.Join(dirPath, entry.Name())
					if err := loadShellsFromDir(subdir); err != nil {
						return err
					}
					continue
				}
				
				if !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "schema.yaml" {
					continue
				}
				
				path := filepath.Join(dirPath, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", path, err)
				}
				
				var shell interfaces.Shell
				if err := yaml.Unmarshal(data, &shell); err != nil {
					return fmt.Errorf("error parsing shell %s: %w", path, err)
				}
				shells = append(shells, &shell)
			}
			return nil
		}
		
		if err := loadShellsFromDir(defaultDir); err != nil {
			return nil, err
		}
		configs = shells
		
	case "language_managers":
		managers := make([]*interfaces.Tool, 0)
		var loadManagersFromDir func(string) error
		loadManagersFromDir = func(dirPath string) error {
			entries, err := l.configFS.ReadDir(dirPath)
			if err != nil {
				return fmt.Errorf("error reading directory %s: %w", dirPath, err)
			}
			
			for _, entry := range entries {
				if entry.IsDir() {
					subdir := filepath.Join(dirPath, entry.Name())
					if err := loadManagersFromDir(subdir); err != nil {
						return err
					}
					continue
				}
				
				if !strings.HasSuffix(entry.Name(), ".yaml") || entry.Name() == "schema.yaml" {
					continue
				}
				
				path := filepath.Join(dirPath, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", path, err)
				}
				
				var tool interfaces.Tool
				if err := yaml.Unmarshal(data, &tool); err != nil {
					return fmt.Errorf("error parsing language manager %s: %w", path, err)
				}
				managers = append(managers, &tool)
			}
			return nil
		}
		
		if err := loadManagersFromDir(defaultDir); err != nil {
			return nil, err
		}
		configs = managers
		
	default:
		return nil, fmt.Errorf("unknown configuration type: %s", dir)
	}
	
	return configs, nil
}

// loadUserConfigs loads configurations from user directory
func (l *Loader) loadUserConfigs(dir string) (interface{}, error) {
	userDir := filepath.Join(l.baseDir, dir)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		return nil, err
	}
	
	var configs interface{}
	switch dir {
	case "tools":
		tools := make([]*interfaces.Tool, 0)
		err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
				return nil
			}
			tool, err := l.loadTool(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			tools = append(tools, tool)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", userDir, err)
		}
		configs = tools
	case "fonts":
		fonts := make([]*interfaces.Font, 0)
		err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
				return nil
			}
			font, err := l.loadFont(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			fonts = append(fonts, font)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", userDir, err)
		}
		configs = fonts
	case "languages":
		languages := make([]*interfaces.Language, 0)
		err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
				return nil
			}
			language, err := l.loadLanguage(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			languages = append(languages, language)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", userDir, err)
		}
		configs = languages
	case "dotfiles":
		dotfiles := make([]*interfaces.Dotfile, 0)
		err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
				return nil
			}
			dotfile, err := l.loadDotfile(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			dotfiles = append(dotfiles, dotfile)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", userDir, err)
		}
		configs = dotfiles
	case "shells":
		shells := make([]*interfaces.Shell, 0)
		err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
				return nil
			}
			shell, err := l.loadShell(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			shells = append(shells, shell)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", userDir, err)
		}
		configs = shells
	case "language_managers":
		managers := make([]*interfaces.Tool, 0)
		err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
				return nil
			}
			tool, err := l.loadTool(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			managers = append(managers, tool)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", userDir, err)
		}
		configs = managers
	default:
		return nil, fmt.Errorf("unknown configuration type: %s", dir)
	}
	
	return configs, nil
}

// mergeToolConfigs merges user tool configs into default configs
func (l *Loader) mergeToolConfigs(defaults, users []*interfaces.Tool) []*interfaces.Tool {
	if users == nil {
		return defaults
	}
	if defaults == nil {
		return users
	}
	
	// Create map of default tools by name
	defaultMap := make(map[string]*interfaces.Tool)
	for _, tool := range defaults {
		defaultMap[tool.Name] = tool
	}
	
	// Merge or append user tools
	result := make([]*interfaces.Tool, 0)
	for _, user := range users {
		if def, exists := defaultMap[user.Name]; exists {
			merged := mergeConfigs(def, user)
			result = append(result, merged)
			delete(defaultMap, user.Name)
		} else {
			result = append(result, user)
		}
	}
	
	// Add remaining defaults
	for _, def := range defaultMap {
		result = append(result, def)
	}
	
	return result
}

// mergeFontConfigs merges user font configs into default configs
func (l *Loader) mergeFontConfigs(defaults, users []*interfaces.Font) []*interfaces.Font {
	if len(users) == 0 {
		return defaults
	}

	// Create a map of default configs by name
	defaultMap := make(map[string]*interfaces.Font)
	for _, def := range defaults {
		defaultMap[def.Name] = def
	}

	// Merge or append user configs
	merged := make([]*interfaces.Font, 0)
	for _, user := range users {
		if def, exists := defaultMap[user.Name]; exists {
			merged = append(merged, mergeConfigs(def, user))
			delete(defaultMap, user.Name)
		} else {
			merged = append(merged, user)
		}
	}

	// Add remaining defaults
	for _, def := range defaultMap {
		merged = append(merged, def)
	}

	return merged
}

// mergeLanguageConfigs merges user language configs into default configs
func (l *Loader) mergeLanguageConfigs(defaults, users []*interfaces.Language) []*interfaces.Language {
	if len(users) == 0 {
		return defaults
	}

	// Create a map of default configs by name
	defaultMap := make(map[string]*interfaces.Language)
	for _, def := range defaults {
		defaultMap[def.Name] = def
	}

	// Merge or append user configs
	merged := make([]*interfaces.Language, 0)
	for _, user := range users {
		if def, exists := defaultMap[user.Name]; exists {
			merged = append(merged, mergeConfigs(def, user))
			delete(defaultMap, user.Name)
		} else {
			merged = append(merged, user)
		}
	}

	// Add remaining defaults
	for _, def := range defaultMap {
		merged = append(merged, def)
	}

	return merged
}

// mergeDotfileConfigs merges dotfile configurations
func (l *Loader) mergeDotfileConfigs(defaults, users []*interfaces.Dotfile) []*interfaces.Dotfile {
	merged := make(map[string]*interfaces.Dotfile)
	for _, d := range defaults {
		merged[d.Name] = d
	}
	for _, u := range users {
		if def, ok := merged[u.Name]; ok {
			// Simple override, or implement more complex merge if needed
			merged[u.Name] = mergeConfigs(def, u)
		} else {
			merged[u.Name] = u
		}
	}
	result := make([]*interfaces.Dotfile, 0, len(merged))
	for _, v := range merged {
		result = append(result, v)
	}
	return result
}

// mergeShellConfigs merges default and user shell configurations
func (l *Loader) mergeShellConfigs(defaults, users []*interfaces.Shell) []*interfaces.Shell {
	merged := make(map[string]*interfaces.Shell)
	for _, s := range defaults {
		merged[s.Name] = s // Assuming Name is unique identifier
	}
	for _, u := range users {
		if def, ok := merged[u.Name]; ok {
			// Simple override for now, can be made more granular
			merged[u.Name] = mergeConfigs(def, u)
		} else {
			merged[u.Name] = u
		}
	}
	result := make([]*interfaces.Shell, 0, len(merged))
	for _, v := range merged {
		result = append(result, v)
	}
	return result
}

// loadTool loads a tool configuration from a file
func (l *Loader) loadTool(path string) (*interfaces.Tool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}
	
	var tool interfaces.Tool
	if err := yaml.Unmarshal(data, &tool); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML from %s: %w", path, err)
	}
	
	return &tool, nil
}

// GetCategories returns a list of categories for a given configuration type
func (l *Loader) GetCategories(configType string) ([]string, error) {
	var dir string
	switch configType {
	case "tools":
		dir = filepath.Join(l.defaultsDir, "tools")
	case "fonts":
		dir = filepath.Join(l.defaultsDir, "fonts")
	case "languages":
		dir = filepath.Join(l.defaultsDir, "languages")
	case "shells":
		dir = filepath.Join(l.defaultsDir, "shells")
	// Note: Dotfiles and Shells might not have categories in the same way
	// but if they do, this can be extended.
	default:
		return nil, fmt.Errorf("unsupported config type for categories: %s", configType)
	}

	categories := make(map[string]bool)
	entries, err := l.configFS.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading default %s directory: %w", configType, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check if it's a direct subdirectory (potential category)
			// or a nested structure (like tools/category/subcategory)
			if configType == "tools" { // Tools can have subcategories
				subEntries, err := l.configFS.ReadDir(filepath.Join(dir, entry.Name()))
				if err != nil {
					// log or handle error, maybe it's not a category dir
					continue
				}
				for _, subEntry := range subEntries {
					if subEntry.IsDir() {
						categories[filepath.Join(entry.Name(), subEntry.Name())] = true
					} else if strings.HasSuffix(subEntry.Name(), ".yaml") && subEntry.Name() != "schema.yaml" {
						// If a .yaml file is directly in a category folder, that folder is a category
						categories[entry.Name()] = true
						break // Found one, no need to check other files in this dir
					}
				}
			} else {
				// For other types, direct subdirectories are categories
				categories[entry.Name()] = true
			}
		} else if strings.HasSuffix(entry.Name(), ".yaml") && entry.Name() != "schema.yaml" && configType == "tools" {
			// If a .yaml tool file is at the root of the 'tools' dir, it has no category (or a default one)
			// This logic might need adjustment based on how uncategorized items are handled.
			// For now, we assume categories are primarily directories.
		}
	}

	result := make([]string, 0, len(categories))
	for cat := range categories {
		result = append(result, cat)
	}
	return result, nil
}

// GetToolsByCategory returns all tools in a specific category
func (l *Loader) GetToolsByCategory(category, subcategory string) ([]*interfaces.Tool, error) {
	dir := filepath.Join(l.baseDir, category, subcategory)
	tools := make([]*interfaces.Tool, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		tool, err := l.loadTool(path)
		if err != nil {
			return fmt.Errorf("error loading tool from %s: %w", path, err)
		}

		tools = append(tools, tool)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
	}

	return tools, nil
}

// GetFonts loads all font configurations
func (l *Loader) GetFonts() ([]*interfaces.Font, error) {
	dir := filepath.Join(l.baseDir, "fonts")
	fonts := make([]*interfaces.Font, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		font, err := l.loadFont(path)
		if err != nil {
			return fmt.Errorf("error loading font from %s: %w", path, err)
		}

		fonts = append(fonts, font)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
	}

	return fonts, nil
}

// loadFont loads a single font configuration from a YAML file
func (l *Loader) loadFont(path string) (*interfaces.Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	var font interfaces.Font
	if err := yaml.Unmarshal(data, &font); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML from %s: %w", path, err)
	}

	return &font, nil
}

// GetLanguages loads all language configurations
func (l *Loader) GetLanguages() ([]*interfaces.Language, error) {
	dir := filepath.Join(l.baseDir, "languages")
	languages := make([]*interfaces.Language, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		language, err := l.loadLanguage(path)
		if err != nil {
			return fmt.Errorf("error loading language from %s: %w", path, err)
		}

		languages = append(languages, language)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
	}

	return languages, nil
}

// loadLanguage loads a single language configuration from a YAML file
func (l *Loader) loadLanguage(path string) (*interfaces.Language, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	var language interfaces.Language
	if err := yaml.Unmarshal(data, &language); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML from %s: %w", path, err)
	}

	return &language, nil
}

// GetDotfiles loads all dotfile configurations
func (l *Loader) GetDotfiles() ([]*interfaces.Dotfile, error) {
	dir := filepath.Join(l.baseDir, "dotfiles")
	dotfiles := make([]*interfaces.Dotfile, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		dotfile, err := l.loadDotfile(path)
		if err != nil {
			return fmt.Errorf("error loading dotfile from %s: %w", path, err)
		}

		dotfiles = append(dotfiles, dotfile)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
	}

	return dotfiles, nil
}

// loadDotfile loads a single dotfile configuration from a YAML file
func (l *Loader) loadDotfile(path string) (*interfaces.Dotfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}
	var dotfile interfaces.Dotfile
	if err := yaml.Unmarshal(data, &dotfile); err != nil {
		return nil, fmt.Errorf("error parsing dotfile %s: %w", path, err)
	}
	return &dotfile, nil
}

// loadShell loads a single shell configuration from a file
func (l *Loader) loadShell(path string) (*interfaces.Shell, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}
	var shell interfaces.Shell
	if err := yaml.Unmarshal(data, &shell); err != nil {
		return nil, fmt.Errorf("error parsing shell %s: %w", path, err)
	}
	return &shell, nil
}

// ExtractDefaults extracts default configurations to the user's config directory
func (l *Loader) ExtractDefaults() error {
	// Create all necessary directories
	dirs := []string{"tools", "fonts", "languages", "dotfiles", "language_managers", "shells"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(l.baseDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", filepath.Join(l.baseDir, dir), err)
		}
	}

	// Walk through the embedded defaults and copy to user config
	err := fs.WalkDir(l.configFS, l.defaultsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and schema files
		if d.IsDir() || strings.HasSuffix(d.Name(), "schema.yaml") {
			return nil
		}

		// Read the file
		data, err := l.configFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Create the target path
		relPath, err := filepath.Rel(l.defaultsDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		targetPath := filepath.Join(l.baseDir, relPath)

		// Ensure the target directory exists
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		// Write the file if it doesn't exist
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			if err := os.WriteFile(targetPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", targetPath, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to extract default configurations: %w", err)
	}

	return nil
} 