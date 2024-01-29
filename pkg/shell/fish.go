package shell

import "shedock/pkg/parsers/apk"

type Fish struct {
	infoCommand                  string
	installShellCommand          string
	commandToFindBuiltinCommands string
}

func NewFish() *Fish {
	return &Fish{
		infoCommand:                  "apk info -a fish",
		installShellCommand:          "apk add fish",
		commandToFindBuiltinCommands: "builtin --names",
	}
}

func (f *Fish) Dependencies(data []byte) ([]*apk.PackageDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.DependsOn()
}

func (f *Fish) Provides(data []byte) ([]*apk.ProviderDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Provides()
}

func (f *Fish) BinaryDependencies(data []byte) ([]string, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Contains()
}

func (f *Fish) InfoCommand() string {
	return f.infoCommand
}

func (f *Fish) InstallShellCommand() string {
	return f.installShellCommand
}

func (f *Fish) CommandToFindBuiltinCommands() string {
	return f.commandToFindBuiltinCommands
}

func (f *Fish) ParseBuiltins(data string) ([]string, error) {
	return []string{}, nil
}
