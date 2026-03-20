package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func init() { Register(&JSONParser{}) }

// JSONParser extracts structure from JSON config files.
type JSONParser struct{}

func (j *JSONParser) Name() string { return "json" }

func (j *JSONParser) Parse(path string) (*ParseResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return &ParseResult{}, nil
	}

	result := &ParseResult{
		Versions: make(map[string]string),
		Metadata: make(map[string]string),
	}

	for k := range raw {
		result.Keys = append(result.Keys, k)
	}

	base := strings.ToLower(strings.TrimSuffix(
		strings.TrimSuffix(baseName(path), ".json"), ".json5"))

	// package.json special handling.
	if strings.HasSuffix(base, "package") {
		extractPackageJSON(raw, result)
	}

	// tsconfig.json special handling.
	if strings.Contains(base, "tsconfig") {
		if co, ok := raw["compilerOptions"]; ok {
			if coMap, ok := co.(map[string]any); ok {
				var opts []string
				for k := range coMap {
					opts = append(opts, k)
				}
				result.Metadata["compiler_options"] = strconv.Itoa(len(opts))
			}
		}
	}

	result.Metadata["key_count"] = strconv.Itoa(len(result.Keys))
	return result, nil
}

func extractPackageJSON(raw map[string]any, result *ParseResult) {
	if name, ok := raw["name"].(string); ok {
		result.Metadata["package_name"] = name
	}
	if ver, ok := raw["version"].(string); ok {
		result.Metadata["package_version"] = ver
	}

	// Extract engine versions.
	if engines, ok := raw["engines"]; ok {
		if engMap, ok := engines.(map[string]any); ok {
			for runtime, ver := range engMap {
				result.Versions[runtime] = fmt.Sprintf("%v", ver)
			}
		}
	}

	// Extract script names.
	if scripts, ok := raw["scripts"]; ok {
		if scriptMap, ok := scripts.(map[string]any); ok {
			var names []string
			for k := range scriptMap {
				names = append(names, k)
			}
			result.Metadata["scripts"] = strings.Join(names, ", ")
			result.Metadata["script_count"] = strconv.Itoa(len(names))
		}
	}

	// Count dependencies.
	countDeps := func(key string) {
		if deps, ok := raw[key]; ok {
			if depMap, ok := deps.(map[string]any); ok {
				result.Metadata[key+"_count"] = strconv.Itoa(len(depMap))
			}
		}
	}
	countDeps("dependencies")
	countDeps("devDependencies")
}

// baseName returns the base name without directory.
func baseName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}
