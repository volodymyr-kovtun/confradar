package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := loadDefaults()
	if err != nil {
		t.Fatalf("loadDefaults() error: %v", err)
	}
	if len(cfg.Categories) == 0 {
		t.Fatal("expected at least one category in defaults")
	}

	// Verify we have the expected categories.
	names := make(map[string]bool)
	for _, cat := range cfg.Categories {
		names[cat.Name] = true
	}

	expected := []string{"Environment", "Docker", "CI/CD", "JavaScript / Node", "Python", "Go"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected category %q in defaults", name)
		}
	}
}

func TestMergeExtraCategories(t *testing.T) {
	base := &Config{
		Categories: []Category{
			{Name: "Docker", Icon: "🐳", Priority: 1},
		},
	}
	overlay := &Config{
		ExtraCategories: []Category{
			{Name: "Custom", Icon: "📁", Priority: 50},
		},
	}

	result := merge(base, overlay)
	if len(result.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(result.Categories))
	}
	if result.Categories[1].Name != "Custom" {
		t.Errorf("expected Custom category, got %q", result.Categories[1].Name)
	}
}

func TestMergeDisableCategories(t *testing.T) {
	base := &Config{
		Categories: []Category{
			{Name: "Docker", Icon: "🐳"},
			{Name: "CI/CD", Icon: "⚙️"},
			{Name: "Linting", Icon: "✨"},
		},
	}
	overlay := &Config{
		DisableCategories: []string{"CI/CD"},
	}

	result := merge(base, overlay)
	if len(result.Categories) != 2 {
		t.Fatalf("expected 2 categories after disabling, got %d", len(result.Categories))
	}
	for _, cat := range result.Categories {
		if cat.Name == "CI/CD" {
			t.Error("CI/CD should have been disabled")
		}
	}
}

func TestMergeOverrideCategories(t *testing.T) {
	base := &Config{
		Categories: []Category{
			{Name: "Docker", Icon: "🐳", Color: "#blue", Priority: 1},
		},
	}
	p := 99
	overlay := &Config{
		OverrideCategories: []CategoryOverride{
			{Name: "Docker", Icon: "🐋", Priority: &p},
		},
	}

	result := merge(base, overlay)
	if result.Categories[0].Icon != "🐋" {
		t.Errorf("expected icon 🐋, got %s", result.Categories[0].Icon)
	}
	if result.Categories[0].Priority != 99 {
		t.Errorf("expected priority 99, got %d", result.Categories[0].Priority)
	}
	if result.Categories[0].Color != "#blue" {
		t.Errorf("color should be preserved, got %s", result.Categories[0].Color)
	}
}

func TestMergeIgnore(t *testing.T) {
	base := &Config{Ignore: []string{"vendor/**"}}
	overlay := &Config{Ignore: []string{"legacy/**", "vendor/**"}}

	result := merge(base, overlay)
	if len(result.Ignore) != 2 {
		t.Fatalf("expected 2 ignore patterns (deduped), got %d", len(result.Ignore))
	}
}

func TestNew(t *testing.T) {
	cfg, err := New(".", CLIFlags{NoConfig: true})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if cfg.Display.Theme != "auto" {
		t.Errorf("expected theme auto, got %s", cfg.Display.Theme)
	}
	if cfg.Output.DefaultFormat != "tree" {
		t.Errorf("expected format tree, got %s", cfg.Output.DefaultFormat)
	}
}
