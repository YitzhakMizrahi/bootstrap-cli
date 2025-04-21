# Shell Management Package

The shell management package provides a platform-agnostic interface for managing shell environments, configurations, and plugins.

## Features

- Multi-shell support (bash, zsh, fish)
- Plugin manager integration
- Shell configuration management
- Environment variable management
- Path management
- Configuration backup and restore

## Usage

```go
import (
    "github.com/YitzhakMizrahi/bootstrap-cli/pkg/core/shell"
    "github.com/YitzhakMizrahi/bootstrap-cli/pkg/core/shell/templates"
)

// Create shell configuration
config := &shell.Config{
    Type:      shell.Zsh,
    PluginMgr: shell.OhMyZsh,
    Plugins:   []string{"git", "docker"},
}

// Create shell manager
manager, err := shell.NewShellManager(config)
if err != nil {
    log.Fatal(err)
}

// Setup shell environment
if err := manager.Setup(config); err != nil {
    log.Fatal(err)
}

// Add configuration
manager.AddAlias("ll", "ls -la")
manager.AddEnvVar("EDITOR", "vim")
manager.AddPath("/usr/local/bin")

// Generate shell configuration
data := &templates.TemplateData{
    Plugins:     config.Plugins,
    Aliases:     config.Aliases,
    EnvVars:     config.EnvVars,
    Path:        config.Path,
    CustomPaths: []string{},
}

if err := manager.GenerateConfig(data); err != nil {
    log.Fatal(err)
}
```

## Shell Types

The package supports multiple shell types:
- `shell.Bash`: Bash shell
- `shell.Zsh`: Zsh shell
- `shell.Fish`: Fish shell

## Plugin Managers

Supported plugin managers:
- `shell.OhMyZsh`: Oh My Zsh for Zsh
- `shell.Antigen`: Antigen for Zsh
- `shell.Zinit`: Zinit for Zsh
- `shell.Fisherman`: Fisherman for Fish

## Configuration Management

The package provides methods for:
- Adding aliases
- Setting environment variables
- Managing PATH entries
- Backing up and restoring configurations

## Template System

Shell configurations are generated using Go templates. The template system supports:
- Plugin sourcing
- Alias definitions
- Environment variable exports
- PATH management
- Custom shell-specific configurations

## Error Handling

All errors are properly wrapped with context using `fmt.Errorf`. Common error cases:
- Shell installation failures
- Plugin manager installation failures
- Configuration file operations
- Permission issues

## Platform Support

The package uses the platform interface to handle:
- Shell detection and installation
- Path resolution
- System-specific operations

## Future Improvements

- [ ] Support for more plugin managers
- [ ] Plugin dependency management
- [ ] Configuration profiles
- [ ] Shell completion generation
- [ ] Theme management
- [ ] Custom plugin repositories 