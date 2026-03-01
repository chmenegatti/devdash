// Package ui — detail_views.go renders K9s-style full-screen detail views.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/chmenegatti/devdash/internal/state"
)

// DetailContent represents the text payload for a detail screen.
type DetailContent struct {
	Title   string
	Crumb   string
	Summary string
	Body    string
}

// ── Shared detail layout ────────────────────────────────────────────────────

func RenderDetailFrame(title string, crumb string, summary string, body string, width, height int) string {
	// Header
	header := HeaderBarStyle.Width(width).Render(
		LogoStyle.Render("⎈ devdash") + "  " +
			lipgloss.NewStyle().Foreground(ColorDim).Render("v0.1.0"),
	)

	// Crumbs
	crumbs := RenderCrumbs("Dashboard", crumb)

	sep := SepStyle.Render(strings.Repeat("─", clamp(width, 0, width)))

	// Body with padding
	bodyStyle := lipgloss.NewStyle().
		Width(width-4).
		Foreground(ColorFg).
		Padding(0, 1)

	// Command bar
	cmdBar := RenderCommandBar([]KeyBinding{
		{Key: "↑/k", Desc: "up"},
		{Key: "↓/j", Desc: "down"},
		{Key: "pgup/pgdn", Desc: "page"},
		{Key: "g/G", Desc: "top/end"},
		{Key: "backspace", Desc: "back"},
		{Key: "q", Desc: "quit"},
	}, width)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		crumbs,
		sep,
		summary,
		sep,
		bodyStyle.Render(body),
		sep,
		cmdBar,
	)
}

// DetailBodySize returns viewport dimensions for detail content.
func DetailBodySize(width, height int) (bodyWidth int, bodyHeight int) {
	bodyWidth = clamp(width-6, 20, 500)
	bodyHeight = clamp(height-9, 3, 200)
	return bodyWidth, bodyHeight
}

// ── Tests detail ────────────────────────────────────────────────────────────

// RenderTestsDetail renders a full-screen view of test results.
func RenderTestsDetail(ds *state.Dashboard, width, height int) string {
	detail := BuildTestsDetail(ds, width)
	return RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, detail.Body, width, height)
}

// BuildTestsDetail builds summary and body content for tests detail view.
func BuildTestsDetail(ds *state.Dashboard, width int) DetailContent {
	var summary string
	switch ds.Tests.Status {
	case state.StatusDone:
		pass := StatusPass.Render("✓ PASS")
		if !ds.Tests.Passed {
			pass = StatusFail.Render("✗ FAIL")
		}
		summary = fmt.Sprintf("  %s  %s  %s",
			pass,
			StatChip("packages", fmt.Sprintf("%d", ds.Tests.Packages)),
			StatChip("duration", ds.Tests.Duration.String()),
		)
	case state.StatusRunning:
		summary = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		summary = "  " + StatusFail.Render("● Error")
	default:
		summary = "  " + StatusIdle.Render("○ No test run yet. Press <t> to run.")
	}

	output := ds.Tests.Output
	if ds.Tests.Err != "" {
		output += "\n" + StatusFail.Render("Error: "+ds.Tests.Err)
	}
	if output == "" {
		output = StatusIdle.Render("(no output)")
	}

	return DetailContent{
		Title:   "Tests",
		Crumb:   "Tests",
		Summary: summary,
		Body:    output,
	}
}

// ── Lint detail ─────────────────────────────────────────────────────────────

// RenderLintDetail renders a full-screen view of lint results.
func RenderLintDetail(ds *state.Dashboard, width, height int) string {
	detail := BuildLintDetail(ds, width)
	return RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, detail.Body, width, height)
}

// BuildLintDetail builds summary and body content for lint detail view.
func BuildLintDetail(ds *state.Dashboard, width int) DetailContent {
	var summary string
	switch ds.Lint.Status {
	case state.StatusDone:
		n := len(ds.Lint.Issues)
		if n == 0 {
			summary = "  " + StatusPass.Render("✓ No lint issues")
		} else {
			summary = "  " + StatusWarn.Render(fmt.Sprintf("▲ %d issues", n))
		}
	case state.StatusRunning:
		summary = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		summary = "  " + StatusFail.Render("● Error: "+ds.Lint.Err)
	default:
		summary = "  " + StatusIdle.Render("○ No lint run yet. Press <l> to run.")
	}

	var body string
	if len(ds.Lint.Issues) > 0 {
		issueStyle := lipgloss.NewStyle().Foreground(ColorWarning)
		var sb strings.Builder
		for i, iss := range ds.Lint.Issues {
			fmt.Fprintf(&sb, "%s %s\n",
				LabelStyle.Render(fmt.Sprintf("%3d.", i+1)),
				issueStyle.Render(iss),
			)
		}
		body = sb.String()
	} else if ds.Lint.Status == state.StatusDone {
		body = StatusPass.Render("✓ Clean — no lint issues found")
	} else {
		body = StatusIdle.Render("(no output)")
	}

	return DetailContent{
		Title:   "Lint",
		Crumb:   "Lint",
		Summary: summary,
		Body:    body,
	}
}

// ── Bench detail ────────────────────────────────────────────────────────────

// RenderBenchDetail renders a full-screen view of benchmark results.
func RenderBenchDetail(ds *state.Dashboard, width, height int) string {
	detail := BuildBenchDetail(ds, width)
	return RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, detail.Body, width, height)
}

