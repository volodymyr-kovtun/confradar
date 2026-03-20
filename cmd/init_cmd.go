package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter .confradar.yml in the current directory",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

const starterConfig = `# confradar configuration
# See: https://github.com/volodymyrkovtun/confradar

# ── EXTENDING PATTERNS ──────────────────────────────────
# Add entirely new categories
# extra_categories:
#   - name: "My Custom Category"
#     icon: "📁"
#     color: "#ff6b6b"
#     priority: 50
#     patterns:
#       - glob: "config/api_*.json"
#         description: "API gateway configs"

# Add patterns to existing built-in categories
# extra_patterns:
#   - category: "Environment"
#     patterns:
#       - glob: ".env.secrets"
#       - glob: "config/.env.*"

# ── OVERRIDING CATEGORIES ───────────────────────────────
# override_categories:
#   - name: "Docker"
#     icon: "🐋"
#     color: "#0db7ed"

# Disable categories entirely (hide from output)
# disable_categories:
#   - "Editor / DX"

# ── IGNORE RULES ────────────────────────────────────────
# ignore:
#   - "legacy/**"
#   - "vendor/**"

# skip_dirs:
#   - "generated"
#   - ".cache"

# ── HEALTH CHECKS ──────────────────────────────────────
# health_checks:
#   - type: "env_sync"
#     source: ".env"
#     targets: [".env.production", ".env.staging"]
#     severity: "warning"
#
#   - type: "key_exists"
#     file: ".env.production"
#     required_keys: ["DATABASE_URL", "REDIS_URL"]
#     severity: "error"
#
#   - type: "file_exists"
#     files: [".env.example", ".editorconfig"]
#     severity: "warning"

# ── DISPLAY PREFERENCES ────────────────────────────────
# display:
#   theme: "auto"
#   redact_env_values: true
#   category_sort: "priority"
#   file_sort: "name"

# ── OUTPUT FORMATS ──────────────────────────────────────
# output:
#   default_format: "tree"

# ── ALIASES ─────────────────────────────────────────────
# aliases:
#   prod-check: "check --health-only --severity error"
#   env-diff: "diff .env .env.production"
`

func runInit(cmd *cobra.Command, args []string) error {
	target := filepath.Join(".", ".confradar.yml")

	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("%s already exists; remove it first to reinitialize", target)
	}

	if err := os.WriteFile(target, []byte(starterConfig), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", target, err)
	}

	fmt.Printf("Created %s — customize it for your project!\n", target)
	return nil
}
