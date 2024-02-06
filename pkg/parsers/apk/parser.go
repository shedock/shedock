package apk

import (
	"regexp"
	"strings"
)

// ApkParser is a parser for information returned by apk info
type ApkParser struct {
	// Data is the raw data returned by apk info -a <package>
	Data []byte
}

// Parse
// fish-3.6.1-r2 depends on:
// bc
// /bin/sh
// so:libc.musl-aarch64.so.1
// so:libgcc_s.so.1
// so:libintl.so.8
// so:libncursesw.so.6
// so:libpcre2-32.so.0
// so:libstdc++.so.6
func (p *ApkParser) DependsOn() ([]*PackageDependency, error) {
	regex := regexp.MustCompile(`depends on:((.|\n[^\n])*)`)
	matches := regex.FindStringSubmatch(string(p.Data))

	if len(matches) < 2 {
		return nil, nil // or return an error
	}

	dependenciesLines := strings.Split(matches[1], "\n")
	var dependencies []*PackageDependency

	for _, line := range dependenciesLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		p := &PackageDependency{}
		if strings.HasPrefix(line, "so:") {
			p.Type = SharedLibrary
			p.Name = strings.TrimPrefix(line, "so:")
		} else {
			p.Type = Binary
			p.Name = line
		}
		dependencies = append(dependencies, p)
	}

	return dependencies, nil
}

// Provides generates the list of commands shipped with a apk package
// fish-3.6.1-r2 provides:
// cmd:fish=3.6.1-r2
// cmd:fish_indent=3.6.1-r2
// cmd:fish_key_reader=3.6.1-r2
func (p *ApkParser) Provides() ([]*ProviderDependency, error) {
	regex := regexp.MustCompile(`provides:((.|\n[^\n])*)`)
	matches := regex.FindStringSubmatch(string(p.Data))

	if len(matches) < 2 {
		return nil, nil // or return an error
	}

	providersLines := strings.Split(matches[1], "\n")
	var providers []*ProviderDependency

	for _, line := range providersLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		p := &ProviderDependency{}
		if strings.HasPrefix(line, "cmd:") {
			binary := strings.TrimPrefix(line, "cmd:")
			binaryComponents := strings.Split(binary, "=")
			p.Version = binaryComponents[1]
			p.Name = binaryComponents[0]
		}
		providers = append(providers, p)
	}

	return providers, nil
}

func (p *ApkParser) Contains() ([]string, error) {
	regex := regexp.MustCompile(`contains:((.|\n[^\n])*)`)
	matches := regex.FindStringSubmatch(string(p.Data))

	if len(matches) < 2 {
		return nil, nil // or return an error
	}

	containsLines := strings.Split(matches[1], "\n")
	var contains []string

	for _, line := range containsLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		contains = append(contains, line)
	}

	return contains, nil
}

// Parse
// fish-3.6.1-r2 installed size:
// 10 MiB
func (p *ApkParser) InstalledSize() {
	// TODO
}
