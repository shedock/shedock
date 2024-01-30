package ldd

import "strings"

type LddParser struct {
	// Data is the raw data returned by ldd <binary>
	Data []byte
}

// Parse parses the output of ldd <binary>
// E.g.:
// / # ldd $(which sh)
//
//	/lib/ld-musl-aarch64.so.1 (0xffffa1ed5000)
//	libc.musl-aarch64.so.1 => /lib/ld-musl-aarch64.so.1 (0xffffa1ed5000)
func (l *LddParser) Parse(output string) []Dependencies {
	lines := strings.Split(output, "\n")
	parsers := []Dependencies{}

	var currentParser *Dependencies
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/bin/") {
			if currentParser != nil {
				parsers = append(parsers, *currentParser)
			}
			currentParser = &Dependencies{
				Binary:       line,
				Dependencies: []Library{},
			}
		} else if strings.HasPrefix(line, "F") || strings.HasPrefix(line, ",") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				currentParser.Dependencies = append(currentParser.Dependencies, Library{
					Name:      parts[1],
					IsSymlink: strings.HasPrefix(line, "F"),
				})
			}
		}
	}

	if currentParser != nil {
		parsers = append(parsers, *currentParser)
	}

	return parsers
}
