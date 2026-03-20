// Package tui implements the interactive terminal UI using Bubble Tea.
package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines a color palette for the TUI.
type Theme struct {
	Name       string
	Background lipgloss.Color
	Foreground lipgloss.Color
	Dim        lipgloss.Color
	Accent     lipgloss.Color
	Border     lipgloss.Color
	Selected   lipgloss.Color
	Error      lipgloss.Color
	Warning    lipgloss.Color
	Success    lipgloss.Color
	Info       lipgloss.Color
}

var themes = map[string]Theme{
	"dark": {
		Name:       "dark",
		Background: lipgloss.Color("#0c0c0c"),
		Foreground: lipgloss.Color("#e0e0e0"),
		Dim:        lipgloss.Color("#666666"),
		Accent:     lipgloss.Color("#06b6d4"),
		Border:     lipgloss.Color("#333333"),
		Selected:   lipgloss.Color("#1a3a4a"),
		Error:      lipgloss.Color("#ef4444"),
		Warning:    lipgloss.Color("#f59e0b"),
		Success:    lipgloss.Color("#22c55e"),
		Info:       lipgloss.Color("#3b82f6"),
	},
	"light": {
		Name:       "light",
		Background: lipgloss.Color("#fafafa"),
		Foreground: lipgloss.Color("#1a1a1a"),
		Dim:        lipgloss.Color("#999999"),
		Accent:     lipgloss.Color("#0891b2"),
		Border:     lipgloss.Color("#d4d4d4"),
		Selected:   lipgloss.Color("#e0f2fe"),
		Error:      lipgloss.Color("#dc2626"),
		Warning:    lipgloss.Color("#d97706"),
		Success:    lipgloss.Color("#16a34a"),
		Info:       lipgloss.Color("#2563eb"),
	},
	"minimal": {
		Name:       "minimal",
		Background: lipgloss.Color(""),
		Foreground: lipgloss.Color("252"),
		Dim:        lipgloss.Color("240"),
		Accent:     lipgloss.Color("75"),
		Border:     lipgloss.Color("236"),
		Selected:   lipgloss.Color("235"),
		Error:      lipgloss.Color("196"),
		Warning:    lipgloss.Color("214"),
		Success:    lipgloss.Color("82"),
		Info:       lipgloss.Color("69"),
	},
}

var themeOrder = []string{"dark", "light", "minimal"}

// GetTheme returns a theme by name, defaulting to dark.
func GetTheme(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["dark"]
}

// NextTheme cycles to the next theme.
func NextTheme(current string) string {
	for i, name := range themeOrder {
		if name == current {
			return themeOrder[(i+1)%len(themeOrder)]
		}
	}
	return themeOrder[0]
}
