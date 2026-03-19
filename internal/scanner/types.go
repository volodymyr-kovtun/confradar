// Package scanner walks a project directory and matches files against config patterns.
package scanner

import (
	"time"
)

// ConfigFile represents a single detected configuration file.
type ConfigFile struct {
	Path        string    `json:"path" yaml:"path"`
	RelPath     string    `json:"rel_path" yaml:"rel_path"`
	Category    string    `json:"category" yaml:"category"`
	PatternName string    `json:"pattern_name" yaml:"pattern_name"`
	Parser      string    `json:"parser" yaml:"parser"`
	Size        int64     `json:"size" yaml:"size"`
	ModTime     time.Time `json:"mod_time" yaml:"mod_time"`
}

// CategoryResult holds all files matched under a single category.
type CategoryResult struct {
	Name     string       `json:"name" yaml:"name"`
	Icon     string       `json:"icon" yaml:"icon"`
	Color    string       `json:"color" yaml:"color"`
	Priority int          `json:"priority" yaml:"priority"`
	Files    []ConfigFile `json:"files" yaml:"files"`
}

// ScanResult is the complete output of a project scan.
type ScanResult struct {
	Root       string                    `json:"root" yaml:"root"`
	Categories map[string]*CategoryResult `json:"categories" yaml:"categories"`
	// Ordered preserves category display order (sorted by priority, then name).
	Ordered    []*CategoryResult `json:"-" yaml:"-"`
	TotalFiles int               `json:"total_files" yaml:"total_files"`
	Duration   time.Duration     `json:"duration" yaml:"duration"`
}
