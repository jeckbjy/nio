// +build !windows

package internal

import "syscall"

func SetNonblock(fd uintptr, nonblocking bool) error {
	return syscall.SetNonblock(int(fd), nonblocking)
}

func CloseOnExec(fd uintptr) {
	syscall.CloseOnExec(int(fd))
}
