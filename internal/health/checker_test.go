package health

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/parser"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

func TestEnvSyncChecker(t *testing.T) {
	dir := t.TempDir()

	// Create source .env with 3 keys.
	os.WriteFile(filepath.Join(dir, ".env"), []byte("A=1\nB=2\nC=3\n"), 0644)
	// Create target missing key C.
	os.WriteFile(filepath.Join(dir, ".env.prod"), []byte("A=1\nB=2\n"), 0644)

	rule := config.HealthCheckRule{
		Type:     "env_sync",
		Source:   ".env",
		Targets:  []string{".env.prod"},
		Severity: "warning",
	}

	ctx := &CheckContext{
		RootPath: dir,
		Files:    map[string]*scanner.ConfigFile{},
		Parsed:   map[string]*parser.ParseResult{},
	}

	checker := &EnvSyncChecker{}
	issues := checker.Check(ctx, rule)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %s", issues[0].Severity)
	}
}

func TestFileExistsChecker(t *testing.T) {
	dir := t.TempDir()

	// Create one file.
	os.WriteFile(filepath.Join(dir, ".env"), []byte("A=1"), 0644)
	// .editorconfig does NOT exist.

	rule := config.HealthCheckRule{
		Type:     "file_exists",
		Files:    []string{".env", ".editorconfig"},
		Severity: "warning",
		Message:  "Required file missing",
	}

	ctx := &CheckContext{
		RootPath: dir,
		Files:    map[string]*scanner.ConfigFile{},
		Parsed:   map[string]*parser.ParseResult{},
	}

	checker := &FileExistsChecker{}
	issues := checker.Check(ctx, rule)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue (.editorconfig missing), got %d", len(issues))
	}
}

func TestFilterBySeverity(t *testing.T) {
	issues := []HealthIssue{
		{Severity: SeverityError, Message: "error1"},
		{Severity: SeverityWarning, Message: "warn1"},
		{Severity: SeverityInfo, Message: "info1"},
	}

	errors := FilterBySeverity(issues, SeverityError)
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors))
	}

	warnings := FilterBySeverity(issues, SeverityWarning)
	if len(warnings) != 2 {
		t.Errorf("expected 2 (errors+warnings), got %d", len(warnings))
	}

	all := FilterBySeverity(issues, "")
	if len(all) != 3 {
		t.Errorf("expected all 3, got %d", len(all))
	}
}
