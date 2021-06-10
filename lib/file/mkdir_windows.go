//+build windows

package file

import (
	"os"
	"path/filepath"
	"syscall"
)

// MkdirAll creates a directory named path, along with any necessary parents.
// Except added os prefix on call to IsPathSeparator and additional comments,
// it is a pure copy of Syncthing's version of os.MkdirAll, which adds check on
// volume name to avoid trying to create the folder \\? for extended length paths.
// For comparsion, see:
// https://github.com/syncthing/syncthing/blob/main/lib/fs/basicfs_windows.go
// https://github.com/golang/go/blob/master/src/os/path.go
func MkdirAll(path string, perm os.FileMode) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{
			Op:   "mkdir",
			Path: path,
			Err:  syscall.ENOTDIR,
		}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		parent := path[0 : j-1]
		if parent != filepath.VolumeName(parent) { // Here is change from os.MkdirAll (avoiding issue trying to create \\? as a directory)
			err = MkdirAll(parent, perm) // Here is another change from os.MkdirAll (not calling fixRootDirectory, it is not necessary)
			if err != nil {
				return err
			}
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = os.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}
