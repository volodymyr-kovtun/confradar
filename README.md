<p align="center">
  <img src="assets/logo.svg" alt="confradar" width="420">
</p>

<p align="center">
  <strong>Instantly see every config file in your project.</strong><br>
  Spot missing env vars. Diff environments. Catch issues before they hit production.
</p>

<p align="center">
  <a href="#install">Install</a> &middot;
  <a href="#quick-start">Quick Start</a> &middot;
  <a href="#commands">Commands</a> &middot;
  <a href="#tui">TUI</a> &middot;
  <a href="#health-checks">Health Checks</a> &middot;
  <a href="#configuration">Configuration</a>
</p>

---

## What is confradar?

confradar scans any software project and presents a unified, categorized view of every configuration file — `.env` files, Docker configs, CI/CD pipelines, build tool configs, linting rules, infrastructure-as-code, and more.

```
$ confradar scan

my-project (18 config files found in 3ms)

🔐 Environment (4)
  ├── .env
  ├── .env.example
  ├── .env.production
  └── .env.staging

🐳 Docker (3)
  ├── Dockerfile
  ├── docker-compose.yml
  └── docker-compose.production.yml

⚙️ CI/CD (1)
  └── .github/workflows/deploy.yml

🌐 Web Server (1)
  └── nginx/default.conf

📦 JavaScript / Node (4)
  ├── .nvmrc
  ├── next.config.js
  ├── package.json
  └── tsconfig.json

✨ Linting / Formatting (2)
  ├── .eslintrc.json
  └── .prettierrc

─── 18 files across 6 categories ───
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
git clone https://github.com/volodymyr-kovtun/confradar.git
cd confradar
go build -o confradar .
sudo mv confradar /usr/local/bin/
```

### With `go install`

```bash
go install github.com/volodymyrkovtun/confradar@latest
```

### Homebrew (after release)

```bash
brew install volodymyr-kovtun/tap/confradar
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

### Flags

```
--format <fmt>       Output format: tree, json, yaml, markdown, table
--no-color           Disable colored output
--config <path>      Use a specific config file
--no-config          Ignore all config files, use built-in defaults only
--severity <level>   Filter health issues: error, warning, info
--health-only        Show only health issues
--max-depth <n>      Limit scan depth
--verbose            Show debug info
--quiet              Suppress output except errors
```

## TUI

Run `confradar` in any project directory to launch the interactive terminal UI.

### Keybindings

**Navigation**

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate tree |
| `Enter` | Expand category / view file |
| `Esc` | Go back / close overlay |
| `Tab` | Switch between panels |
| `1`-`9` | Jump to category by number |
| `q` | Quit |

**Actions**

| Key | Action |
|-----|--------|
| `/` | Search files |
| `d` | Diff two .env files |
| `h` | Toggle health panel |
| `e` | Open in `$EDITOR` |
| `y` | Copy path to clipboard |
| `r` | Rescan project |
| `?` | Show help |

**Display**

| Key | Action |
|-----|--------|
| `c` | Cycle theme (dark / light / minimal) |
| `v` | Toggle .env value redaction |
| `p` | Toggle preview panel |
| `s` | Cycle sort order |

## Health Checks

confradar detects common configuration issues automatically and via custom rules.

### Built-in (automatic)

- `.env.example` out of sync with `.env`
- Dockerfile using `:latest` tag
- Dockerfile missing `USER` directive

### Configurable (via `.confradar.yml`)

| Type | Description |
|------|-------------|
| `env_sync` | Compare keys between .env files |
| `key_exists` | Verify required keys exist |
| `port_conflict` | Detect port conflicts across configs |
| `version_match` | Check runtime version consistency |
| `file_exists` | Verify required files exist |
| `regex_check` | Custom regex pattern matching |
| `dockerfile_best_practices` | Dockerfile analysis |

### Example

```yaml
health_checks:
  - type: "env_sync"
    source: ".env"
    targets: [".env.production"]
    severity: "warning"

  - type: "key_exists"
    file: ".env.production"
    required_keys: ["DATABASE_URL", "SECRET_KEY"]
    severity: "error"
```

### CI/CD Integration

```yaml
# GitHub Actions
- name: Check config health
  run: confradar check --severity error
```

Exit code 1 when error-severity issues are found.

## Configuration

Run `confradar init` to create a `.confradar.yml`:

```yaml
# Add custom patterns
extra_categories:
  - name: "My Configs"
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

# Override display
override_categories:
  - name: "Docker"
    icon: "🐋"

# Hide categories
disable_categories:
  - "Editor / DX"

# Ignore paths
ignore:
  - "vendor/**"
  - "legacy/**"
```

### Three config layers (merged in order)

1. **Built-in defaults** — 119 patterns embedded in the binary
2. **Global** — `~/.config/confradar/config.yml` (all projects)
3. **Project** — `.confradar.yml` in project root (highest priority)

### Environment variable overrides

```bash
CONFRADAR_FORMAT=json confradar scan
CONFRADAR_DISPLAY_THEME=minimal confradar
```

## Output Formats

```bash
confradar scan --format tree .     # default, pretty tree
confradar scan --format json .     # structured JSON
confradar scan --format yaml .     # YAML
confradar scan --format table .    # aligned table
confradar report .                 # Markdown report
```

## Built-in Categories

Environment, Docker, CI/CD, Web Server, Build Tools, JavaScript/Node, Python, Go, Rust, Linting/Formatting, Git, Infrastructure as Code, Testing, Security, Editor/DX, Documentation

## License

MIT
