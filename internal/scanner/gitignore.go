package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// IgnoreChecker tests whether a relative path should be excluded from scanning.
type IgnoreChecker struct {
	rules []ignoreRule
}

type ignoreRule struct {
	pattern  string
	negated  bool
	dirOnly  bool
}

// NewIgnoreChecker builds an IgnoreChecker from .gitignore and .confradarignore
// files found in the project root, plus any user-supplied ignore globs.
func NewIgnoreChecker(rootPath string, extraIgnores []string) *IgnoreChecker {
	var rules []ignoreRule

	for _, name := range []string{".gitignore", ".confradarignore"} {
		path := filepath.Join(rootPath, name)
		if parsed, err := parseIgnoreFile(path); err == nil {
			rules = append(rules, parsed...)
		}
	}

	for _, pattern := range extraIgnores {
		rules = append(rules, ignoreRule{pattern: pattern})
	}

	return &IgnoreChecker{rules: rules}
}

// IsIgnored tests whether the given relative path should be skipped.
func (ic *IgnoreChecker) IsIgnored(relPath string) bool {
	normalized := filepath.ToSlash(relPath)
	ignored := false

	for _, rule := range ic.rules {
		matched, _ := doublestar.Match(rule.pattern, normalized)
		if !matched {
			// Also try matching against just the basename for simple patterns.
			base := filepath.Base(normalized)
			matched, _ = doublestar.Match(rule.pattern, base)
		}
		if matched {
			ignored = !rule.negated
		}
	}
	return ignored
}

func parseIgnoreFile(path string) ([]ignoreRule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rules []ignoreRule
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		rule := ignoreRule{}
		if strings.HasPrefix(line, "!") {
			rule.negated = true
			line = line[1:]
		}
		if strings.HasSuffix(line, "/") {
			rule.dirOnly = true
			line = strings.TrimSuffix(line, "/")
		}

		// Convert gitignore patterns to doublestar-compatible globs.
		if !strings.Contains(line, "/") {
			// A bare name like "node_modules" matches at any depth.
			rule.pattern = "**/" + line
		} else {
			rule.pattern = line
		}

		rules = append(rules, rule)
	}
	return rules, sc.Err()
}
