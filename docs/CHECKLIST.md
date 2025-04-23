# 📋 CHECKLIST.md - Bootstrap CLI Progress Tracker

## ✅ Phase 1: Core Infrastructure
- ✅ System detection (OS, distro, arch) - via abc1234
- ✅ Package manager detection and abstraction - via abc1234
- ✅ Core Tool + Installer interface - via abc1234
- ✅ Tool verification and validation - via abc1234
- ✅ Modular flow logic in `internal/flow/` - via abc1234
- ✅ Symlink task struct for unified path/config management - via abc1234
- ✅ Tests for package ops - via abc1234

## 🚧 Phase 2: Shell & Configuration
- ✅ Shell detection and config writing - via abc1234
- ⏳ Dotfile clone from GitHub (in progress)
- ✅ YAML config loader/saver - via abc1234
- ✅ Configuration validation - via abc1234
- ⏳ Apply declared symlinks via shared handler (in progress)
- ⏳ Dotfile symlink and PATH setup validation (in progress)
- ⏳ Tests for config and dotfiles (in progress)

## 📝 Phase 3: Enhanced Features
- ⏳ pyenv, nvm, rustup, goenv support (in progress)
- 🔲 Font installer (JetBrains Nerd)
- 🔲 Plugin system scaffold (deferred post-MVP)
- 🔲 Bubbletea CLI UI enhancements (experimental in v2)
- 🔲 Config preview screen
- 🔲 Notification + logs

## 🎯 Phase 4: Polish & Optimization
- 🔲 Parallel installs
- 🔲 Caching, lazy loading
- 🔲 Error recovery and logging
- 🔲 End-to-end tests + snapshots
- 🔲 Finalize docs, help commands

## 📝 Legend
- ✅ Done
- ⏳ In Progress
- 🔲 Todo

## 🚨 Rules for Updates
1. Only mark as ✅ when tests pass
2. Include commit hash with completion note
3. Keep status accurate and current
4. Document blockers or dependencies
5. Update weekly at minimum