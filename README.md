<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/github/license/chmenegatti/devdash?style=for-the-badge&color=00d7d7" alt="License">
  <img src="https://img.shields.io/github/stars/chmenegatti/devdash?style=for-the-badge&color=FFD700" alt="Stars">
  <img src="https://img.shields.io/github/issues/chmenegatti/devdash?style=for-the-badge&color=ff5f5f" alt="Issues">
  <img src="https://img.shields.io/github/actions/workflow/status/chmenegatti/devdash/ci.yml?style=for-the-badge&label=CI" alt="CI">
  <img src="https://img.shields.io/codecov/c/github/chmenegatti/devdash?style=for-the-badge&label=Coverage" alt="Coverage">
</p>

<h1 align="center">
  ⎈ devdash
</h1>

<p align="center">
  <strong>A K9s-inspired terminal dashboard for Go developers</strong>
</p>

<p align="center">
  <em>Run tests, check coverage, lint, benchmark, inspect dependencies, and monitor git status — all from one beautiful TUI.</em>
</p>

<p align="center">
  <a href="#-features">Features</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#%EF%B8%8F-keyboard-shortcuts">Shortcuts</a> •
  <a href="#-architecture">Architecture</a> •
  <a href="#-contributing">Contributing</a> •
  <a href="#-license">License</a>
</p>

---

## ✨ Features

| Feature | Key | Description |
|---------|-----|-------------|
| 🧪 **Test Runner** | `t` | Run `go test ./...` with pass/fail, package count & duration |
| 📊 **Coverage** | `c` | Run `go test -cover` with color-coded percentage (🟢 ≥80% 🟡 ≥60% 🔴 <60%) |
| 🔍 **Linter** | `l` | Run `golangci-lint` and display issue count with inline preview |
| ⚡ **Benchmarks** | `b` | Run `go test -bench` with table of iterations & ns/op |
| 📦 **Binary Size** | `s` | Build and measure compiled binary size |
| 🌿 **Git Status** | `g` | Show modified/added/deleted/untracked files with colored indicators |
| 📚 **Dependencies** | `d` | List all module dependencies via `go list -m all` |
| 🔎 **Detail Views** | `Shift+Key` | Full-screen output for tests, lint, benchmarks, git & deps |

### 🎨 Design Philosophy

