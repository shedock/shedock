package shell

import (
	"shedock/pkg/parsers/apk"
	"strings"
)

type Zsh struct {
	infoCommand                  string
	installShellCommand          string
	commandToFindBuiltinCommands string
}

func NewZsh() *Zsh {
	return &Zsh{
		infoCommand:                  "apk info -a zsh",
		installShellCommand:          "apk add zsh",
		commandToFindBuiltinCommands: "zsh -c \"whence -wm '*'\"",
	}
}

func (z *Zsh) Dependencies(data []byte) ([]*apk.PackageDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.DependsOn()
}

func (z *Zsh) Provides(data []byte) ([]*apk.ProviderDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Provides()
}

func (z *Zsh) InfoCommand() string {
	return z.infoCommand
}

func (z *Zsh) InstallShellCommand() string {
	return z.installShellCommand
}

func (z *Zsh) CommandToFindBuiltinCommands() string {
	return z.commandToFindBuiltinCommands
}

func (z *Zsh) BinaryDependencies(data []byte) ([]string, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Contains()
}

func (z *Zsh) ParseBuiltins(data string) ([]string, error) {
	lines := strings.Split(data, "\n")
	var builtins []string

	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			command := parts[0]
			builtins = append(builtins, command)
		}
	}

	return builtins, nil
}
