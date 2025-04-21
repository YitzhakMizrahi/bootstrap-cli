# 🚀 Bootstrap CLI

A modern, interactive command-line tool for setting up and managing development environments across different platforms.

## 🎯 Features

- 🖥️ **System Detection & Setup**
  - Automatic OS detection and configuration
  - Package manager integration
  - System-specific optimizations

- 🛠️ **Development Tools**
  - Modern CLI tools installation
  - Programming language setup
  - Shell configuration
  - Font management

- ⚙️ **Configuration Management**
  - Dotfiles handling
  - Shell customization
  - Tool-specific configs
  - Backup and restore

- 🔌 **Plugin System**
  - Shell plugin management
  - Custom plugin support
  - Version management

## 🏗️ Project Structure

```
bootstrap-cli/
├── cmd/          # Command-line entry points
├── internal/     # Private application code
├── pkg/          # Public libraries
├── scripts/      # Build and maintenance scripts
├── test/         # Test files
└── docs/         # Documentation
```

## 🚦 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Linux (Ubuntu/Debian) or macOS

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/bootstrap-cli.git
cd bootstrap-cli

# Install dependencies
go mod download

# Build the project
make build
```

### Development Setup

```bash
# Run tests
make test

# Run specific tests
go test ./internal/core/...

# Build for development
make build
```

## 🧪 Testing

We use LXC containers for testing to ensure isolation and reproducibility:

```bash
# Create test container
lxc launch ubuntu:22.04 bootstrap-test

# Create snapshot
lxc snapshot bootstrap-test bootstrap-test-snapshot

# Run tests in container
lxc file push build/bin/bootstrap-cli bootstrap-test/usr/local/bin/
lxc exec bootstrap-test -- bootstrap-cli test
```

## 📚 Documentation

- [CLI Overview](docs/CLI_OVERVIEW.md) - Comprehensive feature documentation
- [Project Structure](docs/PROJECT_STRUCTURE.md) - Detailed project architecture
- [Implementation Guide](docs/IMPLEMENTATION.md) - Technical implementation details

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/guides/CONTRIBUTING.md) for details.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔄 Development Status

Current Version: 0.1.0 (In Development)

Status: Active Development - Core Features Implementation

## 📞 Contact

- Issue Tracker: [GitHub Issues](https://github.com/yourusername/bootstrap-cli/issues)
- Source Code: [GitHub Repository](https://github.com/yourusername/bootstrap-cli)

---

## ✅ 3. Optional Nice-to-Haves

| File | Purpose |
|------|---------|
| `LICENSE` | Use MIT if unsure — simple, permissive |
| `Makefile` | Optional shortcut for common tasks (build, run, fmt) |
| `docs/` | Placeholder for later usage or architecture notes |
| `CONTRIBUTING.md` | If you plan to open it up later |
| `starship.toml`, `.zshrc` samples | As optional templates for dotfiles repo setup |

---

