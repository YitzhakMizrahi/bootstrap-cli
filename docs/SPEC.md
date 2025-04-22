# 📘 SPEC.md – Bootstrap CLI Specification

## 🎯 Purpose
Bootstrap CLI is a comprehensive tool for automating the installation and configuration of development tools, environments, and shell preferences across various platforms.

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