package health

import (
	"fmt"
	"path/filepath"

	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
)

func init() { RegisterChecker(&EnvExampleSyncChecker{}) }

// EnvExampleSyncChecker checks that .env.example is in sync with .env.
type EnvExampleSyncChecker struct{}

func (e *EnvExampleSyncChecker) Type() string { return "env_example_sync" }

func (e *EnvExampleSyncChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	envPath := filepath.Join(ctx.RootPath, ".env")
	examplePath := filepath.Join(ctx.RootPath, ".env.example")

	envResult, err := parser.ParseFile("env", envPath)
	if err != nil || envResult == nil {
		return nil
	}
	exampleResult, err := parser.ParseFile("env", examplePath)
	if err != nil || exampleResult == nil {
		return nil
	}

	envKeys := toSet(envResult.Keys)
	exampleKeys := toSet(exampleResult.Keys)

	var issues []HealthIssue

	for key := range envKeys {
		if !exampleKeys[key] {
			issues = append(issues, HealthIssue{
				Severity:  rule.Severity,
				CheckType: "env_example_sync",
				Message:   fmt.Sprintf("Key %q in .env but missing from .env.example", key),
				File:      ".env.example",
			})
		}
	}

	for key := range exampleKeys {
		if !envKeys[key] {
			issues = append(issues, HealthIssue{
				Severity:  SeverityInfo,
				CheckType: "env_example_sync",
				Message:   fmt.Sprintf("Key %q in .env.example but missing from .env", key),
				File:      ".env",
			})
		}
	}

	return issues
}
