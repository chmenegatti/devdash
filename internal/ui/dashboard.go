// Package ui — dashboard.go renders the K9s-inspired dashboard view.
package ui

import (
	"fmt"
	"strings"

	"github.com/cesar/devdash/internal/state"
	"github.com/charmbracelet/lipgloss"
)

// ── Main dashboard renderer ─────────────────────────────────────────────────

// RenderDashboard composes the K9s-style terminal dashboard.
func RenderDashboard(ds *state.Dashboard, width, height int) string {
	// ── Header: logo + info bar ──────────────────────────────────
	header := renderHeader(ds, width)

	// ── Breadcrumbs ──────────────────────────────────────────────
	crumbs := RenderCrumbs("Dashboard", ds.ProjectName)

	// ── Stat tiles row (compact metrics) ─────────────────────────
	statsRow := renderStatsRow(ds, width)

	// ── Main content: responsive uniform grid ────────────────────
	renderers := []func(*state.Dashboard, int) string{
		renderTestsSection,
		renderCoverageSection,
		renderLintSection,
		renderBenchSection,
		renderBinarySection,
		renderGitSection,
		renderDepsSection,
	}

	cols, panelW := dashboardGridConfig(width, height, len(renderers))
	if cols < 1 {
		cols = 1
	}

	columns := make([][]string, cols)
	for i, render := range renderers {
		colIdx := i % cols
		columns[colIdx] = append(columns[colIdx], render(ds, panelW))
	}

	columnViews := make([]string, 0, cols)
	for _, panels := range columns {
		if len(panels) == 0 {
			continue
		}
		columnViews = append(columnViews, lipgloss.JoinVertical(lipgloss.Left, panels...))
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, append([]string{}, columnViews...)...)

	// ── Command bar (bottom) ─────────────────────────────────────
	cmdBar1 := RenderCommandBar([]KeyBinding{
		{Key: "t", Desc: "tests"},
		{Key: "c", Desc: "cover"},
		{Key: "l", Desc: "lint"},
		{Key: "b", Desc: "bench"},
		{Key: "s", Desc: "build"},
		{Key: "g", Desc: "git"},
		{Key: "d", Desc: "deps"},
		{Key: "m", Desc: "report"},
		{Key: "r", Desc: "reset"},
		{Key: "q", Desc: "quit"},
	}, width)
	cmdBar2 := RenderCommandBar([]KeyBinding{
		{Key: "T", Desc: "tests detail"},
		{Key: "L", Desc: "lint detail"},
		{Key: "B", Desc: "bench detail"},
		{Key: "G", Desc: "git detail"},
		{Key: "D", Desc: "deps detail"},
	}, width)

	// ── Compose ──────────────────────────────────────────────────
	sep := SepStyle.Render(strings.Repeat("─", clamp(width, 0, width)))

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		crumbs,
		sep,
		statsRow,
		sep,
		content,
		sep,
		cmdBar1,
		cmdBar2,
	)
}

// dashboardGridConfig chooses a responsive panel grid based on terminal size.
func dashboardGridConfig(termWidth, termHeight, panelCount int) (cols int, panelWidth int) {
	if termWidth <= 0 {
		termWidth = 120
	}
	if termHeight <= 0 {
		termHeight = 40
	}

	cols = 1
	if termWidth >= 170 {
		cols = 3
	} else if termWidth >= 110 {
		cols = 2
	}

	// If height is tight, spread panels across more columns when possible.
	if termHeight < 32 && termWidth >= 170 {
		cols = 3
	}

	if cols > panelCount {
		cols = panelCount
	}

	gap := 1
	panelWidth = (termWidth - (cols-1)*gap) / cols
	if panelWidth < 30 {
		panelWidth = 30
	}
	if panelWidth > 58 {
		panelWidth = 58
	}

	return cols, panelWidth
}

// ── Header ──────────────────────────────────────────────────────────────────

func renderHeader(ds *state.Dashboard, width int) string {
	logo := LogoStyle.Render("⎈ devdash")
	ver := lipgloss.NewStyle().Foreground(ColorDim).Render("v0.1.0")

	left := logo + " " + ver

	// Right side: project path
	right := lipgloss.NewStyle().Foreground(ColorDim).Render(ds.ProjectDir)

	gap := width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	return HeaderBarStyle.Width(width).Render(
		left + strings.Repeat(" ", gap) + right,
	)
}

// ── Top stats row ───────────────────────────────────────────────────────────

