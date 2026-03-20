// Package parser extracts structured metadata from configuration files.
package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Parser extracts metadata from a config file.
type Parser interface {
	Name() string
	Parse(path string) (*ParseResult, error)
}

// ParseResult holds extracted metadata from a parsed config file.
type ParseResult struct {
	Keys     []string          `json:"keys,omitempty"`
	Values   map[string]string `json:"values,omitempty"`
	Ports    []int             `json:"ports,omitempty"`
	Versions map[string]string `json:"versions,omitempty"`
	Services []string          `json:"services,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

var registry = map[string]Parser{}

// Register adds a parser to the global registry.
func Register(p Parser) {
	registry[p.Name()] = p
}

// ForFile returns the appropriate parser for a file based on its parser type tag
// or filename extension.
func ForFile(parserType, path string) Parser {
	if p, ok := registry[parserType]; ok {
		return p
	}

	ext := strings.ToLower(filepath.Ext(path))
	base := strings.ToLower(filepath.Base(path))

	switch {
	case strings.HasPrefix(base, ".env") || base == ".envrc" || base == ".secrets":
		return registry["env"]
	case ext == ".yml" || ext == ".yaml":
		return registry["yaml"]
	case ext == ".json" || ext == ".json5":
		return registry["json"]
	case ext == ".toml":
		return registry["toml"]
	case ext == ".ini" || ext == ".cfg" || ext == ".conf":
		return registry["ini"]
	case base == "dockerfile" || strings.HasPrefix(base, "dockerfile."):
		return registry["dockerfile"]
	case strings.Contains(base, "nginx"):
		return registry["nginx"]
	}

	return nil
}

// ParseFile is a convenience function that finds the right parser and runs it.
func ParseFile(parserType, path string) (*ParseResult, error) {
	p := ForFile(parserType, path)
	if p == nil {
		return &ParseResult{}, nil
	}
	result, err := p.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return result, nil
}
