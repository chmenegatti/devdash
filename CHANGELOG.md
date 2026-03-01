# 📋 Changelog

All notable changes to **devdash** will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- 🧪 **Test Runner** — `go test ./...` with pass/fail status, package count, duration
- 📊 **Coverage** — `go test -cover` with color-coded percentage thresholds
- 🔍 **Linter** — `golangci-lint run` with issue count and inline preview
- ⚡ **Benchmarks** — `go test -bench` with table display (name, iterations, ns/op)
- 📦 **Binary Size** — Build and measure compiled binary size
- 🌿 **Git Status** — `git status --short` with categorised file changes
- 📚 **Dependencies** — `go list -m all` with module listing
- 📝 **Markdown Report Export** — Shortcut `m` generates a full dashboard report (`devdash-report-YYYYMMDD-HHMMSS.md`)
- 🔎 **Detail Views** — Full-screen output for tests (`T`), lint (`L`), benchmarks (`B`), git (`G`), deps (`D`)
- 🎨 **K9s-inspired UI** — Dark theme with cyan accents, breadcrumbs, stat tiles, table layouts, command bar
- ⚡ **Async execution** — All commands run via `tea.Cmd` goroutines, never blocking the UI
- 🧪 **23 unit tests** — Comprehensive parser tests for all modules
- 📄 Open source setup — README, LICENSE (MIT), CONTRIBUTING, CODE_OF_CONDUCT, GitHub templates, CI workflow
- 🪵 **Error Logging** — Persistent diagnostics in `.devdash.log` (failed commands, module/report errors, timestamped entries)

---

## [0.1.0] - 2026-03-01

### Added
- Initial release with full dashboard functionality
- 7 integrated panels: Tests, Coverage, Lint, Benchmarks, Binary Size, Git, Dependencies
- K9s-inspired terminal UI with Bubble Tea + Lip Gloss
- Keyboard-driven navigation with detail views

[Unreleased]: https://github.com/chmenegatti/devdash/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/chmenegatti/devdash/releases/tag/v0.1.0