// BuildBenchDetail builds summary and body content for benchmark detail view.
func BuildBenchDetail(ds *state.Dashboard, width int) DetailContent {
	var summary string
	switch ds.Benchmarks.Status {
	case state.StatusDone:
		n := len(ds.Benchmarks.Entries)
		if n == 0 {
			summary = "  " + StatusIdle.Render("No benchmarks found")
		} else {
			summary = "  " + StatusPass.Render(fmt.Sprintf("● %d benchmarks", n))
		}
	case state.StatusRunning:
		summary = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		summary = "  " + StatusFail.Render("● Error: "+ds.Benchmarks.Err)
	default:
		summary = "  " + StatusIdle.Render("○ No benchmark run yet. Press <b> to run.")
	}

	var body string
	if len(ds.Benchmarks.Entries) > 0 {
		cols := []TableColumn{
			{Header: "#", Width: 5},
			{Header: "NAME", Width: clamp(width-40, 20, 60)},
			{Header: "ITERATIONS", Width: 12},
			{Header: "NS/OP", Width: 12},
		}
		var sb strings.Builder
		sb.WriteString(RenderTableHeader(cols) + "\n")
		for i, e := range ds.Benchmarks.Entries {
			sb.WriteString(RenderTableRow(
				[]string{
					fmt.Sprintf("%d", i+1),
					e.Name,
					fmt.Sprintf("%d", e.Iterations),
					fmt.Sprintf("%.1f", e.NsPerOp),
				}, cols, i%2 == 1) + "\n")
		}
		body = sb.String()
	} else {
		body = StatusIdle.Render("(no output)")
	}

	return DetailContent{
		Title:   "Benchmarks",
		Crumb:   "Benchmarks",
		Summary: summary,
		Body:    body,
	}
}

// ── Deps detail ─────────────────────────────────────────────────────────────

// RenderDepsDetail renders a full-screen view of module dependencies.
func RenderDepsDetail(ds *state.Dashboard, width, height int) string {
	detail := BuildDepsDetail(ds, width)
	return RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, detail.Body, width, height)
}

// BuildDepsDetail builds summary and body content for dependencies detail view.
func BuildDepsDetail(ds *state.Dashboard, width int) DetailContent {
	var summary string
	switch ds.Deps.Status {
	case state.StatusDone:
		summary = "  " + StatusPass.Render(fmt.Sprintf("● %d modules", len(ds.Deps.Deps)))
	case state.StatusRunning:
		summary = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		summary = "  " + StatusFail.Render("● Error: "+ds.Deps.Err)
	default:
		summary = "  " + StatusIdle.Render("○ Press <d> to list dependencies.")
	}

	var body string
	if len(ds.Deps.Deps) > 0 {
		cols := []TableColumn{
			{Header: "#", Width: 5},
			{Header: "MODULE", Width: clamp(width-12, 20, 80)},
		}
		var sb strings.Builder
		sb.WriteString(RenderTableHeader(cols) + "\n")
		for i, d := range ds.Deps.Deps {
			sb.WriteString(RenderTableRow(
				[]string{fmt.Sprintf("%d", i+1), d},
				cols, i%2 == 1) + "\n")
		}
		body = sb.String()
	} else {
		body = StatusIdle.Render("(no output)")
	}

	return DetailContent{
		Title:   "Dependencies",
		Crumb:   "Dependencies",
		Summary: summary,
		Body:    body,
	}
}

// ── Git detail ──────────────────────────────────────────────────────────────

// RenderGitDetail renders a full-screen view of git status.
func RenderGitDetail(ds *state.Dashboard, width, height int) string {
	detail := BuildGitDetail(ds, width)
	return RenderDetailFrame(detail.Title, detail.Crumb, detail.Summary, detail.Body, width, height)
}

// BuildGitDetail builds summary and body content for git detail view.
func BuildGitDetail(ds *state.Dashboard, width int) DetailContent {
	var summary string
	switch ds.Git.Status {
	case state.StatusDone:
		total := len(ds.Git.Modified) + len(ds.Git.Added) + len(ds.Git.Deleted) + len(ds.Git.Other)
		if total == 0 {
			summary = "  " + StatusPass.Render("✓ Working tree clean")
		} else {
			summary = "  " + StatusWarn.Render(fmt.Sprintf("● %d changes", total))
		}
	case state.StatusRunning:
		summary = "  " + StatusWarn.Render("◍ Running…")
	case state.StatusError:
		summary = "  " + StatusFail.Render("● Error: "+ds.Git.Err)
	default:
		summary = "  " + StatusIdle.Render("○ Press <g> to check git status.")
	}

	var sb strings.Builder
	modStyle := lipgloss.NewStyle().Foreground(ColorWarning)
	addStyle := lipgloss.NewStyle().Foreground(ColorSuccess)
	delStyle := lipgloss.NewStyle().Foreground(ColorDanger)
	otherStyle := lipgloss.NewStyle().Foreground(ColorDim)

	renderSection := func(prefix string, label string, files []string, style lipgloss.Style) {
		if len(files) == 0 {
			return
		}
		sb.WriteString(LabelStyle.Render(fmt.Sprintf("  %s (%d)", label, len(files))) + "\n")
		for _, f := range files {
			fmt.Fprintf(&sb, "  %s %s\n", style.Render(prefix), style.Render(f))
		}
		sb.WriteString("\n")
	}

	renderSection("M", "Modified", ds.Git.Modified, modStyle)
	renderSection("A", "Added", ds.Git.Added, addStyle)
	renderSection("D", "Deleted", ds.Git.Deleted, delStyle)
	renderSection("?", "Untracked", ds.Git.Other, otherStyle)

	body := sb.String()
	if body == "" && ds.Git.Status == state.StatusDone {
		body = StatusPass.Render("✓ Working tree clean")
	} else if body == "" {
		body = StatusIdle.Render("(no output)")
	}

	return DetailContent{
		Title:   "Git Status",
		Crumb:   "Git",
		Summary: summary,
		Body:    body,
	}
}
