package health

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
)

func init() { RegisterChecker(&VersionMatchChecker{}) }

// VersionMatchChecker verifies runtime version consistency across files.
type VersionMatchChecker struct{}

func (v *VersionMatchChecker) Type() string { return "version_match" }

func (v *VersionMatchChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	runtime := rule.TypeName

	type versionSource struct {
		version string
		file    string
	}

	var sources []versionSource
	for _, f := range rule.Files {
		ver := extractVersionFromFile(ctx.RootPath, f, runtime)
		if ver != "" {
			sources = append(sources, versionSource{version: ver, file: f})
		}
	}

	if len(sources) < 2 {
		return nil
	}

	// Compare all versions against the first.
	base := sources[0]
	var issues []HealthIssue
	for _, src := range sources[1:] {
		if !versionsMatch(base.version, src.version) {
			issues = append(issues, HealthIssue{
				Severity:  rule.Severity,
				CheckType: "version_match",
				Message:   fmt.Sprintf("%s version mismatch: %s has %s, %s has %s", runtime, base.file, base.version, src.file, src.version),
				File:      src.file,
				Details:   fmt.Sprintf("runtime=%s expected=%s actual=%s", runtime, base.version, src.version),
			})
		}
	}
	return issues
}

func extractVersionFromFile(rootPath, relPath, runtime string) string {
	absPath := filepath.Join(rootPath, relPath)
	base := strings.ToLower(filepath.Base(relPath))

	// Version files that contain just a version number.
	switch base {
	case ".nvmrc", ".node-version", ".python-version", ".ruby-version", ".go-version":
		data, err := os.ReadFile(absPath)
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(data))
	case ".tool-versions":
		data, err := os.ReadFile(absPath)
		if err != nil {
			return ""
		}
		for _, line := range strings.Split(string(data), "\n") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := strings.ToLower(parts[0])
				if name == runtime || (runtime == "node" && name == "nodejs") {
					return parts[1]
				}
			}
		}
		return ""
	}

	// Parse structured files.
	file := &parser.ParseResult{}
	if f, ok := parserForPath(rootPath, relPath); ok {
		pr, err := parser.ParseFile(f, absPath)
		if err == nil && pr != nil {
			file = pr
		}
	}

	if v, ok := file.Versions[runtime]; ok {
		return v
	}
	return ""
}

func parserForPath(rootPath, relPath string) (string, bool) {
	base := strings.ToLower(filepath.Base(relPath))
	switch {
	case base == "dockerfile" || strings.HasPrefix(base, "dockerfile."):
		return "dockerfile", true
	case base == "package.json":
		return "json", true
	case strings.HasSuffix(base, ".toml"):
		return "toml", true
	case strings.HasSuffix(base, ".yml") || strings.HasSuffix(base, ".yaml"):
		return "yaml", true
	}
	return "", false
}

// versionsMatch compares version strings loosely (major version match).
func versionsMatch(a, b string) bool {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")

	// Exact match.
	if a == b {
		return true
	}

	// Major version match (e.g., "20" matches "20.11.0").
	aMajor := strings.SplitN(a, ".", 2)[0]
	bMajor := strings.SplitN(b, ".", 2)[0]
	if aMajor == bMajor && aMajor != "" {
		return true
	}

	// Prefix match with >=, ~, ^.
	for _, prefix := range []string{">=", "~", "^", "> ", "= "} {
		a = strings.TrimPrefix(a, prefix)
		b = strings.TrimPrefix(b, prefix)
	}
	aMajor = strings.SplitN(a, ".", 2)[0]
	bMajor = strings.SplitN(b, ".", 2)[0]
	return aMajor == bMajor && aMajor != ""
}
