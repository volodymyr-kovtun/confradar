package differ

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/volodymyrkovtun/confradar/internal/parser"
)

// Diff compares two .env files and returns the differences.
func Diff(leftPath, rightPath string) (*DiffResult, error) {
	leftResult, err := parser.ParseFile("env", leftPath)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", leftPath, err)
	}
	rightResult, err := parser.ParseFile("env", rightPath)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", rightPath, err)
	}

	result := &DiffResult{
		LeftPath:  leftPath,
		RightPath: rightPath,
	}

	leftKeys := leftResult.Values
	rightKeys := rightResult.Values

	allKeys := make(map[string]bool)
	for k := range leftKeys {
		allKeys[k] = true
	}
	for k := range rightKeys {
		allKeys[k] = true
	}

	sortedKeys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		leftVal, inLeft := leftKeys[key]
		rightVal, inRight := rightKeys[key]

		switch {
		case inLeft && !inRight:
			result.OnlyLeft = append(result.OnlyLeft, DiffKey{Key: key, Value: leftVal})
		case !inLeft && inRight:
			result.OnlyRight = append(result.OnlyRight, DiffKey{Key: key, Value: rightVal})
		case leftVal != rightVal:
			result.Changed = append(result.Changed, DiffPair{Key: key, LeftValue: leftVal, RightValue: rightVal})
		default:
			result.Common = append(result.Common, DiffPair{Key: key, LeftValue: leftVal, RightValue: rightVal})
		}
	}

	return result, nil
}

// RenderDiff writes a color-coded side-by-side comparison.
func RenderDiff(result *DiffResult, w io.Writer, noColor bool, redact bool) {
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	headerStyle := lipgloss.NewStyle().Bold(true)

	if noColor {
		redStyle = lipgloss.NewStyle()
		greenStyle = lipgloss.NewStyle()
		yellowStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
		headerStyle = lipgloss.NewStyle()
	}

	leftName := filepath.Base(result.LeftPath)
	rightName := filepath.Base(result.RightPath)

	fmt.Fprintf(w, "\n%s\n\n",
		headerStyle.Render(fmt.Sprintf("Comparing %s vs %s", leftName, rightName)))

	onlyL, onlyR, changed, common := result.Stats()

	// Only in left.
	if len(result.OnlyLeft) > 0 {
		fmt.Fprintf(w, "%s %s\n",
			redStyle.Render("−"),
			headerStyle.Render(fmt.Sprintf("Only in %s (%d):", leftName, onlyL)))
		for _, dk := range result.OnlyLeft {
			val := dk.Value
			if redact {
				val = "••••••"
			}
			fmt.Fprintf(w, "  %s %s%s\n",
				redStyle.Render("−"),
				dk.Key,
				dimStyle.Render("="+val))
		}
		fmt.Fprintln(w)
	}

	// Only in right.
	if len(result.OnlyRight) > 0 {
		fmt.Fprintf(w, "%s %s\n",
			greenStyle.Render("+"),
			headerStyle.Render(fmt.Sprintf("Only in %s (%d):", rightName, onlyR)))
		for _, dk := range result.OnlyRight {
			val := dk.Value
			if redact {
				val = "••••••"
			}
			fmt.Fprintf(w, "  %s %s%s\n",
				greenStyle.Render("+"),
				dk.Key,
				dimStyle.Render("="+val))
		}
		fmt.Fprintln(w)
	}

	// Changed.
	if len(result.Changed) > 0 {
		fmt.Fprintf(w, "%s %s\n",
			yellowStyle.Render("~"),
			headerStyle.Render(fmt.Sprintf("Different values (%d):", changed)))
		for _, dp := range result.Changed {
			lv, rv := dp.LeftValue, dp.RightValue
			if redact {
				lv, rv = "••••••", "••••••"
			}
			fmt.Fprintf(w, "  %s %s\n", yellowStyle.Render("~"), dp.Key)
			fmt.Fprintf(w, "    %s %s\n", redStyle.Render("−"), dimStyle.Render(lv))
			fmt.Fprintf(w, "    %s %s\n", greenStyle.Render("+"), dimStyle.Render(rv))
		}
		fmt.Fprintln(w)
	}

	// Summary.
	summary := fmt.Sprintf("─── %d common, %d only left, %d only right, %d changed ───",
		common, onlyL, onlyR, changed)
	fmt.Fprintln(w, dimStyle.Render(summary))

}
