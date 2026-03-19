package scanner

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/volodymyrkovtun/confradar/internal/config"
)

// match holds a pre-compiled pattern entry for efficient matching.
type match struct {
	categoryName string
	categoryIcon string
	categoryColor string
	priority     int
	patternName  string
	parser       string
	globs        []string
}

// Matcher tests file paths against all configured patterns.
type Matcher struct {
	matches []match
	ignores []string
}

// NewMatcher creates a Matcher from the effective configuration.
func NewMatcher(cfg *config.Config) *Matcher {
	var matches []match
	for _, cat := range cfg.Categories {
		for _, pat := range cat.Patterns {
			matches = append(matches, match{
				categoryName:  cat.Name,
				categoryIcon:  cat.Icon,
				categoryColor: cat.Color,
				priority:      cat.Priority,
				patternName:   pat.Name,
				parser:        pat.Parser,
				globs:         pat.AllGlobs(),
			})
		}
	}
	return &Matcher{
		matches: matches,
		ignores: cfg.Ignore,
	}
}

// matchResult is returned when a file matches a pattern.
type matchResult struct {
	CategoryName  string
	CategoryIcon  string
	CategoryColor string
	Priority      int
	PatternName   string
	Parser        string
}

// Match tests a relative file path against all patterns and returns
// the first match, or nil if no pattern matches.
func (m *Matcher) Match(relPath string) *matchResult {
	// Normalize to forward slashes for cross-platform glob matching.
	normalized := filepath.ToSlash(relPath)

	for _, ig := range m.ignores {
		if matched, _ := doublestar.Match(ig, normalized); matched {
			return nil
		}
	}

	for _, mt := range m.matches {
		for _, glob := range mt.globs {
			if matched, _ := doublestar.Match(glob, normalized); matched {
				return &matchResult{
					CategoryName:  mt.categoryName,
					CategoryIcon:  mt.categoryIcon,
					CategoryColor: mt.categoryColor,
					Priority:      mt.priority,
					PatternName:   mt.patternName,
					Parser:        mt.parser,
				}
			}
		}
	}
	return nil
}

// IsIgnoredDir checks whether a directory name should be skipped entirely.
func IsIgnoredDir(name string, extraSkipDirs []string) bool {
	// Built-in directories to always skip.
	switch name {
	case ".git", "node_modules", "vendor", "__pycache__",
		".terraform", ".next", ".nuxt", "dist", "build",
		".cache", ".pytest_cache", ".mypy_cache",
		"target", "bin", "obj", ".gradle", ".idea":
		return true
	}

	// Skip hidden directories (starting with .) unless they are known config dirs.
	if strings.HasPrefix(name, ".") {
		switch name {
		case ".github", ".gitlab", ".circleci", ".husky", ".githooks",
			".vscode", ".devcontainer", ".docker", ".cargo",
			".azure-pipelines", ".woodpecker", ".cursor", ".claude",
			".vitepress":
			return false
		default:
			return true
		}
	}

	for _, skip := range extraSkipDirs {
		if name == skip {
			return true
		}
	}
	return false
}
