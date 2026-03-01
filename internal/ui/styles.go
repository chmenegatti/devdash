// Package ui provides styling constants and helpers using Lipgloss.
package ui

import "github.com/charmbracelet/lipgloss"

// Color palette.
var (
	ColorPrimary   = lipgloss.Color("#7D56F4")
	ColorSecondary = lipgloss.Color("#6C6C6C")
	ColorSuccess   = lipgloss.Color("#73D216")
	ColorWarning   = lipgloss.Color("#F5A623")
	ColorDanger    = lipgloss.Color("#FF5F56")
	ColorMuted     = lipgloss.Color("#555555")
	ColorWhite     = lipgloss.Color("#FAFAFA")
	ColorDim       = lipgloss.Color("#888888")
)

// PanelStyle is a reusable bordered card style.
var PanelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorPrimary).
	Padding(0, 1)

// TitleStyle styles the top header bar.
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorWhite).
	Background(ColorPrimary).
	Padding(0, 2).
	MarginBottom(1)

// SubtitleStyle for secondary information.
var SubtitleStyle = lipgloss.NewStyle().
	Foreground(ColorDim).
	MarginBottom(1)

// LabelStyle for field labels inside panels.
var LabelStyle = lipgloss.NewStyle().
	Foreground(ColorSecondary).
	Bold(true)

// ValueStyle for field values.
var ValueStyle = lipgloss.NewStyle().
	Foreground(ColorWhite)

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
	Foreground(ColorMuted).
	Italic(true)

// HelpStyle for the footer help bar.
var HelpStyle = lipgloss.NewStyle().
	Foreground(ColorDim).
	MarginTop(1)

// KeyStyle for a keyboard shortcut key in help text.
var KeyStyle = lipgloss.NewStyle().
	Foreground(ColorPrimary).
	Bold(true)

// SepStyle is a subtle separator.
var SepStyle = lipgloss.NewStyle().
	Foreground(ColorMuted)
