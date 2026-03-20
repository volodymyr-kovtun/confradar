# CLAUDE.md — Development Guide for confradar

## Project Description

confradar is a CLI tool that scans any project directory and presents a unified, categorized view of every configuration file — .env files, Docker configs, CI/CD pipelines, build tool configs, linting rules, and more. It includes an interactive TUI mode (Bubble Tea), structured output formats, health checks, and env file diffing.

## Build

```bash
go build -o confradar .
```

## Test

```bash
go test ./...
```

## Run

```bash
# Launch the interactive TUI (default when stdout is a terminal)
./confradar .

# Scan and print a categorized tree of config files
./confradar scan .

# Run health checks and report issues
./confradar check .

# Diff two .env files side by side
./confradar diff .env .env.prod

# Generate a Markdown/JSON report
./confradar report . --output report.md

# Create a starter .confradar.yml config file
./confradar init

# Show all active patterns (built-in + custom)
./confradar list-patterns

# Print version info
./confradar version
```

### Useful Flags

- `--format <tree|json|yaml|markdown|table>` — output format for scan/report
- `--no-color` — disable colored output
- `--config <path>` — path to a custom .confradar.yml config file
- `--no-config` — ignore all config files, use built-in defaults only
- `--severity <error|warning|info>` — filter health check issues by severity
- `--health-only` — only show health issues
- `--include-hidden` — include dotfiles normally skipped
- `--max-depth <n>` — limit scan depth (0 = unlimited)
- `--verbose` — show debug info
- `--quiet` — suppress all output except errors

## Project Structure

```
.
├── main.go                          # Entrypoint — calls cmd.Execute()
├── cmd/                             # CLI commands (Cobra)
│   ├── root.go                      # Root command, TUI launcher, persistent flags
│   ├── scan.go                      # `scan` subcommand
│   ├── check.go                     # `check` subcommand (health checks)
│   ├── diff.go                      # `diff` subcommand (env file comparison)
│   ├── report.go                    # `report` subcommand (Markdown/JSON reports)
│   ├── init_cmd.go                  # `init` subcommand (scaffold .confradar.yml)
│   ├── list_patterns.go             # `list-patterns` subcommand
│   └── version.go                   # `version` subcommand, build-time vars
├── internal/
│   ├── config/                      # Configuration loading and schema
│   │   ├── config.go                # Config constructor and merging logic
│   │   ├── schema.go                # Config struct definitions (yaml tags)
│   │   ├── loader.go                # File loading and defaults
│   │   └── env_override.go          # Environment variable overrides
│   ├── scanner/                     # Filesystem scanning engine
│   │   ├── scanner.go               # Core scan logic
│   │   ├── matcher.go               # Glob pattern matching
│   │   ├── types.go                 # ScanResult, ConfigFile types
│   │   └── gitignore.go             # .gitignore integration
│   ├── parser/                      # Config file parsers
│   │   ├── parser.go                # Parser registry/interface
│   │   ├── env.go                   # .env parser
│   │   ├── yaml.go                  # YAML parser
│   │   ├── json.go                  # JSON parser
│   │   ├── toml.go                  # TOML parser
│   │   ├── ini.go                   # INI parser
│   │   ├── dockerfile.go            # Dockerfile parser
│   │   ├── nginx.go                 # Nginx config parser
│   │   └── text.go                  # Plain text fallback parser
│   ├── health/                      # Health check engine
│   │   ├── checker.go               # Check runner and auto-checks
│   │   ├── types.go                 # HealthIssue, severity constants
│   │   ├── env_sync.go              # Env file sync check
│   │   ├── env_example_sync.go      # .env.example sync check
│   │   ├── key_exists.go            # Required key existence check
│   │   ├── port_conflict.go         # Port conflict detection
│   │   ├── version_match.go         # Version consistency check
│   │   ├── file_exists.go           # Required file existence check
│   │   ├── regex_check.go           # Regex-based value validation
│   │   └── dockerfile_check.go      # Dockerfile best-practice checks
│   ├── differ/                      # Env file diffing
│   │   ├── env_diff.go              # Diff logic and rendering
│   │   └── types.go                 # Diff result types
│   ├── renderer/                    # Output renderers
│   │   ├── renderer.go              # Renderer interface and factory
│   │   ├── tree.go                  # Tree output
│   │   ├── table.go                 # Table output
│   │   ├── json.go                  # JSON output
│   │   ├── yaml.go                  # YAML output
│   │   └── markdown.go              # Markdown report output
│   └── tui/                         # Interactive TUI (Bubble Tea)
│       ├── app.go                   # Main TUI model
│       ├── tree_view.go             # File tree panel
│       ├── file_viewer.go           # File content viewer panel
│       ├── diff_view.go             # Diff view panel
│       ├── health_panel.go          # Health issues panel
│       ├── search.go                # Search overlay
│       ├── help.go                  # Help overlay
│       ├── statusbar.go             # Status bar component
│       ├── theme.go                 # Color theme definitions
│       └── styles.go                # Lipgloss style helpers
├── patterns/                        # Built-in pattern definitions
│   ├── default.yml                  # Default category/pattern rules
│   └── embed.go                     # go:embed for pattern files
├── confradar                        # Compiled binary (gitignored)
├── go.mod
└── go.sum
```

## Key Conventions

- **Go version**: 1.26 (see go.mod)
- **CLI framework**: [Cobra](https://github.com/spf13/cobra) — all commands live in `cmd/`
- **TUI framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) with [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- **Config serialization**: struct fields use `yaml:"..."` tags; config files are YAML
- **Pattern definitions**: YAML files in `patterns/`, embedded at compile time via `go:embed`
- **Build-time injection**: `cmd.version`, `cmd.commit`, `cmd.date` are set via `-ldflags` at build time
- **Error handling**: commands return errors via `RunE`; root command silences usage on error
- **Output formats**: tree (default), json, yaml, markdown, table — selectable via `--format`
- **No CGO**: builds use `CGO_ENABLED=0` for static binaries

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI command framework |
| `github.com/charmbracelet/bubbletea` | Terminal UI framework |
| `github.com/charmbracelet/lipgloss` | TUI styling and layout |
| `github.com/pelletier/go-toml/v2` | TOML config parsing |
| `gopkg.in/yaml.v3` | YAML config parsing |
| `github.com/bmatcuk/doublestar/v4` | Glob pattern matching |
| `golang.org/x/term` | Terminal detection (isatty) |
