// +build plan9 solaris

package goselect

import (
	"fmt"
	"runtime"
	"syscall"
)

// ErrUnsupported .
var ErrUnsupported = fmt.Errorf("Platofrm %s/%s unsupported", runtime.GOOS, runtime.GOARCH)

func sysSelect(n int, r, w, e *FDSet, timeout *syscall.Timeval) (int, error) {
	return 0, ErrUnsupported
}
