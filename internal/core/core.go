package core

import (
	"fmt"
	"log"
	"regexp"
	"shedock/internal/insights"
	"shedock/internal/instance"
	"shedock/pkg/docker/file"
	apkTypes "shedock/pkg/parsers/apk"
	"shedock/pkg/parsers/ldd"
	shellScriptTypes "shedock/pkg/parsers/shellscript"
	"shedock/pkg/shell"
	"strings"
)

type ImageBuilder struct {
	Script             *shellScriptTypes.Script
	Shell              shell.Shell
	systemBuiltins     []string
	scriptDeps         []shellScriptTypes.Dependency
	cmdOnApk           []string
	shellbuiltns       []string
	cmdNotOnApk        []string
	cmdNotSupported    []string
	sharedLibs         []ldd.Library
	usedSystemBuiltins []string
}

func NewImageBuilder(
	script *shellScriptTypes.Script,
	shell shell.Shell,
) *ImageBuilder {
	return &ImageBuilder{
		Script: script,
		Shell:  shell,
	}
}

func (i *ImageBuilder) Build() error {
	var filteredDeps []string

	err := i.LoadSystemBuiltins()
	if err != nil {
		return err
	}

	err = i.LoadScriptDeps()
	if err != nil {
		return err
	}

	err = i.UsedSystemBuiltins()
	if err != nil {
		return err
	}

	err = i.LoadShellBuiltins()
	if err != nil {
		return err
	}

	// remove not-supported commands from script deps
	// remove shell-builtins and system-builtins from script deps and find what we can get from package manager
	filteredDeps = i.FilterCmdsToInstall()
	log.Println("Filtered dependencies: ", filteredDeps)
	err = i.DependenciesAvailableOnPackageHost(filteredDeps)
	if err != nil {
		return err
	}

	log.Println("Shell builtins: ", i.GetShellBuiltins())
	log.Println("Script deps after filter system builtins: ", i.GetUsedSystemBuiltins())
	log.Println("Script deps after filter shell builtins: ", i.FilterShellBuiltins())
	log.Println("Commands not found on apk: ", i.GetCmdNotOnApk())
	log.Println("Commands not supported in containerized environment: ", i.GetCmdNotSupported())

	err = i.LoadAllSharedLibs()
	if err != nil {
		return err
	}
	cmdOnApk := i.GetCmdOnApk()
	log.Println("Commands available on apk: ", cmdOnApk)
	// commands not available on apk come under not found
	var notFound []string
	for _, dep := range filteredDeps {
		var found bool
		for _, cmd := range cmdOnApk {
			if dep == cmd.Name {
				found = true
				break
			}
		}
		if !found {
			notFound = append(notFound, dep)
		}
	}
	i.UpdateCmdsNotFound(notFound)

	var deps file.Dependencies
	var bins []file.Dependency
	var libs []file.Dependency

	shell, _ := i.Script.GetShell()

	systemBuiltins := i.GetUsedSystemBuiltins()
	externalCommands := i.GetCmdOnApk()
	sharedLibs := i.GetSharedLibs()

	for _, cmd := range systemBuiltins {
		bins = append(bins, file.Dependency{
			FromPath: cmd.Path,
			ToPath:   cmd.Path,
		})
	}

	for _, cmd := range externalCommands {
		bins = append(bins, file.Dependency{
			FromPath: cmd.Path,
			ToPath:   cmd.Path,
		})
	}

	for _, lib := range sharedLibs {
		libs = append(libs, file.Dependency{
			FromPath:   lib.FullPath,
			ToPath:     lib.FullPath,
			Requiredby: lib.DependencyOf,
		})
	}

	deps = file.Dependencies{
		Bin: bins,
		Lib: libs,
	}

	file := &file.Dockerfile{
		DependenciesToInstall: externalCommands,
		Script:                i.Script.ScriptPath,
		ShellPath:             shell,
		Dependencies:          deps,
	}

	_, err = file.Generate()
	if err != nil {
		return err
	}

	return nil
}

// System dependencies are the binaries that are installed in the base image
// E.g. ls, cat, etc.
func (i *ImageBuilder) LoadSystemBuiltins() error {
	container := instance.GetDockerInstance()
	output, err := container.ExecCommand("for binary in /bin/* /usr/bin/*; do echo \"$binary\"; done")
	if err != nil {
		log.Fatalf("%v", err)
	}
	var builtins []string
	// re := regexp.MustCompile(`\\x[0-9A-Fa-f]{2}`)
	for _, builtin := range strings.Split(output, "\n") {
		if builtin != "" {
			filteredPath := strings.Trim(builtin, "\n")
			filteredPath = strings.Map(func(r rune) rune {
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '/' {
					return r
				}
				return -1
			}, filteredPath)
			if filteredPath == "" {
				continue
			}
			builtins = append(builtins, filteredPath)
		}
	}
	i.systemBuiltins = builtins

	return nil
}

