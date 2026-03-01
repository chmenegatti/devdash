# 🔒 Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest `main` | ✅ |
| < 0.1.0 | ❌ |

## Reporting a Vulnerability

If you discover a security vulnerability in devdash, please report it responsibly:

1. **Do NOT** open a public GitHub issue
2. **Email** the maintainers directly or use [GitHub Security Advisories](https://github.com/chmenegatti/devdash/security/advisories/new)
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will acknowledge receipt within **48 hours** and provide a fix timeline within **7 days**.

## Scope

devdash executes shell commands (`go test`, `go build`, `golangci-lint`, `git`) on the local machine. It:

- **Does NOT** send any data over the network
- **Does NOT** require elevated privileges
- **Does** execute commands in the user's project directory

Please report any concerns about command injection or unexpected command execution.
