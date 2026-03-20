package health

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/volodymyrkovtun/confradar/internal/config"
)

func init() { RegisterChecker(&RegexChecker{}) }

// RegexChecker runs a custom regex against a file and flags matches.
type RegexChecker struct{}

func (r *RegexChecker) Type() string { return "regex_check" }

func (r *RegexChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	pattern, err := regexp.Compile(rule.Pattern)
	if err != nil {
		return []HealthIssue{{
			Severity:  SeverityError,
			CheckType: "regex_check",
			Message:   fmt.Sprintf("Invalid regex pattern %q: %v", rule.Pattern, err),
		}}
	}

	filePath := rule.File
	absPath := filepath.Join(ctx.RootPath, filePath)
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil // file doesn't exist, skip silently
	}

	if pattern.Match(data) {
		msg := rule.Message
		if msg == "" {
			msg = fmt.Sprintf("Pattern %q matched in %s", rule.Pattern, filePath)
		}
		return []HealthIssue{{
			Severity:  rule.Severity,
			CheckType: "regex_check",
			Message:   msg,
			File:      filePath,
		}}
	}
	return nil
}
