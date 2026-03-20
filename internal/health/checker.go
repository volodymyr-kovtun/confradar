package health

import (
	"sort"

	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

var checkers = map[string]HealthChecker{}

// RegisterChecker adds a health checker to the registry.
func RegisterChecker(hc HealthChecker) {
	checkers[hc.Type()] = hc
}

// RunChecks executes all health check rules against the scan results.
func RunChecks(rootPath string, result *scanner.ScanResult, rules []config.HealthCheckRule) []HealthIssue {
	if len(rules) == 0 {
		return nil
	}

	// Build file map and parse all detected config files.
	files := make(map[string]*scanner.ConfigFile)
	for _, cat := range result.Categories {
		for i := range cat.Files {
			f := &cat.Files[i]
			files[f.RelPath] = f
		}
	}

	parsed := make(map[string]*parser.ParseResult)
	for relPath, f := range files {
		pr, err := parser.ParseFile(f.Parser, f.Path)
		if err != nil {
			continue
		}
		parsed[relPath] = pr
	}

	ctx := &CheckContext{
		RootPath: rootPath,
		Files:    files,
		Parsed:   parsed,
	}

	var issues []HealthIssue
	for _, rule := range rules {
		hc, ok := checkers[rule.Type]
		if !ok {
			continue
		}
		issues = append(issues, hc.Check(ctx, rule)...)
	}

	// Sort by severity: error > warning > info.
	sort.Slice(issues, func(i, j int) bool {
		return severityRank(issues[i].Severity) < severityRank(issues[j].Severity)
	})

	return issues
}

// RunAutoChecks runs built-in automatic health checks that don't require configuration.
func RunAutoChecks(rootPath string, result *scanner.ScanResult) []HealthIssue {
	files := make(map[string]*scanner.ConfigFile)
	for _, cat := range result.Categories {
		for i := range cat.Files {
			f := &cat.Files[i]
			files[f.RelPath] = f
		}
	}

	parsed := make(map[string]*parser.ParseResult)
	for relPath, f := range files {
		pr, err := parser.ParseFile(f.Parser, f.Path)
		if err != nil {
			continue
		}
		parsed[relPath] = pr
	}

	ctx := &CheckContext{
		RootPath: rootPath,
		Files:    files,
		Parsed:   parsed,
	}

	var issues []HealthIssue

	// Auto-detect .env sync issues.
	issues = append(issues, autoEnvSync(ctx)...)

	// Auto-detect Dockerfile best practices.
	issues = append(issues, autoDockerfileCheck(ctx)...)

	sort.Slice(issues, func(i, j int) bool {
		return severityRank(issues[i].Severity) < severityRank(issues[j].Severity)
	})

	return issues
}

func autoEnvSync(ctx *CheckContext) []HealthIssue {
	// Find .env and .env.example pairs.
	envParsed := ctx.Parsed[".env"]
	exampleParsed := ctx.Parsed[".env.example"]
	if envParsed == nil || exampleParsed == nil {
		return nil
	}

	envKeys := toSet(envParsed.Keys)
	exampleKeys := toSet(exampleParsed.Keys)

	var issues []HealthIssue
	for key := range envKeys {
		if !exampleKeys[key] {
			issues = append(issues, HealthIssue{
				Severity:  SeverityWarning,
				CheckType: "env_example_sync",
				Message:   key + " exists in .env but missing from .env.example",
				File:      ".env.example",
			})
		}
	}
	return issues
}

func autoDockerfileCheck(ctx *CheckContext) []HealthIssue {
	var issues []HealthIssue
	for relPath, pr := range ctx.Parsed {
		f := ctx.Files[relPath]
		if f == nil || f.Parser != "dockerfile" {
			continue
		}
		if pr.Metadata["uses_latest"] == "true" {
			issues = append(issues, HealthIssue{
				Severity:  SeverityWarning,
				CheckType: "dockerfile_best_practices",
				Message:   "Using :latest tag is risky for reproducible builds",
				File:      relPath,
			})
		}
	}
	return issues
}

func severityRank(s string) int {
	switch s {
	case SeverityError:
		return 0
	case SeverityWarning:
		return 1
	case SeverityInfo:
		return 2
	default:
		return 3
	}
}

func toSet(items []string) map[string]bool {
	s := make(map[string]bool, len(items))
	for _, item := range items {
		s[item] = true
	}
	return s
}

// FilterBySeverity returns issues at or above the given severity level.
func FilterBySeverity(issues []HealthIssue, minSeverity string) []HealthIssue {
	if minSeverity == "" {
		return issues
	}
	maxRank := severityRank(minSeverity)
	var filtered []HealthIssue
	for _, issue := range issues {
		if severityRank(issue.Severity) <= maxRank {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}
