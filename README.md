# 🛍️ Bootstrap CLI

**Bootstrap CLI** is a modular, interactive command-line tool for setting up a fully personalized development environment across platforms (Linux, macOS, WSL).

Built for developers who value speed, reproducibility, and elegance — without fluff.

---

## ✨ Features (Phase 1 & 2)

- 🔍 **System Detection**  
  Detect OS, architecture, distro, and system resources  
- 📦 **Package Manager Abstraction**  
  Unified logic for `apt`, `dnf`, `pacman`, `brew`, etc.  
- 🛠️ **Core Tool Installer**  
  Modern CLI tools (e.g. `bat`, `fzf`, `ripgrep`, `starship`)  
- 👋 **Shell & Config Setup**  
  Zsh, Bash, Fish, dotfiles sync, and config templates  
- 🐍 **Language Version Managers**  
  `nvm`, `pyenv`, `goenv`, `rustup` with custom version support  
- 🔧 **Dotfiles Management**  
  Local import, GitHub sync, backup & restore  

_Plugin system, remote config sync, i18n, and accessibility coming in later phases._

---

## 📁 Project Structure

```bash
bootstrap-cli/
├── cmd/                # CLI entrypoints (Cobra commands)
├── internal/           # Business logic: system, packages, install, etc.
├── pkg/                # (Future) public plugins or modules
├── test/               # Integration tests, fixtures
├── docs/               # Specs, plans, architecture
├── scripts/            # Build/test helpers
├── main.go             # Entrypoint
└── cursor-prompt.md    # Prompt used for Cursor AI guidance
```

---

## ⚙️ Getting Started

### Prerequisites
- Go 1.21+
- Git
- Linux/macOS/WSL (tested on Ubuntu 22.04+)

### Install & Build

```bash
git clone https://github.com/YitzhakMizrahi/bootstrap-cli.git
cd bootstrap-cli
go mod tidy
make build
```

### Run

```bash
./bootstrap-cli init
```

---

## 🧪 Testing (LXC Method)

```bash
lxc launch ubuntu:22.04 bootstrap-test
lxc snapshot bootstrap-test clean

# Push binary to test container
lxc file push build/bin/bootstrap-cli bootstrap-test/usr/local/bin/
lxc exec bootstrap-test -- bootstrap-cli init
```

---

## 📚 Documentation

| Topic | Link |
|-------|------|
| 💡 Specification | [`docs/SPEC.md`](docs/SPEC.md) |
| 🧱 Implementation Plan | [`docs/IMPLEMENTATION.md`](docs/IMPLEMENTATION.md) |
| 🦮 Project Structure | [`docs/PROJECT_STRUCTURE.md`](docs/PROJECT_STRUCTURE.md) |
| 🎨 Wireframes | [`docs/WIREFRAMES.md`](docs/WIREFRAMES.md) |
| 📓 Decisions | [`docs/DECISIONS.md`](docs/DECISIONS.md) |

---

## 🛠 Dev Workflow

```bash
make build        # Compile CLI binary
make test         # Run all unit tests
go run main.go    # Run directly
```

To build for release, check `scripts/release/` and `Makefile` targets.

---

## 🚧 Development Status

| Status | Version |
|--------|---------|
| Core Implementation | In Progress |
| Latest Version | `v0.1.0-alpha` |

---

## 🤝 Contributing

We're currently focusing on internal use, but feel free to open issues or forks.  
A formal CONTRIBUTING guide will be added in Phase 3.

---

## 📝 License

MIT License – see [`LICENSE`](LICENSE) file.

---

