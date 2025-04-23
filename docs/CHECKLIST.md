# ✅ CHECKLIST.md – Bootstrap CLI Verification Tracker

This file tracks implementation milestones. Only update entries **after passing tests** for the relevant section.

---

## Phase 1: Core Infrastructure

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Create docs                      | ✅ Done | created docs/ directory       |
| `GetSystemInfo()`                | ✅ Done | Implemented as Detect() in system.go |
| OS/distro/arch/kernel detection  | ✅ Done | Full implementation with tests |
| Package manager abstraction      | ✅ Done | Interface + apt implementation |
| Tool install struct + interface  | ✅ Done | Core Tool + Installer interface implemented |
| Tool verification logic          | ✅ Done | Verification package complete |
| Core tool install (git, curl...) | ✅ Done | Implemented with tests passing |
| Modular flow structure scaffold  | ✅ Done | internal/flow/ created with init/install stubs |

---

## Phase 2: Shell & Config

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Shell detection (zsh/bash/fish) | ✅ Done | Implemented with tests passing |
| Shell config writer              | ⬛ Todo |                               |
| Dotfiles clone from GitHub       | ⬛ Todo | MVP supports only GitHub cloning |
| YAML config loader/saver         | ⬛ Todo |                               |
| Configuration validation         | ⬛ Todo |                               |
| Template rendering logic         | ⬛ Todo | Minimal/dev/sysadmin variants |
| Dotfile validation and symlink test | ⬛ Todo |                               |

---

## Phase 3: Enhanced Features

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Language installers (nvm/pyenv)  | ⬛ Todo | Tests passing for 4 runtimes  |
| Font installer (JetBrains Nerd)  | ⬛ Todo |                               |
| Plugin system scaffold           | ⬛ Todo | Deferred to post-MVP          |
| TUI: Bubbletea base setup        | ⬛ Todo | Optional/experimental in v2   |
| Config preview screen            | ⬛ Todo |                               |

---

## Phase 4: Polish & Optimization

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Parallel installs                | ⬛ Todo |                               |
| Caching + lazy loading           | ⬛ Todo |                               |
| End-to-end tests via LXC         | ⬛ Todo |                               |
| Finalize docs + help             | ⬛ Todo | README, CLI --help, module doc comments |

---

## CLI Commands

| Command     | Status  | Notes                                |
|-------------|---------|--------------------------------------|
| `up`        | ⬛ Todo | Orchestrates full setup flow         |
| `init`      | ⬛ Todo | Interactive prompt-based setup       |
| `detect`    | ⬛ Todo | Print system info                    |
| `install`   | ⬛ Todo | Install tools from config            |
| `dotfiles`  | ⬛ Todo | GitHub clone only                    |
| `shell`     | ⬛ Todo | Shell install and config             |
| `languages` | ⬛ Todo | Runtime installers                   |
| `font`      | ⬛ Todo | Nerd font install                    |
| `validate`  | ⬛ Todo | Run post-install validation          |
| `config`    | ⬛ Todo | View or export config                |
| `version`   | ⬛ Todo | Print CLI version                    |

---

## Update Rules
- ✅ Only update status if unit/integration tests pass
- ➕ Add commit hash/PR ref in Notes if useful
- ⛔ Do not mark incomplete work as Done
- ✍️ Keep this doc aligned with `IMPLEMENTATION.md`