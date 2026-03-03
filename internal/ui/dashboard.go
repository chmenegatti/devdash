// Package ui — dashboard.go renders the K9s-inspired dashboard view.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/chmenegatti/devdash/internal/state"
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
	noticeLine := renderNotice(ds, width)

	// ── Main content: responsive uniform grid ────────────────────
	// Calculate available height for the grid
	// Header(1) + Crumbs(1) + Sep(1) + Stats(1) + Notice(0 or 1) + Sep(1) + Grid(H) + Sep(1) + Bar1(1) + Bar2(1)
	fixedElementsHeight := 1 + 1 + 1 + 1 + 1 + 1 + 1
	if noticeLine != "" {
		fixedElementsHeight++
	}

	gridHeight := height - fixedElementsHeight
	if gridHeight < 10 {
		gridHeight = 10 // minimum fallback
	}

	renderers := []func(*state.Dashboard, int, int) string{
		renderTestsSection,
		renderCoverageSection,
		renderLintSection,
		renderBenchSection,
		renderBinarySection,
		renderGitSection,
		renderDepsSection,
		renderProfileSection,
	}

	cols, rows, panelW, panelH := dashboardGridConfig(width, gridHeight, len(renderers))

	columns := make([][]string, cols)
	for i, render := range renderers {
		colIdx := i % cols

		// Calculate h avoiding extrapolation due to integer division, minus borders
		h := panelH
		rowIdx := i / cols
		if rowIdx == rows-1 {
			// for the last row, give it the remainder, but safely subtract borderHeight
			h += (gridHeight - (rows * (panelH + 2)))
		}

		columns[colIdx] = append(columns[colIdx], render(ds, panelW, h))
	}

	columnViews := make([]string, 0, cols)
	for _, panels := range columns {
		if len(panels) == 0 {
			continue
		}
		columnViews = append(columnViews, lipgloss.JoinVertical(lipgloss.Left, panels...))
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, append([]string{}, columnViews...)...)

	// Determine if we need to pad the content area to push the footer down
	contentHeight := lipgloss.Height(content)
	if contentHeight < gridHeight {
		padding := strings.Repeat("\n", gridHeight-contentHeight)
		content = lipgloss.JoinVertical(lipgloss.Left, content, padding)
	}

	// ── Command bar (bottom) ─────────────────────────────────────
	cmdBar1 := RenderCommandBar([]KeyBinding{
		{Key: "t", Desc: "tests"},
		{Key: "c", Desc: "cover"},
		{Key: "l", Desc: "lint"},
		{Key: "b", Desc: "bench"},
		{Key: "s", Desc: "build"},
		{Key: "g", Desc: "git"},
		{Key: "d", Desc: "deps"},
		{Key: "p", Desc: "profile"},
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
		{Key: "P", Desc: "profile detail"},
	}, width)

	// ── Compose ──────────────────────────────────────────────────
	sep := SepStyle.Render(strings.Repeat("─", clamp(width, 0, width)))

	var topParts []string
	topParts = append(topParts, header, crumbs, sep, statsRow)
	if noticeLine != "" {
		topParts = append(topParts, noticeLine)
	}
	topParts = append(topParts, sep, content, sep, cmdBar1, cmdBar2)

	return lipgloss.JoinVertical(lipgloss.Left, topParts...)
}

func renderNotice(ds *state.Dashboard, width int) string {
	if strings.TrimSpace(ds.Notice) == "" {
		return ""
	}

	style := lipgloss.NewStyle().
		Width(width).
		Padding(0, 1).
		Foreground(ColorFg)

	if strings.HasPrefix(ds.Notice, "❌") {
		style = style.Foreground(ColorDanger)
	} else if strings.HasPrefix(ds.Notice, "📝") || strings.HasPrefix(ds.Notice, "🛠️") {
		style = style.Foreground(ColorAccent)
	}

	return style.Render(ds.Notice)
}

// dashboardGridConfig chooses a responsive panel grid based on terminal size.
func dashboardGridConfig(termWidth, gridHeight, panelCount int) (cols, rows, panelWidth, panelHeight int) {
	if termWidth <= 0 {
		termWidth = 120
	}
	if gridHeight <= 0 {
		gridHeight = 40
	}

	// Determine Columns
	cols = 1
	if termWidth >= 170 {
		cols = 3
	} else if termWidth >= 110 {
		cols = 2
	}

	if cols > panelCount {
		cols = panelCount
	}

	// Calculate Rows needed
	rows = panelCount / cols
	if panelCount%cols != 0 {
		rows++
	}

	// Calculate Widths
	gap := 1
	panelWidth = (termWidth - (cols-1)*gap) / cols
	if panelWidth < 30 {
		panelWidth = 30
	}

	// Calculate Heights
	// The SectionBorder uses RoundedBorder which adds 2 to the height (top+bottom).
	// Let's explicitly subtract the border size and any gap to ensure it fits.
	borderHeight := 2
	totalBorderHeight := rows * borderHeight
	panelHeight = (gridHeight - totalBorderHeight) / rows

	if panelHeight < 5 {
		panelHeight = 5
	}

	return cols, rows, panelWidth, panelHeight
}

