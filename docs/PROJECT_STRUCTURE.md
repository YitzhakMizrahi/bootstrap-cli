# 🧮 PROJECT_STRUCTURE.md – Bootstrap CLI

## 🧱 High-Level Directory Layout
```
bootstrap-cli/
├── cmd/                # CLI commands via Cobra
├── internal/           # Main logic (system, install, shell, flow...)
│   ├── system/         # OS/arch/distro detection
│   ├── packages/       # Package manager detection + abstraction
│   ├── install/        # Tool installation logic
│   ├── shell/          # Shell selection, .rc file writing
│   ├── dotfiles/       # GitHub clone only (MVP scope)
│   ├── config/         # YAML config load/save/validate
│   ├── flow/           # Guided CLI flows (init, install, shell, etc.)
│   ├── ui/             # Prompt modules, spinners, selections
│   ├── symlinks/       # Shared symlink + PATH config logic
│   ├── utils/          # Logger, paths, validations
│   ├── interfaces/     # Shared interfaces (ToolInstaller, ShellManager...)
│   └── testutil/       # Reusable mocks, stubs, and helpers for unit tests
├── pkg/                # Optional public packages and templates
│   ├── templates/      # Static config templates (e.g. .zshrc, starship.toml)
│   ├── plugin/         # Optional plugin loader (post-MVP)
│   └── i18n/           # Language packs (future)
├── test/               # Integration + e2e tests
│   ├── integration/    # Real install test (via LXC)
│   ├── fixtures/       # Static test data (YAML configs, test plans)
│   └── e2e/            # Simulated full user flow test (init → up)
├── docs/               # Specifications + guides
├── scripts/            # Helper scripts (build, test, release)
├── .github/            # CI/CD config, issue templates
├── main.go             # Entrypoint
└── Makefile            # Build shortcuts
```

## 🔧 CI/Linting
- `.golangci.yml` – Includes: gofmt, golint, govet, errcheck
- `.github/workflows/` – `test.yml`, `release.yml`

## 📜 Developer Workflow
```bash
git clone https://github.com/YitzhakMizrahi/bootstrap-cli.git
cd bootstrap-cli
make deps
make build
make test
```

To run manually (instead of Makefile):
```bash
go build -o build/bin/bootstrap-cli main.go
```

## 📆 Makefile Targets
```makefile
build:
	go build -o build/bin/bootstrap-cli main.go

test:
	go test ./...

clean:
	rm -rf build/

lint:
	golangci-lint run ./...

run: build
	./build/bin/bootstrap-cli

deps:
	go mod download
	go mod tidy

build-lxc:
	GOOS=linux GOARCH=amd64 go build -o build/bin/bootstrap-cli-linux-amd64 main.go

validate:
	./build/bin/bootstrap-cli validate

all: lint test build

release:
	./scripts/release/create-release.sh
```

## 🤮 LXC Testing Tips
```bash
# Launch fresh Ubuntu LXC container
lxc launch ubuntu:22.04 bootstrap-test

# Create a clean snapshot to restore from if needed
lxc snapshot bootstrap-test clean-setup

# Push the compiled binary with executable permissions
lxc file push build/bin/bootstrap-cli bootstrap-test/home/devuser/bootstrap-cli --mode=755

# Run the CLI interactively from inside the container
lxc exec bootstrap-test -- su - devuser
```

---

## 🔑 Cobra Command Structure
```
cmd/
├── root.go         # Root command
├── up.go           # Run entire flow (new)
├── init.go         # Interactive setup
├── detect.go       # System info
├── install.go      # Tool install
├── dotfiles.go     # GitHub clone only
├── shell.go        # Shell setup
├── languages.go    # Runtime installs
├── font.go         # Font installer
├── validate.go     # Validate setup
├── config.go       # View/export config
├── version.go      # Print version
```

