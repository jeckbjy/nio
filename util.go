package nio

import (
	"fmt"
	"os"
	"syscall"
)

// ifile describes an object that has ability to return os.File.
type ifile interface {
	// File returns a copy of object's file descriptor.
	File() (*os.File, error)
}

func getFd(conn interface{}) (uintptr, error) {
	if i, ok := conn.(ifile); ok {
		f, err := i.File()
		if err != nil {
			return 0, err
		}

		return f.Fd(), nil
	}

	return 0, fmt.Errorf("bad file descriptor:%+v", conn)
}

func isTemporaryError(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}

	return errno.Temporary()
}
