package renderer

import (
	"io"

	"github.com/volodymyrkovtun/confradar/internal/scanner"
	"gopkg.in/yaml.v3"
)

// YAMLRenderer outputs scan results as YAML.
type YAMLRenderer struct{}

type yamlOutput struct {
	Root       string         `yaml:"root"`
	TotalFiles int            `yaml:"total_files"`
	Duration   string         `yaml:"duration"`
	Categories []yamlCategory `yaml:"categories"`
}

type yamlCategory struct {
	Name  string     `yaml:"name"`
	Icon  string     `yaml:"icon"`
	Count int        `yaml:"count"`
	Files []yamlFile `yaml:"files"`
}

type yamlFile struct {
	RelPath     string `yaml:"rel_path"`
	PatternName string `yaml:"pattern_name"`
	Parser      string `yaml:"parser"`
	Size        int64  `yaml:"size"`
}

// Render writes YAML output to the writer.
func (y *YAMLRenderer) Render(result *scanner.ScanResult, w io.Writer) error {
	out := yamlOutput{
		Root:       result.Root,
		TotalFiles: result.TotalFiles,
		Duration:   result.Duration.String(),
	}

	for _, cat := range result.Ordered {
		yc := yamlCategory{
			Name:  cat.Name,
			Icon:  cat.Icon,
			Count: len(cat.Files),
		}
		for _, f := range cat.Files {
			yc.Files = append(yc.Files, yamlFile{
				RelPath:     f.RelPath,
				PatternName: f.PatternName,
				Parser:      f.Parser,
				Size:        f.Size,
			})
		}
		out.Categories = append(out.Categories, yc)
	}

	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(out)
}
