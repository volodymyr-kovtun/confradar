package renderer

import (
	"encoding/json"
	"io"

	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// JSONRenderer outputs scan results as JSON.
type JSONRenderer struct{}

type jsonOutput struct {
	Root       string           `json:"root"`
	TotalFiles int              `json:"total_files"`
	Duration   string           `json:"duration"`
	Categories []jsonCategory   `json:"categories"`
}

type jsonCategory struct {
	Name     string     `json:"name"`
	Icon     string     `json:"icon"`
	Color    string     `json:"color"`
	Priority int        `json:"priority"`
	Count    int        `json:"count"`
	Files    []jsonFile `json:"files"`
}

type jsonFile struct {
	Path        string `json:"path"`
	RelPath     string `json:"rel_path"`
	PatternName string `json:"pattern_name"`
	Parser      string `json:"parser"`
	Size        int64  `json:"size"`
}

// Render writes JSON output to the writer.
func (j *JSONRenderer) Render(result *scanner.ScanResult, w io.Writer) error {
	out := jsonOutput{
		Root:       result.Root,
		TotalFiles: result.TotalFiles,
		Duration:   result.Duration.String(),
	}

	for _, cat := range result.Ordered {
		jc := jsonCategory{
			Name:     cat.Name,
			Icon:     cat.Icon,
			Color:    cat.Color,
			Priority: cat.Priority,
			Count:    len(cat.Files),
		}
		for _, f := range cat.Files {
			jc.Files = append(jc.Files, jsonFile{
				Path:        f.Path,
				RelPath:     f.RelPath,
				PatternName: f.PatternName,
				Parser:      f.Parser,
				Size:        f.Size,
			})
		}
		out.Categories = append(out.Categories, jc)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
