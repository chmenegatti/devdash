// Package ui — components.go provides reusable rendering helpers for panels.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderPanel renders a bordered card with a title and body content.
func RenderPanel(title, body string, width int) string {
	titleBar := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render(title)

	content := titleBar + "\n" + body

	return PanelStyle.
		Width(width).
		Render(content)
}

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

// RenderHelp builds the keyboard shortcut help bar.
func RenderHelp(bindings []KeyBinding) string {
	parts := make([]string, len(bindings))
	for i, b := range bindings {
		parts[i] = fmt.Sprintf("%s %s",
			KeyStyle.Render(b.Key),
			lipgloss.NewStyle().Foreground(ColorDim).Render(b.Desc),
		)
	}
	return HelpStyle.Render(strings.Join(parts, "  │  "))
}

// KeyBinding maps a key to its description.
type KeyBinding struct {
	Key  string
	Desc string
}
