package shell

import (
	"errors"
	"shedock/pkg/parsers/apk"
)

type Shell interface {
	CommandToFindBuiltinCommands() string
	Dependencies(data []byte) ([]*apk.PackageDependency, error)
	Provides(data []byte) ([]*apk.ProviderDependency, error)
	InfoCommand() string
	InstallShellCommand() string
	BinaryDependencies(data []byte) ([]string, error)
	ParseBuiltins(data string) ([]string, error)
}

// Factory function
func NewShell(shellType string) (Shell, error) {
	switch shellType {
	case ZshShell:
		return NewZsh(), nil
	case BashShell:
		return NewBash(), nil
	case FishShell:
		return NewFish(), nil
	case PwshShell:
		return NewPwsh(), nil
	default:
		return nil, errors.New("invalid shell type or not implemented yet")
	}
}
