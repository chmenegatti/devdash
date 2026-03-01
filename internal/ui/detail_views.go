// Package ui — detail_views.go renders full-screen detail views for module outputs.
package ui

import (
	"fmt"
	"strings"

	"github.com/cesar/devdash/internal/state"
	"github.com/charmbracelet/lipgloss"
)

// detailHeader renders the top bar for a detail view.
func detailHeader(title string, width int) string {
	bar := TitleStyle.Render(fmt.Sprintf("  %s  ", title))
	return bar
}

// detailFooter renders the bottom help bar for a detail view.
func detailFooter() string {
	return RenderHelp([]KeyBinding{
		{Key: "backspace", Desc: "back to dashboard"},
		{Key: "q", Desc: "quit"},
	})
}

// detailBody wraps raw output in a styled full-width container.
func detailBody(content string, width, height int) string {
	style := lipgloss.NewStyle().
		Width(width - 4).
		Foreground(ColorWhite)

	return style.Render(content)
}

// RenderTestsDetail renders a full-screen view of test results.
func RenderTestsDetail(ds *state.Dashboard, width, height int) string {
	header := detailHeader("Tests — Full Output", width)

	// Summary line
	var summary string
	if ds.Tests.Status == state.StatusDone {
		passLabel := "PASS"
		passStyle := StatusPass
		if !ds.Tests.Passed {
			passLabel = "FAIL"
			passStyle = StatusFail
		}
		summary = fmt.Sprintf("%s  %s  %s",
			RenderStatusField("Status", passLabel, passStyle),
			RenderField("Packages", fmt.Sprintf("%d", ds.Tests.Packages)),
			RenderField("Duration", ds.Tests.Duration.String()),
		)
	} else if ds.Tests.Status == state.StatusRunning {
		summary = StatusWarn.Render("Running…")
	} else if ds.Tests.Status == state.StatusError {
		summary = RenderStatusField("Status", "Error", StatusFail)
	} else {
		summary = StatusIdle.Render("No test run yet. Press t to run tests.")
	}

	// Raw output
	output := ds.Tests.Output
	if ds.Tests.Err != "" {
		output += "\n" + StatusFail.Render("Error: "+ds.Tests.Err)
	}
	if output == "" {
		output = StatusIdle.Render("(no output)")
	}

	separator := SepStyle.Render(strings.Repeat("─", min(width-2, 80)))

	body := detailBody(output, width, height)
	footer := detailFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		summary,
		separator,
		"",
		body,
		"",
		footer,
	)
}

// RenderLintDetail renders a full-screen view of lint results.
func RenderLintDetail(ds *state.Dashboard, width, height int) string {
	header := detailHeader("Lint — Full Output", width)

	// Summary line
	var summary string
	if ds.Lint.Status == state.StatusDone {
		n := len(ds.Lint.Issues)
		if n == 0 {
			summary = RenderStatusField("Status", "OK — no issues", StatusPass)
		} else {
			summary = RenderStatusField("Issues", fmt.Sprintf("%d", n), StatusWarn)
		}
	} else if ds.Lint.Status == state.StatusRunning {
		summary = StatusWarn.Render("Running…")
	} else if ds.Lint.Status == state.StatusError {
		summary = RenderStatusField("Status", "Error", StatusFail) +
			"  " + StatusFail.Render(ds.Lint.Err)
	} else {
		summary = StatusIdle.Render("No lint run yet. Press l to run lint.")
	}

	separator := SepStyle.Render(strings.Repeat("─", min(width-2, 80)))

	// Full issue list
	var issueLines string
	if len(ds.Lint.Issues) > 0 {
		var sb strings.Builder
		issueStyle := lipgloss.NewStyle().Foreground(ColorWarning)
		for i, iss := range ds.Lint.Issues {
			sb.WriteString(fmt.Sprintf("%s %s\n",
				LabelStyle.Render(fmt.Sprintf("%3d.", i+1)),
				issueStyle.Render(iss),
			))
		}
		issueLines = sb.String()
	} else if ds.Lint.Status == state.StatusDone {
		issueLines = StatusPass.Render("✓ Clean — no lint issues found")
	} else {
		issueLines = StatusIdle.Render("(no output)")
	}

	body := detailBody(issueLines, width, height)
	footer := detailFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		summary,
		separator,
		"",
		body,
		"",
		footer,
	)
}

// RenderBenchDetail renders a full-screen view of benchmark results.
func RenderBenchDetail(ds *state.Dashboard, width, height int) string {
	header := detailHeader("Benchmarks — Full Output", width)

	var summary string
	if ds.Benchmarks.Status == state.StatusDone {
		n := len(ds.Benchmarks.Entries)
		if n == 0 {
			summary = StatusIdle.Render("No benchmarks found")
		} else {
			summary = RenderField("Benchmarks", fmt.Sprintf("%d", n))
		}
	} else if ds.Benchmarks.Status == state.StatusRunning {
		summary = StatusWarn.Render("Running…")
	} else if ds.Benchmarks.Status == state.StatusError {
		summary = RenderStatusField("Status", "Error", StatusFail) +
			"  " + StatusFail.Render(ds.Benchmarks.Err)
	} else {
		summary = StatusIdle.Render("No benchmark run yet. Press b to run benchmarks.")
	}

	separator := SepStyle.Render(strings.Repeat("─", min(width-2, 80)))

	var body string
	if len(ds.Benchmarks.Entries) > 0 {
		nameStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
		valStyle := lipgloss.NewStyle().Foreground(ColorWhite)
		dimStyle := lipgloss.NewStyle().Foreground(ColorDim)

		var sb strings.Builder
		for i, e := range ds.Benchmarks.Entries {
			sb.WriteString(fmt.Sprintf("%s  %s  %s  %s\n",
				LabelStyle.Render(fmt.Sprintf("%3d.", i+1)),
				nameStyle.Render(e.Name),
				valStyle.Render(fmt.Sprintf("%d iters", e.Iterations)),
				dimStyle.Render(fmt.Sprintf("%.1f ns/op", e.NsPerOp)),
			))
		}
		body = sb.String()
	} else {
		body = StatusIdle.Render("(no output)")
	}

	footer := detailFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		summary,
		separator,
		"",
		detailBody(body, width, height),
		"",
		footer,
	)
}

// min returns the smaller of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
