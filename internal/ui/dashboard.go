// Package ui — dashboard.go renders the full dashboard view from state.
package ui

import (
	"fmt"

	"github.com/cesar/devdash/internal/state"
	"github.com/charmbracelet/lipgloss"
)

const defaultPanelWidth = 36

// RenderDashboard composes the full terminal view from the dashboard state.
func RenderDashboard(ds *state.Dashboard, width, height int) string {
	// ── Header ──────────────────────────────────────────────────
	header := TitleStyle.Render("  Go Developer Dashboard  ")
	projectLine := SubtitleStyle.Render(
		fmt.Sprintf("Project: %s  (%s)", ds.ProjectName, ds.ProjectDir),
	)

	// ── Panels ─────────────────────────────────────────────────
	pw := panelWidth(width)

	testsPanel := renderTestsPanel(ds, pw)
	coveragePanel := renderCoveragePanel(ds, pw)
	lintPanel := renderLintPanel(ds, pw)
	benchPanel := renderBenchPanel(ds, pw)
	binaryPanel := renderBinaryPanel(ds, pw)
	gitPanel := renderGitPanel(ds, pw)
	depsPanel := renderDepsPanel(ds, pw)

	// Arrange panels in two columns
	leftCol := lipgloss.JoinVertical(lipgloss.Left,
		testsPanel,
		coveragePanel,
		lintPanel,
		benchPanel,
	)
	rightCol := lipgloss.JoinVertical(lipgloss.Left,
		binaryPanel,
		gitPanel,
		depsPanel,
	)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "  ", rightCol)

	// ── Help bar ───────────────────────────────────────────────
	help := RenderHelp([]KeyBinding{
		{Key: "t", Desc: "tests"},
		{Key: "c", Desc: "coverage"},
		{Key: "l", Desc: "lint"},
		{Key: "b", Desc: "benchmarks"},
		{Key: "s", Desc: "binary size"},
		{Key: "g", Desc: "git status"},
		{Key: "d", Desc: "deps"},
		{Key: "r", Desc: "refresh"},
		{Key: "q", Desc: "quit"},
	})
	help2 := RenderHelp([]KeyBinding{
		{Key: "T", Desc: "tests detail"},
		{Key: "L", Desc: "lint detail"},
		{Key: "B", Desc: "bench detail"},
		{Key: "G", Desc: "git detail"},
		{Key: "D", Desc: "deps detail"},
	})

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		projectLine,
		columns,
		help,
		help2,
	)
}

// panelWidth adapts to terminal width, falling back to a default.
func panelWidth(termWidth int) int {
	if termWidth <= 0 {
		return defaultPanelWidth
	}
	// Two columns + 2-char gap → each column gets roughly half
	pw := (termWidth - 4) / 2
	if pw < 28 {
		pw = 28
	}
	if pw > 50 {
		pw = 50
	}
	return pw
}

// ── Individual panel renderers ──────────────────────────────────────────────

func renderTestsPanel(ds *state.Dashboard, w int) string {
	body := statusLine(ds.Tests.Status)
	if ds.Tests.Status == state.StatusDone {
		pass := "PASS"
		style := StatusPass
		if !ds.Tests.Passed {
			pass = "FAIL"
			style = StatusFail
		}
		body = RenderStatusField("Status", pass, style) + "\n" +
			RenderField("Packages", fmt.Sprintf("%d", ds.Tests.Packages)) + "\n" +
			RenderField("Duration", ds.Tests.Duration.String())
	} else if ds.Tests.Status == state.StatusError {
		body = RenderStatusField("Status", "Error", StatusFail) + "\n" +
			StatusFail.Render(truncate(ds.Tests.Err, 80))
	}
	return RenderPanel("Tests", body, w)
}

func renderCoveragePanel(ds *state.Dashboard, w int) string {
	body := statusLine(ds.Coverage.Status)
	if ds.Coverage.Status == state.StatusDone {
		pct := fmt.Sprintf("%.1f%%", ds.Coverage.Percentage)
		style := coverageStyle(ds.Coverage.Percentage)
		body = RenderStatusField("Coverage", pct, style)
	} else if ds.Coverage.Status == state.StatusError {
		body = RenderStatusField("Status", "Error", StatusFail)
	}
	return RenderPanel("Coverage", body, w)
}

