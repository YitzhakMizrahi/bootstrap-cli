# 📓 DECISIONS.md – Bootstrap CLI Architectural Decisions

This doc captures high-level technical and strategic decisions made during the project.

---

### ✅ 1. Language: Go
- Chosen for portability, static binaries, speed
- Easy CLI development with `cobra`, `survey`, `promptui`

### ✅ 2. Structure: DDD-lite + Modular
- `internal/` is domain-based: system, install, ui, config, flow
- `pkg/` reserved for public modules and templates
- Clear separation of CLI (`cmd/`) and logic (`internal/`)

### ✅ 3. Configuration: YAML over JSON
- Easier for dev users to edit
- Used with `viper` for flexibility

### ✅ 4. Prompts: `promptui` / `survey` for MVP, `Bubbletea` in v2
- `promptui`/`survey` used for basic multi-selects, validations, and config
- `Bubbletea`, `Bubbles`, and `Lip Gloss` reserved for advanced UI in v2

### ✅ 5. Shells: Zsh (Default), Bash, Fish
- Prioritize Zsh due to ecosystem and plugin support
- Provide fallback logic if not found

### ✅ 6. Fonts: Nerd Fonts (JetBrains Mono)
- JetBrains Mono Nerd + Fira Code Nerd preselected
- User can override in config

### ✅ 7. Dotfiles: GitHub Clone Only (MVP)
- MVP only supports cloning dotfiles from GitHub
- Backup, sync, and restore are deferred to post-MVP

### ✅ 8. CLI Commands
- `bootstrap-cli up` is the main entrypoint for full flow
- All sub-commands (`detect`, `shell`, `install`, etc.) are modular

### ✅ 9. Test Environments: LXC
- Ubuntu container for fast, reproducible install tests
- Snapshot for rollback testing

### ✅ 10. Error Strategy
- Friendly CLI error messages
- Internal error logging for debug
- Recovery suggestion on crash

---

### 🟡 Pending Decisions
- Remote config sync method (S3? Git?)
- Plugin architecture (WASM? Go interfaces?)
- TUI wrapper style (Bubbletea for advanced flow?)
- Dotfile sync strategy beyond GitHub clone

---

### 🧩 Open Topics
- Optional telemetry/metrics?
- Visual theme customization?
- Accessibility defaults (reduced motion, contrast)?