Inspired by [K9s](https://k9scli.io/) — the legendary Kubernetes TUI — **devdash** brings the same sleek, dark-themed, keyboard-driven experience to your Go development workflow:

- 🖤 **Dark theme** with cyan/teal accents
- 📍 **Breadcrumb navigation** between views
- 📊 **Stat tiles** with colored status dots (●/◍/○)
- 📋 **Table layouts** with alternating row highlights
- ⌨️ **Command bar** with discoverable key hints
- ⚡ **Async execution** — UI never blocks while commands run

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21 or later
- **golangci-lint** (optional, for lint panel) — [install guide](https://golangci-lint.run/welcome/install/)
- **Git** (for git status panel)

### Install from source

```bash
go install github.com/chmenegatti/devdash/cmd/dashboard@latest
```

### Or clone and build

```bash
git clone https://github.com/chmenegatti/devdash.git
cd devdash
go build -o devdash ./cmd/dashboard
./devdash
```

### Run directly

```bash
# Navigate to any Go project, then:
devdash
```

> 💡 **Tip:** devdash auto-detects the project from your current working directory.

---

## ⌨️ Keyboard Shortcuts

### Dashboard View

| Key | Action |
|-----|--------|
| `t` | Run tests |
| `c` | Run coverage |
| `l` | Run linter |
| `b` | Run benchmarks |
| `s` | Measure binary size |
| `g` | Check git status |
| `d` | List dependencies |
| `r` | Reset all panels |
| `q` | Quit |

### Detail Views

| Key | Action |
|-----|--------|
| `T` | Full test output |
| `L` | Full lint output |
| `B` | Full benchmark table |
| `G` | Full git status |
| `D` | Full dependency list |
| `Backspace` | Back to dashboard |

---

## 🏗 Architecture

```
devdash/
├── cmd/
│   └── dashboard/          # 🚀 Application entrypoint
│       └── main.go
├── internal/
│   ├── app/                # 🎮 Bubble Tea model (Update/View/Cmd)
│   │   ├── app.go          #     Central model, key dispatch, async wiring
│   │   └── layout.go       #     Layout helpers
│   ├── models/             # 📁 Project detection
│   │   └── project.go
│   ├── modules/            # ⚙️  Feature modules (one per panel)
│   │   ├── tests.go        #     go test runner & parser
│   │   ├── coverage.go     #     go test -cover parser
│   │   ├── lint.go         #     golangci-lint runner & parser
│   │   ├── benchmarks.go   #     go test -bench parser
│   │   ├── binary.go       #     Binary size measurement
│   │   ├── deps.go         #     go list -m all parser
│   │   └── gitstatus.go    #     git status --short parser
│   ├── services/           # 🔧 Shell command abstraction
│   │   ├── exec.go         #     RunCommand wrapper
│   │   └── parser.go       #     Line parsing utilities
│   ├── state/              # 💾 Centralized state management
│   │   └── state.go        #     Dashboard struct + result types
│   └── ui/                 # 🎨 K9s-inspired rendering
│       ├── styles.go       #     Color palette & Lipgloss styles
│       ├── components.go   #     Logo, crumbs, tables, command bar
│       ├── dashboard.go    #     Main dashboard composition
│       └── detail_views.go #     Full-screen detail renderers
├── go.mod
├── go.sum
├── LICENSE
├── README.md
├── CONTRIBUTING.md
├── CHANGELOG.md
└── CODE_OF_CONDUCT.md
```

### Design Patterns

- **Bubble Tea (Elm Architecture)** — Model → Update → View with pure rendering
- **Async Commands** — All shell operations run via `tea.Cmd` goroutines, never blocking the UI
- **Layered Design** — `app` → `ui` → `modules` → `services` → `state`
- **Pure Parsers** — Each module's parser takes a `CommandResult` and returns typed state — easy to unit test

### Tech Stack

| Library | Purpose |
|---------|---------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Terminal UI framework (Elm-style) |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | Styling, layout & colors |
| [Bubbles](https://github.com/charmbracelet/bubbles) | UI components (available for extensions) |

---

## 🧪 Testing

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific module tests
go test ./internal/modules/ -v -run TestParse

# Short mode (skip integration tests)
go test ./... -short
```

Currently **23 unit tests** covering all module parsers plus integration tests for binary size measurement.

---

## 🗺️ Roadmap

We welcome contributions for any of these planned features! See [CONTRIBUTING.md](CONTRIBUTING.md).

- [ ] 🔄 **Auto-refresh** — Periodic re-run with configurable interval
- [ ] 📜 **Scrollable panels** — Scroll through long outputs in detail views
- [ ] 🎛️ **Config file** — `.devdash.yaml` for custom panel layout, colors, and shortcuts
- [ ] 📈 **Flame graphs** — pprof integration with inline visualization
- [ ] 🐳 **Docker support** — Build & run inside containers
- [ ] 🔌 **Plugin system** — Custom panels via Go plugins or external scripts
- [ ] 🌐 **Remote mode** — Monitor CI/CD pipelines via SSH or API
- [ ] 📋 **Clipboard** — Copy panel output to clipboard
- [ ] 🎯 **Focused test run** — Run a single test function by name
- [ ] 📊 **History** — Track test/coverage trends over time
- [ ] 🔔 **Notifications** — Desktop alerts on test failure
- [ ] 🖥️ **Resizable panes** — Drag-to-resize panel layout

---

## 🤝 Contributing

We love contributions! Whether it's a bug fix, new feature, documentation improvement, or just a typo — every PR matters.

Please read our [Contributing Guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) before getting started.

```bash
# Fork & clone
git clone https://github.com/<your-user>/devdash.git
cd devdash

# Create a branch
git checkout -b feature/amazing-feature

# Make changes, then test
go test ./... -v
go vet ./...

# Commit & push
git add -A
git commit -m "feat: add amazing feature"
git push origin feature/amazing-feature
```

Then open a Pull Request 🚀

---

## 📖 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [**Charm**](https://charm.sh/) — For the incredible Bubble Tea & Lip Gloss libraries
- [**K9s**](https://k9scli.io/) — For the design inspiration
- [**golangci-lint**](https://golangci-lint.run/) — For the Go linting ecosystem
- All our [contributors](https://github.com/chmenegatti/devdash/graphs/contributors) ❤️

---

<p align="center">
  <strong>⭐ If you find devdash useful, give it a star!</strong>
</p>

<p align="center">
  Made with ❤️ and Go
</p>
