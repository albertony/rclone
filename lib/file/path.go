//+build !windows

package file

import (
	"path/filepath"

	"github.com/rclone/rclone/lib/encoder"
)

// UNCPath converts an absolute Windows path to a UNC long path.
//
// It does nothing on non windows platforms
func UNCPath(l string) string {
	return l
}

// NonUNCPath removes UNC path prefix if present.
//
// It does nothing on non windows platforms
func NonUNCPath(l string) string {
	return l
}

// LocalPath converts a standard path into local path.
//
// The returned path is always absolute and encoded according to
// supplied encoder, typically encoder.OS.
func LocalPath(s string, noUNC bool, enc encoder.MultiEncoder) string {
	if !filepath.IsAbs(s) {
		s2, err := filepath.Abs(s)
		if err == nil {
			s = s2
		}
	}
	s = enc.FromStandardPath(s)
	return s
}
