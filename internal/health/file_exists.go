package health

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/volodymyrkovtun/confradar/internal/config"
)

func init() { RegisterChecker(&FileExistsChecker{}) }

// FileExistsChecker verifies that required files exist in the project.
type FileExistsChecker struct{}

func (f *FileExistsChecker) Type() string { return "file_exists" }

func (f *FileExistsChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	var issues []HealthIssue
	msg := rule.Message
	if msg == "" {
		msg = "Required file missing"
	}

	for _, file := range rule.Files {
		absPath := filepath.Join(ctx.RootPath, file)
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			issues = append(issues, HealthIssue{
				Severity:  rule.Severity,
				CheckType: "file_exists",
				Message:   fmt.Sprintf("%s: %s", msg, file),
				File:      file,
			})
		}
	}
	return issues
}
