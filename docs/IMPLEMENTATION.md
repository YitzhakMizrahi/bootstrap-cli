# 🛠 IMPLEMENTATION.md – Bootstrap CLI Implementation Plan

## ✅ Phase Breakdown

### 📦 Phase 1: Core Infrastructure [COMPLETED]
- [x] System detection (OS, distro, arch)
- [x] Package manager detection and abstraction
- [x] Core Tool + Installer interface
- [x] Tool verification and validation
- [x] Modular flow logic in `internal/flow/`
- [x] Symlink task struct for unified path/config management
- [x] Tests for package ops

### 🐚 Phase 2: Shell & Configuration [IN PROGRESS]
- [x] Shell detection and config writing
- [ ] Dotfile clone from GitHub (in progress)
- [x] YAML config loader/saver
- [x] Configuration validation
- [ ] Apply declared symlinks via shared handler (in progress)
- [ ] Dotfile symlink and PATH setup validation (in progress)
- [ ] Tests for config and dotfiles (in progress)

### 📚 Phase 3: Enhanced Features [PLANNED]
- [ ] Language runtime support
  - [ ] pyenv for Python
  - [ ] nvm for Node.js
  - [ ] rustup for Rust
  - [ ] goenv for Go
- [ ] Font installer (JetBrains Nerd)
- [ ] Plugin system scaffold (deferred post-MVP)
- [ ] Bubbletea CLI UI enhancements (experimental in v2)
- [ ] Config preview screen
- [ ] Notification + logs

### 🚀 Phase 4: Polish & Optimization [PLANNED]
- [ ] Parallel installs
- [ ] Caching, lazy loading
- [ ] Error recovery and logging
- [ ] End-to-end tests + snapshots
- [ ] Finalize docs, help commands

---

## 🔍 Testing Strategy
- Unit tests per internal package
- Integration tests in LXC container
- Fixture-based test data
- Test coverage targets:
  - Core packages: 80%+
  - UI/Flow packages: 70%+
  - Utils: 90%+

## 📦 Config File Schema
```yaml
shell: zsh
fonts:
  - JetBrains Mono Nerd
languages:
  node: lts
  python: 3.11
  go: latest
  rust: stable
tools:
  - git
  - bat
  - lsd
```

## 🧩 Core Interfaces
### PackageManager
- Handles system package operations
- Supports apt, dnf, pacman, homebrew
- Methods: Install, Remove, Update, IsInstalled

### ToolInstaller
- Manages tool installation and verification
- Supports system-specific package names
- Handles post-install configuration

### ShellManager
- Detects and configures shells
- Manages shell configuration files
- Supports bash, zsh, fish

### DotfilesManager
- Handles dotfile operations
- Supports GitHub repository cloning
- Manages symlinks and backups

## 🔗 Symlink & Path Setup
Bootstrap CLI ensures installed tools and runtime environments are ready out-of-the-box by:

- Defining `SymlinkTask` structs from all flows (dotfiles, shell, languages)
- Applying symlinks via a central `ApplySymlinks([]SymlinkTask)` handler
- Creating symlinks from cloned dotfiles to standard locations (e.g., `~/.zshrc`)
- Appending PATH and language init commands to `.zshrc`/`.bashrc` as needed
- Verifying correctness via shell validation step (`command -v`, `which`)

## 🔮 Future Enhancements
- Notification history
- Remote config sync + backup
- Rollback support
- i18n + accessibility
- GUI launcher wrapper

## 📌 Notes
- Each module testable in isolation
- Use Go interfaces for pluggable logic
- Use promptui/survey for MVP prompts
- Bubbletea reserved for v2 or advanced flows
- Refer to `PROJECT_STRUCTURE.md` for folders

