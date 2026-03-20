package health

import (
	"fmt"
	"path/filepath"

	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
)

func init() { RegisterChecker(&KeyExistsChecker{}) }

// KeyExistsChecker verifies that specific required keys exist in a file.
type KeyExistsChecker struct{}

func (k *KeyExistsChecker) Type() string { return "key_exists" }

func (k *KeyExistsChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	filePath := rule.File
	result, err := parser.ParseFile("env", filepath.Join(ctx.RootPath, filePath))
	if err != nil || result == nil {
		return []HealthIssue{{
			Severity:  rule.Severity,
			CheckType: "key_exists",
			Message:   fmt.Sprintf("Cannot read file %s", filePath),
			File:      filePath,
		}}
	}

	existingKeys := toSet(result.Keys)
	var issues []HealthIssue
	for _, required := range rule.RequiredKeys {
		if !existingKeys[required] {
			issues = append(issues, HealthIssue{
				Severity:  rule.Severity,
				CheckType: "key_exists",
				Message:   fmt.Sprintf("Required key %q missing from %s", required, filePath),
				File:      filePath,
			})
		}
	}
	return issues
}
