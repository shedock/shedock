package parsers

import "strings"

type Library struct {
	Name      string
	IsSymlink bool
}

type Dependencies struct {
	Binary       string
	Dependencies []Library
}

func ParseLddOutput(output string) []Dependencies {
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
