package parser

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func init() { Register(&NginxParser{}) }

// NginxParser extracts key directives from nginx config files.
type NginxParser struct{}

func (n *NginxParser) Name() string { return "nginx" }

var (
	listenRe    = regexp.MustCompile(`listen\s+(\d+)`)
	serverRe    = regexp.MustCompile(`server_name\s+(.+?);`)
	proxyPassRe = regexp.MustCompile(`proxy_pass\s+(.+?);`)
	upstreamRe  = regexp.MustCompile(`upstream\s+(\w+)`)
)

func (n *NginxParser) Parse(path string) (*ParseResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := &ParseResult{
		Metadata: make(map[string]string),
	}

	var serverNames []string
	var upstreams []string

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if m := listenRe.FindStringSubmatch(line); m != nil {
			if port, err := strconv.Atoi(m[1]); err == nil {
				result.Ports = append(result.Ports, port)
			}
		}

		if m := serverRe.FindStringSubmatch(line); m != nil {
			names := strings.Fields(m[1])
			serverNames = append(serverNames, names...)
		}

		if m := proxyPassRe.FindStringSubmatch(line); m != nil {
			result.Keys = append(result.Keys, "proxy_pass:"+m[1])
		}

		if m := upstreamRe.FindStringSubmatch(line); m != nil {
			upstreams = append(upstreams, m[1])
		}
	}

	result.Services = upstreams
	if len(serverNames) > 0 {
		result.Metadata["server_names"] = strings.Join(serverNames, ", ")
	}
	if len(upstreams) > 0 {
		result.Metadata["upstream_count"] = strconv.Itoa(len(upstreams))
	}

	return result, sc.Err()
}
