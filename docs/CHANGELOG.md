# Changelog

## [Unreleased]

### Added
- New TUI-based interface using Bubble Tea
- Modular screen system with separate screens for each step
- Base selector component for consistent selection UI
- Progress bar components for installation feedback
- Unified styling system with consistent theme
- Configuration-driven tool, font, and language selection
- Two-phase initialization process: `init` and `up` commands

### Changed
- Split initialization into two commands:
  - `init`: Handles first-time setup (config extraction, env vars)
  - `up`: Runs the interactive TUI for installation
- Moved to a more modular architecture with clear separation of concerns
- Improved error handling and user feedback
- Enhanced configuration loading with default/user config merging

### Removed
- Old CLI-based interface
- Direct package installation without user confirmation
- Hardcoded tool and font configurations

## Architecture Changes

### Command Structure
- `init`: First-time setup and configuration
- `up`: Interactive TUI for installation

### Package Organization
- `cmd/`: Command entrypoints
- `internal/`
  - `config/`: Configuration loading and merging
  - `interfaces/`: Core type definitions
  - `ui/`
    - `app/`: Main application model
    - `components/`: Reusable UI components
    - `screens/`: Individual screen implementations
    - `styles/`: UI styling and theming
    - `utils/`: UI utilities
  - `install/`: Installation logic
  - `packages/`: Package management

### UI Components
- Base selector for consistent selection interfaces
- Progress bars for installation feedback
- Styled components using lipgloss
- Screen-based navigation 