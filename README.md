# confradar

**Instantly see every config file in your project. Spot missing env vars. Diff environments. Catch issues before they hit production.**

confradar scans any software project directory and presents a unified, categorized view of every configuration file — `.env` files, Docker configs, CI/CD pipelines, build tool configs, linting rules, infrastructure-as-code, and more.

```
$ confradar scan

Subero (12 config files found in 10ms)

🔐 Environment (5)
  ├── .env.example
  ├── backend/.env.example
  ├── frontend/.env.example
  ├── .env.production.example
  └── .env.staging.example

🐳 Docker (3)
  ├── compose.ghcr.yml
  ├── compose.production.yml
  └── compose.yml

⚙️ CI/CD (1)
  └── .github/workflows/manual-deploy.yml

📝 Git (1)
  └── .gitignore

💻 Editor / DX (2)
  ├── CLAUDE.md
  └── .claude/settings.local.json

─── 12 files across 5 categories ───
```

## Features

- **119 built-in patterns** across 16 categories — works instantly on any project
- **Interactive TUI** — navigate your config tree, preview files, diff environments
- **Health checks** — detect missing env vars, port conflicts, version mismatches, Dockerfile issues
- **Env diffing** — compare `.env` vs `.env.production` side-by-side with redacted values
- **Multiple output formats** — tree, JSON, YAML, Markdown, table
- **Fully customizable** — add patterns, categories, health checks via `.confradar.yml`
- **Zero config required** — just run `confradar` in any project

## Install

### From source (requires Go 1.22+)

```bash
git clone https://github.com/volodymyrkovtun/confradar.git
cd confradar
go build -o confradar .
sudo mv confradar /usr/local/bin/  # optional: make available globally
```

### With `go install`

```bash
go install github.com/volodymyrkovtun/confradar@latest
```

### Homebrew (after release with GoReleaser)

```bash
brew install volodymyrkovtun/tap/confradar
```

## Quick Start

```bash
# Scan current directory
confradar scan

# Launch interactive TUI
confradar

# Check for health issues (great for CI)
confradar check

# Diff two .env files
confradar diff .env .env.production

# Generate a Markdown report
confradar report --output CONFIGS.md
```

## Commands

| Command | Description |
|---------|-------------|
| `confradar` | Launch TUI (or scan if not a terminal) |
| `confradar scan [path]` | Print categorized config tree |
| `confradar diff <file1> <file2>` | Compare two .env files side-by-side |
| `confradar check [path]` | Run health checks (exit 1 on errors) |
| `confradar report [path]` | Generate Markdown/JSON report |
| `confradar list-patterns` | Show all 119 active patterns |
| `confradar init` | Create starter `.confradar.yml` |
| `confradar version` | Print version info |

## Flags

| Flag | Description |
|------|-------------|
| `--format <fmt>` | Output format: `tree`, `json`, `yaml`, `markdown`, `table` |
| `--no-color` | Disable colored output |
| `--config <path>` | Use a specific config file |
| `--no-config` | Ignore all config files, use built-in defaults only |
| `--severity <level>` | Filter health issues: `error`, `warning`, `info` |
| `--health-only` | Show only health issues |
| `--max-depth <n>` | Limit scan depth |
| `--verbose` | Show debug info |
| `--quiet` | Suppress output except errors |

## TUI Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate tree |
| `Enter` | Expand category / view file |
| `Esc` | Go back |
| `Tab` | Switch panels |
| `/` | Search files |
| `d` | Diff two .env files |
| `h` | Toggle health panel |
| `?` | Show help |
| `c` | Cycle theme (dark/light/minimal) |
| `v` | Toggle value redaction |
| `e` | Open in $EDITOR |
| `y` | Copy path to clipboard |
| `r` | Rescan |
| `q` | Quit |

## Configuration

Run `confradar init` to create a `.confradar.yml` in your project:

```yaml
# Add custom patterns
extra_categories:
  - name: "My Custom Category"
    icon: "📁"
    color: "#ff6b6b"
    priority: 50
    patterns:
      - glob: "config/*.json"
        description: "Custom config files"

# Add patterns to existing categories
extra_patterns:
  - category: "Environment"
    patterns:
      - glob: ".env.secrets"

# Override category display
override_categories:
  - name: "Docker"
    icon: "🐋"

# Hide categories
disable_categories:
  - "Editor / DX"

# Health checks
health_checks:
  - type: "env_sync"
    source: ".env"
    targets: [".env.production", ".env.staging"]
    severity: "warning"

  - type: "key_exists"
    file: ".env.production"
    required_keys: ["DATABASE_URL", "SECRET_KEY"]
    severity: "error"

  - type: "file_exists"
    files: [".env.example", ".editorconfig"]
    severity: "warning"

  - type: "port_conflict"
    files: ["docker-compose.yml", ".env"]
    severity: "warning"
```

## Health Check Types

| Type | Description |
|------|-------------|
| `env_sync` | Compare keys between .env files |
| `key_exists` | Verify required keys exist |
| `port_conflict` | Detect port conflicts across configs |
| `version_match` | Check runtime version consistency |
| `file_exists` | Verify required files exist |
| `regex_check` | Custom regex pattern matching |
| `env_example_sync` | Check .env.example is up to date |
| `dockerfile_best_practices` | :latest tags, missing USER, etc. |

## CI/CD Integration

Use `confradar check` in your CI pipeline:

```yaml
# GitHub Actions
- name: Check config health
  run: confradar check --severity error
```

Environment variable overrides (no config file needed):

```bash
CONFRADAR_FORMAT=json confradar scan
CONFRADAR_DISPLAY_THEME=minimal confradar
```

## Built-in Categories

Environment, Docker, CI/CD, Web Server, Build Tools, JavaScript/Node, Python, Go, Rust, .NET, Linting/Formatting, Git, Infrastructure as Code, Testing, Security, Editor/DX, Documentation

## License

MIT