func renderStatsRow(ds *state.Dashboard, width int) string {
	chips := []string{}

	// Tests status
	chips = append(chips, statWithDot("Tests", ds.Tests.Status, func() string {
		if ds.Tests.Passed {
			return "PASS"
		}
		return "FAIL"
	}))

	// Coverage
	chips = append(chips, statWithDot("Cover", ds.Coverage.Status, func() string {
		return fmt.Sprintf("%.0f%%", ds.Coverage.Percentage)
	}))

	// Lint
	chips = append(chips, statWithDot("Lint", ds.Lint.Status, func() string {
		n := len(ds.Lint.Issues)
		if n == 0 {
			return "OK"
		}
		return fmt.Sprintf("%d", n)
	}))

	// Binary
	chips = append(chips, statWithDot("Binary", ds.Binary.Status, func() string {
		return formatBytes(ds.Binary.Size)
	}))

	// Benchmarks
	chips = append(chips, statWithDot("Bench", ds.Benchmarks.Status, func() string {
		return fmt.Sprintf("%d", len(ds.Benchmarks.Entries))
	}))

	// Git
	chips = append(chips, statWithDot("Git", ds.Git.Status, func() string {
		total := len(ds.Git.Modified) + len(ds.Git.Added) + len(ds.Git.Deleted) + len(ds.Git.Other)
		if total == 0 {
			return "clean"
		}
		return fmt.Sprintf("%d changes", total)
	}))

	// Deps
	chips = append(chips, statWithDot("Deps", ds.Deps.Status, func() string {
		return fmt.Sprintf("%d", len(ds.Deps.Deps))
	}))

	return "  " + strings.Join(chips, "    ")
}

func statWithDot(label string, s state.Status, valueFn func() string) string {
	dot := "○"
	dotStyle := lipgloss.NewStyle().Foreground(ColorSubtle)
	valStr := "–"

	switch s {
	case state.StatusRunning:
		dot = "◍"
		dotStyle = lipgloss.NewStyle().Foreground(ColorWarning)
		valStr = "…"
	case state.StatusDone:
		dot = "●"
		dotStyle = lipgloss.NewStyle().Foreground(ColorSuccess)
		valStr = valueFn()
	case state.StatusError:
		dot = "●"
		dotStyle = lipgloss.NewStyle().Foreground(ColorDanger)
		valStr = "err"
	}

	return fmt.Sprintf("%s %s %s",
		dotStyle.Render(dot),
		InfoStyle.Render(label),
		InfoValueStyle.Render(valStr),
	)
}

// ── Tests section ───────────────────────────────────────────────────────────

func renderTestsSection(ds *state.Dashboard, w int) string {
	var body string
	switch ds.Tests.Status {
	case state.StatusDone:
		pass := StatusPass.Render("✓ PASS")
		if !ds.Tests.Passed {
			pass = StatusFail.Render("✗ FAIL")
		}
		testsValue := fmt.Sprintf("%d", ds.Tests.TestCases)
		if ds.Tests.TestCases == 0 {
			testsValue = "-"
		}
		body = fmt.Sprintf("  %s  %s  %s",
			pass,
			StatChip("tests", testsValue),
			StatChip("dur", ds.Tests.Duration.String()),
		)
		if ds.Tests.FailedTests > 0 {
			body += "\n  " + StatusFail.Render(fmt.Sprintf("%d failed", ds.Tests.FailedTests))
		}
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error: "+truncate(ds.Tests.Err, w-10))
	default:
		body = "  " + StatusIdle.Render("○ idle — press <t>")
	}
	return RenderSection("Tests", body, w)
}

func renderCoverageSection(ds *state.Dashboard, w int) string {
	var body string
	switch ds.Coverage.Status {
	case state.StatusDone:
		pct := ds.Coverage.Percentage
		style := StatusFail
		if pct >= 80 {
			style = StatusPass
		} else if pct >= 60 {
			style = StatusWarn
		}
		body = fmt.Sprintf("  %s %s",
			style.Render(fmt.Sprintf("%.1f%%", pct)),
			StatChip("target", "80%+"),
		)
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error: "+truncate(ds.Coverage.Err, w-10))
	default:
		body = "  " + StatusIdle.Render("○ idle — press <c>")
	}
	return RenderSection("Coverage", body, w)
}

// ── Lint section ────────────────────────────────────────────────────────────

func renderLintSection(ds *state.Dashboard, w int) string {
	var body string
	switch ds.Lint.Status {
	case state.StatusDone:
		n := len(ds.Lint.Issues)
		if n == 0 {
			body = "  " + StatusPass.Render("✓ No issues")
		} else {
			header := StatusWarn.Render(fmt.Sprintf("  ▲ %d issues", n))
			var lines []string
			for i, iss := range ds.Lint.Issues {
				if i >= 4 {
					lines = append(lines,
						lipgloss.NewStyle().Foreground(ColorDim).Render(
							fmt.Sprintf("    … and %d more", n-4)),
					)
					break
				}
				lines = append(lines,
					lipgloss.NewStyle().Foreground(ColorWarning).Render(
						fmt.Sprintf("  %d. %s", i+1, truncate(iss, w-8))),
				)
			}
			body = header + "\n" + strings.Join(lines, "\n")
		}
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error")
	default:
		body = "  " + StatusIdle.Render("○ idle — press <l>")
	}
	return RenderSection("Lint", body, w)
}

