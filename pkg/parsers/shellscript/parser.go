package shellscript

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
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
	variableDel := make(map[string]*syntax.Assign)
	uniqDeps := make(map[string]Dependency)
	var allArgs []string

	alphaRegex, err := regexp.Compile("^[a-z]{2,}$")
	if err != nil {
		return nil, err
	}

	// Walk the AST and classify nodes
	syntax.Walk(file, func(node syntax.Node) bool {
		switch n := node.(type) {
		// case *syntax.CmdSubst:
		// case *syntax.Word:
		case *syntax.CallExpr:
			// fmt.Println("CallExpr", reflect.TypeOf(n))
			if len(n.Args) > 0 {
				word := n.Args[0]
				// fmt.Println(word.Parts)
				for _, part := range word.Parts {
					// fmt.Println(part, reflect.TypeOf(part))
					if lit, ok := part.(*syntax.Lit); ok {
						// add to uniqDeps if not already present
						if _, ok := uniqDeps[lit.Value]; !ok {
							uniqDeps[lit.Value] = Dependency{
								Name: lit.Value,
							}
						}

						// collect arguments
						args := []string{}
						for _, arg := range n.Args[1:] {
							// fmt.Println(arg.Parts, reflect.TypeOf(arg.Parts))
							for _, argPart := range arg.Parts {
								if argLit, ok := argPart.(*syntax.Lit); ok {
									args = append(args, argLit.Value)
								}
							}
						}
						// make sure the args are appended to the dependency if it already exists
						if dep, ok := uniqDeps[lit.Value]; ok {
							// TODO only add unique args
							dep.Args = append(dep.Args, args...)
							uniqDeps[lit.Value] = dep
							allArgs = append(allArgs, args...)
						}
						// // add arguments to the dependency
						// uniqDeps[lit.Value] = Dependency{
						// 	Name: lit.Value,
						// 	Args: args,
						// }
					}
				}
			}
		case *syntax.FuncDecl:
			funcDel[n.Name.Value] = n
		case *syntax.Word:
			if len(n.Parts) > 0 {
				for _, part := range n.Parts {
					if lit, ok := part.(*syntax.Lit); ok {
						if alphaRegex.MatchString(lit.Value) {
							if _, ok := uniqDeps[lit.Value]; !ok {
								uniqDeps[lit.Value] = Dependency{
									Name: lit.Value,
								}
							}
						}
					}
				}
			}
		case *syntax.Assign:
			// get variable declarations, so that we can remove them from the finalDeps later
			variableDel[n.Name.Value] = n
		case *syntax.CmdSubst:
			// TODO handle command substitution
			fmt.Println("CmdSubst", reflect.TypeOf(n))
		}
		return true
	})

	// var matchvar []string
	for _, d := range uniqDeps {
		// filter out function declarations
		if _, ok := funcDel[d.Name]; !ok {
			// filter out variable declarations
			if _, ok := variableDel[d.Name]; !ok {
				finalDeps = append(finalDeps, d)
			}
		}
	}

	// // remove dep from uniqDeps that are also variables
	// for k := range variableDel {
	// 	if _, ok := uniqDeps[k]; ok {
	// 		matchvar = append(matchvar, k)
	// 	}
	// }
	// fmt.Println("matchvar", matchvar)

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
