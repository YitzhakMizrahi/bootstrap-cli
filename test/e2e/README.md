# End-to-End Testing with LXC

This directory contains end-to-end tests that run in LXC containers to verify the behavior of bootstrap-cli in isolated environments.

## Test Environment

We use Ubuntu 22.04 LXC containers (`bootstrap-test`) for testing. This provides:
- Isolated testing environment
- Snapshot/restore capabilities
- Real system package management
- Privilege operations testing

## Setup

1. Create base container:
```bash
lxc launch ubuntu:22.04 bootstrap-test
```

2. Create initial snapshot (clean state):
```bash
lxc snapshot bootstrap-test clean
```

3. Push binary to container:
```bash
lxc file push ./bootstrap-cli bootstrap-test/home/devuser/bootstrap-cli --mode=755
```

## Test Scenarios

### 1. Package Management Tests
```bash
# Start from clean snapshot
lxc restore bootstrap-test clean

# Run package installation test
lxc exec bootstrap-test -- /home/devuser/bootstrap-cli pkg install git

# Verify installation
lxc exec bootstrap-test -- which git
```

### 2. Shell Configuration Tests
```bash
# Start from clean snapshot
lxc restore bootstrap-test clean

# Run shell setup test
lxc exec bootstrap-test -- /home/devuser/bootstrap-cli shell setup --type=zsh

# Verify configuration
lxc exec bootstrap-test -- cat /home/devuser/.zshrc
```

### 3. Full Environment Setup Test
```bash
# Start from clean snapshot
lxc restore bootstrap-test clean

# Run full setup
lxc exec bootstrap-test -- /home/devuser/bootstrap-cli setup --config=/path/to/test/config.yaml

# Run verification checks
./test/e2e/verify.sh bootstrap-test
```

## Creating Test Snapshots

Different test scenarios might need different base states. Here's how to manage them:

1. **Clean State**:
   ```bash
   lxc snapshot bootstrap-test clean
   ```

2. **With Development Tools**:
   ```bash
   lxc restore bootstrap-test clean
   lxc exec bootstrap-test -- apt install -y build-essential git
   lxc snapshot bootstrap-test dev-tools
   ```

3. **With Shell Setup**:
   ```bash
   lxc restore bootstrap-test clean
   lxc exec bootstrap-test -- /home/devuser/bootstrap-cli shell setup
   lxc snapshot bootstrap-test shell-setup
   ```

## Test Scripts

The `scripts/` directory contains helper scripts for running tests:

- `run_e2e_tests.sh`: Runs all end-to-end tests
- `verify.sh`: Verifies the state of a container after tests
- `cleanup.sh`: Cleans up test containers and snapshots

## Adding New Tests

1. Create a new test script in `test/e2e/tests/`
2. Update `run_e2e_tests.sh` to include your test
3. Add verification steps to `verify.sh`
4. Document any new snapshots needed

## Best Practices

1. **Always start from a known state**:
   ```bash
   lxc restore bootstrap-test clean
   ```

2. **Verify changes**:
   ```bash
   # Check file contents
   lxc exec bootstrap-test -- cat /path/to/file
   
   # Check installed packages
   lxc exec bootstrap-test -- dpkg -l | grep package-name
   
   # Check processes
   lxc exec bootstrap-test -- ps aux | grep process-name
   ```

3. **Clean up after tests**:
   ```bash
   # Remove test artifacts
   lxc exec bootstrap-test -- rm -rf /path/to/artifacts
   
   # Restore clean state
   lxc restore bootstrap-test clean
   ```

4. **Document test requirements**:
   - Required snapshots
   - Expected initial state
   - Expected final state
   - Verification steps 