package confman

import (
	"path/filepath"
	"unsafe"
)

// A Path represents a joinable file path.
type Path string

// Joins a list of strings to this Path, separating them with the platform-specific separator.
func (p Path) Join(elem ...string) Path {
	parts := make([]string, 0, len(elem)+1)
	parts = append(parts, string(p))
	parts = append(parts, elem...)

	return Path(filepath.Join(parts...))
}

// Joins a list of Paths to this Path, separating them with the platform-specific separator.
func (p Path) JoinP(elem ...Path) Path {
	elp := unsafe.SliceData(elem)
	sp := (*string)(unsafe.Pointer(elp))
	return p.Join(unsafe.Slice(sp, len(elem))...)
}
