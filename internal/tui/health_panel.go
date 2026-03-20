package tui

import (
	"fmt"
	"strings"

	"github.com/volodymyrkovtun/confradar/internal/health"
)

// HealthPanel displays health check issues.
type HealthPanel struct {
	issues  []health.HealthIssue
	cursor  int
	offset  int
	height  int
	visible bool
}

// NewHealthPanel creates a health panel.
func NewHealthPanel() HealthPanel {
	return HealthPanel{height: 8}
}

// SetIssues updates the issues list.
func (hp *HealthPanel) SetIssues(issues []health.HealthIssue) {
	hp.issues = issues
}

// Toggle shows/hides the panel.
func (hp *HealthPanel) Toggle() {
	hp.visible = !hp.visible
}

// IsVisible returns whether the panel is shown.
func (hp *HealthPanel) IsVisible() bool {
	return hp.visible && len(hp.issues) > 0
}

// SetHeight sets the visible height.
func (hp *HealthPanel) SetHeight(h int) {
	hp.height = h
	if hp.height < 1 {
		hp.height = 1
	}
}

// MoveUp moves the cursor up.
func (hp *HealthPanel) MoveUp() {
	hp.cursor--
	if hp.cursor < 0 {
		hp.cursor = 0
	}
	if hp.cursor < hp.offset {
		hp.offset = hp.cursor
	}
}

// MoveDown moves the cursor down.
func (hp *HealthPanel) MoveDown() {
	hp.cursor++
	if hp.cursor >= len(hp.issues) {
		hp.cursor = len(hp.issues) - 1
	}
	if hp.cursor >= hp.offset+hp.height {
		hp.offset = hp.cursor - hp.height + 1
	}
}

// SelectedFile returns the file associated with the selected issue.
func (hp *HealthPanel) SelectedFile() string {
	if hp.cursor >= len(hp.issues) {
		return ""
	}
	return hp.issues[hp.cursor].File
}

// View renders the health panel.
func (hp *HealthPanel) View(styles Styles, width int) string {
	if !hp.visible || len(hp.issues) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(styles.Accent.Render(fmt.Sprintf(" Health Issues (%d)", len(hp.issues))))
	b.WriteByte('\n')

	end := hp.offset + hp.height - 1
	if end > len(hp.issues) {
		end = len(hp.issues)
	}

	for i := hp.offset; i < end; i++ {
		issue := hp.issues[i]
		isCurrent := i == hp.cursor

		var icon string
		var styledIcon string
		switch issue.Severity {
		case health.SeverityError:
			icon = "✗"
			styledIcon = styles.Error.Render(icon)
		case health.SeverityWarning:
			icon = "!"
			styledIcon = styles.Warning.Render(icon)
		default:
			icon = "i"
			styledIcon = styles.Info.Render(icon)
		}

		file := ""
		if issue.File != "" {
			file = styles.Dim.Render(" [" + issue.File + "]")
		}

		line := fmt.Sprintf(" %s %s%s", styledIcon, issue.Message, file)
		if isCurrent {
			line = styles.Selected.Render(padRight(line, width))
		}

		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String()
}

// Summary returns a short summary string for the status bar.
func (hp *HealthPanel) Summary(styles Styles) string {
	if len(hp.issues) == 0 {
		return styles.Success.Render("✓ healthy")
	}

	errors, warnings := 0, 0
	for _, issue := range hp.issues {
		switch issue.Severity {
		case health.SeverityError:
			errors++
		case health.SeverityWarning:
			warnings++
		}
	}

	var parts []string
	if errors > 0 {
		parts = append(parts, styles.Error.Render(fmt.Sprintf("%d errors", errors)))
	}
	if warnings > 0 {
		parts = append(parts, styles.Warning.Render(fmt.Sprintf("%d warnings", warnings)))
	}
	return strings.Join(parts, " ")
}
