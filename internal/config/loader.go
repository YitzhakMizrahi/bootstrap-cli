package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"gopkg.in/yaml.v3"
)

//go:embed defaults/tools/modern/*.yaml
//go:embed defaults/tools/schema.yaml
//go:embed defaults/dotfiles/shell/*.yaml
//go:embed defaults/languages/*.yaml
//go:embed defaults/language_managers/*.yaml
var defaultConfigs embed.FS

// ConfigLoader handles loading and parsing configuration files
type ConfigLoader struct {
	baseDir     string // User config directory
	defaultsDir string // Embedded defaults directory
	configFS    embed.FS
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(baseDir string) *ConfigLoader {
	loader := &ConfigLoader{
		baseDir:     baseDir,
		defaultsDir: "defaults",
		configFS:    defaultConfigs,
	}
	
	// Extract embedded configs on initialization
	if err := loader.ExtractEmbeddedConfigs(); err != nil {
		// Log error but don't fail - we'll try to use user configs as fallback
		fmt.Printf("Warning: Failed to extract embedded configs: %v\n", err)
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
func (l *ConfigLoader) LoadTools() ([]*interfaces.Tool, error) {
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
func (l *ConfigLoader) LoadFonts() ([]*install.Font, error) {
	configs, err := l.loadConfigsFromDir("fonts")
	if err != nil {
		return nil, err
	}
	fonts, ok := configs.([]*install.Font)
	if !ok {
		return nil, fmt.Errorf("failed to convert configs to fonts")
	}
	return fonts, nil
}

// LoadLanguages loads all language configurations
func (l *ConfigLoader) LoadLanguages() ([]*install.Language, error) {
	configs, err := l.loadConfigsFromDir("languages")
	if err != nil {
		return nil, err
	}
	languages, ok := configs.([]*install.Language)
	if !ok {
		return nil, fmt.Errorf("failed to convert configs to languages")
	}
	return languages, nil
}

// LoadDotfiles loads all dotfile configurations
func (l *ConfigLoader) LoadDotfiles() ([]*interfaces.Dotfile, error) {
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

// LoadLanguageManagers loads all language manager configurations
func (l *ConfigLoader) LoadLanguageManagers() ([]*interfaces.Tool, error) {
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
func (l *ConfigLoader) loadConfigsFromDir(dir string) (interface{}, error) {
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
		configs = l.mergeToolConfigs(defaultConfigs.([]*interfaces.Tool), userConfigs.([]*interfaces.Tool))
	case "fonts":
		configs = l.mergeFontConfigs(defaultConfigs.([]*install.Font), userConfigs.([]*install.Font))
	case "languages":
		configs = l.mergeLanguageConfigs(defaultConfigs.([]*install.Language), userConfigs.([]*install.Language))
	case "dotfiles":
		configs = l.mergeDotfileConfigs(defaultConfigs.([]*interfaces.Dotfile), userConfigs.([]*interfaces.Dotfile))
	default:
		return nil, fmt.Errorf("unknown configuration type: %s", dir)
	}
	
	return configs, nil
}

// loadDefaultConfigs loads configurations from embedded defaults
func (l *ConfigLoader) loadDefaultConfigs(dir string) (interface{}, error) {
	defaultDir := filepath.Join(l.defaultsDir, dir)
	entries, err := l.configFS.ReadDir(defaultDir)
	if err != nil {
		return nil, fmt.Errorf("error reading default directory %s: %w", defaultDir, err)
	}
	
	var configs interface{}
	switch dir {
	case "tools":
		tools := make([]*interfaces.Tool, 0)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
				path := filepath.Join(defaultDir, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("error reading default file %s: %w", path, err)
				}
				
				var tool interfaces.Tool
				if err := yaml.Unmarshal(data, &tool); err != nil {
					return nil, fmt.Errorf("error parsing tool %s: %w", path, err)
				}
				tools = append(tools, &tool)
			}
		}
		configs = tools
	case "fonts":
		fonts := make([]*install.Font, 0)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
				path := filepath.Join(defaultDir, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("error reading default file %s: %w", path, err)
				}
				
				var font install.Font
				if err := yaml.Unmarshal(data, &font); err != nil {
					return nil, fmt.Errorf("error parsing font %s: %w", path, err)
				}
				fonts = append(fonts, &font)
			}
		}
		configs = fonts
	case "languages":
		languages := make([]*install.Language, 0)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
				path := filepath.Join(defaultDir, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("error reading default file %s: %w", path, err)
				}
				
				var language install.Language
				if err := yaml.Unmarshal(data, &language); err != nil {
					return nil, fmt.Errorf("error parsing language %s: %w", path, err)
				}
				languages = append(languages, &language)
			}
		}
		configs = languages
	case "dotfiles":
		dotfiles := make([]*interfaces.Dotfile, 0)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
				path := filepath.Join(defaultDir, entry.Name())
				data, err := l.configFS.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("error reading default file %s: %w", path, err)
				}
				
				var dotfile interfaces.Dotfile
				if err := yaml.Unmarshal(data, &dotfile); err != nil {
					return nil, fmt.Errorf("error parsing dotfile %s: %w", path, err)
				}
				dotfiles = append(dotfiles, &dotfile)
			}
		}
		configs = dotfiles
	}
	
	return configs, nil
}

