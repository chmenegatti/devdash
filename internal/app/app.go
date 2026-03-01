// Package app implements the Bubble Tea application model.
package app

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/chmenegatti/devdash/internal/logs"
	"github.com/chmenegatti/devdash/internal/modules"
	"github.com/chmenegatti/devdash/internal/state"
	"github.com/chmenegatti/devdash/internal/ui"
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

// profileResultMsg carries the result of a completed CPU profile run.
type profileResultMsg struct {
	result state.ProfileResult
}

// reportResultMsg carries the result of markdown report generation.
type reportResultMsg struct {
	path string
	err  error
}

// viewMode represents which screen is currently displayed.
type viewMode int

const (
	viewDashboard     viewMode = iota // Main dashboard
	viewTestsDetail                   // Full test output
	viewLintDetail                    // Full lint output
	viewBenchDetail                   // Full benchmark output
	viewDepsDetail                    // Full dependency list
	viewGitDetail                     // Full git status
	viewProfileDetail                 // Full inline flamegraph
)

// Model is the top-level Bubble Tea model for the dashboard.
type Model struct {
	state          *state.Dashboard
	width          int
	height         int
	ready          bool
	view           viewMode
	detailViewport viewport.Model
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
		m.syncDetailViewport()
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)

	// ── Async result messages ──────────────────────────────────
	case testsResultMsg:
		m.state.Tests = msg.result
		logModuleError("tests", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case coverageResultMsg:
		m.state.Coverage = msg.result
		logModuleError("coverage", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case lintResultMsg:
		m.state.Lint = msg.result
		logModuleError("lint", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case benchResultMsg:
		m.state.Benchmarks = msg.result
		logModuleError("benchmarks", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case binaryResultMsg:
		m.state.Binary = msg.result
		logModuleError("binary", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case depsResultMsg:
		m.state.Deps = msg.result
		logModuleError("dependencies", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case gitResultMsg:
		m.state.Git = msg.result
		logModuleError("git", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case profileResultMsg:
		m.state.Profile = msg.result
		logModuleError("profile", msg.result.Status, msg.result.Err)
		m.syncDetailViewport()
		return m, nil
	case reportResultMsg:
		if msg.err != nil {
			logs.Errorf("report generation failed: %v", msg.err)
			m.state.Notice = "❌ Erro ao gerar relatório: " + msg.err.Error()
			return m, nil
		}
		logs.Infof("report generated: %s", msg.path)
		fileName := filepath.Base(msg.path)
		m.state.Notice = "📝 Relatório gerado: " + fileName + " (" + msg.path + ")"
		return m, nil
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
		detail := m.currentDetailContent()
		return ui.RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, m.detailViewport.View(), m.state.Version, m.width, m.height)
	case viewLintDetail:
		detail := m.currentDetailContent()
		return ui.RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, m.detailViewport.View(), m.state.Version, m.width, m.height)
	case viewBenchDetail:
		detail := m.currentDetailContent()
		return ui.RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, m.detailViewport.View(), m.state.Version, m.width, m.height)
	case viewDepsDetail:
		detail := m.currentDetailContent()
		return ui.RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, m.detailViewport.View(), m.state.Version, m.width, m.height)
	case viewGitDetail:
		detail := m.currentDetailContent()
		return ui.RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, m.detailViewport.View(), m.state.Version, m.width, m.height)
	case viewProfileDetail:
		detail := m.currentDetailContent()
		return ui.RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, m.detailViewport.View(), m.state.Version, m.width, m.height)
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
		var cmd tea.Cmd
		m.detailViewport, cmd = m.detailViewport.Update(msg)
		return m, cmd
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
		m.syncDetailViewport()
		return m, nil
	case "c":
		m.state.Coverage.Status = state.StatusRunning
		return m, m.runCoverageCmd()
	case "l":
		m.state.Lint.Status = state.StatusRunning
		return m, m.runLintCmd()
	case "L":
		m.view = viewLintDetail
		m.syncDetailViewport()
		return m, nil
	case "b":
		m.state.Benchmarks.Status = state.StatusRunning
		return m, m.runBenchCmd()
	case "B":
		m.view = viewBenchDetail
		m.syncDetailViewport()
		return m, nil
	case "s":
		m.state.Binary.Status = state.StatusRunning
		return m, m.runBinaryCmd()
	case "g":
		m.state.Git.Status = state.StatusRunning
		return m, m.runGitCmd()
	case "G":
		m.view = viewGitDetail
		m.syncDetailViewport()
		return m, nil
	case "d":
		m.state.Deps.Status = state.StatusRunning
		return m, m.runDepsCmd()
	case "D":
		m.view = viewDepsDetail
		m.syncDetailViewport()
		return m, nil
	case "p":
		m.state.Profile.Status = state.StatusRunning
		return m, m.runProfileCmd()
	case "P":
		m.view = viewProfileDetail
		m.syncDetailViewport()
		return m, nil
	case "m":
		m.state.Notice = "🛠️ Gerando relatório Markdown..."
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
		m.state.Profile = state.ProfileResult{}
		m.state.Notice = ""
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

// runProfileCmd returns a tea.Cmd that runs CPU profiling and produces an inline flamegraph.
func (m Model) runProfileCmd() tea.Cmd {
	dir := m.state.ProjectDir
	return func() tea.Msg {
		result := modules.RunCPUProfile(dir)
		return profileResultMsg{result: result}
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

func logModuleError(module string, status state.Status, errText string) {
	errText = strings.TrimSpace(errText)
	if status != state.StatusError && errText == "" {
		return
	}
	if errText == "" {
		errText = "unknown error"
	}
	logs.Errorf("module=%s status=%s err=%s", module, status.String(), errText)
}

func (m *Model) syncDetailViewport() {
	if m.view == viewDashboard || m.width <= 0 || m.height <= 0 {
		return
	}

	detail := m.currentDetailContent()
	bodyWidth, bodyHeight := ui.DetailBodySize(m.width, m.height)

	if m.detailViewport.Width == 0 || m.detailViewport.Height == 0 {
		m.detailViewport = viewport.New(bodyWidth, bodyHeight)
	} else {
		m.detailViewport.Width = bodyWidth
		m.detailViewport.Height = bodyHeight
	}

	m.detailViewport.SetContent(detail.Body)
}

func (m Model) currentDetailContent() ui.DetailContent {
	switch m.view {
	case viewTestsDetail:
		return ui.BuildTestsDetail(m.state, m.width)
	case viewLintDetail:
		return ui.BuildLintDetail(m.state, m.width)
	case viewBenchDetail:
		return ui.BuildBenchDetail(m.state, m.width)
	case viewDepsDetail:
		return ui.BuildDepsDetail(m.state, m.width)
	case viewGitDetail:
		return ui.BuildGitDetail(m.state, m.width)
	case viewProfileDetail:
		return ui.BuildProfileDetail(m.state, m.width)
	default:
		return ui.DetailContent{}
	}
}
