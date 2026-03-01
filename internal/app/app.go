// Package app implements the Bubble Tea application model.
package app

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/cesar/devdash/internal/modules"
	"github.com/cesar/devdash/internal/state"
	"github.com/cesar/devdash/internal/ui"
)

// ── Messages ────────────────────────────────────────────────────────────────

// testsResultMsg carries the result of a completed test run back to Update.
type testsResultMsg struct {
	result state.TestsResult
}

// coverageResultMsg carries the result of a completed coverage run.
type coverageResultMsg struct {
	result state.CoverageResult
}

// lintResultMsg carries the result of a completed lint run.
type lintResultMsg struct {
	result state.LintResult
}

// benchResultMsg carries the result of a completed benchmark run.
type benchResultMsg struct {
	result state.BenchmarkResult
}

// binaryResultMsg carries the result of a completed binary size check.
type binaryResultMsg struct {
	result state.BinaryResult
}

// depsResultMsg carries the result of a completed dependency listing.
type depsResultMsg struct {
	result state.DepsResult
}

// gitResultMsg carries the result of a completed git status check.
type gitResultMsg struct {
	result state.GitResult
}

// reportResultMsg carries the result of markdown report generation.
type reportResultMsg struct {
	path string
	err  error
}

// viewMode represents which screen is currently displayed.
type viewMode int

const (
	viewDashboard   viewMode = iota // Main dashboard
	viewTestsDetail                 // Full test output
	viewLintDetail                  // Full lint output
	viewBenchDetail                 // Full benchmark output
	viewDepsDetail                  // Full dependency list
	viewGitDetail                   // Full git status
)

// Model is the top-level Bubble Tea model for the dashboard.
type Model struct {
	state  *state.Dashboard
	width  int
	height int
	ready  bool
	view   viewMode
}

// New creates a new application Model.
func New(ds *state.Dashboard) Model {
	return Model{
		state: ds,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)

	// ── Async result messages ──────────────────────────────────
	case testsResultMsg:
		m.state.Tests = msg.result
		return m, nil
	case coverageResultMsg:
		m.state.Coverage = msg.result
		return m, nil
	case lintResultMsg:
		m.state.Lint = msg.result
		return m, nil
	case benchResultMsg:
		m.state.Benchmarks = msg.result
		return m, nil
	case binaryResultMsg:
		m.state.Binary = msg.result
		return m, nil
	case depsResultMsg:
		m.state.Deps = msg.result
		return m, nil
	case gitResultMsg:
		m.state.Git = msg.result
		return m, nil
	case reportResultMsg:
		if msg.err != nil {
			return m, tea.Printf("❌ Erro ao gerar relatório: %v", msg.err)
		}
		fileName := filepath.Base(msg.path)
		return m, tea.Printf("📝 Relatório gerado: %s (%s)", fileName, msg.path)
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}
	switch m.view {
	case viewTestsDetail:
		return ui.RenderTestsDetail(m.state, m.width, m.height)
	case viewLintDetail:
		return ui.RenderLintDetail(m.state, m.width, m.height)
	case viewBenchDetail:
		return ui.RenderBenchDetail(m.state, m.width, m.height)
	case viewDepsDetail:
		return ui.RenderDepsDetail(m.state, m.width, m.height)
	case viewGitDetail:
		return ui.RenderGitDetail(m.state, m.width, m.height)
	default:
		return ui.RenderDashboard(m.state, m.width, m.height)
	}
}

// handleKey processes keyboard input and returns updated model + commands.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys — work in any view
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "backspace":
		m.view = viewDashboard
		return m, nil
	}

	// View-specific keys
	switch m.view {
	case viewDashboard:
		return m.handleDashboardKey(msg)
	default:
		// In detail views, only global keys apply
		return m, nil
	}
}