// Shell builtins are the binaries that are installed by the shell
// E.g. [[, readarray, zmodload, etc.
func (i *ImageBuilder) LoadShellBuiltins() error {
	container := instance.GetDockerInstance()
	_, err := container.ExecCommand(i.Shell.InstallShellCommand())
	if err != nil {
		log.Fatalf("%v", err)
	}

	output, err := container.ExecCommand(i.Shell.CommandToFindBuiltinCommands())
	if err != nil {
		log.Fatalf("%v", err)
	}

	builtins, err := i.Shell.ParseBuiltins(output)
	if err != nil {
		log.Fatalf("Failed to parse builtins: %v", err)
	}
	i.shellbuiltns = builtins
	return nil
}

// Shell libraries are the libraries that are required by the shell. Also called as "shared libraries"
// E.g. libreadline, libncurses, etc.
func (i *ImageBuilder) GetShellLibraries() ([]*apkTypes.PackageDependency, error) {
	container := instance.GetDockerInstance()
	// Execute command in the container and get the output
	output, err := container.ExecCommand(i.Shell.InfoCommand())
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Parse the output
	dependencies, err := i.Shell.Dependencies([]byte(output))
	if err != nil {
		log.Fatalf("Failed to parse dependencies: %v", err)
	}

	return dependencies, nil
}

func (i *ImageBuilder) GetShellBinaryBuiltins() ([]string, error) {
	container := instance.GetDockerInstance()
	_, err := container.ExecCommand(i.Shell.InstallShellCommand())
	if err != nil {
		return []string{}, fmt.Errorf("%v", err)
	}

	// Execute the info command again to get the stuff installed by the shell
	output, err := container.ExecCommand(i.Shell.InfoCommand())
	if err != nil {
		return []string{}, fmt.Errorf("%v", err)
	}

	// Parse the output
	binDeps, err := i.Shell.BinaryDependencies([]byte(output))
	if err != nil {
		return []string{}, fmt.Errorf("failed to parse dependencies: %v", err)
	}

	return binDeps, nil
}

// Find dependencies that can be installed by the package manager (for now apk)
func (i *ImageBuilder) DependenciesAvailableOnPackageHost(deps []string) error {
	container := instance.GetDockerInstance()
	var depsOnApk []string
	var err error

	for _, dep := range deps {
		// Execute the info command again to get the stuff installed by the shell
		output, err := container.ExecCommand(fmt.Sprintf("apk info -a %s", dep))
		if err != nil {
			return err
		}

		// Parse the output
		parser := apkTypes.ApkParser{Data: []byte(output)}
		dependencies, err := parser.Provides()
		if err != nil {
			return err
		}

		for _, dependency := range dependencies {
			if dependency.Name == dep {
				depsOnApk = append(depsOnApk, dependency.Name)
			}
		}
	}

	i.cmdOnApk = depsOnApk
	return err
}

func (i *ImageBuilder) LoadAllSharedLibs() error {
	// Combine libraries from
	// 1. System Builtins used
	// 2. Shell used
	// 3. Packages on APK used
	container := instance.GetDockerInstance()

	var deps []string
	var sl []ldd.Library
	uniqueLibs := make(map[string]bool)

	deps = append(deps, i.usedSystemBuiltins...)

	// install all packages that are available on apk
	for _, dep := range i.cmdOnApk {
		_, err := container.ExecCommand(fmt.Sprintf("apk add %s", dep))
		if err != nil {
			return err
		}
	}

	deps = append(deps, i.cmdOnApk...)
	shellname, err := i.Script.GetShell()
	if err != nil {
		return err
	}

	// the shell is also a dependency
	deps = append(deps, shellname)

	for _, dep := range deps {
		depList := make(map[string][]string)

		// Find the shared libraries required by the binary
		output, err := container.ExecCommand(fmt.Sprintf("ldd $(which %s)", dep))
		if err != nil {
			return err
		}
		lddParser := ldd.LddParser{Data: []byte(output)}
		libs := lddParser.Parse()

		if len(libs) > 0 {
			for _, lib := range libs {
				if _, exists := uniqueLibs[lib.SoName]; !exists {
					depList[lib.SoName] = append(depList[lib.SoName], dep)
					lib.DependencyOf = depList[lib.SoName]

					sl = append(sl, *lib)
					uniqueLibs[lib.SoName] = true
				} else {
					// update DependencyOf of the shared library
					for i, l := range sl {
						if l.SoName == lib.SoName {
							sl[i].DependencyOf = append(sl[i].DependencyOf, dep)
						}
					}
				}
			}
		}
	}

	i.sharedLibs = sl

	return nil
}

