package goselect

import (
	"syscall"
	"time"
)

// Select wraps syscall.Select with Go types
func Select(n int, r, w, e *FDSet, timeout time.Duration) (int, error) {
	var timeval *syscall.Timeval
	if timeout >= 0 {
		t := syscall.NsecToTimeval(timeout.Nanoseconds())
		timeval = &t
	}

	return sysSelect(n, r, w, e, timeval)
}
