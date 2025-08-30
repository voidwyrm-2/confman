package confman

import "os"

// Mkdir creates a new subdirectory then returns the absolute path to that subdirectory.
//
// [os.ModeDir] is masked onto perm internally, and does not need to be specified.
func (c *Config) Mkdir(name string, perm os.FileMode) (Path, error) {
	path := c.child(name)
	return Path(path), os.Mkdir(path, perm|os.ModeDir)
}

// MkdirAll creates a subdirectory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// subdirectories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
//
// [os.ModeDir] is masked onto perm internally, and does not need to be specified.
func (c *Config) MkdirAll(name string, perm os.FileMode) (Path, error) {
	path := c.child(name)
	return Path(path), os.MkdirAll(path, perm|os.ModeDir)
}
