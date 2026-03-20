package tui

import (
	"fmt"
	"strings"
)

// SearchOverlay manages the fuzzy search input.
type SearchOverlay struct {
	query   string
	active  bool
	cursor  int
}

// Activate shows the search overlay.
func (s *SearchOverlay) Activate() {
	s.active = true
	s.query = ""
	s.cursor = 0
}

// Deactivate hides the search overlay.
func (s *SearchOverlay) Deactivate() {
	s.active = false
}

// IsActive returns whether search is active.
func (s *SearchOverlay) IsActive() bool {
	return s.active
}

// Query returns the current search query.
func (s *SearchOverlay) Query() string {
	return s.query
}

// InsertChar adds a character to the query.
func (s *SearchOverlay) InsertChar(ch rune) {
	s.query = s.query[:s.cursor] + string(ch) + s.query[s.cursor:]
	s.cursor++
}

// Backspace removes the character before the cursor.
func (s *SearchOverlay) Backspace() {
	if s.cursor > 0 {
		s.query = s.query[:s.cursor-1] + s.query[s.cursor:]
		s.cursor--
	}
}

// Clear resets the query.
func (s *SearchOverlay) Clear() {
	s.query = ""
	s.cursor = 0
}

// View renders the search bar.
func (s *SearchOverlay) View(styles Styles, width int) string {
	if !s.active {
		return ""
	}

	prompt := styles.SearchPrompt.Render(" / ")
	query := s.query
	cursorChar := "█"

	line := fmt.Sprintf("%s%s%s", prompt, query, styles.Accent.Render(cursorChar))
	return line + strings.Repeat(" ", max(0, width-len(prompt)-len(query)-1))
}

