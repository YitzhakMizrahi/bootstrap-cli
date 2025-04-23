# 🔄 INTERFACES.md - Bootstrap CLI Interface Documentation

## 📦 Package Management

### PackageManager Interface
```go
type PackageManager interface {
    Install(packageName string) error
    IsInstalled(packageName string) bool
    GetName() string
    IsAvailable() bool
    Update() error
    Upgrade() error
    Remove(packageName string) error
    GetVersion(packageName string) (string, error)
    ListInstalled() ([]string, error)
}
```

**Purpose**: Provides a unified interface for system package managers (apt, dnf, pacman, homebrew).

**Methods**:
- `Install`: Installs a package by name
- `IsInstalled`: Checks if a package is installed
- `GetName`: Returns package manager name (apt, brew, dnf, pacman)
- `IsAvailable`: Checks if package manager is available on system
- `Update`: Updates package list
- `Upgrade`: Upgrades all packages
- `Remove`: Removes a package
- `GetVersion`: Gets version of installed package
- `ListInstalled`: Lists all installed packages

## 🛠 Tool Management

### Tool Struct
```go
type Tool struct {
    Name        string
    Description string
    Category    string
    Tags        []string
    
    PackageName  string
    PackageNames struct {
        APT    string
        Brew   string
        DNF    string
        Pacman string
    }

    Version       string
    VerifyCommand string
    PostInstall   []struct {
        Command     string
        Description string
    }

    ShellConfig struct {
        Aliases   map[string]string
        Functions map[string]string
        Env       map[string]string
    }

    Files []struct {
        Source      string
        Destination string
        Type        string
        Permissions int
        Content     string
    }
}
```

### ToolInstaller Interface
```go
type ToolInstaller interface {
    Install(tool *Tool) error
    Verify(tool *Tool) error
    IsInstalled(tool *Tool) bool
}
```

**Purpose**: Manages tool installation, verification, and configuration.

## 🐚 Shell Management

### ShellManager Interface
```go
type ShellManager interface {
    DetectCurrent() (*ShellInfo, error)
    ListAvailable() ([]*ShellInfo, error)
    IsInstalled(shell ShellType) bool
    GetInfo(shell ShellType) (*ShellInfo, error)
    ConfigureShell(config *ShellConfig) error
}
```

**Purpose**: Handles shell detection, configuration, and management.

### ShellInfo Struct
```go
type ShellInfo struct {
    Current     string
    Available   []string
    DefaultPath string
    Type        string
    Path        string
    Version     string
    IsDefault   bool
    IsAvailable bool
    ConfigFiles []string
}
```

### ShellConfig Struct
```go
type ShellConfig struct {
    Aliases   map[string]string
    Exports   map[string]string
    Functions map[string]string
    Path      []string
    Source    []string
}
```

## 📄 Dotfiles Management

### DotfilesManager Interface
```go
type DotfilesManager interface {
    Initialize() error
    CloneUserRepo(repoURL string) error
    ApplyDotfile(config *Dotfile) error
    BackupFile(filePath string, suffix string) error
    CreateSymlink(source, destination string) error
    WriteContentFile(content []byte, destination string) error
    ApplyShellConfig(config *ShellConfig) error
}
```

**Purpose**: Manages dotfile operations, including cloning, symlinking, and configuration.

### Dotfile Struct
```go
type Dotfile struct {
    Name            string
    Description     string
    Category        string
    Tags            []string
    Files           []DotfileFile
    Dependencies    []string
    ShellConfig     ShellConfig
    PostInstall     []string
    RequiresRestart bool
    SourceRepo      string
    BaseDir         string
    SymlinkStrategy SymlinkStrategy
}
```

### DotfileFile Struct
```go
type DotfileFile struct {
    Source       string
    Destination  string
    Operation    DotfileOperation
    Backup       bool
    BackupSuffix string
    Content      string
}
```

## 🔄 Constants and Types

### Shell Types
```go
type ShellType string

const (
    BashShell ShellType = "bash"
    ZshShell  ShellType = "zsh"
    FishShell ShellType = "fish"
)
```

### Package Manager Types
```go
type PackageManagerType string

const (
    APT      PackageManagerType = "apt"
    DNF      PackageManagerType = "dnf"
    Pacman   PackageManagerType = "pacman"
    Homebrew PackageManagerType = "brew"
)
```

### Dotfile Operations
```go
type DotfileOperation string

const (
    Create  DotfileOperation = "create"
    Update  DotfileOperation = "update"
    Delete  DotfileOperation = "delete"
    Symlink DotfileOperation = "symlink"
)
```

### Symlink Strategies
```go
type SymlinkStrategy int

const (
    SymlinkToHome SymlinkStrategy = iota
    SymlinkToDotfiles
    CopyToHome
)
```

## 📝 Best Practices

1. **Interface Location**
   - All interfaces MUST be defined in `internal/interfaces/`
   - Check existing interfaces before creating new ones
   - Keep interface implementations in respective packages

2. **Testing**
   - Mock implementations should be in `internal/testutil/`
   - Each interface should have corresponding tests
   - Use interfaces to enable dependency injection in tests

3. **Error Handling**
   - Use custom error types for domain-specific errors
   - Wrap errors with context using `fmt.Errorf`
   - Return early from functions when encountering errors

4. **Documentation**
   - All exported types must have godoc comments
   - Document breaking changes in interfaces
   - Keep interface documentation up to date

5. **Implementation Guidelines**
   - Follow interface segregation principle
   - Keep interfaces small and focused
   - Use composition over inheritance
   - Implement only required methods 