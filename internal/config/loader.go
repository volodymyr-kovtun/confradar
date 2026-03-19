package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/volodymyrkovtun/confradar/patterns"
	"gopkg.in/yaml.v3"
)

// Load reads and merges all configuration layers:
//  1. Built-in defaults (embedded in binary)
//  2. Global user config (~/.config/confradar/config.yml)
//  3. Project config (.confradar.yml in project root)
//
// The flags parameter controls overrides from CLI arguments.
func Load(projectPath string, flags CLIFlags) (*Config, error) {
	base, err := loadDefaults()
	if err != nil {
		return nil, fmt.Errorf("loading built-in defaults: %w", err)
	}

	if flags.NoConfig {
		return base, nil
	}

	if flags.ConfigPath != "" {
		overlay, err := loadFile(flags.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("loading config %s: %w", flags.ConfigPath, err)
		}
		base = merge(base, overlay)
		return base, nil
	}

	if global, err := loadGlobal(); err == nil {
		base = merge(base, global)
	}

	if project, err := loadProject(projectPath); err == nil {
		base = merge(base, project)
	}

	return base, nil
}

func loadDefaults() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(patterns.DefaultYAML, &cfg); err != nil {
		return nil, fmt.Errorf("parsing default patterns: %w", err)
	}
	return &cfg, nil
}

func loadGlobal() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return loadFile(filepath.Join(home, ".config", "confradar", "config.yml"))
}

func loadProject(projectPath string) (*Config, error) {
	return loadFile(filepath.Join(projectPath, ".confradar.yml"))
}

func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return &cfg, nil
}

// merge combines a base config with an overlay, following the documented merge strategy.
func merge(base, overlay *Config) *Config {
	result := *base

	// extra_categories: append to existing list
	if len(overlay.ExtraCategories) > 0 {
		result.Categories = append(result.Categories, overlay.ExtraCategories...)
	}
	if len(overlay.Categories) > 0 {
		result.Categories = append(result.Categories, overlay.Categories...)
	}

	// extra_patterns: append patterns to matching categories
	for _, ep := range overlay.ExtraPatterns {
		for i, cat := range result.Categories {
			if cat.Name == ep.Category {
				result.Categories[i].Patterns = append(result.Categories[i].Patterns, ep.Patterns...)
				break
			}
		}
	}

	// override_categories: match by name, replace non-zero fields
	for _, oc := range overlay.OverrideCategories {
		for i, cat := range result.Categories {
			if cat.Name == oc.Name {
				if oc.Icon != "" {
					result.Categories[i].Icon = oc.Icon
				}
				if oc.Color != "" {
					result.Categories[i].Color = oc.Color
				}
				if oc.Priority != nil {
					result.Categories[i].Priority = *oc.Priority
				}
				break
			}
		}
	}

	// disable_categories: remove by name
	if len(overlay.DisableCategories) > 0 {
		disabled := make(map[string]bool, len(overlay.DisableCategories))
		for _, name := range overlay.DisableCategories {
			disabled[name] = true
		}
		filtered := result.Categories[:0]
		for _, cat := range result.Categories {
			if !disabled[cat.Name] {
				filtered = append(filtered, cat)
			}
		}
		result.Categories = filtered
	}

	// ignore: union of all globs
	if len(overlay.Ignore) > 0 {
		seen := make(map[string]bool)
		for _, g := range result.Ignore {
			seen[g] = true
		}
		for _, g := range overlay.Ignore {
			if !seen[g] {
				result.Ignore = append(result.Ignore, g)
			}
		}
	}

	// skip_dirs: union
	if len(overlay.SkipDirs) > 0 {
		seen := make(map[string]bool)
		for _, d := range result.SkipDirs {
			seen[d] = true
		}
		for _, d := range overlay.SkipDirs {
			if !seen[d] {
				result.SkipDirs = append(result.SkipDirs, d)
			}
		}
	}

	// health_checks: append
	result.HealthChecks = append(result.HealthChecks, overlay.HealthChecks...)

	// display: deep merge (non-zero overlay fields overwrite)
	mergeDisplay(&result.Display, &overlay.Display)

	// output: deep merge
	mergeOutput(&result.Output, &overlay.Output)

	// aliases: merge maps
	if len(overlay.Aliases) > 0 {
		if result.Aliases == nil {
			result.Aliases = make(map[string]string)
		}
		for k, v := range overlay.Aliases {
			result.Aliases[k] = v
		}
	}

	return &result
}

func mergeDisplay(base, overlay *DisplayConfig) {
	if overlay.Theme != "" {
		base.Theme = overlay.Theme
	}
	if overlay.ShowPreview != nil {
		base.ShowPreview = overlay.ShowPreview
	}
	if overlay.PreviewLines > 0 {
		base.PreviewLines = overlay.PreviewLines
	}
	if overlay.ShowMetadata != nil {
		base.ShowMetadata = overlay.ShowMetadata
	}
	if overlay.ShowHiddenCount != nil {
		base.ShowHiddenCount = overlay.ShowHiddenCount
	}
	if overlay.RedactEnvValues != nil {
		base.RedactEnvValues = overlay.RedactEnvValues
	}
	if overlay.RedactPattern != "" {
		base.RedactPattern = overlay.RedactPattern
	}
	if overlay.CategorySort != "" {
		base.CategorySort = overlay.CategorySort
	}
	if overlay.FileSort != "" {
		base.FileSort = overlay.FileSort
	}
	if overlay.CollapseThreshold > 0 {
		base.CollapseThreshold = overlay.CollapseThreshold
	}
}

func mergeOutput(base, overlay *OutputConfig) {
	if overlay.DefaultFormat != "" {
		base.DefaultFormat = overlay.DefaultFormat
	}
	if overlay.JSON.IncludeFileContents != nil {
		base.JSON.IncludeFileContents = overlay.JSON.IncludeFileContents
	}
	if overlay.JSON.IncludeParsedKeys != nil {
		base.JSON.IncludeParsedKeys = overlay.JSON.IncludeParsedKeys
	}
	if overlay.JSON.IncludeHealthIssues != nil {
		base.JSON.IncludeHealthIssues = overlay.JSON.IncludeHealthIssues
	}
	if overlay.JSON.Pretty != nil {
		base.JSON.Pretty = overlay.JSON.Pretty
	}
	if overlay.Markdown.IncludeTOC != nil {
		base.Markdown.IncludeTOC = overlay.Markdown.IncludeTOC
	}
	if overlay.Markdown.IncludeHealthSummary != nil {
		base.Markdown.IncludeHealthSummary = overlay.Markdown.IncludeHealthSummary
	}
	if overlay.Markdown.IncludeFileList != nil {
		base.Markdown.IncludeFileList = overlay.Markdown.IncludeFileList
	}
	if overlay.Markdown.IncludeEnvDiff != nil {
		base.Markdown.IncludeEnvDiff = overlay.Markdown.IncludeEnvDiff
	}
	if overlay.Markdown.Template != "" {
		base.Markdown.Template = overlay.Markdown.Template
	}
}
