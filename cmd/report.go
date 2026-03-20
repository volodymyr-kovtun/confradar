package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/health"
	"github.com/volodymyrkovtun/confradar/internal/renderer"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

var reportOutput string

var reportCmd = &cobra.Command{
	Use:   "report [path]",
	Short: "Generate a configuration report",
	Long:  "Scan a project and generate a comprehensive Markdown or JSON report.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runReport,
}

func init() {
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "", "write report to file instead of stdout")
	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
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

	// Run health checks.
	issues := health.RunChecks(absPath, result, cfg.HealthChecks)
	autoIssues := health.RunAutoChecks(absPath, result)
	issues = append(issues, autoIssues...)

	format := config.EffectiveFormat(cfg, flags)

	// Determine output writer.
	var w *os.File
	if reportOutput != "" {
		w, err = os.Create(reportOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}

	switch format {
	case "json":
		r := &renderer.JSONRenderer{}
		return r.Render(result, w)
	case "yaml":
		r := &renderer.YAMLRenderer{}
		return r.Render(result, w)
	default: // markdown
		r := &renderer.MarkdownRenderer{
			Issues:     issues,
			IncludeTOC: true,
		}
		return r.Render(result, w)
	}
}
