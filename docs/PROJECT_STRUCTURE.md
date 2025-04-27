# Project Structure

## Overview
Bootstrap CLI is organized into a modular structure that separates concerns and promotes code reuse. The main components are organized as follows:

```
bootstrap-cli/
├── cmd/                    # Command entrypoints
│   ├── init/              # Initialization command
│   └── up/                # Main TUI command
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   │   ├── defaults/      # Default configurations
│   │   ├── loader.go      # Config loader
│   │   └── schema/        # YAML schemas
│   ├── interfaces/        # Core type definitions
│   │   ├── tool.go
│   │   ├── font.go
│   │   └── language.go
│   ├── install/          # Installation logic
│   │   ├── font.go
│   │   ├── language.go
│   │   └── tool.go
│   ├── packages/         # Package management
│   │   ├── detector/     # System detection
│   │   └── factory/      # Package manager factory
│   └── ui/              # User interface
│       ├── app/         # Main application
│       ├── components/  # Reusable components
│       ├── screens/     # Screen implementations
│       ├── styles/      # UI styling
│       └── utils/       # UI utilities
└── docs/               # Documentation
    ├── CHANGELOG.md    # Change history
    ├── DECISIONS.md    # Architecture decisions
    └── SPEC.md         # Specifications
```

## Key Components

### Command Layer (`cmd/`)
- `init/`: Handles first-time setup
  - Configuration extraction
  - Environment setup
- `up/`: Main TUI application
  - Interactive installation flow
  - User selections
  - Progress tracking

### Configuration (`internal/config/`)
- Configuration loading and merging
- Default configurations in YAML
- Schema validation
- User override support

### Core Types (`internal/interfaces/`)
- Tool definitions
- Font specifications
- Language configurations
- Package manager interfaces

### Installation (`internal/install/`)
- Tool installation logic
- Font installation
- Language setup
- Error handling and rollback

### UI Layer (`internal/ui/`)
- Bubble Tea components
- Screen management
- Consistent styling
- Progress indicators

## Configuration Structure

### Tools
```yaml
name: string
description: string
category: string
package_names:
  apt: string[]
  dnf: string[]
  pacman: string[]
install_commands: string[]
verify_commands: string[]
```

### Fonts
```yaml
name: string
description: string
source: string
install_commands: string[]
verify_commands: string[]
```

### Languages
```yaml
name: string
description: string
version: string
manager: string
install_commands: string[]
verify_commands: string[]
```

## Development Workflow

1. Configuration Changes
   - Add/modify YAML in `internal/config/defaults/`
   - Update schemas if needed
   - Test with `init` command

2. UI Changes
   - Modify screens in `internal/ui/screens/`
   - Update components if needed
   - Style changes in `internal/ui/styles/`

3. Installation Logic
   - Update relevant files in `internal/install/`
   - Add error handling and rollback
   - Test with `up` command

4. Testing
   - Run linters: `go vet` and `revive`
   - Run tests: `go test ./...`
   - Manual testing with both commands

```
bootstrap-cli/
├── cmd/                # Command implementations
│   ├── init/          # Initialization command
│   └── up/            # Update command
├── internal/          # Internal packages
│   ├── config/        # Configuration management
│   │   ├── defaults/  # Default configurations
│   │   └── loader.go  # Configuration loader
│   ├── interfaces/    # Shared interfaces (ShellManager...)
│   ├── install/       # Installation logic
│   ├── pipeline/      # Pipeline-based installation
│   ├── packages/      # Package manager implementations
│   ├── ui/            # User interface components
│   └── utils/         # Utility functions
├── docs/              # Documentation
├── test/              # Test files
└── main.go            # Entry point
```

## Key Components

### Command Packages
- `cmd/init`: Handles system initialization
- `cmd/up`: Manages system updates

### Internal Packages
- `config`: Configuration management and loading
- `interfaces`: Shared interfaces and types
- `install`: Installation logic and tools
- `pipeline`: Pipeline-based installation system
- `packages`: Package manager implementations
- `ui`: User interface components
- `utils`: Utility functions

