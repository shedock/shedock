package file

type Dependency struct {
	// FromPath is the path of the dependency in the first layer
	FromPath string
	// ToPath is the path of the dependency in the second layer
	ToPath string
	// Requiredby represents a string containing information about who requires this dependency.
	// E.g
	// # Required by user script: <script_name>
	// # Required by binaries: bash, find, grep, sed, xargs etc.
	Requiredby []string
}

type Dependencies struct {
	// Bin represents the binary command dependencies of the user script
	Bin []Dependency
	// Lib represents the dependencies shared libraries
	Lib []Dependency
}
