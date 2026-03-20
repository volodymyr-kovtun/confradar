// Package health implements the config health check engine.
package health

import (
	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// Severity levels for health issues.
const (
	SeverityError   = "error"
	SeverityWarning = "warning"
	SeverityInfo    = "info"
)

// HealthIssue represents a single detected problem.
type HealthIssue struct {
	Severity  string `json:"severity" yaml:"severity"`
	CheckType string `json:"check_type" yaml:"check_type"`
	Message   string `json:"message" yaml:"message"`
	File      string `json:"file,omitempty" yaml:"file,omitempty"`
	Details   string `json:"details,omitempty" yaml:"details,omitempty"`
}

// HealthChecker runs a specific type of health check.
type HealthChecker interface {
	Type() string
	Check(ctx *CheckContext, rule config.HealthCheckRule) []HealthIssue
}

// CheckContext provides all data needed by health checkers.
type CheckContext struct {
	RootPath string
	Files    map[string]*scanner.ConfigFile
	Parsed   map[string]*parser.ParseResult
}