func (i *ImageBuilder) LoadScriptDeps() error {
	scriptDeps, err := i.Script.Dependencies()
	if err != nil {
		return err
	}
	i.scriptDeps = scriptDeps
	return err
}

func (i *ImageBuilder) UsedSystemBuiltins() error {
	scriptDeps := i.GetScriptDeps()
	systemBuiltins := i.GetSystemBuiltins()

	var usedSystemBuiltins []string
	for _, dep := range scriptDeps {
		for _, builtin := range systemBuiltins {
			// get basename of the binary
			baseBuiltIn := strings.Split(builtin, "/")
			builtinName := baseBuiltIn[len(baseBuiltIn)-1]

			if builtinName == dep.Name {
				usedSystemBuiltins = append(usedSystemBuiltins, builtinName)
			}
		}
	}
	i.usedSystemBuiltins = usedSystemBuiltins

	return nil
}

func (i *ImageBuilder) FilterCmdsToInstall() []string {
	var filteredDeps []string

	shellBuiltins := i.GetShellBuiltins()
	scriptDeps := i.GetScriptDeps()
	systemBuiltins := i.GetSystemBuiltins()

	for _, dep := range scriptDeps {
		var found bool

		// Check if the dependency is a shell builtin
		for _, builtin := range shellBuiltins {
			if dep.Name == builtin {
				found = true
				break
			}
		}

		// Check if the dependency is a system builtin
		for _, builtin := range systemBuiltins {
			if strings.Contains(builtin, dep.Name) {
				found = true
				break
			}
		}

		// Check if the dependency cannot be used in a containerized environment
		for _, cmd := range insights.NOT_SUPPORTED_COMMANDS {
			if dep.Name == cmd {
				found = true
				i.cmdNotSupported = append(i.cmdNotSupported, dep.Name)
				break
			}
		}

		if !found {
			// its not a builtin, so we need to install it
			filteredDeps = append(filteredDeps, dep.Name)
		}
	}
	return filteredDeps
}

// WIP
func (i *ImageBuilder) FilterShellBuiltins() []string {
	var filteredDeps []string

	shellBuiltins := i.GetShellBuiltins()
	scriptDeps := i.GetScriptDeps()

	for _, dep := range scriptDeps {
		var found bool
		for _, builtin := range shellBuiltins {
			if strings.Contains(builtin, dep.Name) {
				found = true
				break
			}
		}
		if !found {
			filteredDeps = append(filteredDeps, dep.Name)
		}
	}
	return filteredDeps
}

func (i *ImageBuilder) GetSharedLibs() []ldd.Library {
	return i.sharedLibs
}

func (i *ImageBuilder) GetUsedSystemBuiltins() []apkTypes.Package {
	var deps []apkTypes.Package
	container := instance.GetDockerInstance()

	r := regexp.MustCompile(`(/[^:\s]+)`)

	for _, dep := range i.usedSystemBuiltins {
		// find location of the binary
		output, err := container.ExecCommand(fmt.Sprintf("which %s", dep))
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Extract the file path from the output
		matches := r.FindStringSubmatch(output)
		if len(matches) < 2 {
			log.Fatalf("GetUsedSystemBuiltins: Could not extract file path from output: %s", output)
		}
		depPath := matches[1]

		deps = append(deps, apkTypes.Package{
			Name: dep,
			Path: depPath,
		})
	}
	return deps
	// return i.usedSystemBuiltins
}

func (i *ImageBuilder) GetSystemBuiltins() []string {
	return i.systemBuiltins
}

func (i *ImageBuilder) GetScriptDeps() []shellScriptTypes.Dependency {
	return i.scriptDeps
}

func (i *ImageBuilder) GetShellBuiltins() []string {
	return i.shellbuiltns
}

func (i *ImageBuilder) GetCmdOnApk() []apkTypes.Package {
	var deps []apkTypes.Package
	container := instance.GetDockerInstance()

	r := regexp.MustCompile(`(/[^:\s]+)`)

	for _, dep := range i.cmdOnApk {
		// find location of the binary
		output, err := container.ExecCommand(fmt.Sprintf("which %s", dep))
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Extract the file path from the output
		matches := r.FindStringSubmatch(output)

		if len(matches) >= 1 {
			depPath := matches[0]

			deps = append(deps, apkTypes.Package{
				Name: dep,
				Path: depPath,
			})
		}
	}
	return deps
}

func (i *ImageBuilder) GetCmdNotOnApk() []string {
	return i.cmdNotOnApk
}

func (i *ImageBuilder) GetCmdNotSupported() []string {
	return i.cmdNotSupported
}

func (i *ImageBuilder) UpdateCmdsNotFound(cmds []string) {
	i.cmdNotOnApk = cmds
}
