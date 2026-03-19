package config

import "sort"

// New loads the effective configuration for the given project path and CLI flags.
func New(projectPath string, flags CLIFlags) (*Config, error) {
	cfg, err := Load(projectPath, flags)
	if err != nil {
		return nil, err
	}

	// Apply environment variable overrides.
	ApplyEnvOverrides(cfg)

	// Apply display defaults.
	if cfg.Display.Theme == "" {
		cfg.Display.Theme = "auto"
	}
	if cfg.Display.PreviewLines == 0 {
		cfg.Display.PreviewLines = 5
	}
	if cfg.Display.RedactPattern == "" {
		cfg.Display.RedactPattern = "••••••"
	}
	if cfg.Display.CategorySort == "" {
		cfg.Display.CategorySort = "priority"
	}
	if cfg.Display.FileSort == "" {
		cfg.Display.FileSort = "name"
	}
	if cfg.Display.CollapseThreshold == 0 {
		cfg.Display.CollapseThreshold = 15
	}

	// Apply output defaults.
	if cfg.Output.DefaultFormat == "" {
		cfg.Output.DefaultFormat = "tree"
	}

	// Sort categories by priority.
	sort.Slice(cfg.Categories, func(i, j int) bool {
		if cfg.Categories[i].Priority == cfg.Categories[j].Priority {
			return cfg.Categories[i].Name < cfg.Categories[j].Name
		}
		return cfg.Categories[i].Priority < cfg.Categories[j].Priority
	})

	return cfg, nil
}

// EffectiveFormat returns the output format to use, with CLI flag taking precedence.
func EffectiveFormat(cfg *Config, flags CLIFlags) string {
	if flags.Format != "" {
		return flags.Format
	}
	return cfg.Output.DefaultFormat
}
