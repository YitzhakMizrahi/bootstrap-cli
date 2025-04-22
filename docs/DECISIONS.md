# 📓 DECISIONS.md – Bootstrap CLI Architectural Decisions

This doc captures high-level technical and strategic decisions made during the project.

---

### ✅ 1. Language: Go
- Chosen for portability, static binaries, speed
- Easy CLI development with `cobra`, `survey`

### ✅ 2. Structure: DDD-lite + Modular
- `internal/` is domain-based: system, install, ui, config
- `pkg/` reserved for public modules (plugins, templates)
- Clear separation of CLI (`cmd/`) and logic (`internal/`)

### ✅ 3. Configuration: YAML over JSON
- Easier for dev users to edit
- Used with `viper` for flexibility

### ✅ 4. Prompts: `Bubbletea + Bubbles + Lip Gloss`
- Excellent for interactive CLI UX
- Supports all major input types, styling

### ✅ 5. Shells: Zsh (Default), Bash, Fish
- Prioritize Zsh due to ecosystem and plugin support
- Provide fallback logic if not found

### ✅ 6. Fonts: Nerd Fonts (JetBrains Mono)
- JetBrains Mono Nerd + Fira Code Nerd preselected
- User can override in config

### ✅ 7. Test Environments: LXC
- Ubuntu container for fast, reproducible install tests
- Snapshot for rollback testing

### ✅ 8. Error Strategy
- Friendly CLI error messages
- Internal error logging for debug
- Recovery suggestion on crash

---

### 🟡 Pending Decisions
- Remote config sync method (S3? Git?)
- Plugin architecture (WASM? Go interfaces?)
- Dotfiles symlink manager vs full copy
- TUI wrapper (Bubbletea? Custom rendering?)

---

### 🧩 Open Topics
- Optional telemetry/metrics?
- Visual theme customization?
- Accessibility defaults (reduced motion, contrast)?

