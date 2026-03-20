package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/volodymyrkovtun/confradar/internal/differ"
)

var diffCmd = &cobra.Command{
	Use:   "diff <file1> <file2>",
	Short: "Diff two .env files side by side",
	Long:  "Compare two .env files and show missing, extra, and changed keys.",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	result, err := differ.Diff(args[0], args[1])
	if err != nil {
		return fmt.Errorf("diffing files: %w", err)
	}

	redact := true // Default to redacting values for security.
	differ.RenderDiff(result, os.Stdout, flags.NoColor, redact)
	return nil
}
