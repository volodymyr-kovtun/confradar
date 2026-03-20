package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/volodymyrkovtun/confradar/internal/config"
)

var listPatternsCmd = &cobra.Command{
	Use:   "list-patterns",
	Short: "Show all active patterns (built-in + custom)",
	RunE:  runListPatterns,
}

func init() {
	rootCmd.AddCommand(listPatternsCmd)
}

func runListPatterns(cmd *cobra.Command, args []string) error {
	path := "."
	cfg, err := config.New(path, flags)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	headerStyle := lipgloss.NewStyle().Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	if flags.NoColor {
		headerStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
	}

	totalPatterns := 0
	for _, cat := range cfg.Categories {
		catHeader := fmt.Sprintf("%s %s", cat.Icon, cat.Name)
		fmt.Fprintf(os.Stdout, "\n%s %s\n",
			headerStyle.Render(catHeader),
			dimStyle.Render(fmt.Sprintf("(%d patterns)", len(cat.Patterns))),
		)

		for i, pat := range cat.Patterns {
			isLast := i == len(cat.Patterns)-1
			prefix := "├── "
			if isLast {
				prefix = "└── "
			}

			globs := pat.AllGlobs()
			globStr := strings.Join(globs, ", ")

			desc := ""
			if pat.Description != "" {
				desc = dimStyle.Render(" — " + pat.Description)
			}

			fmt.Fprintf(os.Stdout, "  %s%s %s%s\n",
				dimStyle.Render(prefix),
				pat.Name,
				dimStyle.Render("["+globStr+"]"),
				desc,
			)
		}
		totalPatterns += len(cat.Patterns)
	}

	fmt.Fprintf(os.Stdout, "\n%s\n",
		dimStyle.Render(fmt.Sprintf("─── %d patterns across %d categories ───",
			totalPatterns, len(cfg.Categories))),
	)
	return nil
}
