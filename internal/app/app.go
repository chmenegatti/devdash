// Package app implements the Bubble Tea application model.
package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cesar/devdash/internal/state"
	"github.com/cesar/devdash/internal/ui"
)

// Model is the top-level Bubble Tea model for the dashboard.
type Model struct {
	state  *state.Dashboard
	width  int
	height int
	ready  bool
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
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}
	return ui.RenderDashboard(m.state, m.width, m.height)
}

// handleKey processes keyboard input and returns updated model + commands.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	// Phase 1 placeholders - will wire up real commands in later phases.
	case "t":
		m.state.Tests.Status = state.StatusRunning
		return m, nil
	case "c":
		m.state.Coverage.Status = state.StatusRunning
		return m, nil
	case "l":
		m.state.Lint.Status = state.StatusRunning
		return m, nil
	case "b":
		m.state.Benchmarks.Status = state.StatusRunning
		return m, nil
	case "g":
		m.state.Git.Status = state.StatusRunning
		return m, nil
	case "d":
		m.state.Deps.Status = state.StatusRunning
		return m, nil
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
