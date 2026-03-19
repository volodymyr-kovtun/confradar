package config

import (
	"os"
	"strconv"
	"strings"
)

// ApplyEnvOverrides reads CONFRADAR_* environment variables and applies them
// to the configuration. This allows CI/CD pipelines to override settings
// without a config file.
func ApplyEnvOverrides(cfg *Config) {
	if v := os.Getenv("CONFRADAR_DISPLAY_THEME"); v != "" {
		cfg.Display.Theme = v
	}
	if v := os.Getenv("CONFRADAR_DISPLAY_REDACT_ENV_VALUES"); v != "" {
		b := v == "true" || v == "1"
		cfg.Display.RedactEnvValues = &b
	}
	if v := os.Getenv("CONFRADAR_DISPLAY_CATEGORY_SORT"); v != "" {
		cfg.Display.CategorySort = v
	}
	if v := os.Getenv("CONFRADAR_DISPLAY_FILE_SORT"); v != "" {
		cfg.Display.FileSort = v
	}
	if v := os.Getenv("CONFRADAR_OUTPUT_DEFAULT_FORMAT"); v != "" {
		cfg.Output.DefaultFormat = v
	}
	if v := os.Getenv("CONFRADAR_FORMAT"); v != "" {
		cfg.Output.DefaultFormat = v
	}
	if v := os.Getenv("CONFRADAR_MAX_DEPTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Display.CollapseThreshold = n
		}
	}
	if v := os.Getenv("CONFRADAR_IGNORE"); v != "" {
		extra := strings.Split(v, ",")
		for _, g := range extra {
			g = strings.TrimSpace(g)
			if g != "" {
				cfg.Ignore = append(cfg.Ignore, g)
			}
		}
	}
}
