package tui

import (
	"fmt"
	"strings"

	"github.com/volodymyrkovtun/confradar/internal/differ"
)

// DiffView displays a side-by-side .env diff.
type DiffView struct {
	result  *differ.DiffResult
	visible bool
	offset  int
	height  int
	redact  bool
}

// NewDiffView creates a diff view.
func NewDiffView() DiffView {
	return DiffView{height: 20, redact: true}
}

// Show displays a diff result.
func (dv *DiffView) Show(result *differ.DiffResult) {
	dv.result = result
	dv.visible = true
	dv.offset = 0
}

// Hide closes the diff view.
func (dv *DiffView) Hide() {
	dv.visible = false
}

// IsVisible returns whether the diff view is shown.
func (dv *DiffView) IsVisible() bool {
	return dv.visible && dv.result != nil
}

// SetHeight sets the visible height.
func (dv *DiffView) SetHeight(h int) {
	dv.height = h
}

// ScrollUp scrolls up.
func (dv *DiffView) ScrollUp(n int) {
	dv.offset -= n
	if dv.offset < 0 {
		dv.offset = 0
	}
}

// ScrollDown scrolls down.
func (dv *DiffView) ScrollDown(n int) {
	dv.offset += n
}

// ToggleRedact toggles value redaction.
func (dv *DiffView) ToggleRedact() {
	dv.redact = !dv.redact
}

// View renders the diff.
func (dv *DiffView) View(styles Styles, width int) string {
	if dv.result == nil {
		return ""
	}

	var b strings.Builder
	r := dv.result

	b.WriteString(styles.Accent.Render(fmt.Sprintf(" Diff: %s vs %s",
		shortPath(r.LeftPath), shortPath(r.RightPath))))
	b.WriteByte('\n')
	b.WriteByte('\n')

	var lines []string

	// Only in left.
	if len(r.OnlyLeft) > 0 {
		lines = append(lines, styles.Error.Render(fmt.Sprintf(" ─ Only in %s (%d):", shortPath(r.LeftPath), len(r.OnlyLeft))))
		for _, dk := range r.OnlyLeft {
			val := dk.Value
			if dv.redact {
				val = "••••••"
			}
			lines = append(lines, fmt.Sprintf("   %s %s=%s", styles.Error.Render("−"), dk.Key, styles.Dim.Render(val)))
		}
		lines = append(lines, "")
	}

	// Only in right.
	if len(r.OnlyRight) > 0 {
		lines = append(lines, styles.Success.Render(fmt.Sprintf(" ─ Only in %s (%d):", shortPath(r.RightPath), len(r.OnlyRight))))
		for _, dk := range r.OnlyRight {
			val := dk.Value
			if dv.redact {
				val = "••••••"
			}
			lines = append(lines, fmt.Sprintf("   %s %s=%s", styles.Success.Render("+"), dk.Key, styles.Dim.Render(val)))
		}
		lines = append(lines, "")
	}

	// Changed.
	if len(r.Changed) > 0 {
		lines = append(lines, styles.Warning.Render(fmt.Sprintf(" ─ Different values (%d):", len(r.Changed))))
		for _, dp := range r.Changed {
			lv, rv := dp.LeftValue, dp.RightValue
			if dv.redact {
				lv, rv = "••••••", "••••••"
			}
			lines = append(lines, fmt.Sprintf("   %s %s", styles.Warning.Render("~"), dp.Key))
			lines = append(lines, fmt.Sprintf("     %s %s", styles.Error.Render("−"), styles.Dim.Render(lv)))
			lines = append(lines, fmt.Sprintf("     %s %s", styles.Success.Render("+"), styles.Dim.Render(rv)))
		}
		lines = append(lines, "")
	}

	// Common count.
	lines = append(lines, styles.Dim.Render(fmt.Sprintf(" %d keys in common", len(r.Common))))

	// Apply scrolling.
	end := dv.offset + dv.height - 2
	if end > len(lines) {
		end = len(lines)
	}
	start := dv.offset
	if start > len(lines) {
		start = len(lines)
	}

	for i := start; i < end; i++ {
		b.WriteString(lines[i])
		b.WriteByte('\n')
	}

	return b.String()
}

func shortPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) <= 2 {
		return path
	}
	return parts[len(parts)-1]
}
