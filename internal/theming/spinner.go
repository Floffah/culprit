package theming

import (
	"charm.land/huh/v2/spinner"
	"charm.land/lipgloss/v2"
)

type spinnerTheme struct{}

func (t spinnerTheme) Theme(isDark bool) *spinner.Styles {
	lightDark := lipgloss.LightDark(isDark)
	title := lightDark(
		lipgloss.Color("#00020A"),
		lipgloss.Color("#FFFDF5"),
	)
	return &spinner.Styles{
		Spinner: lipgloss.NewStyle().Foreground(lipgloss.Color("#F780E2")).PaddingLeft(3),
		Title:   lipgloss.NewStyle().Foreground(title),
	}
}

func SpinnerTheme() spinner.Theme {
	return &spinnerTheme{}
}
