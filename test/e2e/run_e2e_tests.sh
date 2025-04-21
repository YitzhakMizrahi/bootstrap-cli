#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
CONTAINER="bootstrap-test"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTS_DIR="$SCRIPT_DIR/tests"
BINARY_PATH="../../bootstrap-cli"

# Ensure binary exists and is up to date
echo -e "${BLUE}Building bootstrap-cli...${NC}"
(cd ../.. && go build -o bootstrap-cli)

# Push binary to container
echo -e "${BLUE}Pushing binary to container...${NC}"
lxc file push $BINARY_PATH $CONTAINER/home/devuser/bootstrap-cli --mode=755

# Run all test scripts
echo -e "${BLUE}Running end-to-end tests...${NC}\n"

FAILED_TESTS=()
TOTAL_TESTS=0
PASSED_TESTS=0

for test in "$TESTS_DIR"/*.sh; do
    if [ -f "$test" ]; then
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        test_name=$(basename "$test")
        
        echo -e "${BLUE}Running test: $test_name${NC}"
        if bash "$test"; then
            PASSED_TESTS=$((PASSED_TESTS + 1))
            echo -e "${GREEN}✓ Test passed: $test_name${NC}\n"
        else
            FAILED_TESTS+=("$test_name")
            echo -e "${RED}✗ Test failed: $test_name${NC}\n"
        fi
    fi
done

# Print summary
echo -e "${BLUE}Test Summary:${NC}"
echo -e "Total tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}${#FAILED_TESTS[@]}${NC}"

if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    echo -e "\n${RED}Failed tests:${NC}"
    for test in "${FAILED_TESTS[@]}"; do
        echo -e "${RED}- $test${NC}"
    done
    exit 1
fi

echo -e "\n${GREEN}All tests passed!${NC}" 