package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// TableRenderer outputs scan results as an aligned table.
type TableRenderer struct {
	NoColor bool
}

// Render writes a table to the writer.
func (t *TableRenderer) Render(result *scanner.ScanResult, w io.Writer) error {
	if len(result.Ordered) == 0 {
		fmt.Fprintln(w, "No configuration files found.")
		return nil
	}

	headerStyle := lipgloss.NewStyle().Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	if t.NoColor {
		headerStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
	}

	// Calculate column widths.
	catWidth, fileWidth, patternWidth := 10, 10, 10
	for _, cat := range result.Ordered {
		if len(cat.Name) > catWidth {
			catWidth = len(cat.Name)
		}
		for _, f := range cat.Files {
			if len(f.RelPath) > fileWidth {
				fileWidth = len(f.RelPath)
			}
			if len(f.PatternName) > patternWidth {
				patternWidth = len(f.PatternName)
			}
		}
	}

	// Cap widths.
	if fileWidth > 60 {
		fileWidth = 60
	}
	if patternWidth > 30 {
		patternWidth = 30
	}
	if catWidth > 25 {
		catWidth = 25
	}

	// Header.
	header := fmt.Sprintf("%-*s  %-*s  %-*s  %s",
		catWidth, "CATEGORY",
		fileWidth, "FILE",
		patternWidth, "PATTERN",
		"SIZE",
	)
	fmt.Fprintln(w, headerStyle.Render(header))
	fmt.Fprintln(w, dimStyle.Render(strings.Repeat("─", catWidth+fileWidth+patternWidth+12)))

	for _, cat := range result.Ordered {
		for _, f := range cat.Files {
			relPath := f.RelPath
			if len(relPath) > fileWidth {
				relPath = "…" + relPath[len(relPath)-fileWidth+1:]
			}
			patName := f.PatternName
			if len(patName) > patternWidth {
				patName = patName[:patternWidth-1] + "…"
			}
			catName := cat.Icon + " " + cat.Name
			if len(catName) > catWidth {
				catName = catName[:catWidth-1] + "…"
			}

			fmt.Fprintf(w, "%-*s  %-*s  %-*s  %s\n",
				catWidth, catName,
				fileWidth, relPath,
				patternWidth, patName,
				FormatSize(f.Size),
			)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s\n",
		dimStyle.Render(fmt.Sprintf("%d files, %d categories, scanned in %s",
			result.TotalFiles, len(result.Ordered),
			result.Duration.Round(1e6))),
	)

	return nil
}
