package tui

import "github.com/charmbracelet/lipgloss"

// Styles holds pre-computed lipgloss styles derived from a Theme.
type Styles struct {
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	Normal        lipgloss.Style
	Dim           lipgloss.Style
	Accent        lipgloss.Style
	Selected      lipgloss.Style
	CategoryName  lipgloss.Style
	FileName      lipgloss.Style
	Border        lipgloss.Style
	StatusBar     lipgloss.Style
	Error         lipgloss.Style
	Warning       lipgloss.Style
	Success       lipgloss.Style
	Info          lipgloss.Style
	Help          lipgloss.Style
	SearchPrompt  lipgloss.Style
}

// NewStyles creates a Styles set from a Theme.
func NewStyles(theme Theme) Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Foreground),
		Subtitle: lipgloss.NewStyle().
			Foreground(theme.Dim),
		Normal: lipgloss.NewStyle().
			Foreground(theme.Foreground),
		Dim: lipgloss.NewStyle().
			Foreground(theme.Dim),
		Accent: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true),
		Selected: lipgloss.NewStyle().
			Background(theme.Selected).
			Foreground(theme.Foreground).
			Bold(true),
		CategoryName: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true),
		FileName: lipgloss.NewStyle().
			Foreground(theme.Foreground),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border),
		StatusBar: lipgloss.NewStyle().
			Foreground(theme.Dim).
			Background(theme.Border),
		Error: lipgloss.NewStyle().
			Foreground(theme.Error),
		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning),
		Success: lipgloss.NewStyle().
			Foreground(theme.Success),
		Info: lipgloss.NewStyle().
			Foreground(theme.Info),
		Help: lipgloss.NewStyle().
			Foreground(theme.Dim),
		SearchPrompt: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true),
	}
}
