package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/health"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Run health checks and report issues",
	Long:  "Scan a project, run all configured health checks, and report issues. Exits with code 1 if any error-severity issues are found.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	cfg, err := config.New(path, flags)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	result, err := scanner.Scan(path, cfg)
	if err != nil {
		return fmt.Errorf("scanning: %w", err)
	}

	// Run configured health checks.
	issues := health.RunChecks(absPath, result, cfg.HealthChecks)

	// Also run automatic checks.
	autoIssues := health.RunAutoChecks(absPath, result)
	issues = append(issues, autoIssues...)

	// Filter by severity.
	if flags.Severity != "" {
		issues = health.FilterBySeverity(issues, flags.Severity)
	}

	// Render issues.
	renderIssues(issues)

	// Exit with code 1 if any errors.
	for _, issue := range issues {
		if issue.Severity == health.SeverityError {
			os.Exit(1)
		}
	}

	return nil
}

func renderIssues(issues []health.HealthIssue) {
	if len(issues) == 0 {
		greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true)
		fmt.Println(greenStyle.Render("✓ No health issues found"))
		return
	}

	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3b82f6"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	headerStyle := lipgloss.NewStyle().Bold(true)

	if flags.NoColor {
		errorStyle = lipgloss.NewStyle()
		warnStyle = lipgloss.NewStyle()
		infoStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
		headerStyle = lipgloss.NewStyle()
	}

	errors, warnings, infos := 0, 0, 0
	for _, issue := range issues {
		switch issue.Severity {
		case health.SeverityError:
			errors++
		case health.SeverityWarning:
			warnings++
		case health.SeverityInfo:
			infos++
		}
	}

	fmt.Printf("\n%s\n\n",
		headerStyle.Render(fmt.Sprintf("Health Check Results (%d issues)", len(issues))))

	for _, issue := range issues {
		var icon, styledSeverity string
		switch issue.Severity {
		case health.SeverityError:
			icon = errorStyle.Render("✗")
			styledSeverity = errorStyle.Render("ERROR")
		case health.SeverityWarning:
			icon = warnStyle.Render("!")
			styledSeverity = warnStyle.Render("WARN ")
		default:
			icon = infoStyle.Render("i")
			styledSeverity = infoStyle.Render("INFO ")
		}

		file := ""
		if issue.File != "" {
			file = dimStyle.Render(" [" + issue.File + "]")
		}

		fmt.Printf("  %s %s %s%s\n", icon, styledSeverity, issue.Message, file)
	}

	fmt.Printf("\n%s\n",
		dimStyle.Render(fmt.Sprintf("─── %d errors, %d warnings, %d info ───",
			errors, warnings, infos)))
}
