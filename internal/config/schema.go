// Package config handles loading, merging, and validating confradar configuration.
package config

// Config is the top-level configuration after merging all three layers
// (built-in defaults, global user config, project config).
type Config struct {
	Categories         []Category        `yaml:"categories"`
	ExtraCategories    []Category        `yaml:"extra_categories"`
	OverrideCategories []CategoryOverride `yaml:"override_categories"`
	DisableCategories  []string           `yaml:"disable_categories"`
	ExtraPatterns      []ExtraPattern     `yaml:"extra_patterns"`
	Ignore             []string           `yaml:"ignore"`
	SkipDirs           []string           `yaml:"skip_dirs"`
	HealthChecks       []HealthCheckRule  `yaml:"health_checks"`
	Display            DisplayConfig      `yaml:"display"`
	Output             OutputConfig       `yaml:"output"`
	Aliases            map[string]string  `yaml:"aliases"`
}

// Category groups related config file patterns under a common name with display properties.
type Category struct {
	Name     string    `yaml:"name"`
	Icon     string    `yaml:"icon"`
	Color    string    `yaml:"color"`
	Priority int       `yaml:"priority"`
	Patterns []Pattern `yaml:"patterns"`
}

// CategoryOverride selectively replaces fields on a built-in category matched by name.
type CategoryOverride struct {
	Name     string `yaml:"name"`
	Icon     string `yaml:"icon,omitempty"`
	Color    string `yaml:"color,omitempty"`
	Priority *int   `yaml:"priority,omitempty"`
}

// Pattern defines a glob-based rule for detecting a config file.
type Pattern struct {
	Name        string   `yaml:"name"`
	Glob        string   `yaml:"glob,omitempty"`
	Globs       []string `yaml:"globs,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Parser      string   `yaml:"parser,omitempty"`
	MaxDepth    int      `yaml:"max_depth,omitempty"`
}

// AllGlobs returns every glob string for this pattern, combining both
// the singular Glob field and the plural Globs field.
func (p Pattern) AllGlobs() []string {
	var out []string
	if p.Glob != "" {
		out = append(out, p.Glob)
	}
	out = append(out, p.Globs...)
	return out
}

// ExtraPattern adds patterns to an existing category by name.
type ExtraPattern struct {
	Category string    `yaml:"category"`
	Patterns []Pattern `yaml:"patterns"`
}

// HealthCheckRule defines a single health check to run during scanning.
type HealthCheckRule struct {
	Type         string   `yaml:"type"`
	Source       string   `yaml:"source,omitempty"`
	Targets      []string `yaml:"targets,omitempty"`
	File         string   `yaml:"file,omitempty"`
	Files        []string `yaml:"files,omitempty"`
	RequiredKeys []string `yaml:"required_keys,omitempty"`
	IgnoreKeys   []string `yaml:"ignore_keys,omitempty"`
	Severity     string   `yaml:"severity"`
	TypeName     string   `yaml:"type_name,omitempty"`
	Pattern      string   `yaml:"pattern,omitempty"`
	Message      string   `yaml:"message,omitempty"`
}

// DisplayConfig controls how confradar renders output.
type DisplayConfig struct {
	Theme              string `yaml:"theme"`
	ShowPreview        *bool  `yaml:"show_preview,omitempty"`
	PreviewLines       int    `yaml:"preview_lines,omitempty"`
	ShowMetadata       *bool  `yaml:"show_metadata,omitempty"`
	ShowHiddenCount    *bool  `yaml:"show_hidden_count,omitempty"`
	RedactEnvValues    *bool  `yaml:"redact_env_values,omitempty"`
	RedactPattern      string `yaml:"redact_pattern,omitempty"`
	CategorySort       string `yaml:"category_sort,omitempty"`
	FileSort           string `yaml:"file_sort,omitempty"`
	CollapseThreshold  int    `yaml:"collapse_threshold,omitempty"`
}

// OutputConfig controls the format and options for CLI output.
type OutputConfig struct {
	DefaultFormat string         `yaml:"default_format,omitempty"`
	JSON          JSONConfig     `yaml:"json,omitempty"`
	Markdown      MarkdownConfig `yaml:"markdown,omitempty"`
}

// JSONConfig controls JSON output behavior.
type JSONConfig struct {
	IncludeFileContents *bool `yaml:"include_file_contents,omitempty"`
	IncludeParsedKeys   *bool `yaml:"include_parsed_keys,omitempty"`
	IncludeHealthIssues *bool `yaml:"include_health_issues,omitempty"`
	Pretty              *bool `yaml:"pretty,omitempty"`
}

// MarkdownConfig controls markdown report output.
type MarkdownConfig struct {
	IncludeTOC           *bool  `yaml:"include_toc,omitempty"`
	IncludeHealthSummary *bool  `yaml:"include_health_summary,omitempty"`
	IncludeFileList      *bool  `yaml:"include_file_list,omitempty"`
	IncludeEnvDiff       *bool  `yaml:"include_env_diff,omitempty"`
	Template             string `yaml:"template,omitempty"`
}

// CLIFlags holds values parsed from command-line flags.
type CLIFlags struct {
	ConfigPath    string
	NoConfig      bool
	Format        string
	NoColor       bool
	Verbose       bool
	Quiet         bool
	Severity      string
	HealthOnly    bool
	IncludeHidden bool
	MaxDepth      int
}

// BoolVal returns the value of a *bool pointer, falling back to a default.
func BoolVal(p *bool, def bool) bool {
	if p != nil {
		return *p
	}
	return def
}

// BoolPtr returns a pointer to a bool value.
func BoolPtr(v bool) *bool {
	return &v
}
