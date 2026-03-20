package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/renderer"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan and print a categorized tree of config files",
	Long:  "Scan a project directory and display all detected configuration files grouped by category.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	cfg, err := config.New(path, flags)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	format := config.EffectiveFormat(cfg, flags)

	var result *scanner.ScanResult
	if flags.MaxDepth > 0 {
		result, err = scanner.ScanWithDepth(path, cfg, flags.MaxDepth)
	} else {
		result, err = scanner.Scan(path, cfg)
	}
	if err != nil {
		return fmt.Errorf("scanning: %w", err)
	}

	r, err := renderer.New(format, flags.NoColor)
	if err != nil {
		return err
	}

	return r.Render(result, os.Stdout)
}
