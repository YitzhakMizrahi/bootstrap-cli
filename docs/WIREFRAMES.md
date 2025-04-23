# 🗄 WIREFRAMES.md – Bootstrap CLI Screens

This document outlines the core screens and interactions of the guided CLI setup. This is meant to support the `init` and `up` flows.

Note: MVP uses `promptui` or `survey` for all interactive elements. `bubbletea` UI planned for future (v2).

---

## 🌀 Flow Overview
```
Welcome
  ↳ System Detection (auto-detect)
     ↳ Dotfiles GitHub Clone (optional)
     ↳ Shell Selection
     ↳ Font Installer (Optional)
     ↳ Tool Selection (Multi-Select)
     ↳ Language Runtimes (nvm, pyenv...)
     ↳ Install Progress
     ↳ Validation & Finish
```

---

## 📢 CLI UI Components (MVP)
- Checkbox-style Multi-Selects (Tools, Shells)
- Input (GitHub dotfiles URL, optional)
- Dropdown/Select (Shell)
- Confirmation (yes/no)
- Spinners during install phases
- Progress bars (optional MVP)

---

## 🌐 Screens with Visual Wireframes

### Welcome Screen
```
+------------------------------------------+
|   ✨ Bootstrap CLI ✨                     |
|   Setup your dev machine with ease       |
|                                          |
|   [Y] Start  [N] Exit                    |
+------------------------------------------+
```

### System Detection (Auto)
```
+------------------------------------------+
|  System Info Detected:                   |
|  OS: Ubuntu 22.04                        |
|  Arch: x86_64                            |
|  RAM: 8 GB   | Disk: 256 GB              |
|  Package Manager: apt                    |
|                                          |
|  [Enter] Continue                        |
+------------------------------------------+
```

### Dotfiles (GitHub Clone - Optional)
```
+------------------------------------------+
|  Clone dotfiles from GitHub?             |
|  ( ) Yes                                 |
|  ( ) No                                  |
|                                          |
|  If Yes, enter repo URL:                 |
|  > https://github.com/yourname/dotfiles |
|                                          |
|  [Enter] Confirm                         |
+------------------------------------------+
```

### Shell Setup
```
+------------------------------------------+
|  Choose your shell:                      |
|  ( ) zsh                                 |
|  ( ) bash                                |
|  ( ) fish                                |
|                                          |
|  [Enter] Confirm                         |
+------------------------------------------+
```

### Font Installer (Optional)
```
+------------------------------------------+
|  Install JetBrains Mono Nerd Font?       |
|  ( ) Yes                                 |
|  ( ) No                                  |
|                                          |
|  [Enter] Confirm                         |
+------------------------------------------+
```

### Tool Selection
```
+------------------------------------------+
|  Select tools to install:                |
|  [x] git     [x] curl     [x] bat        |
|  [ ] lsd     [x] fzf      [ ] zoxide     |
|                                          |
|  [Space] Toggle  [Enter] Confirm         |
+------------------------------------------+
```

### Language Runtimes
```
+------------------------------------------+
|  Install runtimes:                       |
|  Node.js (nvm):      [lts]               |
|  Python (pyenv):     [3.11]              |
|  Go (goenv):         [latest]            |
|  Rust (rustup):      [stable]            |
|                                          |
|  [Enter] Continue                        |
+------------------------------------------+
```

### Install Progress
```
+------------------------------------------+
|  Installing tools...                     |
|  - git installed ✓                    |
|  - curl installed ✓                   |
|  - bat installed ✓                    |
|  - fzf installed …                     |
|                                          |
|  [spinner] Working...                    |
+------------------------------------------+
```

### Final Validation
```
+------------------------------------------+
|  Validation Results:                     |
|  - Shell setup: OK                       |
|  - Tools installed: OK                   |
|  - Language runtimes: OK                 |
|  - Paths and symlinks: Configured        |
|                                          |
|  ✅ All systems go!                     |
|  [Enter] Finish                          |
+------------------------------------------+
```

---

## 🧠 Notes
- MVP CLI UI powered by `promptui` or `survey`
- `bubbletea`, `lipgloss`, and full TUI flow saved for v2
- Screens should be composable like prompt modules (functions with props)
- Final setup ensures required binaries are symlinked and paths configured so installed tools work out of the box

