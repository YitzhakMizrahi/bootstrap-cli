# Bootstrap CLI Development Issues Tracker

## Current Status

### 1. 📦 Fix embedded-defaults extraction
- [x] Implement `ExtractEmbeddedConfigs()` in `internal/config/loader.go`
  - [x] Write files from `defaults/...` to `<baseDir>/defaults/...`
  - [x] Handle all subdirectories properly
  - [x] Add proper error handling
- [x] Expand `//go:embed` to include all defaults
- [x] Remove double-call to `ExtractEmbeddedConfigs()`
  - [x] Keep only the call in `main.go`
  - [x] Remove redundant calls

### 2. 🔀 Harden the merge logic
- [x] Add runtime checks on type assertions
- [x] Extend switch statement to cover all subfolders
- [x] Improve error handling for type assertions
- [x] Add support for `language_managers` directory

### 3. 🛠 Unify installation strategy
- [x] Decide between `interfaces.ToolInstaller` and `pipeline.Installer`
  - [x] Choose `pipeline.Installer` for its advanced features
  - [x] Document decision in code comments
- [x] Update code based on chosen strategy
  - [x] Remove `interfaces.ToolInstaller` interface
  - [x] Update documentation to reflect pipeline usage
  - [x] Clean up project structure docs
- [x] Remove dead code
  - [x] Remove unused `ToolInstaller` implementations
  - [x] Clean up any remaining references
- [ ] Implement retry/timeouts and rollback support
  - [ ] Add retry configuration to pipeline steps
  - [ ] Implement rollback for failed installations
  - [ ] Add timeout configuration for long-running operations

### 4. 🎨 Drive the UI from configs
- [ ] Replace hard-coded lists with loader calls
  - [ ] Update `cmd/init`
  - [ ] Update `cmd/up`
- [ ] Update UI prompts to use configuration data
  - [ ] Fonts
  - [ ] Languages
  - [ ] Language managers
  - [ ] Dotfiles

### 5. 🗂 Implement dotfiles MVP
- [ ] Create `DotfilesManager` implementation
  - [ ] Clone user's repo URL
  - [ ] Handle file operations (backup, symlink, copy)
  - [ ] Implement post-install hooks
  - [ ] Handle restart requirements
- [ ] Wire into installation phases
  - [ ] Add to `cmd/up`
  - [ ] Add to `cmd/init`

### 6. 🔄 Rollback & retry
- [ ] Implement rollback functionality
  - [ ] Add rollback to installation steps
  - [ ] Test rollback scenarios
- [ ] Add retry mechanism
  - [ ] Implement for network-dependent steps
  - [ ] Add timeout handling
  - [ ] Test retry scenarios

### 7. ✅ Bolster test coverage
- [ ] Add unit tests
  - [ ] Config merging
  - [ ] Loader fallbacks
  - [ ] Bad YAML handling
  - [ ] `ExtractEmbeddedConfigs`
- [ ] Add pipeline tests
  - [ ] Step execution
  - [ ] Rollback behavior
  - [ ] Timeout handling
- [ ] Add integration tests
  - [ ] Full "init → up" flow
  - [ ] LXC container testing
  - [ ] Idempotency verification

## Progress Tracking
- Total Issues: 7
- Completed Issues: 2 (Issue #1 - Fix embedded-defaults extraction, Issue #2 - Harden the merge logic)
- In Progress: 1 (Issue #3 - Unify installation strategy)
- Remaining Issues: 4

## Notes
- Last Updated: [Current Date]
- Next Focus: Complete installation strategy unification (Issue #3)
- Priority: High - Fix core functionality before moving to features 