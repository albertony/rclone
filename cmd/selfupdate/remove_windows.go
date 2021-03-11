// +build windows

package selfupdate

import (
	"golang.org/x/sys/windows"
)

func removeOnReboot(path string) error {
	pathp, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(pathp, nil, windows.MOVEFILE_DELAY_UNTIL_REBOOT)
}
