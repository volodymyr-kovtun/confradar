package parser

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func init() { Register(&DockerfileParser{}) }

// DockerfileParser extracts metadata from Dockerfiles.
type DockerfileParser struct{}

func (d *DockerfileParser) Name() string { return "dockerfile" }

var (
	fromRe   = regexp.MustCompile(`(?i)^FROM\s+(.+?)(\s+AS\s+\w+)?$`)
	exposeRe = regexp.MustCompile(`(?i)^EXPOSE\s+(.+)$`)
	envRe    = regexp.MustCompile(`(?i)^ENV\s+(\w+)[=\s]+(.*)$`)
)

func (d *DockerfileParser) Parse(path string) (*ParseResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := &ParseResult{
		Values:   make(map[string]string),
		Versions: make(map[string]string),
		Metadata: make(map[string]string),
	}

	var stages []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		// FROM directives.
		if m := fromRe.FindStringSubmatch(line); m != nil {
			image := strings.TrimSpace(m[1])
			stages = append(stages, image)

			// Extract version from image tag.
			if parts := strings.SplitN(image, ":", 2); len(parts) == 2 {
				imageName := parts[0]
				tag := parts[1]

				// Detect runtime versions.
				switch {
				case strings.Contains(imageName, "node"):
					result.Versions["node"] = extractVersion(tag)
				case strings.Contains(imageName, "python"):
					result.Versions["python"] = extractVersion(tag)
				case strings.Contains(imageName, "ruby"):
					result.Versions["ruby"] = extractVersion(tag)
				case strings.Contains(imageName, "golang") || strings.Contains(imageName, "go"):
					result.Versions["go"] = extractVersion(tag)
				}

				if tag == "latest" {
					result.Metadata["uses_latest"] = "true"
				}
			}
		}

		// EXPOSE directives.
		if m := exposeRe.FindStringSubmatch(line); m != nil {
			for _, portStr := range strings.Fields(m[1]) {
				portStr = strings.TrimSuffix(portStr, "/tcp")
				portStr = strings.TrimSuffix(portStr, "/udp")
				if port, err := strconv.Atoi(portStr); err == nil {
					result.Ports = append(result.Ports, port)
				}
			}
		}

		// ENV directives.
		if m := envRe.FindStringSubmatch(line); m != nil {
			key := m[1]
			val := strings.TrimSpace(m[2])
			result.Keys = append(result.Keys, key)
			result.Values[key] = val
		}
	}

	result.Services = stages
	result.Metadata["stage_count"] = strconv.Itoa(len(stages))
	return result, sc.Err()
}

// extractVersion tries to pull a semver-like string from a Docker tag.
func extractVersion(tag string) string {
	// Tags often look like "20-alpine", "3.12-slim", "1.22.0"
	parts := strings.SplitN(tag, "-", 2)
	return parts[0]
}
