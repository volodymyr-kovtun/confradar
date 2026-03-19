// Package cmd defines all CLI commands for confradar.
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/tui"
	"golang.org/x/term"
)

var flags config.CLIFlags

var rootCmd = &cobra.Command{
	Use:   "confradar [path]",
	Short: "Instantly see every config file in your project",
	Long: `confradar scans any project directory and presents a unified,
categorized view of every configuration file — .env files, Docker configs,
CI/CD pipelines, build tool configs, linting rules, and more.

Run without arguments to launch the interactive TUI, or use subcommands
for specific tasks.`,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runRoot,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&flags.ConfigPath, "config", "", "path to config file")
	pf.BoolVar(&flags.NoConfig, "no-config", false, "ignore all config files, use built-in defaults only")
	pf.StringVar(&flags.Format, "format", "", "output format: tree, json, yaml, markdown, table")
	pf.BoolVar(&flags.NoColor, "no-color", false, "disable colored output")
	pf.BoolVar(&flags.Verbose, "verbose", false, "show debug info")
	pf.BoolVar(&flags.Quiet, "quiet", false, "suppress all output except errors")
	pf.StringVar(&flags.Severity, "severity", "", "filter health issues: error, warning, info")
	pf.BoolVar(&flags.HealthOnly, "health-only", false, "only show health issues")
	pf.BoolVar(&flags.IncludeHidden, "include-hidden", false, "include dotfiles normally skipped")
	pf.IntVar(&flags.MaxDepth, "max-depth", 0, "limit scan depth (0 = unlimited)")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func runRoot(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// If a specific format is requested or stdout is not a TTY, use scan mode.
	if flags.Format != "" || !isTerminal() {
		return runScan(cmd, args)
	}

	// Launch TUI.
	cfg, err := config.New(path, flags)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	app := tui.NewApp(path, cfg)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running TUI: %w", err)
	}
	return nil
}

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
