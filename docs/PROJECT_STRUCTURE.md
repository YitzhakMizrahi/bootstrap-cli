# 🧮 PROJECT_STRUCTURE.md – Bootstrap CLI

## 🧱 High-Level Directory Layout
```
bootstrap-cli/
├── cmd/                # CLI commands via Cobra
├── internal/           # Main logic (system, install, shell...)
├── pkg/                # Optional public packages
├── test/               # Integration + e2e tests
├── docs/               # Specifications + guides
├── scripts/            # Helper scripts (build, test, release)
├── .github/            # CI/CD config, issue templates
├── main.go             # Entrypoint
└── Makefile            # Build shortcuts
```

## 📂 internal/
Organized by domain (loosely DDD-inspired):
- `system/` – OS/arch/distro detection
- `packages/` – Package manager detection + abstraction
- `install/` – Tool installation logic
- `shell/` – Shell selection, .rc file writing
- `dotfiles/` – GitHub clone, local sync, restore
- `config/` – YAML config load/save/validate
- `ui/` – CLI components (spinners, forms, prompts)
- `utils/` – Logger, paths, input validation

## 📂 pkg/
For public APIs or pluggable modules:
- `plugin/` – Optional plugin loader, future extension
- `config/` – Shell/editor/tool config templates
- `i18n/` – Language packs

## 🧪 test/
- `integration/` – Real install test (via LXC)
- `fixtures/` – Static test data
- `e2e/` – Simulated user flow test

## 🔧 CI/Linting
- `.golangci.yml` – Includes: gofmt, golint, govet, errcheck
- `.github/workflows/` – `test.yml`, `release.yml`

## 📜 Developer Workflow
```bash
git clone https://github.com/YitzhakMizrahi/bootstrap-cli.git
cd bootstrap-cli
go mod tidy
make build
make test
```

## 📦 Makefile Targets
```makefile
build:
	go build -o build/bin/bootstrap-cli cmd/init.go

test:
	go test ./...

clean:
	rm -rf build/

release:
	./scripts/release/create-release.sh
```

## 🧪 LXC Testing Tips
```bash
lxc launch ubuntu:22.04 bootstrap-test
lxc snapshot bootstrap-test clean
lxc file push build/bin/bootstrap-cli bootstrap-test/usr/local/bin/
lxc exec bootstrap-test -- bootstrap-cli init
```