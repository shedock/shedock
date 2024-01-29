package core

import (
	"log"
	"shedock/internal/instance"
	apkTypes "shedock/pkg/parsers/apk"
	shellScriptTypes "shedock/pkg/parsers/shellscript"
	"shedock/pkg/shell"
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
		for _, builtin := range shellbuiltns {
			if dep.Name == builtin {
				found = true
				break
			}
		}
		if !found {
			filteredDeps = append(filteredDeps, dep.Name)
		}
	}
	return nil
}

func (i *ImageBuilder) getSysteDependencies() ([]string, error) {
	return []string{}, nil
}

func (i *ImageBuilder) getShellBuiltins() ([]string, error) {
	container := instance.GetDockerInstance()
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
