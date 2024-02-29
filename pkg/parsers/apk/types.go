package apk

type PackageDependency struct {
	Name string
	Type DependencyType
}

type ProviderDependency struct {
	// Name of the command that the dependency provides
	Name string
	// Version of the command that the dependency provides
	Version string
}

type Package struct {
	// Name of the binary package
	Name string
	// Path of the package once installed
	Path string
}

type DependencyType string

const (
	Binary        DependencyType = "binary"
	SharedLibrary DependencyType = "shared library"
)

func (d DependencyType) IsBinary() bool {
	return d == Binary
}

func (d DependencyType) IsSharedLibrary() bool {
	return d == SharedLibrary
}
