package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/volodymyrkovtun/confradar/internal/config"
)

// Scan walks the project at rootPath and returns all detected config files
// grouped by category. It performs a single-pass directory walk.
func Scan(rootPath string, cfg *config.Config) (*ScanResult, error) {
	start := time.Now()
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	matcher := NewMatcher(cfg)
	ignoreChecker := NewIgnoreChecker(absRoot, cfg.Ignore)

	result := &ScanResult{
		Root:       absRoot,
		Categories: make(map[string]*CategoryResult),
	}

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries gracefully
		}

		name := d.Name()

		if d.IsDir() {
			if path == absRoot {
				return nil
			}

			if IsIgnoredDir(name, cfg.SkipDirs) {
				return fs.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			return nil
		}

		if ignoreChecker.IsIgnored(relPath) {
			return nil
		}

		mr := matcher.Match(relPath)
		if mr == nil {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		cf := ConfigFile{
			Path:        path,
			RelPath:     filepath.ToSlash(relPath),
			Category:    mr.CategoryName,
			PatternName: mr.PatternName,
			Parser:      mr.Parser,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
		}

		cat, exists := result.Categories[mr.CategoryName]
		if !exists {
			cat = &CategoryResult{
				Name:     mr.CategoryName,
				Icon:     mr.CategoryIcon,
				Color:    mr.CategoryColor,
				Priority: mr.Priority,
			}
			result.Categories[mr.CategoryName] = cat
		}
		cat.Files = append(cat.Files, cf)
		result.TotalFiles++

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	// Sort files within each category.
	for _, cat := range result.Categories {
		sortFiles(cat.Files, cfg.Display.FileSort)
	}

	// Build ordered category list.
	result.Ordered = buildOrdered(result.Categories, cfg.Display.CategorySort)
	result.Duration = time.Since(start)

	return result, nil
}

// ScanWithDepth is like Scan but enforces a maximum directory depth.
func ScanWithDepth(rootPath string, cfg *config.Config, maxDepth int) (*ScanResult, error) {
	start := time.Now()
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	matcher := NewMatcher(cfg)
	ignoreChecker := NewIgnoreChecker(absRoot, cfg.Ignore)

	result := &ScanResult{
		Root:       absRoot,
		Categories: make(map[string]*CategoryResult),
	}

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()
		relPath, relErr := filepath.Rel(absRoot, path)
		if relErr != nil {
			return nil
		}

		// Calculate depth.
		depth := 0
		if relPath != "." {
			depth = strings.Count(filepath.ToSlash(relPath), "/") + 1
		}

		if d.IsDir() {
			if path == absRoot {
				return nil
			}
			if IsIgnoredDir(name, cfg.SkipDirs) {
				return fs.SkipDir
			}
			if maxDepth > 0 && depth >= maxDepth {
				return fs.SkipDir
			}
			return nil
		}

		if ignoreChecker.IsIgnored(relPath) {
			return nil
		}

		mr := matcher.Match(relPath)
		if mr == nil {
			return nil
		}

		info, infoErr := d.Info()
		if infoErr != nil {
			return nil
		}

		cf := ConfigFile{
			Path:        path,
			RelPath:     filepath.ToSlash(relPath),
			Category:    mr.CategoryName,
			PatternName: mr.PatternName,
			Parser:      mr.Parser,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
		}

		cat, exists := result.Categories[mr.CategoryName]
		if !exists {
			cat = &CategoryResult{
				Name:     mr.CategoryName,
				Icon:     mr.CategoryIcon,
				Color:    mr.CategoryColor,
				Priority: mr.Priority,
			}
			result.Categories[mr.CategoryName] = cat
		}
		cat.Files = append(cat.Files, cf)
		result.TotalFiles++

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	for _, cat := range result.Categories {
		sortFiles(cat.Files, cfg.Display.FileSort)
	}
	result.Ordered = buildOrdered(result.Categories, cfg.Display.CategorySort)
	result.Duration = time.Since(start)

	return result, nil
}

func sortFiles(files []ConfigFile, sortBy string) {
	switch sortBy {
	case "path":
		sort.Slice(files, func(i, j int) bool { return files[i].RelPath < files[j].RelPath })
	case "size":
		sort.Slice(files, func(i, j int) bool { return files[i].Size < files[j].Size })
	case "modified":
		sort.Slice(files, func(i, j int) bool { return files[i].ModTime.Before(files[j].ModTime) })
	default: // "name"
		sort.Slice(files, func(i, j int) bool {
			return filepath.Base(files[i].RelPath) < filepath.Base(files[j].RelPath)
		})
	}
}

func buildOrdered(cats map[string]*CategoryResult, sortBy string) []*CategoryResult {
	ordered := make([]*CategoryResult, 0, len(cats))
	for _, cat := range cats {
		ordered = append(ordered, cat)
	}

	switch sortBy {
	case "name":
		sort.Slice(ordered, func(i, j int) bool { return ordered[i].Name < ordered[j].Name })
	case "file_count":
		sort.Slice(ordered, func(i, j int) bool { return len(ordered[i].Files) > len(ordered[j].Files) })
	default: // "priority"
		sort.Slice(ordered, func(i, j int) bool {
			if ordered[i].Priority == ordered[j].Priority {
				return ordered[i].Name < ordered[j].Name
			}
			return ordered[i].Priority < ordered[j].Priority
		})
	}
	return ordered
}
