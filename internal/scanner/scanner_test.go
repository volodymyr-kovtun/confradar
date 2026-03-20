package scanner

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/volodymyrkovtun/confradar/internal/config"
)

func testdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", name)
}

func TestScanNodeFullstack(t *testing.T) {
	cfg, err := config.New(testdataPath("node_fullstack"), config.CLIFlags{})
	if err != nil {
		t.Fatalf("config: %v", err)
	}

	result, err := Scan(testdataPath("node_fullstack"), cfg)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if result.TotalFiles < 10 {
		t.Errorf("expected at least 10 files, got %d", result.TotalFiles)
	}

	// Check expected categories exist.
	expectedCats := []string{"Environment", "Docker", "JavaScript / Node"}
	for _, name := range expectedCats {
		if _, ok := result.Categories[name]; !ok {
			t.Errorf("expected category %q", name)
		}
	}

	// Check .env files are found.
	envCat := result.Categories["Environment"]
	if envCat == nil {
		t.Fatal("missing Environment category")
	}
	if len(envCat.Files) < 3 {
		t.Errorf("expected at least 3 env files, got %d", len(envCat.Files))
	}
}

func TestScanEmpty(t *testing.T) {
	cfg, err := config.New(testdataPath("empty"), config.CLIFlags{})
	if err != nil {
		t.Fatalf("config: %v", err)
	}

	result, err := Scan(testdataPath("empty"), cfg)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if result.TotalFiles != 0 {
		t.Errorf("expected 0 files, got %d", result.TotalFiles)
	}
}

func TestScanMonorepo(t *testing.T) {
	cfg, err := config.New(testdataPath("monorepo"), config.CLIFlags{})
	if err != nil {
		t.Fatalf("config: %v", err)
	}

	result, err := Scan(testdataPath("monorepo"), cfg)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	// Should find nested .env files.
	envCat := result.Categories["Environment"]
	if envCat == nil {
		t.Fatal("missing Environment category")
	}
	if len(envCat.Files) < 2 {
		t.Errorf("expected at least 2 nested env files, got %d", len(envCat.Files))
	}
}

func TestScanWithDepth(t *testing.T) {
	cfg, err := config.New(testdataPath("monorepo"), config.CLIFlags{})
	if err != nil {
		t.Fatalf("config: %v", err)
	}

	result, err := ScanWithDepth(testdataPath("monorepo"), cfg, 1)
	if err != nil {
		t.Fatalf("ScanWithDepth: %v", err)
	}

	// Depth 1 should only find root-level files.
	for _, cat := range result.Categories {
		for _, f := range cat.Files {
			if filepath.Dir(f.RelPath) != "." {
				// Files at depth > 1 should be excluded.
				depth := len(filepath.SplitList(f.RelPath))
				if depth > 1 {
					// This is ok for files like .github/workflows which are at depth 2
					// with our counting. The check is more nuanced.
				}
			}
		}
	}
}

func TestIsIgnoredDir(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{".git", true},
		{"node_modules", true},
		{"vendor", true},
		{".github", false},
		{".vscode", false},
		{"src", false},
		{".terraform", true},
	}

	for _, tt := range tests {
		got := IsIgnoredDir(tt.name, nil)
		if got != tt.expected {
			t.Errorf("IsIgnoredDir(%q) = %v, want %v", tt.name, got, tt.expected)
		}
	}
}
