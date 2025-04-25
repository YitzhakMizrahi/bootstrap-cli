# TROUBLESHOOTING.md

## Common Issues & Resolutions

| Symptom                                      | Possible Cause                               | Resolution                                                                                  |
|----------------------------------------------|----------------------------------------------|---------------------------------------------------------------------------------------------|
| `Failed to extract embedded configs: ...`    | Missing or incorrect `go:embed` pattern       | Ensure `//go:embed defaults/**` is declared and paths match directory structure.            |
| `panic: interface conversion`                | Bad type assertion in config merge           | Add type checks with `ok` on assertions and return errors instead of panicking.            |
| `root privileges required`                  | Running install without `sudo`                | Rerun command with `sudo`; consider auto-relaunch with elevated rights.                     |
| `failed to detect package manager`          | Unsupported distro or missing binary         | Verify distro support; ensure `apt`, `dnf`, `pacman`, or `brew` are installed and on PATH.  |
| `step timed out after`                       | Long-running install action                  | Increase `Timeout` for that step or optimize the installation command.                     |
| `rollback failed: ...`                      | Rollback hook error                           | Inspect rollback function and add idempotent cleanup for partial-state scenarios.           |

## Debugging Techniques

1. **Enable Debug Logging**  
   ```bash
   bootstrap-cli --debug init
   ```
2. **Inspect Temporary Configs**  
   ```bash
   ls $TMPDIR/bootstrap-cli-config/defaults
   ```
3. **Manual Step Execution**  
   - Copy the failing command from logs and run it directly to capture verbose errors.  
   - Use `strace` or `ltrace` for low-level diagnostics.
4. **Pipeline Progress**  
   ```go
   fmt.Println(installer.GetProgress())
   ```

## Integration Test Failures

- Ensure LXC container has correct user permissions and network access.  
- Clean up leftover containers: `lxc delete --force bootstrap-test`.  
- Snapshot rollback: `lxc restore bootstrap-test clean-setup`.

## Reporting Bugs

1. Collect logs: `~/.cache/bootstrap-cli/logs/*.log`  
2. Reproduce with `--debug` and include log file in issue.  
3. Open an issue at `https://github.com/YitzhakMizrahi/bootstrap-cli/issues` with:  
   - Steps to reproduce  
   - Version (`bootstrap-cli version`)  
   - OS/distro info (`uname -a`)

