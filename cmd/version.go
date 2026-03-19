package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build-time variables injected via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of confradar",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("confradar %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