// ── Header ──────────────────────────────────────────────────────────────────

func renderHeader(ds *state.Dashboard, width int) string {
	logo := LogoStyle.Render("⎈ devdash")
	ver := lipgloss.NewStyle().Foreground(ColorDim).Render(ds.Version)

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

	// Profile
	chips = append(chips, statWithDot("Profile", ds.Profile.Status, func() string {
		if ds.Profile.TotalSamples > 0 {
			return fmt.Sprintf("%d", ds.Profile.TotalSamples)
		}
		return "ok"
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

func renderTestsSection(ds *state.Dashboard, w, h int) string {
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
	return RenderSection("Tests", body, w, h)
}

func renderCoverageSection(ds *state.Dashboard, w, h int) string {
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
	return RenderSection("Coverage", body, w, h)
}

// ── Lint section ────────────────────────────────────────────────────────────

func renderLintSection(ds *state.Dashboard, w, h int) string {
	var body string
	switch ds.Lint.Status {
	case state.StatusDone:
		n := len(ds.Lint.Issues)
		if n == 0 {
			body = "  " + StatusPass.Render("✓ No issues")
		} else {
			header := StatusWarn.Render(fmt.Sprintf("  ▲ %d issues", n))
			var lines []string

			// Calculate how many issues we can show based on panel height (-3 for borders/header)
			maxLines := h - 3
			if maxLines < 1 {
				maxLines = 1
			}

			for i, iss := range ds.Lint.Issues {
				if i >= maxLines-1 && i < len(ds.Lint.Issues)-1 {
					lines = append(lines,
						lipgloss.NewStyle().Foreground(ColorDim).Render(
							fmt.Sprintf("    … and %d more", n-i)),
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
	return RenderSection("Lint", body, w, h)
}

// ── Benchmarks section ──────────────────────────────────────────────────────

func renderBenchSection(ds *state.Dashboard, w, h int) string {
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

			// Dynamic row count based on height
			maxRows := h - 4
			if maxRows < 1 {
				maxRows = 1
			}

			for i, e := range ds.Benchmarks.Entries {
				if i >= maxRows {
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
	return RenderSection("Benchmarks", body, w, h)
}

func renderBinarySection(ds *state.Dashboard, w, h int) string {
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
	return RenderSection("Binary", body, w, h)
}

// ── Git section ─────────────────────────────────────────────────────────────

func renderGitSection(ds *state.Dashboard, w, h int) string {
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
	return RenderSection("Git", body, w, h)
}

// ── Dependencies section ────────────────────────────────────────────────────

func renderDepsSection(ds *state.Dashboard, w, h int) string {
	var body string
	switch ds.Deps.Status {
	case state.StatusDone:
		n := len(ds.Deps.Deps)
		summary := StatusPass.Render(fmt.Sprintf("  %d modules", n))
		var lines []string

		maxLines := h - 3
		if maxLines < 1 {
			maxLines = 1
		}

		for i, d := range ds.Deps.Deps {
			if i >= maxLines-1 && i < len(ds.Deps.Deps)-1 {
				lines = append(lines,
					lipgloss.NewStyle().Foreground(ColorDim).Render(
						fmt.Sprintf("  … and %d more", n-i)),
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
	return RenderSection("Dependencies", body, w, h)
}

func renderProfileSection(ds *state.Dashboard, w, h int) string {
	var body string
	switch ds.Profile.Status {
	case state.StatusDone:
		header := fmt.Sprintf("  %s %s",
			StatusPass.Render("✓ flamegraph ready"),
			StatChip("pkg", truncate(ds.Profile.TargetPackage, w-24)),
		)
		preview := strings.Split(strings.TrimSpace(ds.Profile.Flamegraph), "\n")

		maxLines := h - 4
		if maxLines < 1 {
			maxLines = 1
		}

		if len(preview) > maxLines {
			preview = preview[:maxLines]
		}
		if len(preview) == 0 {
			body = header
		} else {
			body = header + "\n" + strings.Join(preview, "\n")
		}
	case state.StatusRunning:
		body = "  " + StatusWarn.Render("◍ Profiling with go test -cpuprofile…")
	case state.StatusError:
		body = "  " + StatusFail.Render("● Error: "+truncate(ds.Profile.Err, w-10))
	default:
		body = "  " + StatusIdle.Render("○ idle — press <p>")
	}
	return RenderSection("CPU Profile", body, w, h)
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
