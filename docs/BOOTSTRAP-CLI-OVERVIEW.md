# 🧭 Project Overview: bootstrap-cli

## 🔥 Vision
`bootstrap-cli` is a modular, interactive command-line tool designed to **bootstrap a personalized development environment** quickly and consistently across platforms. Inspired by Powerlevel10k-style onboarding, it allows users to:

- Select shells, plugin managers, prompts, editors
- Choose CLI tools and programming languages
- Symlink dotfiles intelligently
- Install language environments via `pyenv`, `goenv`, `nvm`, etc.
- Operate seamlessly in both native systems and sandboxed test containers

It’s built for **full-stack developers**, **dotfiles enthusiasts**, and **DevOps tinkerers** who want reproducible, elegant, minimal setups.

---

## 🧱 Technologies Used

| Layer                  | Tool / Framework               |
|------------------------|--------------------------------|
| Language               | Go                             |
| CLI Framework          | Cobra                          |
| Interactive Prompts    | bubbletea                      |
| Config Serialization   | YAML via gopkg.in/yaml.v3      |
| Testing Environment    | LXD + Ubuntu containers        |
| Plugin Management      | Zinit (for zsh)                |
| Dotfiles               | Symlinks via Go logic          |

---

## 🗂 Project Folder Structure

```
bootstrap-cli/
├── cmd/               # Cobra commands: init, install, link
│   ├── init.go
│   ├── install.go
│   ├── link.go
│   └── up.go          # NEW: orchestrates init → install → link
│
├── config/            # YAML config handling (save/load)
│   └── config.go
│
├── installer/         # CLI, languages, editor installers
│   └── cli.go
│
├── prompts/           # Interactive prompt logic
│   └── init_prompts.go
│
├── symlink/           # Dotfiles symlinking logic
│   └── linker.go
│
├── types/             # Shared types (UserConfig, ToolOption)
│   └── types.go
│
└── main.go            # CLI entry point
```

---

## 🧪 Testing Architecture

### Goal:
Validate the CLI tool in **isolated container environments** without affecting the host.

### Approach:
- Use **LXD** containers with Ubuntu 22.04 (`bootstrap-test`)
- Push compiled binary from WSL → container via:
  ```bash
  lxc file push ./bootstrap-cli bootstrap-test/home/devuser/bootstrap-cli --mode=755
  ```
- Inside container:
  ```bash
  ./bootstrap-cli init
  ./bootstrap-cli install
  ./bootstrap-cli link
  ```
- Snapshots created for state rollback:
  ```bash
  lxc snapshot bootstrap-test clean-setup
  lxc restore bootstrap-test clean-setup
  ```

### Networking Fixes Used:
- Manual NAT masquerade:
  ```bash
  sudo iptables -t nat -A POSTROUTING -s 10.85.50.0/24 ! -d 10.85.50.0/24 -j MASQUERADE
  ```
- UFW Rules:
  ```bash
  sudo ufw allow in on lxdbr0
  sudo ufw allow out on lxdbr0
  sudo ufw route allow in on lxdbr0 out on eth0
  sudo ufw route allow in on eth0 out on lxdbr0
  ```

---

## ✅ Key Features So Far

- Interactive setup wizard (`bootstrap-cli init`)
- Config persistence in `~/.config/bootstrap/config.yaml`
- CLI installs via Homebrew: lsd, bat, fzf, lazygit, etc.
- Language bootstrapping: Python (pyenv), Node (nvm), Go, Rust
- Editor setup with options like LazyVim / AstroNvim
- Dotfile symlinking with optional backup and relative/absolute modes
- Dev mode and backup flags

---

## 🧩 Planned Enhancements

- 🔒 Lock select/multiselect prompts to prevent free typing
- 🧹 Better input validation and default fallback handling
- 🙏 Handle Ctrl+C (SIGINT) gracefully in the wizard
- 💄 Prompt polish: icons, links, descriptions like tauri create
- ⚙️ Settings auto-detection or minimal .bootstraprc for headless mode
- 💬 Dynamic prompt filtering based on OS/platform
- 🧠 Intelligent fallback when dotfiles not found
- 📦 Package manager installs beyond Homebrew (apt, dnf, etc.)
- 🧪 CI testing with snapshot restoration

---

## 🧑‍💻 Dev Workflow

- Code on host (WSL2, Go v1.24.2)
- Build binary locally:
  ```bash
  go build -o bootstrap-cli
  ```
- Push to container and test in isolation
- Use LXD snapshots to save/restore test states

---

## 💬 Philosophy

Keep it:
- Minimal
- Modular
- Hackable
- Reproducible

CLI should feel like a guided wizard — smooth, fast, low-friction.

---
