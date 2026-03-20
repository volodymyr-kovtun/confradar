package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// treeItem represents a single item in the tree view (category header or file).
type treeItem struct {
	isCategory bool
	category   *scanner.CategoryResult
	file       *scanner.ConfigFile
	expanded   bool
}

// TreeView manages the navigable config file tree.
type TreeView struct {
	items    []treeItem
	cursor   int
	offset   int
	height   int
	filter   string
	allItems []treeItem // unfiltered
}

// NewTreeView builds the tree from scan results.
func NewTreeView(result *scanner.ScanResult) TreeView {
	var items []treeItem
	for _, cat := range result.Ordered {
		items = append(items, treeItem{
			isCategory: true,
			category:   cat,
			expanded:   true,
		})
		for i := range cat.Files {
			items = append(items, treeItem{
				file: &cat.Files[i],
			})
		}
	}
	return TreeView{
		items:    items,
		allItems: items,
		height:   20,
	}
}

// SetHeight updates the visible height.
func (tv *TreeView) SetHeight(h int) {
	tv.height = h
	if tv.height < 1 {
		tv.height = 1
	}
}

// MoveUp moves the cursor up, skipping hidden items.
func (tv *TreeView) MoveUp() {
	visible := tv.visibleItems()
	if len(visible) == 0 {
		return
	}
	tv.cursor--
	if tv.cursor < 0 {
		tv.cursor = 0
	}
	tv.adjustOffset()
}

// MoveDown moves the cursor down.
func (tv *TreeView) MoveDown() {
	visible := tv.visibleItems()
	if len(visible) == 0 {
		return
	}
	tv.cursor++
	if tv.cursor >= len(visible) {
		tv.cursor = len(visible) - 1
	}
	tv.adjustOffset()
}

// Toggle expands/collapses the current category.
func (tv *TreeView) Toggle() {
	visible := tv.visibleItems()
	if tv.cursor >= len(visible) {
		return
	}
	idx := visible[tv.cursor]
	item := &tv.items[idx]
	if item.isCategory {
		item.expanded = !item.expanded
	}
}

// Selected returns the currently selected file, or nil if a category is selected.
func (tv *TreeView) Selected() *scanner.ConfigFile {
	visible := tv.visibleItems()
	if tv.cursor >= len(visible) {
		return nil
	}
	item := tv.items[visible[tv.cursor]]
	return item.file
}

// SelectedCategory returns the currently selected category, or nil.
func (tv *TreeView) SelectedCategory() *scanner.CategoryResult {
	visible := tv.visibleItems()
	if tv.cursor >= len(visible) {
		return nil
	}
	item := tv.items[visible[tv.cursor]]
	if item.isCategory {
		return item.category
	}
	return nil
}

// JumpToCategory jumps to the nth category (1-indexed).
func (tv *TreeView) JumpToCategory(n int) {
	visible := tv.visibleItems()
	catNum := 0
	for i, idx := range visible {
		if tv.items[idx].isCategory {
			catNum++
			if catNum == n {
				tv.cursor = i
				tv.adjustOffset()
				return
			}
		}
	}
}

// SetFilter applies a search filter to the tree.
func (tv *TreeView) SetFilter(query string) {
	tv.filter = strings.ToLower(query)
	tv.cursor = 0
	tv.offset = 0
}

// ClearFilter removes the search filter.
func (tv *TreeView) ClearFilter() {
	tv.filter = ""
	tv.cursor = 0
	tv.offset = 0
}

// visibleItems returns indices into tv.items for currently visible items.
func (tv *TreeView) visibleItems() []int {
	var visible []int
	var currentCat *treeItem

	for i := range tv.items {
		item := &tv.items[i]
		if item.isCategory {
			if tv.filter != "" {
				// Show category if any of its files match.
				if tv.categoryHasMatch(i) {
					visible = append(visible, i)
				}
			} else {
				visible = append(visible, i)
			}
			currentCat = item
			continue
		}

		// File item.
		if currentCat != nil && !currentCat.expanded {
			continue
		}

		if tv.filter != "" {
			if !tv.fileMatches(item.file) {
				continue
			}
		}

		visible = append(visible, i)
	}
	return visible
}

func (tv *TreeView) categoryHasMatch(catIdx int) bool {
	for i := catIdx + 1; i < len(tv.items); i++ {
		if tv.items[i].isCategory {
			break
		}
		if tv.fileMatches(tv.items[i].file) {
			return true
		}
	}
	return false
}

func (tv *TreeView) fileMatches(f *scanner.ConfigFile) bool {
	if f == nil {
		return false
	}
	q := tv.filter
	return strings.Contains(strings.ToLower(f.RelPath), q) ||
		strings.Contains(strings.ToLower(f.PatternName), q) ||
		strings.Contains(strings.ToLower(f.Category), q)
}

func (tv *TreeView) adjustOffset() {
	if tv.cursor < tv.offset {
		tv.offset = tv.cursor
	}
	if tv.cursor >= tv.offset+tv.height {
		tv.offset = tv.cursor - tv.height + 1
	}
}

// View renders the tree.
func (tv *TreeView) View(styles Styles, width int) string {
	visible := tv.visibleItems()
	if len(visible) == 0 {
		return styles.Dim.Render("  No matching files")
	}

	var b strings.Builder
	end := tv.offset + tv.height
	if end > len(visible) {
		end = len(visible)
	}

	for vi := tv.offset; vi < end; vi++ {
		idx := visible[vi]
		item := tv.items[idx]
		isCurrent := vi == tv.cursor

		var line string
		if item.isCategory {
			arrow := "▼"
			if !item.expanded {
				arrow = "▶"
			}
			catLine := fmt.Sprintf(" %s %s %s (%d)",
				arrow, item.category.Icon, item.category.Name, len(item.category.Files))
			if isCurrent {
				line = styles.Selected.Render(padRight(catLine, width))
			} else {
				line = styles.CategoryName.Render(catLine)
			}
		} else {
			name := filepath.Base(item.file.RelPath)
			dir := filepath.Dir(item.file.RelPath)
			prefix := "   "

			var fileLine string
			if dir != "." {
				fileLine = fmt.Sprintf("%s %s/%s", prefix, dir, name)
			} else {
				fileLine = fmt.Sprintf("%s %s", prefix, name)
			}

			if isCurrent {
				line = styles.Selected.Render(padRight(fileLine, width))
			} else {
				line = styles.FileName.Render(fileLine)
			}
		}

		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String()
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
