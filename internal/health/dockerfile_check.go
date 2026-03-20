package health

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/volodymyrkovtun/confradar/internal/config"
)

func init() { RegisterChecker(&DockerfileCheckChecker{}) }

// DockerfileCheckChecker checks Dockerfile best practices.
type DockerfileCheckChecker struct{}

func (d *DockerfileCheckChecker) Type() string { return "dockerfile_best_practices" }

func (d *DockerfileCheckChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	var issues []HealthIssue

	for relPath, f := range ctx.Files {
		if f.Parser != "dockerfile" {
			continue
		}

		pr := ctx.Parsed[relPath]
		if pr == nil {
			continue
		}

		// Check for :latest tag.
		if pr.Metadata["uses_latest"] == "true" {
			issues = append(issues, HealthIssue{
				Severity:  rule.Severity,
				CheckType: "dockerfile_best_practices",
				Message:   fmt.Sprintf("Using :latest tag in %s is risky for reproducible builds", relPath),
				File:      relPath,
			})
		}

		// Check for missing .dockerignore.
		dir := filepath.Dir(filepath.Join(ctx.RootPath, relPath))
		dockerignorePath := filepath.Join(dir, ".dockerignore")
		if _, err := os.Stat(dockerignorePath); os.IsNotExist(err) {
			// Also check project root.
			rootIgnore := filepath.Join(ctx.RootPath, ".dockerignore")
			if _, err := os.Stat(rootIgnore); os.IsNotExist(err) {
				issues = append(issues, HealthIssue{
					Severity:  SeverityInfo,
					CheckType: "dockerfile_best_practices",
					Message:   fmt.Sprintf("No .dockerignore found for %s", relPath),
					File:      relPath,
				})
			}
		}

		// Check for USER directive (running as root).
		absPath := filepath.Join(ctx.RootPath, relPath)
		data, err := os.ReadFile(absPath)
		if err != nil {
			continue
		}
		content := string(data)
		upper := strings.ToUpper(content)
		if !strings.HasPrefix(upper, "USER ") && !strings.Contains(upper, "\nUSER ") {
			issues = append(issues, HealthIssue{
				Severity:  SeverityInfo,
				CheckType: "dockerfile_best_practices",
				Message:   fmt.Sprintf("No USER directive in %s — container will run as root", relPath),
				File:      relPath,
			})
		}
	}

	return issues
}
