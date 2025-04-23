package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// ConfigLoader handles loading and parsing configuration files
type ConfigLoader struct {
	baseDir string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(baseDir string) *ConfigLoader {
	return &ConfigLoader{
		baseDir: baseDir,
	}
}

// LoadTools loads all tool configurations
func (l *ConfigLoader) LoadTools() ([]*install.Tool, error) {
	configs, err := l.loadConfigsFromDir("tools")
	if err != nil {
		return nil, err
	}
	tools, ok := configs.([]*install.Tool)
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
		return nil, fmt.Errorf("invalid dotfile configuration")
	}

	return dotfiles, nil
}

// loadConfigsFromDir loads all configurations from a specific directory
func (l *ConfigLoader) loadConfigsFromDir(dir string) (interface{}, error) {
	dirPath := filepath.Join(l.baseDir, dir)
	var configs interface{}

	switch dir {
	case "tools":
		tools := make([]*install.Tool, 0)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
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
			return nil, fmt.Errorf("error walking directory %s: %w", dirPath, err)
		}
		configs = tools

	case "fonts":
		fonts := make([]*install.Font, 0)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
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
			return nil, fmt.Errorf("error walking directory %s: %w", dirPath, err)
		}
		configs = fonts

	case "languages":
		languages := make([]*install.Language, 0)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
				return nil
			}
			lang, err := l.loadLanguage(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			languages = append(languages, lang)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", dirPath, err)
		}
		configs = languages

	case "dotfiles":
		dotfiles := make([]*interfaces.Dotfile, 0)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
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
			return nil, fmt.Errorf("error walking directory %s: %w", dirPath, err)
		}
		configs = dotfiles

	default:
		return nil, fmt.Errorf("unknown configuration type: %s", dir)
	}

	return configs, nil
}

// loadConfigFile loads a single configuration file
func (l *ConfigLoader) loadConfigFile(path string) (*install.Tool, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var tool install.Tool
	if err := v.Unmarshal(&tool); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
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
func (l *ConfigLoader) GetToolsByCategory(category, subcategory string) ([]*install.Tool, error) {
	dir := filepath.Join(l.baseDir, category, subcategory)
	tools := make([]*install.Tool, 0)

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

// loadTool loads a single tool configuration from a YAML file
func (l *ConfigLoader) loadTool(path string) (*install.Tool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	var tool install.Tool
	if err := yaml.Unmarshal(data, &tool); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML from %s: %w", path, err)
	}

	return &tool, nil
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

// GetDotfiles returns all dotfile configurations
func (l *ConfigLoader) GetDotfiles() ([]*interfaces.Dotfile, error) {
	dotfiles := make([]*interfaces.Dotfile, 0)
	
	files, err := filepath.Glob(filepath.Join(l.baseDir, "dotfiles", "*.yaml"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		dotfile, err := l.loadDotfile(file)
		if err != nil {
			return nil, err
		}
		dotfiles = append(dotfiles, dotfile)
	}

	return dotfiles, nil
}

// loadDotfile loads a single dotfile configuration
func (l *ConfigLoader) loadDotfile(path string) (*interfaces.Dotfile, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dotfile interfaces.Dotfile
	if err := yaml.Unmarshal(data, &dotfile); err != nil {
		return nil, err
	}

	return &dotfile, nil
} 