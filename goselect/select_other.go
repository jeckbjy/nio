// +build !linux,!windows,!plan9,!solaris

package goselect

import "syscall"

func sysSelect(n int, r, w, e *FDSet, timeout *syscall.Timeval) (int, error) {
	// The Go syscall.Select for the BSD unixes is buggy. It
	// returns only the error and not the number of active file
	// descriptors. To cope with this, we return "nfd" as the
	// number of active file-descriptors. This can cause
	// significant performance degradation but there's nothing
	// else we can do.
	err := syscall.Select(n, (*syscall.FdSet)(r), (*syscall.FdSet)(w), (*syscall.FdSet)(e), timeout)
	if err != nil {
		return 0, err
	}
	return n, nil
}
