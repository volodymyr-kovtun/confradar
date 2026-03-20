package parser

func init() { Register(&TextParser{}) }

// TextParser is a no-op parser for files that don't need structural extraction.
type TextParser struct{}

func (t *TextParser) Name() string { return "text" }

func (t *TextParser) Parse(path string) (*ParseResult, error) {
	return &ParseResult{
		Metadata: map[string]string{},
	}, nil
}
