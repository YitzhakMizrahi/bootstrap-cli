## User Journey Wireframes

### Color Coding Legend
```
┌─────────────────────────────────────────┐
│           Color Coding Legend           │
│                                         │
│  🔵 Blue: Primary actions               │
│  🟢 Green: Success/Completion           │
│  🔴 Red: Errors/Warnings                │
│  🟡 Yellow: Information/Notes           │
│  ⚪ White: Background                   │
│  ⚫ Black: Text                         │
└─────────────────────────────────────────┘
```

### 1. Initial Launch & System Detection
```
┌─────────────────────────────────────────┐
│           Bootstrap CLI                 │
│                                         │
│  ██████╗  ██████╗  ██████╗ ████████╗███████╗████████╗██████╗  █████╗ ██████╗ 
│  ██╔══██╗██╔═══██╗██╔═══██╗╚══██╔══╝██╔════╝╚══██╔══╝██╔══██╗██╔══██╗██╔══██╗
│  ██████╔╝██║   ██║██║   ██║   ██║   ███████╗   ██║   ██████╔╝███████║██████╔╝
│  ██╔══██╗██║   ██║██║   ██║   ██║   ╚════██║   ██║   ██╔══██╗██╔══██║██╔═══╝ 
│  ██████╔╝╚██████╔╝╚██████╔╝   ██║   ███████║   ██║   ██║  ██║██║  ██║██║     
│  ╚═════╝  ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     
│                                         │
│  🔵 Welcome to Bootstrap CLI!           │
│                                         │
│  System Information:                    │
│  🟡 Detected System: Ubuntu 22.04       │
│  🟡 Architecture: x86_64                │
│  🟡 Kernel: 5.15.0-56-generic           │
│  🟡 Memory: 16GB                        │
│  🟡 Disk Space: 256GB free              │
│                                         │
│  🔵 [Start Setup]                       │
│  ⚪ [Skip Detection]                    │
└─────────────────────────────────────────┘
```

### 2. Core Development Tools
```
┌─────────────────────────────────────────┐
│      Core Development Tools             │
│                                         │
│  Essential Tools:                       │
│  ☑ git                                  │
│  ☑ curl                                 │
│  ☑ wget                                 │
│  ☑ build-essential                      │
│                                         │
│  Modern CLI Tools:                      │
│  ☑ bat (better cat)                     │
│  ☑ lsd (better ls)                      │
│  ☑ ripgrep (better grep)                │
│  ☑ fd (better find)                     │
│  ☑ fzf (fuzzy finder)                   │
│                                         │
│  System Tools:                          │
│  ☑ htop/btop                            │
│  ☑ tree                                 │
│  ☑ tldr                                 │
│  ☑ neofetch                             │
│                                         │
│  [Install Selected]                     │
└─────────────────────────────────────────┘
```

### 3. Shell & Terminal Setup
```
┌─────────────────────────────────────────┐
│      Shell & Terminal Setup             │
│                                         │
│  Choose Shell:                          │
│  ○ Zsh (Recommended)                    │
│  ○ Bash                                 │
│  ○ Fish                                 │
│                                         │
│  Terminal Enhancements:                 │
│  ☑ starship (modern prompt)             │
│  ☑ zoxide (smart cd)                    │
│  ☑ direnv (env management)              │
│  ☑ navi (interactive cheatsheets)       │
│                                         │
│  Fonts:                                 │
│  ☑ Nerd Fonts Complete                  │
│  ☑ JetBrains Mono Nerd                  │
│  ☑ Fira Code Nerd                       │
│                                         │
│  [Continue]                             │
└─────────────────────────────────────────┘
```

### 4. Programming Languages
```
┌─────────────────────────────────────────┐
│      Programming Languages              │
│                                         │
│  Choose languages to install:           │
│                                         │
│  Node.js:                               │
│  ○ Latest LTS                           │
│  ○ Latest Current                       │
│  ○ Custom Version                       │
│                                         │
│  Python:                                │
│  ○ System Python                        │
│  ○ pyenv (multiple versions)            │
│                                         │
│  Go:                                    │
│  ○ Latest Stable                        │
│  ○ Custom Version                       │
│                                         │
│  Rust:                                  │
│  ○ Latest Stable                        │
│  ○ Custom Version                       │
│                                         │
│  [Install Selected]                     │
└─────────────────────────────────────────┘
```

### 5. Dotfiles Management
```
┌─────────────────────────────────────────┐
│      Dotfiles Management                │
│                                         │
│  Choose dotfiles setup:                 │
│                                         │
│  ○ Start fresh                          │
│  ○ Clone from GitHub                    │
│    Repository URL: [________________]   │
│                                         │
│  Dotfiles to manage:                    │
│  ☑ .zshrc/.bashrc/.fishrc               │
│  ☑ .gitconfig                           │
│  ☑ .vimrc/.nvim                         │
│  ☑ .tmux.conf                           │
│  ☑ SSH config                           │
│                                         │
│  [Continue]                             │
└─────────────────────────────────────────┘
```

