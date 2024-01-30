package core

import (
	"fmt"
	"log"
	"shedock/internal/instance"
	apkTypes "shedock/pkg/parsers/apk"
	shellScriptTypes "shedock/pkg/parsers/shellscript"
	"shedock/pkg/shell"
	"strings"
)

type ImageBuilder struct {
	Script *shellScriptTypes.Script
	Shell  shell.Shell
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

	systemBuiltins, err := i.getSystemBuiltins()
	if err != nil {
		return err
	}

	shellType, err := i.Script.GetShell()
	if err != nil {
		return err
	}
	_, err = shell.NewShell(shellType)
	if err != nil {
		log.Fatalf("Failed to get shell: %v", err)
	}

	scriptDeps, err := i.Script.Dependencies()
	if err != nil {
		return err
	}

	shellbuiltns, err := i.getShellBuiltins()
	if err != nil {
		return err
	}

	for _, dep := range scriptDeps {
		var found bool

		// Check if the dependency is a shell builtin
		for _, builtin := range shellbuiltns {
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

		if !found {
			filteredDeps = append(filteredDeps, dep.Name)
		}
	}

	log.Println(filteredDeps)
	cmdOnApk, err := i.filterDependenciesFromPackageHost(filteredDeps)
	if err != nil {
		return err
	}
	log.Println(cmdOnApk)
	return nil
}

// System dependencies are the binaries that are installed in the base image
// E.g. ls, cat, etc.
func (i *ImageBuilder) getSystemBuiltins() ([]string, error) {
	container := instance.GetDockerInstance()
	output, err := container.ExecCommand("for binary in /bin/* /usr/bin/*; do echo \"$binary\"; done")
	if err != nil {
		log.Fatalf("%v", err)
	}
	var builtins []string
	builtins = append(builtins, strings.Split(output, "\n")...)
	// Parse the output

	return builtins, nil
}

// Shell builtins are the binaries that are installed by the shell
// E.g. [[, readarray, zmodload, etc.
func (i *ImageBuilder) getShellBuiltins() ([]string, error) {
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
	return builtins, nil
}

// Shell libraries are the libraries that are required by the shell. Also called as "shared libraries"
// E.g. libreadline, libncurses, etc.
func (i *ImageBuilder) getShellLibraries() ([]*apkTypes.PackageDependency, error) {
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

func (i *ImageBuilder) getShellBinaryBuiltins() ([]string, error) {
	container := instance.GetDockerInstance()
	_, err := container.ExecCommand(i.Shell.InstallShellCommand())
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Execute the info command again to get the stuff installed by the shell
	output, err := container.ExecCommand(i.Shell.InfoCommand())
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Parse the output
	binDeps, err := i.Shell.BinaryDependencies([]byte(output))
	if err != nil {
		log.Fatalf("Failed to parse dependencies: %v", err)
	}

	return binDeps, nil
}

// Find dependencies that can be installed by the package manager (for now apk)
func (i *ImageBuilder) filterDependenciesFromPackageHost(deps []string) ([]string, error) {
	container := instance.GetDockerInstance()
	var depsOnApk []string
	var err error

	for _, dep := range deps {
		// Execute the info command again to get the stuff installed by the shell
		output, err := container.ExecCommand(fmt.Sprintf("apk info -a %s", dep))
		if err != nil {
			return []string{}, err
		}

		// Parse the output
		parser := apkTypes.ApkParser{Data: []byte(output)}
		dependencies, err := parser.Provides()
		if err != nil {
			return []string{}, err
		}

		for _, dependency := range dependencies {
			depsOnApk = append(depsOnApk, dependency.Name)
		}
	}

	return depsOnApk, err
}
