// Package ui — components.go provides reusable rendering helpers for panels.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Logo ────────────────────────────────────────────────────────────────────

const devdashLogo = `╺┳╸╻ ╻┏━┓   ┏━╸┏━┓   ╺┳┓┏━┓┏━┓╻ ╻
 ┃ ┣━┫┣╸    ┃╺┓┃ ┃    ┃┃┣━┫┗━┓┣━┫
 ╹ ╹ ╹┗━╸   ┗━┛┗━┛   ╺┻┛╹ ╹┗━┛╹ ╹`

// RenderLogo returns the styled ASCII logo.
func RenderLogo() string {
	return LogoStyle.Render(devdashLogo)
}

// ── Breadcrumbs ─────────────────────────────────────────────────────────────

// RenderCrumbs renders K9s-style breadcrumb navigation.
func RenderCrumbs(items ...string) string {
	parts := make([]string, len(items))
	for i, item := range items {
		parts[i] = CrumbStyle.Render(item)
	}
	sep := CrumbSepStyle.Render(" ▸ ")
	return strings.Join(parts, sep)
}

// ── Stat chips (header info bar) ────────────────────────────────────────────

// StatChip renders a "label:value" pair for the header stats bar.
func StatChip(label, value string) string {
	return fmt.Sprintf("%s:%s",
		InfoStyle.Render(label),
		InfoValueStyle.Render(value),
	)
}

// ── Status dot ──────────────────────────────────────────────────────────────

// StatusDot returns a colored bullet indicator: ● ● ● ○
func StatusDot(label string, style lipgloss.Style) string {
	return style.Render("●") + " " + style.Render(label)
}

// ── Section panel ───────────────────────────────────────────────────────────

// RenderSection renders a K9s-style bordered section with a title label.
func RenderSection(title, body string, width int, height ...int) string {
	titleBar := SectionTitleStyle.Render("╸" + title + "╺")
	content := titleBar + "\n" + body

	style := SectionBorder.Width(width)
	if len(height) > 0 && height[0] > 0 {
		style = style.Height(height[0])
	}
	return style.Render(content)
}

// RenderPanel is kept for backward compat. Delegates to RenderSection.
func RenderPanel(title, body string, width int) string {
	return RenderSection(title, body, width)
}

// ── Field rendering ─────────────────────────────────────────────────────────

// RenderField formats a label: value line.
func RenderField(label, value string) string {
	return fmt.Sprintf("%s %s",
		LabelStyle.Render(label+":"),
		ValueStyle.Render(value),
	)
}

// RenderStatusField renders a label with a colored status value.
func RenderStatusField(label, value string, style lipgloss.Style) string {
	return fmt.Sprintf("%s %s",
		LabelStyle.Render(label+":"),
		style.Render(value),
	)
}

// ── Table helpers ───────────────────────────────────────────────────────────

// TableColumn describes a column in a simple text table.
type TableColumn struct {
	Header string
	Width  int
}

// RenderTableHeader renders column headers with a separator line below.
func RenderTableHeader(cols []TableColumn) string {
	cells := make([]string, len(cols))
	for i, c := range cols {
		cells[i] = TableHeaderStyle.
			Width(c.Width).
			Render(c.Header)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

// RenderTableRow renders a single data row, optionally with alt-row bg.
func RenderTableRow(values []string, cols []TableColumn, alt bool) string {
	style := TableRowStyle
	if alt {
		style = TableRowAltStyle
	}
	cells := make([]string, len(cols))
	for i, c := range cols {
		val := ""
		if i < len(values) {
			val = values[i]
		}
		cells[i] = style.Width(c.Width).Render(val)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

// ── Command bar (bottom) ────────────────────────────────────────────────────

// KeyBinding maps a key to its description.
type KeyBinding struct {
	Key  string
	Desc string
}

// RenderHelp builds the K9s-style command bar.
func RenderHelp(bindings []KeyBinding) string {
	parts := make([]string, len(bindings))
	for i, b := range bindings {
		parts[i] = lipgloss.JoinHorizontal(
			lipgloss.Center,
			HotkeyBoxStyle.Render(b.Key),
			HotkeyDescStyle.Render(b.Desc),
		)
	}
	return HelpStyle.Render(strings.Join(parts, "   "))
}

// RenderCommandBar renders a full-width bottom bar with key hints.
func RenderCommandBar(bindings []KeyBinding, width int) string {
	inner := RenderHelp(bindings)
	bar := lipgloss.NewStyle().
		Width(width).
		Background(ColorHeaderBg).
		Padding(0, 1).
		Render(inner)
	return bar
}
