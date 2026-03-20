package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// StatusBar renders the bottom status bar.
type StatusBar struct {
	result *scanner.ScanResult
	theme  string
}

// NewStatusBar creates a status bar.
func NewStatusBar() StatusBar {
	return StatusBar{theme: "dark"}
}

// SetResult updates scan results.
func (sb *StatusBar) SetResult(result *scanner.ScanResult) {
	sb.result = result
}

// SetTheme updates the current theme name.
func (sb *StatusBar) SetTheme(name string) {
	sb.theme = name
}

// View renders the status bar.
func (sb *StatusBar) View(styles Styles, healthSummary string, width int) string {
	if sb.result == nil {
		return ""
	}

	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a1a")).
		Foreground(lipgloss.Color("#999999")).
		Width(width)

	projectName := filepath.Base(sb.result.Root)

	left := fmt.Sprintf(" %s │ %d files │ %d categories │ %s",
		projectName,
		sb.result.TotalFiles,
		len(sb.result.Ordered),
		sb.result.Duration.Round(1e6),
	)

	right := fmt.Sprintf("%s │ theme: %s │ ? help ", healthSummary, sb.theme)

	gap := width - len(stripAnsi(left)) - len(stripAnsi(right))
	if gap < 0 {
		gap = 0
	}

	return barStyle.Render(left + strings.Repeat(" ", gap) + right)
}

// stripAnsi removes ANSI escape codes for length calculation.
func stripAnsi(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
