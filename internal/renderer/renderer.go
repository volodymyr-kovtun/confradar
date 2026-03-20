// Package renderer provides output formatters for scan results.
package renderer

import (
	"fmt"
	"io"

	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// Renderer writes scan results to an output stream.
type Renderer interface {
	Render(result *scanner.ScanResult, w io.Writer) error
}

// New returns a Renderer for the given format name.
func New(format string, noColor bool) (Renderer, error) {
	switch format {
	case "tree", "":
		return &TreeRenderer{NoColor: noColor}, nil
	case "json":
		return &JSONRenderer{}, nil
	case "yaml":
		return &YAMLRenderer{}, nil
	case "table":
		return &TableRenderer{NoColor: noColor}, nil
	default:
		return nil, fmt.Errorf("unknown format: %q (valid: tree, json, yaml, table)", format)
	}
}
