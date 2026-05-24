package theming

import "charm.land/lipgloss/v2"

func MutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6b7280")).
		MarginLeft(4).
		Padding(0, 1)
}

func ScriptWrittenStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#064e3b")).
		Padding(0, 1).
		MarginLeft(5)
}

func ScriptRunScriptBeforeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		MarginLeft(5)
}

func ScriptRunScriptAfterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Italic(true)
}

func ScriptRunCommandStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Italic(true)
}

func DangerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#dc2626")).
		Bold(true).
		MarginLeft(4).
		Width(60).
		Padding(0, 1)
}
