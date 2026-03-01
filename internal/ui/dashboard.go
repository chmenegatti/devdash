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

	// ── Main content: two-column panels ──────────────────────────
	pw := panelWidth(width)

	leftCol := lipgloss.JoinVertical(lipgloss.Left,
		renderTestsSection(ds, pw),
		renderLintSection(ds, pw),
		renderBenchSection(ds, pw),
	)
	rightCol := lipgloss.JoinVertical(lipgloss.Left,
		renderGitSection(ds, pw),
		renderDepsSection(ds, pw),
	)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, " ", rightCol)

	// ── Command bar (bottom) ─────────────────────────────────────
	cmdBar1 := RenderCommandBar([]KeyBinding{
		{Key: "t", Desc: "tests"},
		{Key: "c", Desc: "cover"},
		{Key: "l", Desc: "lint"},
		{Key: "b", Desc: "bench"},
		{Key: "s", Desc: "build"},
		{Key: "g", Desc: "git"},
		{Key: "d", Desc: "deps"},
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

// panelWidth adapts to terminal width for two-column layout.
func panelWidth(termWidth int) int {
	if termWidth <= 0 {
		return 40
	}
	pw := (termWidth - 3) / 2
	if pw < 30 {
		pw = 30
	}
	if pw > 60 {
		pw = 60
	}
	return pw
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
		body = fmt.Sprintf("  %s  %s  %s",
			pass,
			StatChip("pkg", fmt.Sprintf("%d", ds.Tests.Packages)),
			StatChip("dur", ds.Tests.Duration.String()),
		)
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error: "+truncate(ds.Tests.Err, w-10))
	default:
		body = "  " + StatusIdle.Render("○ idle — press <t>")
	}
	return RenderSection("Tests", body, w)
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
