package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvParser(t *testing.T) {
	content := `# Database config
DATABASE_URL=postgres://localhost/myapp
REDIS_URL=redis://localhost:6379
PORT=3000
SECRET_KEY="my secret key"
EMPTY_VAR=
# This is a comment
DEBUG=true
`
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &EnvParser{}
	result, err := p.Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if len(result.Keys) != 6 {
		t.Errorf("expected 6 keys, got %d: %v", len(result.Keys), result.Keys)
	}

	if result.Values["DATABASE_URL"] != "postgres://localhost/myapp" {
		t.Errorf("DATABASE_URL = %q", result.Values["DATABASE_URL"])
	}

	if result.Values["SECRET_KEY"] != "my secret key" {
		t.Errorf("SECRET_KEY should have quotes stripped, got %q", result.Values["SECRET_KEY"])
	}

	if result.Values["EMPTY_VAR"] != "" {
		t.Errorf("EMPTY_VAR should be empty, got %q", result.Values["EMPTY_VAR"])
	}

	// PORT=3000 should be detected as a port.
	if len(result.Ports) != 1 || result.Ports[0] != 3000 {
		t.Errorf("expected port 3000, got %v", result.Ports)
	}
}

func TestEnvParserEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	p := &EnvParser{}
	result, err := p.Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if len(result.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(result.Keys))
	}
}
