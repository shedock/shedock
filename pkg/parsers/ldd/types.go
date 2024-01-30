package ldd

type Library struct {
	Name      string
	IsSymlink bool
}

type Dependencies struct {
	Binary       string
	Dependencies []Library
}
