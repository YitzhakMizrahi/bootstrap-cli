# 📓 DECISIONS.md – Bootstrap CLI Architectural Decisions

This doc captures high-level technical and strategic decisions made during the project.

---

### ✅ 1. Language: Go
- Chosen for portability, static binaries, speed
- Easy CLI development with `cobra`, `promptui`, and `bubbletea`

### ✅ 2. Structure: DDD-lite + Modular
- `internal/` is domain-based: system, install, shell, flow, etc.
- `pkg/` reserved for public modules (templates, optional plugins)
- Clear separation of CLI (`cmd/`) and logic (`internal/`)
- Reusable shared interfaces live in `internal/interfaces/`

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

### ✅ 9. Interface Definitions
- Shared interfaces (e.g., `ToolInstaller`, `ShellManager`, `SymlinkApplier`) are defined once in `internal/interfaces/`
- Prevents duplication and improves code reuse across packages

### ✅ 10. Testing Strategy
- **Unit tests** live alongside logic files in the form of `*_test.go`
- **Integration tests** live under `test/integration/` and test cross-module or CLI-level flows
- **End-to-end tests** live under `test/e2e/` and simulate full user journeys via `init`, `up`, etc.
- **Fixtures** (e.g., example YAML config files) live in `test/fixtures/`
- **Mocks and fakes** used across modules live in `internal/testutil/` for shared use across packages

### ✅ 11. Declarative Extensibility
- Bootstrap CLI is designed to be **configuration-driven**
- Adding tools, language runtimes, fonts, or prompts should be possible by editing config files or templates — not by writing new Go code
- All installable items should be registered through centralized metadata or schema (e.g., a `tools.yaml` or Go struct registry)
- Encourages scalability and lowers barrier to contribution or user customization

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

