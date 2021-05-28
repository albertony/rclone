//+build windows

package file

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rclone/rclone/lib/encoder"
)

// Pattern to match a windows absolute path: "c:\" and similar
var isAbsWinDrive = regexp.MustCompile(`^[a-zA-Z]\:\\`)

// UNCPath converts an absolute Windows path to a UNC long path.
//
// It does nothing on non windows platforms
func UNCPath(l string) string {
	// If prefix is "\\", we already have a UNC path or server.
	if strings.HasPrefix(l, `\\`) {
		// If already long path, just keep it
		if strings.HasPrefix(l, `\\?\`) {
			return l
		}

		// Trim "\\" from path and add UNC prefix.
		return `\\?\UNC\` + strings.TrimPrefix(l, `\\`)
	}
	if isAbsWinDrive.MatchString(l) {
		return `\\?\` + l
	}
	return l
}

// NonUNCPath removes UNC path prefix if present.
//
// It does nothing on non windows platforms
func NonUNCPath(l string) string {
	// If prefix is "\\", we already have a UNC path or server.
	if strings.HasPrefix(l, `\\?\UNC\`) {
		return l[8:]
	}
	if strings.HasPrefix(l, `\\?\`) {
		return l[4:]
	}
	return l
}

// LocalPath converts a standard path into local path.
//
// The returned path is using always absolute, using os path separator,
// and encoded according to a supplied encoder - typically encoder.OS.
// By default also returned as UNC long path.
func LocalPath(s string, noUNC bool, enc encoder.MultiEncoder) string {
	// Ensure absolute and clean path (avoid multiple separators etc)
	if !filepath.IsAbs(s) {
		s2, err := filepath.Abs(s)
		if err == nil {
			s = s2
		}
	} else {
		s = filepath.Clean(s)
	}
	// Encode the absolute path, but must avoid encoding the ':' after drive letters.
	s = filepath.ToSlash(s)
	vol := filepath.VolumeName(s)
	s = vol + enc.FromStandardPath(s[len(vol):])
	s = filepath.FromSlash(s)
	// Optionally add UNC prefix
	if !noUNC {
		s = UNCPath(s)
	}
	return s
}
