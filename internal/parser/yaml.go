package parser

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func init() { Register(&YAMLParser{}) }

// YAMLParser extracts structure from YAML config files.
type YAMLParser struct{}

func (y *YAMLParser) Name() string { return "yaml" }

func (y *YAMLParser) Parse(path string) (*ParseResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return &ParseResult{}, nil // gracefully skip unparseable files
	}

	result := &ParseResult{
		Versions: make(map[string]string),
		Metadata: make(map[string]string),
	}

	// Extract top-level keys.
	for k := range raw {
		result.Keys = append(result.Keys, k)
	}

	// Extract services from docker-compose style.
	if services, ok := raw["services"]; ok {
		if svcMap, ok := services.(map[string]any); ok {
			for name := range svcMap {
				result.Services = append(result.Services, name)
			}
			result.Metadata["service_count"] = strconv.Itoa(len(svcMap))
		}
	}

	// Extract version field.
	if v, ok := raw["version"]; ok {
		result.Versions["compose"] = fmt.Sprintf("%v", v)
	}

	// Extract ports from services.
	result.Ports = extractYAMLPorts(raw)

	result.Metadata["key_count"] = strconv.Itoa(len(result.Keys))
	return result, nil
}

func extractYAMLPorts(raw map[string]any) []int {
	var ports []int
	services, ok := raw["services"]
	if !ok {
		return ports
	}
	svcMap, ok := services.(map[string]any)
	if !ok {
		return ports
	}

	for _, svc := range svcMap {
		svcDef, ok := svc.(map[string]any)
		if !ok {
			continue
		}
		portsVal, ok := svcDef["ports"]
		if !ok {
			continue
		}
		portsList, ok := portsVal.([]any)
		if !ok {
			continue
		}
		for _, p := range portsList {
			switch v := p.(type) {
			case string:
				port := parsePortMapping(v)
				if port > 0 {
					ports = append(ports, port)
				}
			case int:
				ports = append(ports, v)
			case float64:
				ports = append(ports, int(v))
			}
		}
	}
	return ports
}

func parsePortMapping(s string) int {
	// Handle "8080:80", "127.0.0.1:8080:80", "8080"
	parts := splitLast(s, ":")
	if len(parts) < 2 {
		if port, err := strconv.Atoi(s); err == nil {
			return port
		}
		return 0
	}
	// Host port is the second-to-last part.
	hostPart := parts[0]
	// Could be "127.0.0.1:8080" or just "8080".
	if idx := lastIndexByte(hostPart, ':'); idx >= 0 {
		hostPart = hostPart[idx+1:]
	}
	port, _ := strconv.Atoi(hostPart)
	return port
}

func splitLast(s, sep string) []string {
	idx := lastIndexStr(s, sep)
	if idx < 0 {
		return []string{s}
	}
	return []string{s[:idx], s[idx+len(sep):]}
}

func lastIndexStr(s, substr string) int {
	for i := len(s) - len(substr); i >= 0; i-- {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func lastIndexByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}
