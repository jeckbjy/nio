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
	_, err := syscall.Write(p.wfd, []byte{0, 0, 0, 0, 0, 0, 0, 1})
	return err
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
			fd := int(ev.Fd)

			if fd == p.wfd {
				continue
			}

			sk := s.getSelectionKey(uintptr(fd))
			if sk == nil {
				// close socket?
				continue
			}

			sk.reset()

			if ev.Events&(syscall.EPOLLIN|syscall.EPOLLERR|syscall.EPOLLHUP) != 0 {
				sk.setReadyIn()
			}

			if ev.Events&(syscall.EPOLLOUT|syscall.EPOLLERR|syscall.EPOLLHUP) != 0 {
				sk.setReadyOut()
			}

			if cb != nil {
				cb(sk)
			} else {
				s.readyKeys = append(s.readyKeys, sk)
			}
		}

		return nil
	}

	return nil
}

func (p *epoll) Add(fd uintptr, ops int) error {
	ev := &syscall.EpollEvent{Events: toEpollEvents(ops), Fd: int32(fd)}
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_ADD, int(fd), ev)
}

func (p *epoll) Delete(fd uintptr, ops int) error {
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_DEL, int(fd), nil)
}

func (p *epoll) Modify(fd uintptr, old, ops int) error {
	ev := &syscall.EpollEvent{Events: toEpollEvents(ops), Fd: int32(fd)}
	return syscall.EpollCtl(p.efd, syscall.EPOLL_CTL_MOD, int(fd), ev)
}

func toEpollEvents(ops int) uint32 {
	events := syscall.EPOLLET | syscall.EPOLLPRI

	if ops&op_IN != 0 {
		events |= syscall.EPOLLIN
	}

	if ops&op_OUT != 0 {
		events |= syscall.EPOLLOUT
	}

	return uint32(events)
}