### 6. Tool Configurations
```
┌─────────────────────────────────────────┐
│      Tool Configurations                │
│                                         │
│  Shell Configurations:                  │
│  ☑ .zshrc/.bashrc/.fishrc               │
│  ☑ .zshenv/.bashenv/.fish_env           │
│  ☑ .zprofile/.bash_profile              │
│  ☑ .zlogin/.bash_login                  │
│                                         │
│  Starship Configuration:                │
│  ☑ starship.toml                        │
│    • Prompt customization               │
│    • Module configuration               │
│    • Transient prompt                   │
│                                         │
│  Git Configuration:                     │
│  ☑ .gitconfig                           │
│  ☑ .gitignore_global                    │
│  ☑ .gitattributes                       │
│  ☑ SSH config (~/.ssh/config)           │
│                                         │
│  [Configure]                            │
└─────────────────────────────────────────┘
```

### 7. Configuration Templates
```
┌─────────────────────────────────────────┐
│      Configuration Templates            │
│                                         │
│  Choose template style:                 │
│                                         │
│  Shell Templates:                       │
│  ○ Minimal                              │
│  ○ Developer                            │
│  ○ System Administrator                 │
│  ○ Data Scientist                       │
│                                         │
│  Editor Templates:                      │
│  ○ Basic                                │
│  ○ IDE-like                             │
│  ○ Minimalist                           │
│                                         │
│  [Apply Template]                       │
└─────────────────────────────────────────┘
```

### 8. Installation Progress (Enhanced)
```
┌─────────────────────────────────────────┐
│      Installation Progress              │
│                                         │
│  Core Tools:                            │
│  • git                                  │
│  • curl                                 │ 
│  • wget                                 │
│  • build-essential                      │
│                                         │
│  [███████████████████████] 100%         │
│                                         │
│  Shell & Terminal:                      │
│  • zsh                                  │
│  • starship                             │
│  • zoxide                               │
│                                         │
│  [███████████████████████] 100%         │
│                                         │
│  Programming Languages:                 │
│  • Node.js 18.17.0                      │
│  • Python 3.10.12                       │
│  • Go 1.21.0                            │
│                                         │
│  [███████████████████████] 100%         │
│                                         │
│  Configuration:                         │
│  • .zshrc                               │
│  • starship.toml                        │
│  • .gitconfig                           │
│                                         │
│  [███████████████████████] 100%         │
│                                         │
│  🔵 [Continue]                          │
│  ⚪ [View Log]                          │
└─────────────────────────────────────────┘
```

### 9. Configuration Validation
```
┌─────────────────────────────────────────┐
│      Configuration Validation           │
│                                         │
│  Validating configurations...           │
│                                         │
│  Shell Configs:                         │
│  ✓ .zshrc syntax valid                  │
│  ✓ .zshenv syntax valid                 │
│  ✓ .zprofile syntax valid               │
│                                         │
│  Tool Configs:                          │
│  ✓ starship.toml valid                  │
│  ✓ .gitconfig valid                     │
│  ✓ .editorconfig valid                  │
│                                         │
│  [Fix Issues] [Continue]                │
└─────────────────────────────────────────┘
```

### 10. Completion & Next Steps
```
┌─────────────────────────────────────────┐
│      Setup Complete!                    │
│                                         │
│  Your development environment is        │
│  ready to use!                          │
│                                         │
│  Next steps:                            │
│  1. Restart your terminal               │
│  2. Review your dotfiles                │
│  3. Customize your configuration        │
│  4. Install additional tools as needed  │
│                                         │
│  Documentation:                         │
│  • View setup log                       │
│  • Access configuration guide           │
│  • Troubleshooting guide                │
│                                         │
│  [Finish]                               │
└─────────────────────────────────────────┘
```

### 11. Flowchart with Decision Points

┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Initial Launch │────▶│ System Detection│────▶│ Package Manager │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                         │
                                                         ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Completion     │◀────│  Configuration  │◀────│  Core Tools     │
│  & Next Steps   │     │  Validation     │     │  Installation   │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        ▲                       ▲                        │
        │                       │                        ▼
        │                       │              ┌─────────────────┐
        │                       │              │  Shell &        │
        │                       │              │  Terminal Setup │
        │                       │              └─────────────────┘
        │                       │                        │
        │                       │                        ▼
        │                       │              ┌─────────────────┐
        │                       │              │  Programming    │
        │                       │              │  Languages      │
        │                       │              └─────────────────┘
        │                       │                        │
        │                       │                        ▼
        │                       │              ┌─────────────────┐
        │                       │              │  Dotfiles       │
        │                       │              │  Management     │
        │                       │              └─────────────────┘
        │                       │                        │
        │                       │                        ▼
        │                       │              ┌─────────────────┐
        │                       │              │  Tool           │
        │                       │              │  Configurations │
        │                       │              └─────────────────┘
        │                       │                        │
        │                       │                        ▼
        │                       │              ┌─────────────────┐
        │                       │              │  Configuration  │
        │                       │              │  Templates      │
        │                       │              └─────────────────┘
        │                       │                        │
        │                       │                        ▼
        │                       │              ┌─────────────────┐
        │                       │              │  Installation   │
        │                       │              │  Progress       │
        │                       │              └─────────────────┘
        │                       │                        │
        │                       └────────────────────────┘
        │
        └────────────────────────────────────────────────┘
                                                         