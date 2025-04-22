# 🛠 IMPLEMENTATION.md – Bootstrap CLI Implementation Plan

## ✅ Phase Breakdown

### 📦 Phase 1: Core Infrastructure
- [ ] System detection (OS, distro, arch)
- [ ] Package manager detection and abstraction
- [ ] Core tool struct + install logic
- [ ] Tool verification and validation
- [ ] Tests for package ops

### 🐚 Phase 2: Shell & Configuration
- [ ] Shell detection and config writing
- [ ] Dotfile management (clone, sync, backup)
- [ ] Load/save config files (YAML)
- [ ] Shell plugin support
- [ ] Initial config templates
- [ ] Tests for config and dotfiles

### 📚 Phase 3: Enhanced Features
- [ ] pyenv, nvm, rustup, goenv support
- [ ] Font installer (Nerd Fonts)
- [ ] Plugin system scaffold
- [ ] Advanced UI (forms, feedback)
- [ ] Notification + logs

### 🚀 Phase 4: Polish & Optimization
- [ ] Parallel installs
- [ ] Caching, lazy loading
- [ ] Error recovery and logging
- [ ] End-to-end tests + snapshots
- [ ] Finalize docs, help commands

## 🔍 Testing Strategy
- `go test ./...`
- Unit tests per internal package
- Integration tests in LXC container
- Fixture-based test data

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

## 🧩 Future Enhancements
- Notification history
- Voice or TUI layer
- Automatic rollback
- i18n + accessibility

## 📌 Notes
- Keep each module testable in isolation
- Use Go interfaces for mocking CLI behavior
- Prefer minimal deps: `cobra`, `survey`, `viper`

Refer to `docs/PROJECT_STRUCTURE.md` for folder layout.