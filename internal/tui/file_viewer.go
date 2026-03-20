package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxFileSize = 64 * 1024 // 64KB
const maxLines = 1000

// FileViewer displays file contents with line numbers.
type FileViewer struct {
	lines    []string
	filePath string
	offset   int
	height   int
	redact   bool
}

// NewFileViewer creates a file viewer.
func NewFileViewer() FileViewer {
	return FileViewer{height: 20, redact: true}
}

// SetHeight sets the visible height.
func (fv *FileViewer) SetHeight(h int) {
	fv.height = h
	if fv.height < 1 {
		fv.height = 1
	}
}

// LoadFile reads a file for display.
func (fv *FileViewer) LoadFile(path string) error {
	fv.filePath = path
	fv.offset = 0

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.Size() > maxFileSize {
		fv.lines = []string{
			fmt.Sprintf("File too large (%d bytes). Showing first %d bytes.", info.Size(), maxFileSize),
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(data) > maxFileSize {
		data = data[:maxFileSize]
	}

	allLines := strings.Split(string(data), "\n")
	if len(allLines) > maxLines {
		allLines = allLines[:maxLines]
		allLines = append(allLines, fmt.Sprintf("... truncated at %d lines", maxLines))
	}

	fv.lines = allLines
	return nil
}

// ScrollUp scrolls the viewer up.
func (fv *FileViewer) ScrollUp(n int) {
	fv.offset -= n
	if fv.offset < 0 {
		fv.offset = 0
	}
}

// ScrollDown scrolls the viewer down.
func (fv *FileViewer) ScrollDown(n int) {
	fv.offset += n
	max := len(fv.lines) - fv.height
	if max < 0 {
		max = 0
	}
	if fv.offset > max {
		fv.offset = max
	}
}

// ToggleRedact toggles value redaction for .env files.
func (fv *FileViewer) ToggleRedact() {
	fv.redact = !fv.redact
}

// IsRedacting returns whether values are being redacted.
func (fv *FileViewer) IsRedacting() bool {
	return fv.redact
}

// View renders the file contents.
func (fv *FileViewer) View(styles Styles, width int) string {
	if len(fv.lines) == 0 {
		return styles.Dim.Render("  Select a file to preview")
	}

	isEnv := isEnvFile(fv.filePath)
	var b strings.Builder

	// Header.
	header := fmt.Sprintf(" %s (%d lines)", filepath.Base(fv.filePath), len(fv.lines))
	b.WriteString(styles.Accent.Render(header))
	b.WriteByte('\n')

	end := fv.offset + fv.height - 1
	if end > len(fv.lines) {
		end = len(fv.lines)
	}

	lineNumWidth := len(fmt.Sprintf("%d", end))

	for i := fv.offset; i < end; i++ {
		line := fv.lines[i]

		// Redact .env values if needed.
		if isEnv && fv.redact {
			line = redactEnvLine(line)
		}

		// Truncate long lines.
		maxLineWidth := width - lineNumWidth - 4
		if maxLineWidth > 0 && len(line) > maxLineWidth {
			line = line[:maxLineWidth-1] + "…"
		}

		lineNum := fmt.Sprintf("%*d", lineNumWidth, i+1)
		b.WriteString(styles.Dim.Render(lineNum+" │ "))
		b.WriteString(styles.Normal.Render(line))
		b.WriteByte('\n')
	}

	return b.String()
}

func isEnvFile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return strings.HasPrefix(base, ".env") || base == ".envrc" || base == ".secrets"
}

func redactEnvLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return line
	}
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return line
	}
	return line[:idx+1] + "••••••"
}
