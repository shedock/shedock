package shell

import "shedock/pkg/parsers/apk"

type Pwsh struct {
	infoCommand                  string
	installShellCommand          string
	commandToFindBuiltinCommands string
}

func NewPwsh() *Pwsh {
	return &Pwsh{
		infoCommand:                  "apk info -a powershell",
		installShellCommand:          "apk add powershell",
		commandToFindBuiltinCommands: "Get-Command",
	}
}

func (p *Pwsh) Dependencies(data []byte) ([]*apk.PackageDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.DependsOn()
}

func (p *Pwsh) Provides(data []byte) ([]*apk.ProviderDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Provides()
}

func (p *Pwsh) BinaryDependencies(data []byte) ([]string, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Contains()
}

func (p *Pwsh) InfoCommand() string {
	return p.infoCommand
}

func (p *Pwsh) InstallShellCommand() string {
	return p.installShellCommand
}

func (p *Pwsh) CommandToFindBuiltinCommands() string {
	return p.commandToFindBuiltinCommands
}

func (p *Pwsh) ParseBuiltins(data string) ([]string, error) {
	return []string{}, nil
}
