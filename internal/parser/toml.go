package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pelletier/go-toml/v2"
)

func init() { Register(&TOMLParser{}) }

// TOMLParser extracts structure from TOML config files.
type TOMLParser struct{}

func (t *TOMLParser) Name() string { return "toml" }

func (t *TOMLParser) Parse(path string) (*ParseResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return &ParseResult{}, nil
	}

	result := &ParseResult{
		Versions: make(map[string]string),
		Metadata: make(map[string]string),
	}

	for k := range raw {
		result.Keys = append(result.Keys, k)
	}

	// pyproject.toml: extract project version.
	if project, ok := raw["project"]; ok {
		if pm, ok := project.(map[string]any); ok {
			if v, ok := pm["version"]; ok {
				result.Versions["python_project"] = toString(v)
			}
			if name, ok := pm["name"]; ok {
				result.Metadata["project_name"] = toString(name)
			}
			if req, ok := pm["requires-python"]; ok {
				result.Versions["python"] = toString(req)
			}
		}
	}

	// Cargo.toml: extract package info.
	if pkg, ok := raw["package"]; ok {
		if pm, ok := pkg.(map[string]any); ok {
			if v, ok := pm["version"]; ok {
				result.Versions["cargo"] = toString(v)
			}
			if name, ok := pm["name"]; ok {
				result.Metadata["package_name"] = toString(name)
			}
			if edition, ok := pm["edition"]; ok {
				result.Versions["rust_edition"] = toString(edition)
			}
		}
	}

	result.Metadata["key_count"] = strconv.Itoa(len(result.Keys))
	return result, nil
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
