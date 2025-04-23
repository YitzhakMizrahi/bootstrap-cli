.PHONY: build test clean lint run deps build-lxc all release validate deploy-lxc

# Build the application
build:
	go build -o build/bin/bootstrap-cli main.go

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
build-lxc:
	GOOS=linux GOARCH=amd64 go build -o build/bin/bootstrap-cli-linux-amd64 main.go

# Validate setup
validate:
	./build/bin/bootstrap-cli validate

# Default target
all: lint test build

# Create a release
release:
	./scripts/release/create-release.sh
