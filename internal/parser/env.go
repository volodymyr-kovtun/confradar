package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func init() { Register(&EnvParser{}) }

// EnvParser extracts key-value pairs from .env files.
type EnvParser struct{}

func (e *EnvParser) Name() string { return "env" }

func (e *EnvParser) Parse(path string) (*ParseResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := &ParseResult{
		Values:   make(map[string]string),
		Metadata: make(map[string]string),
	}

	var ports []int
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// Strip surrounding quotes.
		val = stripQuotes(val)

		// Remove inline comments.
		if !strings.HasPrefix(val, "'") {
			if ci := strings.Index(val, " #"); ci >= 0 {
				val = strings.TrimSpace(val[:ci])
			}
		}

		result.Keys = append(result.Keys, key)
		result.Values[key] = val

		// Detect port values.
		if isPortKey(key) {
			if port, err := strconv.Atoi(val); err == nil && port > 0 && port <= 65535 {
				ports = append(ports, port)
			}
		}
	}

	result.Ports = ports
	result.Metadata["key_count"] = strconv.Itoa(len(result.Keys))
	return result, sc.Err()
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func isPortKey(key string) bool {
	k := strings.ToUpper(key)
	return strings.Contains(k, "PORT") || strings.HasSuffix(k, "_P") || k == "LISTEN"
}
