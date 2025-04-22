# 📘 SPEC.md – Bootstrap CLI Specification

## 🎯 Project Vision
Bootstrap CLI aims to be the standard tool for setting up development environments across Linux, macOS, and WSL. It provides an interactive, guided experience for installing and configuring development tools, making it easy for developers to replicate their preferred environment across different machines.

## 🎯 Core Objectives
1. **Guided Setup Experience**
   - Interactive CLI interface for selecting tools and configurations
   - Clear categorization of tools (essential, modern CLI, system tools)
   - Step-by-step wizard with progress tracking
   - Validation at each step

2. **Cross-Platform Compatibility**
   - Support for Ubuntu/Debian, Fedora, Arch Linux
   - macOS support via Homebrew
   - WSL2 compatibility
   - Consistent experience across platforms

3. **Modular Architecture**
   - Package manager abstraction
   - Pluggable tool definitions
   - Template-based configurations
   - Extensible shell support

4. **Reproducible Environments**
   - Config export/import
   - Version-controlled dotfiles
   - Validation and verification
   - Idempotent operations

## ❌ Non-Goals
- GUI application development
- System-wide configuration management
- Container orchestration
- Remote system management

## 🎯 Success Criteria
1. **Usability**
   - Complete setup in under 10 minutes
   - No manual intervention needed
   - Clear error messages and recovery
   - Comprehensive help documentation

2. **Reliability**
   - 100% success rate for supported platforms
   - Validation for all installed tools
   - No system corruption possible
   - Clean rollback on failure

3. **Performance**
   - Parallel installation where possible
   - Caching of downloaded assets
   - Minimal memory footprint
   - Quick startup time

## 💡 Core Features

### 1. System Detection
- OS type, distro, architecture, kernel, RAM, disk
- Package manager detection (`apt`, `dnf`, `pacman`, `brew`, etc.)

### 2. Core Tool Installation
- Default tool detection (git, curl, etc.)
- Modern CLI tools (bat, fzf, lsd, ripgrep, etc.)

### 3. Shell & Terminal Setup
- Shell selection (zsh, bash, fish)
- Prompt configuration (Starship)
- Font install (JetBrains Mono, Nerd Fonts)

### 4. Language Managers
- Node (nvm), Python (pyenv), Go (goenv), Rust (rustup)

### 5. Dotfiles Management
- GitHub clone, or fresh start
- Sync, backup, override modes

### 6. Configuration Management
- Templates (minimal, dev, sysadmin, data sci)
- Syntax validation and testing
- Custom overrides and export

### 7. UI/UX Components
- Interactive CLI with:
  - spinners
  - progress bars
  - step indicators
  - validation prompts

## 🔁 User Journey
1. Launch → detect system
2. Choose tools → install
3. Select shell + setup config
4. Manage dotfiles + fonts
5. Install languages
6. Validate config + complete

## 📊 Wireframes Summary
- Welcome Screen
- Package Manager Select
- Core Tools
- Shell + Fonts
- Language Versions
- Dotfiles
- Config Templates
- Install Progress
- Validation & Finish

## 🧪 Validation Goals
- Config syntax check
- Tool availability after install
- Dotfile path + symlink check

## 🧭 Templates
- Minimal (bare essentials)
- Developer (git, nvm, zsh, tools)
- System Admin (network tools, tmux, btop)
- Data Scientist (Python, Jupyter, etc.)

## 🧩 Future Ideas
- Plugin system
- Remote sync and backup
- Notification/logging layer
- GUI launcher wrapper

_See `docs/WIREFRAMES.md` for screen-by-screen layout._