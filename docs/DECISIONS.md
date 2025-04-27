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

### ✅ 12. Pipeline-Based Installation Architecture
- Adopted a pipeline-based approach for tool installation and configuration
- Each tool installation is broken down into discrete, verifiable steps
- Key components:
  - `Tool` struct with platform-specific installation strategies
  - `InstallationContext` for managing state and environment
  - `VerifyStrategy` for robust post-install verification
  - `InstallStrategy` for flexible, platform-aware installation steps
- Benefits:
  - Better error handling and recovery
  - Platform-specific installation paths
  - Consistent verification across tools
  - Modular and testable installation steps
  - Clear separation between tool definition and installation logic

### ✅ 13. Shell Configuration Management
- Shell configuration is treated as a first-class concern
- Structured approach using `ShellConfig` type:
  - Aliases
  - Functions
  - Environment variables
- Enables consistent shell setup across different shells and platforms

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

## 1. Split Commands (2024-04-25)

### Context
The original design had a single command that handled both initialization and installation. This created confusion about the state of the system and made it difficult to handle first-time setup.

### Decision
Split the functionality into two distinct commands:
- `init`: Handles first-time setup
- `up`: Runs the interactive TUI

### Consequences
- Clearer separation of concerns
- Better user experience with explicit initialization
- More maintainable codebase
- Slight increase in complexity due to state management

## 2. Bubble Tea TUI (2024-04-25)

### Context
The original CLI used simple prompts and lacked visual feedback. Users had difficulty understanding the installation progress and their choices.

### Decision
Adopt the Bubble Tea framework for a full-screen TUI with:
- Interactive selection screens
- Progress indicators
- Consistent styling
- Keyboard navigation

### Consequences
- More intuitive user interface
- Better visual feedback
- Increased development complexity
- Need for careful state management

## 3. Configuration-Driven Design (2024-04-25)

### Context
Tool and font configurations were initially hardcoded, making it difficult to add or modify options.

### Decision
Move all configurations to YAML files with:
- Schema validation
- Default configurations
- User override support
- Embedded defaults

### Consequences
- Easier to add new tools and fonts
- Better maintainability
- Need for careful schema design
- Additional complexity in configuration loading

## 4. Component-Based UI (2024-04-25)

### Context
UI code was scattered and inconsistent, making it hard to maintain a consistent look and feel.

### Decision
Create reusable UI components:
- Base selector for lists
- Progress indicators
- Consistent styling system
- Screen-based navigation

### Consequences
- More consistent UI
- Better code reuse
- Easier to add new features
- Need for careful component design

## 5. Installation Pipeline (2024-04-25)

### Context
Installation process was linear and lacked proper error handling and rollback capabilities.

### Decision
Implement a pipeline-based installation system:
- Sequential or parallel installation
- Error handling at each step
- Rollback capabilities
- Progress tracking

### Consequences
- More robust installation process
- Better error recovery
- Increased complexity
- Need for careful testing

## Future Considerations

1. **Package Manager Abstraction**
   - Consider more flexible package manager support
   - Add support for more package managers
   - Improve version handling

2. **Configuration Management**
   - Consider remote configuration support
   - Add configuration validation
   - Improve override mechanics

3. **Testing Strategy**
   - Add more unit tests
   - Implement integration tests
   - Add UI testing

4. **Error Handling**
   - Improve error messages
   - Add detailed logging
   - Implement better recovery strategies

