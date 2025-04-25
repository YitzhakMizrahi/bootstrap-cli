# SECURITY.md

## Security Overview
Bootstrap CLI must be secure by default and follow best practices to protect users and systems.

### Principle of Least Privilege
- Run non-critical operations as an unprivileged user.
- Only request elevated (`sudo`) privileges when absolutely necessary (e.g., package installation).
- Immediately drop elevated rights after performing privileged actions.

### Input Validation & Sanitization
- Validate and sanitize all user-supplied inputs (URLs, file paths, config values).
- Reject or escape dangerous patterns (e.g., shell metacharacters) before execution.
- Avoid using `os/exec` with unsanitized strings; prefer building argument slices.

### Secure Defaults
- Set file permissions conservatively (e.g., `0644` for config files, `0755` for executables).
- Disable automatic remote code execution unless explicitly approved by the user.
- Ship minimal embedded defaults (no secrets or credentials in repos).

### Dependency Management
- Pin critical dependencies in `go.mod` and run `go mod tidy` regularly.
- Audit external dependencies for known vulnerabilities.
- Use `govulncheck` or similar tools in CI to scan for CVEs.

### Logging & Audit
- Log security-relevant events (e.g., privilege elevation, config changes) with timestamps.
- Avoid logging sensitive data (passwords, tokens, user secrets).
- Provide a `--log-file` flag for users to redirect logs securely.

### Update & Patch Strategy
- Encourage users to update to the latest Bootstrap CLI version.
- Sign releases or provide checksums to verify binary integrity.
- Document the update process in `docs/SECURITY.md`.

