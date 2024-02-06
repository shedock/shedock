package ldd

import (
	"strings"
)

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
func (l *LddParser) Parse() []Library {
	output := string(l.Data)

	lines := strings.Split(output, "\n")

	var parsers []Library

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=>")
		if len(parts) == 2 {
			soName := strings.TrimSpace(parts[0])
			soNameParts := strings.Fields(soName)
			soName = soNameParts[len(soNameParts)-1]
			fullPathWithAddress := strings.TrimSpace(parts[1])
			fullPath := strings.Split(fullPathWithAddress, " ")[0]
			parsers = append(parsers, Library{
				SoName:   soName,
				FullPath: fullPath,
			})

		}
	}

	return parsers
}
