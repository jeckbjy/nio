// +build windows

package internal

import (
	"syscall"
)

func SetNonblock(fd uintptr, nonblocking bool) error {
	return syscall.SetNonblock(syscall.Handle(fd), nonblocking)
}

func CloseOnExec(fd uintptr) {
	syscall.CloseOnExec(syscall.Handle(fd))
}
