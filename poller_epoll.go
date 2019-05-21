// +build linux

package nio

import (
	"fmt"
	"syscall"
)

func newPoller() poller {
	return &epoll{}
}

// https://medium.com/@copyconstruct/the-method-to-epolls-madness-d9d2d6378642
type epoll struct {
	efd    int // epoll fd
	wfd    int // wakeup fd
	events []syscall.EpollEvent
}

func epoll_create() (int, error) {
	fd, err := syscall.EpollCreate1(0)
	if err == nil {
		return fd, nil
	}

	return syscall.EpollCreate(1024)
}

func (p *epoll) Open() error {
	fd, err := epoll_create()
	if err != nil {
		return err
	}

	r0, _, e0 := syscall.Syscall(syscall.SYS_EVENTFD2, 0, 0, 0)
	if e0 != 0 {
		syscall.Close(fd)
		return fmt.Errorf("create eventfd fail")
	}

	syscall.CloseOnExec(fd)

	p.events = make([]syscall.EpollEvent, maxEventNum)
	p.efd = fd
	p.wfd = int(r0)
	return nil
}

func (p *epoll) Close() error {
	if err := syscall.Close(p.wfd); err != nil {
		return err
	}

	return syscall.Close(p.efd)
}

func (p *epoll) Wakeup() error {
	return nil
}

func (p *epoll) Wait(s *Selector, cb SelectCB, msec int) error {
	for {
		n, err := syscall.EpollWait(p.efd, p.events, msec)
		if err != nil {
			if isTemporaryError(err) {
				continue
			}

			return err
		}

		for i := 0; i < n; i++ {
			ev := &p.events[i]
			fd := ev.Fd
			sk := s.getSelectionKey(uintptr(fd))
			if sk == nil {
				// close socket?
				continue
			}

			sk.reset()

			// check error?

			if ev.Events&(syscall.EPOLLIN|syscall.EPOLLERR|syscall.EPOLLHUP) != 0 {
				if sk.isInterests(OP_ACCEPT) {
					sk.ready |= OP_ACCEPT
				} else if sk.isInterests(OP_READ) {
					sk.ready |= OP_READ
				}
			}

			if ev.Events&(syscall.EPOLLOUT|syscall.EPOLLERR|syscall.EPOLLHUP) != 0 {
				if sk.isInterests(OP_WRITE) {
					sk.ready |= OP_WRITE
				} else if sk.isInterests(OP_CONNECT) {
					sk.ready |= OP_CONNECT
				}
			}

			if cb != nil {
				cb(sk)
			} else {
				s.readyKeys = append(s.readyKeys, sk)
			}
		}

		break
	}

	return nil
}

func (p *epoll) Add(fd uintptr, ops int) error {
	ev := &syscall.EpollEvent{Events: toEpollEvents(ops), Fd: int32(fd)}
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_ADD, int(fd), ev)
}

func (p *epoll) Modify(fd uintptr, ops int) error {
	ev := &syscall.EpollEvent{Events: toEpollEvents(ops), Fd: int32(fd)}
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_MOD, int(fd), ev)
}

func (p *epoll) Delete(fd uintptr, ops int) error {
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_DEL, int(fd), nil)
}

func toEpollEvents(ops int) uint32 {
	events := syscall.EPOLLET | syscall.EPOLLPRI

	if ops&(OP_ACCEPT|OP_READ) != 0 {
		events |= syscall.EPOLLIN
	}

	if ops&(OP_WRITE|OP_CONNECT) != 0 {
		events |= syscall.EPOLLOUT
	}

	return uint32(events)
}
