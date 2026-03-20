package health

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
)

func init() { RegisterChecker(&EnvSyncChecker{}) }

// EnvSyncChecker compares keys between .env files.
type EnvSyncChecker struct{}

func (e *EnvSyncChecker) Type() string { return "env_sync" }

func (e *EnvSyncChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	sourcePath := rule.Source
	sourceResult, err := parser.ParseFile("env", filepath.Join(ctx.RootPath, sourcePath))
	if err != nil || sourceResult == nil {
		return nil
	}

	sourceKeys := toSet(sourceResult.Keys)
	ignoreSet := buildIgnoreSet(rule.IgnoreKeys)

	var issues []HealthIssue
	for _, target := range rule.Targets {
		targetResult, err := parser.ParseFile("env", filepath.Join(ctx.RootPath, target))
		if err != nil || targetResult == nil {
			issues = append(issues, HealthIssue{
				Severity:  rule.Severity,
				CheckType: "env_sync",
				Message:   fmt.Sprintf("Cannot read target file %s", target),
				File:      target,
			})
			continue
		}

		targetKeys := toSet(targetResult.Keys)

		for key := range sourceKeys {
			if ignoreSet[key] || matchesAnyGlob(key, rule.IgnoreKeys) {
				continue
			}
			if !targetKeys[key] {
				issues = append(issues, HealthIssue{
					Severity:  rule.Severity,
					CheckType: "env_sync",
					Message:   fmt.Sprintf("Key %q exists in %s but missing from %s", key, sourcePath, target),
					File:      target,
					Details:   fmt.Sprintf("source=%s target=%s key=%s", sourcePath, target, key),
				})
			}
		}
	}
	return issues
}

func buildIgnoreSet(keys []string) map[string]bool {
	set := make(map[string]bool, len(keys))
	for _, k := range keys {
		if !strings.Contains(k, "*") {
			set[k] = true
		}
	}
	return set
}

func matchesAnyGlob(key string, patterns []string) bool {
	for _, p := range patterns {
		if !strings.Contains(p, "*") {
			continue
		}
		matched, _ := filepath.Match(p, key)
		if matched {
			return true
		}
	}
	return false
}
