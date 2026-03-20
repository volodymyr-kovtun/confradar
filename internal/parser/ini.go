package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func init() { Register(&INIParser{}) }

// INIParser extracts key-value pairs from INI/cfg style files.
type INIParser struct{}

func (i *INIParser) Name() string { return "ini" }

func (i *INIParser) Parse(path string) (*ParseResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := &ParseResult{
		Values:   make(map[string]string),
		Metadata: make(map[string]string),
	}

	section := ""
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			idx = strings.IndexByte(line, ':')
		}
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		fullKey := key
		if section != "" {
			fullKey = section + "." + key
		}

		result.Keys = append(result.Keys, fullKey)
		result.Values[fullKey] = val
	}

	result.Metadata["key_count"] = strconv.Itoa(len(result.Keys))
	return result, sc.Err()
}
