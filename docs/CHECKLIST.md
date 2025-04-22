# ✅ CHECKLIST.md – Bootstrap CLI Verification Tracker

This file tracks implementation milestones. Only update entries **after passing tests** for the relevant section.

---

## Phase 1: Core Infrastructure

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Create docs                      | ✅ Done | created docs/ directory       |
| `GetSystemInfo()`                | ⬜ Todo | Covered by unit tests         |
| OS/distro/arch/kernel detection  | ⬜ Todo | Verified with mock + real sys |
| Package manager abstraction      | ⬜ Todo | apt/dnf/pacman implemented    |
| Tool install struct + interface  | ⬜ Todo |                               |
| Tool verification logic          | ⬜ Todo |                               |
| Core tool install (git, curl...) | ⬜ Todo |                               |

---

## Phase 2: Shell & Config

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Shell detection (zsh/bash/fish) | ⬜ Todo |                               |
| Shell config writer              | ⬜ Todo |                               |
| Dotfiles clone/import/backup     | ⬜ Todo |                               |
| YAML config loader/saver         | ⬜ Todo |                               |
| Configuration validation         | ⬜ Todo |                               |

---

## Phase 3: Enhanced Features

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Language installers (nvm/pyenv)  | ⬜ Todo | Tests passing for 4 runtimes  |
| Font installer (JetBrains Nerd)  | ⬜ Todo |                               |
| Plugin system scaffold           | ⬜ Todo |                               |
| TUI: Bubbletea base setup        | ⬜ Todo |                               |
| Config preview screen            | ⬜ Todo |                               |

---

## Phase 4: Polish & Optimization

| Feature                          | Status  | Notes                         |
|----------------------------------|---------|-------------------------------|
| Parallel installs                | ⬜ Todo |                               |
| Caching + lazy loading           | ⬜ Todo |                               |
| End-to-end tests via LXC         | ⬜ Todo |                               |
| Finalize docs + help             | ⬜ Todo |                               |

---

## Update Rules
- ✅ Only update status if unit/integration tests pass
- ➕ Add commit hash/PR ref in Notes if useful
- 🛑 Do not mark incomplete work as Done
- ✍️ Keep this doc aligned with `IMPLEMENTATION.md`

---

