# 🤝 Contributing to devdash

First off, thank you for considering contributing to **devdash**! 🎉

Every contribution matters — whether it's a bug fix, new feature, documentation improvement, or even a typo correction. This document provides guidelines to help you contribute effectively.

---

## 📋 Table of Contents

- [Code of Conduct](#-code-of-conduct)
- [How Can I Contribute?](#-how-can-i-contribute)
- [Development Setup](#-development-setup)
- [Project Structure](#-project-structure)
- [Coding Guidelines](#-coding-guidelines)
- [Commit Convention](#-commit-convention)
- [Pull Request Process](#-pull-request-process)
- [Reporting Bugs](#-reporting-bugs)
- [Suggesting Features](#-suggesting-features)

---

## 📜 Code of Conduct

This project follows our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

---

## 💡 How Can I Contribute?

### 🐛 Fix a Bug
Check out the [open issues](https://github.com/chmenegatti/devdash/issues?q=is%3Aissue+is%3Aopen+label%3Abug) labeled `bug`.

### ✨ Add a Feature
Look at the [Roadmap](README.md#%EF%B8%8F-roadmap) or issues labeled [`enhancement`](https://github.com/chmenegatti/devdash/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement).

### 📖 Improve Documentation
Documentation improvements are always welcome — README updates, code comments, examples, etc.

### 🧪 Add Tests
Help us improve test coverage! Each module in `internal/modules/` has a corresponding `_test.go` file.

---

## 🛠 Development Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.21+ | Language runtime |
| Git | any | Version control |
| golangci-lint | latest | Linting (optional) |

### Getting Started

```bash
# 1. Fork the repository on GitHub

# 2. Clone your fork
git clone https://github.com/<your-username>/devdash.git
cd devdash

# 3. Add upstream remote
git remote add upstream https://github.com/chmenegatti/devdash.git

# 4. Install dependencies
go mod download

# 5. Verify everything works
go build ./...
go test ./...
go vet ./...
```

### Running the App

```bash
# Build and run
go run ./cmd/dashboard

# Or build a binary
go build -o devdash ./cmd/dashboard
./devdash
```

---

## 📁 Project Structure

```
internal/
├── app/        # Bubble Tea model — the "controller"
├── models/     # Data models (project detection)
├── modules/    # Feature modules — one per dashboard panel
├── services/   # Shell command abstraction layer
├── state/      # Centralized state container
└── ui/         # Rendering layer — styles, components, views
```

### Adding a New Module

1. **Create** `internal/modules/yourmodule.go` with a `RunYourModule(dir string) state.YourResult` function
2. **Add** the result type to `internal/state/state.go`
3. **Wire** a message type and async command in `internal/app/app.go`
4. **Render** a panel in `internal/ui/dashboard.go` and optionally a detail view in `internal/ui/detail_views.go`
5. **Test** — add `internal/modules/yourmodule_test.go` with parser unit tests
6. **Document** — update the README keyboard shortcuts table

---

## 📝 Coding Guidelines

### Go Style

- Follow standard [Go conventions](https://go.dev/doc/effective_go)
- Run `gofmt` before committing (or use `goimports`)
- All exports must have doc comments
- Keep functions small and focused

### Architecture Rules

- **Modules** must be pure: take `services.CommandResult`, return typed state — no side effects
- **UI functions** must be pure renderers: take state, return string — no I/O
- **Services** are the only layer that touches `os/exec`
- **State** is read/write only through the Bubble Tea Update cycle

### Testing

```bash
# Run all tests
go test ./... -v

# Run with race detection
go test ./... -race

# Run specific tests
go test ./internal/modules/ -v -run TestParse

# Check coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

- Every parser function must have unit tests
- Use `services.CommandResult` structs to mock command output
- Integration tests should be guarded with `testing.Short()`

---

## 💬 Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, no code change |
| `refactor` | Code restructuring |
| `test` | Adding or fixing tests |
| `chore` | Build, CI, tooling |
| `perf` | Performance improvement |

### Examples

```
feat(modules): add go vet module with issue parsing
fix(ui): correct panel width calculation on narrow terminals
docs: update README with new keyboard shortcuts
test(modules): add coverage parser edge cases
```

---

## 🔀 Pull Request Process

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

2. **Make your changes** with small, focused commits

3. **Ensure quality**:
   ```bash
   go test ./... -v
   go vet ./...
   gofmt -l .
   ```

4. **Push** to your fork:
   ```bash
   git push origin feature/my-awesome-feature
   ```

5. **Open a Pull Request** against `main` with:
   - Clear title following commit conventions
   - Description of what changed and why
   - Screenshots/recordings for UI changes
   - Link to related issue(s)

6. **Respond to review** — maintainers may request changes

### PR Checklist

- [ ] Code compiles (`go build ./...`)
- [ ] All tests pass (`go test ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] Code is formatted (`gofmt`)
- [ ] New code has tests
- [ ] Documentation updated if needed

---

## 🐛 Reporting Bugs

Use the [Bug Report template](https://github.com/chmenegatti/devdash/issues/new?template=bug_report.md) and include:

- **Go version** (`go version`)
- **OS and terminal** (e.g., macOS + iTerm2, Linux + Alacritty)
- **Steps to reproduce**
- **Expected vs actual behavior**
- **Terminal screenshot** if it's a UI issue

---

## 💡 Suggesting Features

Use the [Feature Request template](https://github.com/chmenegatti/devdash/issues/new?template=feature_request.md) and describe:

- **Problem** — What pain point does this solve?
- **Proposed solution** — How should it work?
- **Alternatives** — What else did you consider?

---

## ❤️ Thank You!

Every contribution makes devdash better for the entire Go community. We appreciate your time and effort!

If you have questions, feel free to open a [Discussion](https://github.com/chmenegatti/devdash/discussions) or reach out to the maintainers.
