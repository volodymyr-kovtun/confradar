package tui

import (
	"strings"
)

// HelpOverlay displays keybinding help.
type HelpOverlay struct {
	visible bool
}

// Toggle shows/hides the help overlay.
func (h *HelpOverlay) Toggle() {
	h.visible = !h.visible
}

// IsVisible returns whether help is shown.
func (h *HelpOverlay) IsVisible() bool {
	return h.visible
}

// Hide closes the help overlay.
func (h *HelpOverlay) Hide() {
	h.visible = false
}

type helpBinding struct {
	key  string
	desc string
}

type helpSection struct {
	title    string
	bindings []helpBinding
}

var helpSections = []helpSection{
	{
		title: "Navigation",
		bindings: []helpBinding{
			{"↑/k", "Move up"},
			{"↓/j", "Move down"},
			{"Enter", "Expand/collapse or view file"},
			{"Esc", "Go back / close overlay"},
			{"Tab", "Switch focus (tree ↔ viewer)"},
			{"1-9", "Jump to category"},
			{"q", "Quit"},
		},
	},
	{
		title: "Actions",
		bindings: []helpBinding{
			{"/", "Search files"},
			{"d", "Diff two .env files"},
			{"h", "Toggle health panel"},
			{"r", "Rescan project"},
			{"e", "Open in $EDITOR"},
			{"y", "Copy file path"},
			{"o", "Open in system app"},
		},
	},
	{
		title: "Display",
		bindings: []helpBinding{
			{"c", "Cycle color theme"},
			{"s", "Cycle sort order"},
			{"p", "Toggle preview panel"},
			{"v", "Toggle .env value redaction"},
			{"?", "Toggle this help"},
		},
	},
}

// View renders the help overlay.
func (h *HelpOverlay) View(styles Styles, width, height int) string {
	if !h.visible {
		return ""
	}

	var b strings.Builder
	b.WriteString(styles.Accent.Render(" Keybindings"))
	b.WriteString("\n\n")

	for _, section := range helpSections {
		b.WriteString(styles.Title.Render(" " + section.title))
		b.WriteByte('\n')
		for _, binding := range section.bindings {
			key := styles.Accent.Render(padRight("  "+binding.key, 14))
			b.WriteString(key)
			b.WriteString(styles.Normal.Render(binding.desc))
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
	}

	b.WriteString(styles.Dim.Render(" Press ? or Esc to close"))

	return b.String()
}
