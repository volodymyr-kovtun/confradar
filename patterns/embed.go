// Package patterns provides the embedded default pattern definitions.
package patterns

import _ "embed"

// DefaultYAML contains the built-in pattern definitions embedded at compile time.
//
//go:embed default.yml
var DefaultYAML []byte
