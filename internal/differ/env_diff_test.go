package differ

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiff(t *testing.T) {
	dir := t.TempDir()

	env1 := `DATABASE_URL=postgres://localhost/app
REDIS_URL=redis://localhost:6379
PORT=3000
SECRET_KEY=abc123
`
	env2 := `DATABASE_URL=postgres://prod/app
REDIS_URL=redis://localhost:6379
PORT=3000
API_KEY=xyz789
`
	path1 := filepath.Join(dir, ".env")
	path2 := filepath.Join(dir, ".env.production")
	os.WriteFile(path1, []byte(env1), 0644)
	os.WriteFile(path2, []byte(env2), 0644)

	result, err := Diff(path1, path2)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}

	// SECRET_KEY only in left.
	if len(result.OnlyLeft) != 1 || result.OnlyLeft[0].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY only in left, got %v", result.OnlyLeft)
	}

	// API_KEY only in right.
	if len(result.OnlyRight) != 1 || result.OnlyRight[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY only in right, got %v", result.OnlyRight)
	}

	// DATABASE_URL changed.
	if len(result.Changed) != 1 || result.Changed[0].Key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL changed, got %v", result.Changed)
	}

	// PORT and REDIS_URL common.
	if len(result.Common) != 2 {
		t.Errorf("expected 2 common keys, got %d", len(result.Common))
	}

	if !result.HasDifferences() {
		t.Error("should have differences")
	}
}
