package health

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
)

func init() { RegisterChecker(&PortConflictChecker{}) }

// PortConflictChecker detects port conflicts across config files.
type PortConflictChecker struct{}

func (p *PortConflictChecker) Type() string { return "port_conflict" }

func (p *PortConflictChecker) Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue {
	type portSource struct {
		port int
		file string
	}

	var sources []portSource
	files := rule.Files
	if len(files) == 0 {
		// Check all parsed files.
		for relPath, pr := range ctx.Parsed {
			for _, port := range pr.Ports {
				sources = append(sources, portSource{port: port, file: relPath})
			}
		}
	} else {
		for _, f := range files {
			absPath := filepath.Join(ctx.RootPath, f)
			file := ctx.Files[f]
			parserType := "env"
			if file != nil {
				parserType = file.Parser
			}
			pr, err := parser.ParseFile(parserType, absPath)
			if err != nil || pr == nil {
				continue
			}
			for _, port := range pr.Ports {
				sources = append(sources, portSource{port: port, file: f})
			}
		}
	}

	// Find duplicates.
	portMap := make(map[int][]string)
	for _, ps := range sources {
		portMap[ps.port] = append(portMap[ps.port], ps.file)
	}

	var issues []HealthIssue
	for port, filenames := range portMap {
		if len(filenames) > 1 {
			issues = append(issues, HealthIssue{
				Severity:  rule.Severity,
				CheckType: "port_conflict",
				Message:   fmt.Sprintf("Port %d used in multiple files: %s", port, strings.Join(filenames, ", ")),
				File:      filenames[0],
				Details:   fmt.Sprintf("port=%d files=%s", port, strings.Join(filenames, ",")),
			})
		}
	}
	return issues
}
