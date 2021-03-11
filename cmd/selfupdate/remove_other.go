// +build !windows

package selfupdate

import "github.com/pkg/errors"

func removeOnReboot(path string) error {
	errors.New("remove on reboot not supported on this platform")
}
