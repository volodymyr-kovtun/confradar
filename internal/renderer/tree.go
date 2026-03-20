package renderer

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// TreeRenderer outputs a pretty-printed categorized tree.
type TreeRenderer struct {
	NoColor bool
}

// Render writes the tree to the given writer.
func (t *TreeRenderer) Render(result *scanner.ScanResult, w io.Writer) error {
	if len(result.Ordered) == 0 {
		fmt.Fprintln(w, "No configuration files found.")
		return nil
	}

	headerStyle := lipgloss.NewStyle().Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	if t.NoColor {
		headerStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
		pathStyle = lipgloss.NewStyle()
	}

	projectName := filepath.Base(result.Root)
	fmt.Fprintf(w, "\n%s %s\n",
		headerStyle.Render(projectName),
		dimStyle.Render(fmt.Sprintf("(%d config files found in %s)", result.TotalFiles, result.Duration.Round(1e6))),
	)
	fmt.Fprintln(w)

	for catIdx, cat := range result.Ordered {
		isLastCat := catIdx == len(result.Ordered)-1

		catHeader := fmt.Sprintf("%s %s", cat.Icon, cat.Name)
		countStr := fmt.Sprintf("(%d)", len(cat.Files))

		var catColor lipgloss.Color
		if cat.Color != "" && !t.NoColor {
			catColor = lipgloss.Color(cat.Color)
			catHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(catColor)
			fmt.Fprintf(w, "%s %s\n", catHeaderStyle.Render(catHeader), dimStyle.Render(countStr))
		} else {
			fmt.Fprintf(w, "%s %s\n", headerStyle.Render(catHeader), dimStyle.Render(countStr))
		}

		for fileIdx, file := range cat.Files {
			isLastFile := fileIdx == len(cat.Files)-1

			var prefix string
			if isLastCat && isLastFile {
				prefix = "  └── "
			} else if isLastFile {
				prefix = "  └── "
			} else {
				prefix = "  ├── "
			}

			dir := filepath.Dir(file.RelPath)
			base := filepath.Base(file.RelPath)

			var line string
			if dir != "." {
				line = fmt.Sprintf("%s%s%s",
					dimStyle.Render(prefix),
					dimStyle.Render(dir+"/"),
					pathStyle.Render(base),
				)
			} else {
				line = fmt.Sprintf("%s%s",
					dimStyle.Render(prefix),
					pathStyle.Render(base),
				)
			}

			fmt.Fprintln(w, line)
		}

		if !isLastCat {
			fmt.Fprintln(w)
		}
	}

	// Summary line.
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s\n",
		dimStyle.Render(fmt.Sprintf("─── %d files across %d categories ───",
			result.TotalFiles, len(result.Ordered))),
	)

	return nil
}

// FormatSize returns a human-readable file size.
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// RepeatStr repeats a string n times.
func RepeatStr(s string, n int) string {
	return strings.Repeat(s, n)
}
