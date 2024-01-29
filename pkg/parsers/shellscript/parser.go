package shellscript

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"mvdan.cc/sh/syntax"
)

// Script represents a shell script
type Script struct {
	// ScriptPath is the path to the shell script
	ScriptPath string
}

// Dependencies returns a list of external dependencies used in the script e.g. curl, wget, etc.
func (s *Script) Dependencies() ([]Dependency, error) {
	f, err := os.Open(s.ScriptPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	parser := syntax.NewParser()
	file, err := parser.Parse(f, s.ScriptPath)
	if err != nil {
		return nil, err
	}

	var finalDeps []Dependency

	funcDel := make(map[string]*syntax.FuncDecl)
	uniqDeps := make(map[string]Dependency)

	// Walk the AST and classify nodes
	syntax.Walk(file, func(node syntax.Node) bool {
		switch n := node.(type) {
		case *syntax.CallExpr:
			if len(n.Args) > 0 {
				word := n.Args[0]
				for _, part := range word.Parts {
					if lit, ok := part.(*syntax.Lit); ok {
						// add to uniqDeps if not already present
						if _, ok := uniqDeps[lit.Value]; !ok {
							uniqDeps[lit.Value] = Dependency{
								Name: lit.Value,
							}
						}
					}
				}
			}
		case *syntax.FuncDecl:
			funcDel[n.Name.Value] = n
		}
		return true
	})

	for _, d := range uniqDeps {
		// check if a command is a func decl
		if _, ok := funcDel[d.Name]; !ok {
			finalDeps = append(finalDeps, d)
		}
	}
	// sort by name
	sort.Slice(finalDeps, func(i, j int) bool {
		return finalDeps[i].Name < finalDeps[j].Name
	})

	return finalDeps, nil
}

// GetShell returns the shell the script is written in e.g. bash, zsh, fish
func (s *Script) GetShell() (string, error) {
	f, err := os.Open(s.ScriptPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#!") {
			// The shebang line typically contains the path to the shell, e.g. "/bin/bash" or "/usr/bin/env bash"
			parts := strings.Split(line, "/")
			possibleCommand := parts[len(parts)-1]
			sheBangElements := strings.Split(possibleCommand, " ")
			if len(sheBangElements) > 1 {
				// TODO handle case "/usr/bin/env bash -e"
				// https://www.shellcheck.net/wiki/SC2096
				return sheBangElements[len(sheBangElements)-1], nil
			}
			return parts[len(parts)-1], nil
		}
	}

	return "", fmt.Errorf("could not determine shell")
}

// Check whether the script is actually a shell script
func (s *Script) Validate() (bool, error) {
	// check if file exists
	if _, err := os.Stat(s.ScriptPath); os.IsNotExist(err) {
		return false, fmt.Errorf("file does not exist")
	}

	shell, err := s.GetShell()
	if err != nil {
		return false, err
	}
	if shell == "" {
		return false, fmt.Errorf("could not determine shell")
	}
	return true, nil
}