// ── Benchmarks section ──────────────────────────────────────────────────────

func renderBenchSection(ds *state.Dashboard, w int) string {
	var body string
	switch ds.Benchmarks.Status {
	case state.StatusDone:
		if len(ds.Benchmarks.Entries) == 0 {
			body = "  " + StatusIdle.Render("No benchmarks found")
		} else {
			cols := []TableColumn{
				{Header: "NAME", Width: w - 28},
				{Header: "ITERS", Width: 10},
				{Header: "NS/OP", Width: 12},
			}
			header := RenderTableHeader(cols)
			var rows []string
			for i, e := range ds.Benchmarks.Entries {
				if i >= 5 {
					break
				}
				rows = append(rows, RenderTableRow(
					[]string{
						truncate(e.Name, w-30),
						fmt.Sprintf("%d", e.Iterations),
						fmt.Sprintf("%.1f", e.NsPerOp),
					}, cols, i%2 == 1))
			}
			body = header + "\n" + strings.Join(rows, "\n")
		}
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error: "+truncate(ds.Benchmarks.Err, w-10))
	default:
		body = "  " + StatusIdle.Render("○ idle — press <b>")
	}
	return RenderSection("Benchmarks", body, w)
}

func renderBinarySection(ds *state.Dashboard, w int) string {
	var body string
	switch ds.Binary.Status {
	case state.StatusDone:
		body = fmt.Sprintf("  %s  %s",
			StatusPass.Render(formatBytes(ds.Binary.Size)),
			StatChip("build", "ok"),
		)
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error: "+truncate(ds.Binary.Err, w-10))
	default:
		body = "  " + StatusIdle.Render("○ idle — press <s>")
	}
	return RenderSection("Binary", body, w)
}

// ── Git section ─────────────────────────────────────────────────────────────

func renderGitSection(ds *state.Dashboard, w int) string {
	var body string
	switch ds.Git.Status {
	case state.StatusDone:
		total := len(ds.Git.Modified) + len(ds.Git.Added) + len(ds.Git.Deleted) + len(ds.Git.Other)
		if total == 0 {
			body = "  " + StatusPass.Render("✓ Working tree clean")
		} else {
			var lines []string
			if len(ds.Git.Modified) > 0 {
				lines = append(lines, fmt.Sprintf("  %s %s",
					StatusWarn.Render(fmt.Sprintf("M %d", len(ds.Git.Modified))),
					lipgloss.NewStyle().Foreground(ColorDim).Render("modified"),
				))
			}
			if len(ds.Git.Added) > 0 {
				lines = append(lines, fmt.Sprintf("  %s %s",
					StatusPass.Render(fmt.Sprintf("A %d", len(ds.Git.Added))),
					lipgloss.NewStyle().Foreground(ColorDim).Render("added"),
				))
			}
			if len(ds.Git.Deleted) > 0 {
				lines = append(lines, fmt.Sprintf("  %s %s",
					StatusFail.Render(fmt.Sprintf("D %d", len(ds.Git.Deleted))),
					lipgloss.NewStyle().Foreground(ColorDim).Render("deleted"),
				))
			}
			if len(ds.Git.Other) > 0 {
				lines = append(lines, fmt.Sprintf("  %s %s",
					lipgloss.NewStyle().Foreground(ColorDim).Bold(true).Render(fmt.Sprintf("? %d", len(ds.Git.Other))),
					lipgloss.NewStyle().Foreground(ColorDim).Render("untracked"),
				))
			}
			body = strings.Join(lines, "\n")
		}
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error: "+truncate(ds.Git.Err, w-10))
	default:
		body = "  " + StatusIdle.Render("○ idle — press <g>")
	}
	return RenderSection("Git", body, w)
}

// ── Dependencies section ────────────────────────────────────────────────────

func renderDepsSection(ds *state.Dashboard, w int) string {
	var body string
	switch ds.Deps.Status {
	case state.StatusDone:
		n := len(ds.Deps.Deps)
		summary := StatusPass.Render(fmt.Sprintf("  %d modules", n))
		var lines []string
		for i, d := range ds.Deps.Deps {
			if i >= 6 {
				lines = append(lines,
					lipgloss.NewStyle().Foreground(ColorDim).Render(
						fmt.Sprintf("  … and %d more", n-6)),
				)
				break
			}
			lines = append(lines,
				lipgloss.NewStyle().Foreground(ColorFg).Render("  "+truncate(d, w-4)),
			)
		}
		body = summary + "\n" + strings.Join(lines, "\n")
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error")
	default:
		body = "  " + StatusIdle.Render("○ idle — press <d>")
	}
	return RenderSection("Dependencies", body, w)
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func truncate(s string, max int) string {
	if max < 4 {
		max = 4
	}
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "…"
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
