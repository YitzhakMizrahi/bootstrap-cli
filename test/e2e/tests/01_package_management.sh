#!/bin/bash
set -e

# Test configuration
CONTAINER="bootstrap-test"
BINARY_PATH="/home/devuser/bootstrap-cli"
TEST_PACKAGE="git"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "Running package management tests..."

# Restore clean snapshot
echo "Restoring clean snapshot..."
lxc restore $CONTAINER clean

# Test package installation
echo "Testing package installation..."
if lxc exec $CONTAINER -- $BINARY_PATH pkg install $TEST_PACKAGE; then
    echo -e "${GREEN}✓ Package installation command succeeded${NC}"
else
    echo -e "${RED}✗ Package installation command failed${NC}"
    exit 1
fi

# Verify package is installed
echo "Verifying package installation..."
if lxc exec $CONTAINER -- which $TEST_PACKAGE >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Package $TEST_PACKAGE is installed${NC}"
else
    echo -e "${RED}✗ Package $TEST_PACKAGE is not installed${NC}"
    exit 1
fi

# Test package is-installed check
echo "Testing package existence check..."
if lxc exec $CONTAINER -- $BINARY_PATH pkg is-installed $TEST_PACKAGE; then
    echo -e "${GREEN}✓ Package existence check succeeded${NC}"
else
    echo -e "${RED}✗ Package existence check failed${NC}"
    exit 1
fi

# Test package uninstallation
echo "Testing package uninstallation..."
if lxc exec $CONTAINER -- $BINARY_PATH pkg uninstall $TEST_PACKAGE; then
    echo -e "${GREEN}✓ Package uninstallation command succeeded${NC}"
else
    echo -e "${RED}✗ Package uninstallation command failed${NC}"
    exit 1
fi

# Verify package is uninstalled
echo "Verifying package uninstallation..."
if ! lxc exec $CONTAINER -- which $TEST_PACKAGE >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Package $TEST_PACKAGE is uninstalled${NC}"
else
    echo -e "${RED}✗ Package $TEST_PACKAGE is still installed${NC}"
    exit 1
fi

echo -e "\n${GREEN}All package management tests passed!${NC}" 