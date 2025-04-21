# Bootstrap CLI Project Structure

```
bootstrap-cli/
├── .github/                    # GitHub specific files
│   ├── workflows/             # GitHub Actions workflows
│   │   ├── test.yml          # CI/CD pipeline
│   │   └── release.yml       # Release automation
│   └── ISSUE_TEMPLATE/       # Issue templates
│
├── cmd/                       # Command-line entry points
│   └── bootstrap/            # Main CLI application
│       └── main.go           # Entry point
│
├── internal/                  # Private application code
│   ├── core/                 # Core functionality
│   │   ├── system/          # System detection and info
│   │   ├── package/         # Package manager operations
│   │   ├── shell/           # Shell management
│   │   └── config/          # Configuration management
│   │
│   ├── install/             # Installation logic
│   │   ├── tools/          # Tool installation
│   │   ├── languages/      # Programming language setup
│   │   ├── fonts/          # Font installation
│   │   └── dotfiles/       # Dotfiles management
│   │
│   ├── ui/                 # User interface components
│   │   ├── progress/       # Progress indicators
│   │   ├── prompts/        # User prompts
│   │   └── display/        # Display formatting
│   │
│   └── utils/              # Utility functions
│       ├── logger/         # Logging utilities
│       ├── validator/      # Input validation
│       └── security/       # Security utilities
│
├── pkg/                     # Public libraries
│   ├── plugin/             # Plugin system
│   │   ├── manager/        # Plugin management
│   │   └── api/           # Plugin API
│   │
│   ├── config/            # Configuration templates
│   │   ├── shell/         # Shell configurations
│   │   ├── editor/        # Editor configurations
│   │   └── tools/         # Tool configurations
│   │
│   └── i18n/              # Internationalization
│       └── locales/       # Language files
│
├── scripts/                # Build and maintenance scripts
│   ├── build/             # Build scripts
│   ├── test/              # Test scripts
│   └── release/           # Release scripts
│
├── test/                  # Test files
│   ├── integration/       # Integration tests
│   ├── e2e/              # End-to-end tests
│   └── fixtures/         # Test fixtures
│
├── docs/                  # Documentation
│   ├── api/              # API documentation
│   ├── guides/           # User guides
│   └── examples/         # Usage examples
│
├── build/                # Build artifacts
│   ├── bin/             # Compiled binaries
│   └── dist/            # Distribution packages
│
├── .gitignore           # Git ignore file
├── .golangci.yml        # GolangCI-Lint configuration
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Makefile            # Build automation
├── README.md           # Project overview
└── LICENSE             # License file
```

## Directory Details

### `.github/`
- Contains GitHub-specific configurations
- CI/CD workflows for automated testing and releases
- Issue and PR templates for standardized contributions

### `cmd/`
- Entry points for the application
- `bootstrap/` contains the main CLI application
- Each subdirectory represents a different executable

### `internal/`
- Private application code not meant for external use
- Organized by functional domains:
  - `core/`: Core system functionality
  - `install/`: Installation and setup logic
  - `ui/`: User interface components
  - `utils/`: Common utilities

### `pkg/`
- Public libraries that can be used by other projects
- Includes:
  - Plugin system
  - Configuration templates
  - Internationalization support

### `scripts/`
- Build and maintenance scripts
- Test automation
- Release management
- Development utilities

### `test/`
- Test files organized by type:
  - Integration tests
  - End-to-end tests
  - Test fixtures and mocks

### `docs/`
- Comprehensive documentation
- API references
- User guides
- Example configurations

### `build/`
- Build artifacts and distribution packages
- Platform-specific binaries
- Package distributions

## Key Files

### `Makefile`
```makefile
.PHONY: build test clean release

# Build targets
build:
	go build -o build/bin/bootstrap-cli cmd/bootstrap/main.go

# Test targets
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf build/

# Release management
release:
	./scripts/release/create-release.sh
```

### `.golangci.yml`
```yaml
linters:
  enable:
    - gofmt
    - golint
    - govet
    - errcheck
    - staticcheck
    - gosimple

run:
  deadline: 5m
```

### `go.mod`
```go
module github.com/YitzhakMizrahi/bootstrap-cli

go 1.21

require (
    github.com/spf13/cobra v1.7.0
    github.com/spf13/viper v1.16.0
    // Add other dependencies
)
```

## Development Workflow

1. **Setup Development Environment**
   ```bash
   # Clone repository
   git clone https://github.com/YitzhakMizrahi/bootstrap-cli.git
   cd bootstrap-cli

   # Install dependencies
   go mod download

   # Build project
   make build
   ```

2. **Running Tests**
   ```bash
   # Run all tests
   make test

   # Run specific test
   go test ./internal/core/...
   ```

3. **Building for Distribution**
   ```bash
   # Build for current platform
   make build

   # Build for all platforms
   make release
   ```

## LXC Container Testing

For testing in LXC containers:

1. **Container Setup**
   ```bash
   # Create Ubuntu container
   lxc launch ubuntu:22.04 bootstrap-test

   # Create snapshot
   lxc snapshot bootstrap-test bootstrap-test-snapshot
   ```

2. **Test Environment**
   ```bash
   # Restore from snapshot for clean test
   lxc restore bootstrap-test bootstrap-test-snapshot

   # Copy test files
   lxc file push build/bin/bootstrap-cli bootstrap-test/usr/local/bin/
   ```

3. **Run Tests**
   ```bash
   # Execute in container
   lxc exec bootstrap-test -- bootstrap-cli test
   ```