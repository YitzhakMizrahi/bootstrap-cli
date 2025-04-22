.PHONY: build test clean release

# Build the application
build:
	go build -o build/bin/bootstrap-cli

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf build/

# Run linter
lint:
	golangci-lint run ./...

# Run the application
run: build
	./build/bin/bootstrap-cli

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build for LXC testing
build-lxc: build
	GOOS=linux GOARCH=amd64 go build -o build/bin/bootstrap-cli-linux-amd64 cmd/bootstrap/main.go

# Default target
all: lint test build

# Create a release
release:
	./scripts/release/create-release.sh 