### Documentation
- `docs/`: Project documentation
  - `INTERFACES.md`: Interface documentation
  - `IMPLEMENTATION.md`: Implementation details
  - `DECISIONS.md`: Architecture decisions

## 🧱 High-Level Directory Layout
```
bootstrap-cli/
├── cmd/                # CLI commands via Cobra
├── internal/           # Main logic (system, install, shell, flow...)
│   ├── system/         # OS/arch/distro detection
│   ├── packages/       # Package manager detection + abstraction
│   ├── install/        # Tool installation logic
│   ├── shell/          # Shell selection, .rc file writing
│   ├── dotfiles/       # GitHub clone only (MVP scope)
│   ├── config/         # YAML config load/save/validate
│   ├── flow/           # Guided CLI flows (init, install, shell, etc.)
│   ├── ui/             # Prompt modules, spinners, selections
│   ├── symlinks/       # Shared symlink + PATH config logic
│   ├── utils/          # Logger, paths, validations
│   ├── interfaces/     # Shared interfaces (ToolInstaller, ShellManager...)
│   └── testutil/       # Reusable mocks, stubs, and helpers for unit tests
├── pkg/                # Optional public packages and templates
│   ├── templates/      # Static config templates (e.g. .zshrc, starship.toml)
│   ├── plugin/         # Optional plugin loader (post-MVP)
│   └── i18n/           # Language packs (future)
├── test/               # Integration + e2e tests
│   ├── integration/    # Real install test (via LXC)
│   ├── fixtures/       # Static test data (YAML configs, test plans)
│   └── e2e/            # Simulated full user flow test (init → up)
├── docs/               # Specifications + guides
├── scripts/            # Helper scripts (build, test, release)
├── .github/            # CI/CD config, issue templates
├── main.go             # Entrypoint
└── Makefile            # Build shortcuts
```

## 🔧 CI/Linting
- `.golangci.yml` – Includes: gofmt, golint, govet, errcheck
- `.github/workflows/` – `test.yml`, `release.yml`

## 📜 Developer Workflow
```bash
git clone https://github.com/YitzhakMizrahi/bootstrap-cli.git
cd bootstrap-cli
make deps
make build
make test
```

To run manually (instead of Makefile):
```bash
go build -o build/bin/bootstrap-cli main.go
```

## 📆 Makefile Targets
```makefile
build:
	go build -o build/bin/bootstrap-cli main.go

test:
	go test ./...

clean:
	rm -rf build/

lint:
	golangci-lint run ./...

run: build
	./build/bin/bootstrap-cli

deps:
	go mod download
	go mod tidy

build-lxc:
	GOOS=linux GOARCH=amd64 go build -o build/bin/bootstrap-cli-linux-amd64 main.go

validate:
	./build/bin/bootstrap-cli validate

all: lint test build

release:
	./scripts/release/create-release.sh
```

## 🤮 LXC Testing Tips
```bash
# Launch fresh Ubuntu LXC container
lxc launch ubuntu:22.04 bootstrap-test

# Create a clean snapshot to restore from if needed
lxc snapshot bootstrap-test clean-setup

# Push the compiled binary with executable permissions
lxc file push build/bin/bootstrap-cli bootstrap-test/home/devuser/bootstrap-cli --mode=755

# Run the CLI interactively from inside the container
lxc exec bootstrap-test -- su - devuser
```

---

## 🔑 Cobra Command Structure
```
cmd/
├── root.go         # Root command
├── up.go           # Run entire flow (new)
├── init.go         # Interactive setup
├── detect.go       # System info
├── install.go      # Tool install
├── dotfiles.go     # GitHub clone only
├── shell.go        # Shell setup
├── languages.go    # Runtime installs
├── font.go         # Font installer
├── validate.go     # Validate setup
├── config.go       # View/export config
├── version.go      # Print version
```