// loadUserConfigs loads configurations from user directory
func (l *ConfigLoader) loadUserConfigs(dir string) (interface{}, error) {
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
		fonts := make([]*install.Font, 0)
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
		languages := make([]*install.Language, 0)
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
	}
	
	return configs, nil
}

// mergeToolConfigs merges user tool configs into default configs
func (l *ConfigLoader) mergeToolConfigs(defaults, users []*interfaces.Tool) []*interfaces.Tool {
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

// mergeFontConfigs merges font configurations
func (l *ConfigLoader) mergeFontConfigs(defaults, users []*install.Font) []*install.Font {
	merged := make([]*install.Font, 0)
	
	// Create a map of fonts by name
	fontMap := make(map[string]*install.Font)
	for _, font := range defaults {
		fontMap[font.Name] = font
	}
	
	// Merge or add user fonts
	for _, font := range users {
		if defaultFont, exists := fontMap[font.Name]; exists {
			merged = append(merged, mergeConfigs(defaultFont, font))
		} else {
			merged = append(merged, font)
		}
	}
	
	// Add remaining default fonts
	for _, font := range defaults {
		if _, exists := fontMap[font.Name]; !exists {
			merged = append(merged, font)
		}
	}
	
	return merged
}

// mergeLanguageConfigs merges language configurations
func (l *ConfigLoader) mergeLanguageConfigs(defaults, users []*install.Language) []*install.Language {
	merged := make([]*install.Language, 0)
	
	// Create a map of languages by name
	langMap := make(map[string]*install.Language)
	for _, lang := range defaults {
		langMap[lang.Name] = lang
	}
	
	// Merge or add user languages
	for _, lang := range users {
		if defaultLang, exists := langMap[lang.Name]; exists {
			merged = append(merged, mergeConfigs(defaultLang, lang))
		} else {
			merged = append(merged, lang)
		}
	}
	
	// Add remaining default languages
	for _, lang := range defaults {
		if _, exists := langMap[lang.Name]; !exists {
			merged = append(merged, lang)
		}
	}
	
	return merged
}

// mergeDotfileConfigs merges dotfile configurations
func (l *ConfigLoader) mergeDotfileConfigs(defaults, users []*interfaces.Dotfile) []*interfaces.Dotfile {
	merged := make([]*interfaces.Dotfile, 0)
	
	// Create a map of dotfiles by name
	dotfileMap := make(map[string]*interfaces.Dotfile)
	for _, dotfile := range defaults {
		dotfileMap[dotfile.Name] = dotfile
	}
	
	// Merge or add user dotfiles
	for _, dotfile := range users {
		if defaultDotfile, exists := dotfileMap[dotfile.Name]; exists {
			merged = append(merged, mergeConfigs(defaultDotfile, dotfile))
		} else {
			merged = append(merged, dotfile)
		}
	}
	
	// Add remaining default dotfiles
	for _, dotfile := range defaults {
		if _, exists := dotfileMap[dotfile.Name]; !exists {
			merged = append(merged, dotfile)
		}
	}
	
	return merged
}

// loadTool loads a tool configuration from a file
func (l *ConfigLoader) loadTool(path string) (*interfaces.Tool, error) {
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

// GetCategories returns all available categories for a given type
func (l *ConfigLoader) GetCategories(configType string) ([]string, error) {
	var categories []string
	dirPath := filepath.Join(l.baseDir, configType)

	// Walk through the directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip files and non-yaml files
		if !info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Get the relative path from the config type directory
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		categories = append(categories, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dirPath, err)
	}

	return categories, nil
}

// GetToolsByCategory returns all tools in a specific category
func (l *ConfigLoader) GetToolsByCategory(category, subcategory string) ([]*interfaces.Tool, error) {
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
func (l *ConfigLoader) GetFonts() ([]*install.Font, error) {
	dir := filepath.Join(l.baseDir, "fonts")
	fonts := make([]*install.Font, 0)

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
func (l *ConfigLoader) loadFont(path string) (*install.Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	var font install.Font
	if err := yaml.Unmarshal(data, &font); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML from %s: %w", path, err)
	}

	return &font, nil
}

// GetLanguages loads all language configurations
func (l *ConfigLoader) GetLanguages() ([]*install.Language, error) {
	dir := filepath.Join(l.baseDir, "languages")
	languages := make([]*install.Language, 0)

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
func (l *ConfigLoader) loadLanguage(path string) (*install.Language, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	var language install.Language
	if err := yaml.Unmarshal(data, &language); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML from %s: %w", path, err)
	}

	return &language, nil
}

// GetDotfiles loads all dotfile configurations
func (l *ConfigLoader) GetDotfiles() ([]*interfaces.Dotfile, error) {
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
func (l *ConfigLoader) loadDotfile(path string) (*interfaces.Dotfile, error) {
	var dotfile interfaces.Dotfile
	data, err := l.configFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading dotfile %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &dotfile); err != nil {
		return nil, fmt.Errorf("error parsing dotfile %s: %w", path, err)
	}
	return &dotfile, nil
} 