// handleDashboardKey handles keys specific to the main dashboard view.
func (m Model) handleDashboardKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "t":
		m.state.Tests.Status = state.StatusRunning
		return m, m.runTestsCmd()
	case "T":
		m.view = viewTestsDetail
		return m, nil
	case "c":
		m.state.Coverage.Status = state.StatusRunning
		return m, m.runCoverageCmd()
	case "l":
		m.state.Lint.Status = state.StatusRunning
		return m, m.runLintCmd()
	case "L":
		m.view = viewLintDetail
		return m, nil
	case "b":
		m.state.Benchmarks.Status = state.StatusRunning
		return m, m.runBenchCmd()
	case "B":
		m.view = viewBenchDetail
		return m, nil
	case "s":
		m.state.Binary.Status = state.StatusRunning
		return m, m.runBinaryCmd()
	case "g":
		m.state.Git.Status = state.StatusRunning
		return m, m.runGitCmd()
	case "G":
		m.view = viewGitDetail
		return m, nil
	case "d":
		m.state.Deps.Status = state.StatusRunning
		return m, m.runDepsCmd()
	case "D":
		m.view = viewDepsDetail
		return m, nil
	case "m":
		return m, m.runReportCmd()
	case "r":
		// refresh - reset all to idle
		m.state.Tests = state.TestsResult{}
		m.state.Coverage = state.CoverageResult{}
		m.state.Lint = state.LintResult{}
		m.state.Benchmarks = state.BenchmarkResult{}
		m.state.Binary = state.BinaryResult{}
		m.state.Deps = state.DepsResult{}
		m.state.Git = state.GitResult{}
		return m, nil
	}
	return m, nil
}

// ── Async commands ──────────────────────────────────────────────────────────

// runTestsCmd returns a tea.Cmd that runs go test asynchronously.
func (m Model) runTestsCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunTests(dir)
		return testsResultMsg{result: result}
	}
}

// runCoverageCmd returns a tea.Cmd that runs go test -cover asynchronously.
func (m Model) runCoverageCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunCoverage(dir)
		return coverageResultMsg{result: result}
	}
}

// runLintCmd returns a tea.Cmd that runs golangci-lint asynchronously.
func (m Model) runLintCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunLint(dir)
		return lintResultMsg{result: result}
	}
}

// runBenchCmd returns a tea.Cmd that runs go test -bench asynchronously.
func (m Model) runBenchCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunBenchmarks(dir)
		return benchResultMsg{result: result}
	}
}

// runBinaryCmd returns a tea.Cmd that builds and measures binary size.
func (m Model) runBinaryCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunBinarySize(dir)
		return binaryResultMsg{result: result}
	}
}

// runDepsCmd returns a tea.Cmd that lists module dependencies.
func (m Model) runDepsCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunDeps(dir)
		return depsResultMsg{result: result}
	}
}

// runGitCmd returns a tea.Cmd that runs git status.
func (m Model) runGitCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunGitStatus(dir)
		return gitResultMsg{result: result}
	}
}

// runReportCmd returns a tea.Cmd that exports a markdown report from current dashboard state.
func (m Model) runReportCmd() tea.Cmd {
	projectDir := m.state.ProjectDir
	snapshot := snapshotDashboard(*m.state)
	return func() tea.Msg {
		path, err := modules.GenerateReportFile(projectDir, snapshot)
		return reportResultMsg{path: path, err: err}
	}
}

func snapshotDashboard(ds state.Dashboard) state.Dashboard {
	s := ds
	s.Lint.Issues = append([]string(nil), ds.Lint.Issues...)
	s.Benchmarks.Entries = append([]state.BenchmarkEntry(nil), ds.Benchmarks.Entries...)
	s.Deps.Deps = append([]string(nil), ds.Deps.Deps...)
	s.Git.Modified = append([]string(nil), ds.Git.Modified...)
	s.Git.Added = append([]string(nil), ds.Git.Added...)
	s.Git.Deleted = append([]string(nil), ds.Git.Deleted...)
	s.Git.Other = append([]string(nil), ds.Git.Other...)
	return s
}
