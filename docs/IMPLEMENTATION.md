# 🛠 IMPLEMENTATION.md – Bootstrap CLI Implementation Plan

## ✅ Phase Breakdown

### 📦 Phase 1: Core Infrastructure
- [ ] System detection (OS, distro, arch)
- [ ] Package manager detection and abstraction
- [ ] Core Tool + Installer interface
- [ ] Tool verification and validation
- [ ] Modular flow logic in `internal/flow/`
- [ ] Symlink task struct for unified path/config management
- [ ] Tests for package ops

### 🐚 Phase 2: Shell & Configuration
- [ ] Shell detection and config writing
- [ ] Dotfile clone from GitHub (optional)
- [ ] YAML config loader/saver
- [ ] Configuration validation
- [ ] Apply declared symlinks via shared handler
- [ ] Dotfile symlink and PATH setup validation
- [ ] Tests for config and dotfiles

### 📚 Phase 3: Enhanced Features
- [ ] pyenv, nvm, rustup, goenv support
- [ ] Font installer (JetBrains Nerd)
- [ ] Plugin system scaffold (deferred post-MVP)
- [ ] Bubbletea CLI UI enhancements (experimental in v2)
- [ ] Config preview screen
- [ ] Notification + logs

### 🚀 Phase 4: Polish & Optimization
- [ ] Parallel installs
- [ ] Caching, lazy loading
- [ ] Error recovery and logging
- [ ] End-to-end tests + snapshots
- [ ] Finalize docs, help commands

---

## 🔍 Testing Strategy
- `go test ./...`
- Unit tests per internal package
- Integration tests in LXC container
- Fixture-based test data

---

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

---

## 🧩 CLI Commands (MVP)
- `bootstrap-cli up` - full flow orchestration
- `bootstrap-cli init` - interactive setup
- `bootstrap-cli detect` - detect system
- `bootstrap-cli dotfiles` - clone from GitHub
- `bootstrap-cli shell` - install & configure shell
- `bootstrap-cli install` - install core tools
- `bootstrap-cli languages` - install runtimes
- `bootstrap-cli font` - install fonts
- `bootstrap-cli validate` - validate setup
- `bootstrap-cli config` - view/export config
- `bootstrap-cli version` - show version

---

## 🔗 Symlink & Path Setup
Bootstrap CLI ensures installed tools and runtime environments are ready out-of-the-box by:

- Defining `SymlinkTask` structs from all flows (dotfiles, shell, languages)
- Applying symlinks via a central `ApplySymlinks([]SymlinkTask)` handler
- Creating symlinks from cloned dotfiles to standard locations (e.g., `~/.zshrc`)
- Appending PATH and language init commands to `.zshrc`/`.bashrc` as needed
- Verifying correctness via shell validation step (`command -v`, `which`)

---

## 🔮 Future Enhancements
- Notification history
- Remote config sync + backup
- Rollback support
- i18n + accessibility
- GUI launcher wrapper

---

## 📌 Notes
- Each module testable in isolation
- Use Go interfaces for pluggable logic
- Use promptui/survey for MVP prompts
- Bubbletea reserved for v2 or advanced flows
- Refer to `PROJECT_STRUCTURE.md` for folders