func renderLintPanel(ds *state.Dashboard, w int) string {
	body := statusLine(ds.Lint.Status)
	if ds.Lint.Status == state.StatusDone {
		n := len(ds.Lint.Issues)
		statusText := "OK"
		style := StatusPass
		if n > 0 {
			statusText = fmt.Sprintf("%d issues", n)
			style = StatusWarn
		}
		body = RenderStatusField("Status", statusText, style)
		// Show first 5 issues
		for i, iss := range ds.Lint.Issues {
			if i >= 5 {
				body += "\n" + StatusWarn.Render(fmt.Sprintf("  … and %d more", n-5))
				break
			}
			body += "\n  " + lipgloss.NewStyle().Foreground(ColorDim).Render(truncate(iss, w-6))
		}
	}
	return RenderPanel("Lint", body, w)
}

func renderBenchPanel(ds *state.Dashboard, w int) string {
	body := statusLine(ds.Benchmarks.Status)
	if ds.Benchmarks.Status == state.StatusDone {
		if len(ds.Benchmarks.Entries) == 0 {
			body = StatusIdle.Render("No benchmarks found")
		} else {
			lines := ""
			for i, e := range ds.Benchmarks.Entries {
				if i >= 5 {
					break
				}
				lines += fmt.Sprintf("%s  %d  %.0f ns/op\n",
					ValueStyle.Render(e.Name),
					e.Iterations,
					e.NsPerOp,
				)
			}
			body = lines
		}
	}
	return RenderPanel("Benchmarks", body, w)
}

func renderBinaryPanel(ds *state.Dashboard, w int) string {
	body := statusLine(ds.Binary.Status)
	if ds.Binary.Status == state.StatusDone {
		body = RenderField("Size", formatBytes(ds.Binary.Size))
	} else if ds.Binary.Status == state.StatusError {
		body = RenderStatusField("Status", "Error", StatusFail)
	}
	return RenderPanel("Binary", body, w)
}

func renderGitPanel(ds *state.Dashboard, w int) string {
	body := statusLine(ds.Git.Status)
	if ds.Git.Status == state.StatusDone {
		total := len(ds.Git.Modified) + len(ds.Git.Added) + len(ds.Git.Deleted) + len(ds.Git.Other)
		if total == 0 {
			body = StatusPass.Render("Clean")
		} else {
			body = ""
			if len(ds.Git.Modified) > 0 {
				body += RenderField("Modified", fmt.Sprintf("%d", len(ds.Git.Modified))) + "\n"
			}
			if len(ds.Git.Added) > 0 {
				body += RenderField("Added", fmt.Sprintf("%d", len(ds.Git.Added))) + "\n"
			}
			if len(ds.Git.Deleted) > 0 {
				body += RenderField("Deleted", fmt.Sprintf("%d", len(ds.Git.Deleted))) + "\n"
			}
			if len(ds.Git.Other) > 0 {
				body += RenderField("Untracked", fmt.Sprintf("%d", len(ds.Git.Other))) + "\n"
			}
			// Trim trailing newline
			if len(body) > 0 && body[len(body)-1] == '\n' {
				body = body[:len(body)-1]
			}
		}
	}
	return RenderPanel("Git", body, w)
}

func renderDepsPanel(ds *state.Dashboard, w int) string {
	body := statusLine(ds.Deps.Status)
	if ds.Deps.Status == state.StatusDone {
		n := len(ds.Deps.Deps)
		body = RenderField("Modules", fmt.Sprintf("%d", n))
		for i, d := range ds.Deps.Deps {
			if i >= 6 {
				body += "\n" + lipgloss.NewStyle().Foreground(ColorDim).Render(
					fmt.Sprintf("  … and %d more", n-6),
				)
				break
			}
			body += "\n  " + lipgloss.NewStyle().Foreground(ColorDim).Render(truncate(d, w-6))
		}
	}
	return RenderPanel("Dependencies", body, w)
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func statusLine(s state.Status) string {
	switch s {
	case state.StatusRunning:
		return StatusWarn.Render("Running…")
	case state.StatusError:
		return StatusFail.Render("Error")
	default:
		return StatusIdle.Render("—")
	}
}

func coverageStyle(pct float64) lipgloss.Style {
	switch {
	case pct >= 80:
		return StatusPass
	case pct >= 60:
		return StatusWarn
	default:
		return StatusFail
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max < 4 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func formatBytes(b int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case b >= gb:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(gb))
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(mb))
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(kb))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
