package shell

import (
	"shedock/pkg/parsers/apk"
	"strings"
)

type Bash struct {
	infoCommand                  string
	installShellCommand          string
	commandToFindBuiltinCommands string
}

func NewBash() *Bash {
	return &Bash{
		infoCommand:                  "apk info -a bash",
		installShellCommand:          "apk add bash",
		commandToFindBuiltinCommands: "bash -c \"compgen -b\"",
	}
}

func (b *Bash) Dependencies(data []byte) ([]*apk.PackageDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.DependsOn()
}

func (b *Bash) Provides(data []byte) ([]*apk.ProviderDependency, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Provides()
}

func (b *Bash) BinaryDependencies(data []byte) ([]string, error) {
	apkParser := apk.ApkParser{Data: data}
	return apkParser.Contains()
}

func (b *Bash) InfoCommand() string {
	return b.infoCommand
}

func (b *Bash) InstallShellCommand() string {
	return b.installShellCommand
}

func (b *Bash) CommandToFindBuiltinCommands() string {
	return b.commandToFindBuiltinCommands
}

func (b *Bash) ParseBuiltins(data string) ([]string, error) {
	var builtins []string

	// scanner := bufio.NewScanner(strings.NewReader(data))
	// if scanner.Scan() {
	// 	line := scanner.Text()
	// 	builtins = append(builtins, line)
	// }

	// if err := scanner.Err(); err != nil {
	// 	return nil, err
	// }

	builtins = strings.Split(data, "\n")

	return builtins, nil
}
