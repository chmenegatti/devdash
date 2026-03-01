// Package ui provides styling constants and helpers using Lipgloss.
// Uses a K9s-inspired dark theme with cyan/blue accents.
package ui

import "github.com/charmbracelet/lipgloss"

// ── K9s-inspired color palette ──────────────────────────────────────────────

var (
	// Core accent — K9s signature cyan/teal
	ColorAccent    = lipgloss.Color("#00d7d7")
	ColorAccentDim = lipgloss.Color("#008b8b")

	// Logo gradient
	ColorLogo1 = lipgloss.Color("#5fafff")
	ColorLogo2 = lipgloss.Color("#00d7d7")

	// Semantic
	ColorSuccess = lipgloss.Color("#00d787")
	ColorWarning = lipgloss.Color("#d7af00")
	ColorDanger  = lipgloss.Color("#ff5f5f")

	// Neutrals
	ColorWhite     = lipgloss.Color("#e4e4e4")
	ColorFg        = lipgloss.Color("#bcbcbc")
	ColorDim       = lipgloss.Color("#6c6c6c")
	ColorSubtle    = lipgloss.Color("#444444")
	ColorBorder    = lipgloss.Color("#3a3a3a")
	ColorBg        = lipgloss.Color("#1c1c1c")
	ColorHeaderBg  = lipgloss.Color("#262626")
	ColorSelected  = lipgloss.Color("#303030")
	ColorCrumbBg   = lipgloss.Color("#303030")
	ColorStatusBar = lipgloss.Color("#1c1c1c")

	// Backward compat aliases
	ColorPrimary   = ColorAccent
	ColorSecondary = ColorDim
	ColorMuted     = ColorSubtle
)

// ── Header ──────────────────────────────────────────────────────────────────

// LogoStyle renders the "devdash" logo text.
var LogoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorLogo2)

// HeaderBarStyle is the full-width top bar background.
var HeaderBarStyle = lipgloss.NewStyle().
	Background(ColorHeaderBg).
	Foreground(ColorFg).
	Padding(0, 1)

// CrumbStyle renders breadcrumb segments.
var CrumbStyle = lipgloss.NewStyle().
	Background(ColorCrumbBg).
	Foreground(ColorAccent).
	Bold(true).
	Padding(0, 1)

// CrumbSepStyle for the separator between crumbs.
var CrumbSepStyle = lipgloss.NewStyle().
	Foreground(ColorDim)

// InfoStyle for inline key:value stats on the header line.
var InfoStyle = lipgloss.NewStyle().
	Foreground(ColorDim)

// InfoValueStyle for stat values in the header.
var InfoValueStyle = lipgloss.NewStyle().
	Foreground(ColorFg)

// ── Section / Panel ─────────────────────────────────────────────────────────

// SectionBorder uses a thin line box.
var SectionBorder = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(ColorBorder).
	Padding(0, 1)

// SectionTitleStyle renders a panel header label.
var SectionTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorAccent)

// ── Table ───────────────────────────────────────────────────────────────────

// TableHeaderStyle for column headers inside panels.
var TableHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorAccentDim).
	BorderBottom(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(ColorBorder)

// TableRowStyle for regular data rows.
var TableRowStyle = lipgloss.NewStyle().
	Foreground(ColorFg)

// TableRowAltStyle for alternating rows (subtle highlight).
var TableRowAltStyle = lipgloss.NewStyle().
	Foreground(ColorFg).
	Background(ColorSelected)

// ── Field labels & values ───────────────────────────────────────────────────

// LabelStyle for field labels inside panels.
var LabelStyle = lipgloss.NewStyle().
	Foreground(ColorDim).
	Bold(true)

// ValueStyle for field values.
var ValueStyle = lipgloss.NewStyle().
	Foreground(ColorWhite)

// ── Status badges ───────────────────────────────────────────────────────────

// StatusPass styles pass/ok values.
var StatusPass = lipgloss.NewStyle().
	Foreground(ColorSuccess).
	Bold(true)

// StatusFail styles fail/error values.
var StatusFail = lipgloss.NewStyle().
	Foreground(ColorDanger).
	Bold(true)

// StatusWarn styles warning values.
var StatusWarn = lipgloss.NewStyle().
	Foreground(ColorWarning).
	Bold(true)

// StatusIdle styles idle/placeholder values.
var StatusIdle = lipgloss.NewStyle().
	Foreground(ColorSubtle).
	Italic(true)

// ── Help / Command bar ──────────────────────────────────────────────────────

// HelpStyle for the command bar at the bottom.
var HelpStyle = lipgloss.NewStyle().
	Foreground(ColorDim)

// KeyStyle for a keyboard shortcut key.
var KeyStyle = lipgloss.NewStyle().
	Foreground(ColorAccent).
	Bold(true)

// ── Separators ──────────────────────────────────────────────────────────────

// SepStyle for horizontal rule lines.
var SepStyle = lipgloss.NewStyle().
	Foreground(ColorBorder)

// ── Backward compatibility (used by old component code) ─────────────────────

// PanelStyle = SectionBorder (alias for backward compat)
var PanelStyle = SectionBorder

// TitleStyle (alias — used by detail_views)
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorBg).
	Background(ColorAccent).
	Padding(0, 2)

// SubtitleStyle for secondary information.
var SubtitleStyle = lipgloss.NewStyle().
	Foreground(ColorDim).
	MarginBottom(0)
