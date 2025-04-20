.PHONY: build install test clean lxd-push lxd-test help

BINARY_NAME=bootstrap-cli
GOPATH=$(shell go env GOPATH)
GO_FILES=$(shell find . -name "*.go" -type f)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags="-X 'github.com/YitzhakMizrahi/bootstrap-cli/cmd.Version=$(VERSION)'"

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY_NAME) .

install: build ## Install the binary to GOPATH/bin
	cp $(BINARY_NAME) $(GOPATH)/bin/

test: ## Run the tests
	go test -v ./...

clean: ## Clean build artifacts
	go clean
	rm -f $(BINARY_NAME)

lxd-push: build ## Push the binary to the LXD test container
	lxc file push ./$(BINARY_NAME) bootstrap-test/home/ubuntu/$(BINARY_NAME) --mode=755

lxd-test: lxd-push ## Run the binary in the LXD test container
	lxc exec bootstrap-test -- sudo -u ubuntu -i ./$(BINARY_NAME) up

lxd-create: ## Create a test container
	lxc launch ubuntu:22.04 bootstrap-test
	lxc exec bootstrap-test -- apt-get update
	lxc exec bootstrap-test -- apt-get install -y git curl

lxd-snap: ## Create a snapshot of the test container
	lxc snapshot bootstrap-test clean-setup

lxd-restore: ## Restore the test container to the clean-setup snapshot
	lxc restore bootstrap-test clean-setup