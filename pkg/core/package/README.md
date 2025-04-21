# Package Manager

The package manager package provides a platform-agnostic interface for managing system packages. It abstracts away the differences between various package managers (apt, yum, pacman, etc.) and provides a consistent API for package management operations.

## Usage

```go
import "github.com/YitzhakMizrahi/bootstrap-cli/pkg/core/package"

// Create a new package manager
manager := pkg.NewManager()

// Install a package
err := manager.Install("git")
if err != nil {
    log.Fatal(err)
}

// Check if a package is installed
if manager.IsInstalled("git") {
    fmt.Println("Git is installed")
}

// Uninstall a package
err = manager.Uninstall("git")
if err != nil {
    log.Fatal(err)
}
```

## Features

- Platform-agnostic package management
- Automatic privilege elevation when needed
- Support for multiple package managers:
  - apt (Debian/Ubuntu)
  - yum/dnf (Red Hat/CentOS/Fedora)
  - pacman (Arch Linux)
  - More to come...

## Interface

The package provides a `Manager` interface with the following methods:

```go
type Manager interface {
    Install(name string) error
    Uninstall(name string) error
    IsInstalled(name string) bool
    GetPackageManager() string
}
```

## Implementation Details

The package manager uses platform-specific implementations internally to perform the actual package management operations. This is handled through the `internal/platform` package, which provides implementations for different operating systems.

## Error Handling

All errors are properly wrapped with context using `fmt.Errorf`. Common error cases include:
- Failed privilege elevation
- Unsupported package manager
- Package installation/uninstallation failures

## Future Improvements

- [ ] Support for macOS (Homebrew)
- [ ] Package version management
- [ ] Dependency resolution
- [ ] Batch operations
- [ ] Package search functionality 