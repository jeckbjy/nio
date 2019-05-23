// +build darwin dragonfly freebsd netbsd openbsd

package nio

import (
	"log"
	"syscall"
)

func newPoller() poller {
	return &kqueue{}
}

// http://eradman.com/posts/kqueue-tcp.html
type kqueue struct {
	kfd    int
	events []syscall.Kevent_t
}

func (p *kqueue) Open() error {
	fd, err := syscall.Kqueue()
	if err != nil {
		return err
	}

	changes := []syscall.Kevent_t{{Ident: 0, Filter: syscall.EVFILT_USER, Flags: syscall.EV_ADD | syscall.EV_CLEAR}}
	_, err = syscall.Kevent(fd, changes, nil, nil)
	if err != nil {
		syscall.Close(fd)
		return err
	}

	syscall.CloseOnExec(fd)
	p.kfd = fd
	p.events = make([]syscall.Kevent_t, maxEventNum)
	return nil
}

func (p *kqueue) Close() error {
	return syscall.Close(p.kfd)
}

func (p *kqueue) Wakeup() error {
	changes := []syscall.Kevent_t{{Ident: 0, Filter: syscall.EVFILT_USER, Fflags: syscall.NOTE_TRIGGER}}
	_, err := syscall.Kevent(p.kfd, changes, nil, nil)
	return err
}

func (p *kqueue) Wait(s *Selector, cb SelectCB, msec int) error {
	for {
		n, err := syscall.Kevent(p.kfd, nil, p.events, nil)
		if err != nil {
			// check temporary
			if isTemporaryError(err) {
				continue
			}

			return err
		}

		for i := 0; i < n; i++ {
			ev := &p.events[i]
			fd := uintptr(ev.Ident)
			sk := s.keys[fd]
			if sk == nil {
				continue
			}

			sk.reset()

			if ev.Flags&(syscall.EV_ERROR|syscall.EV_EOF) != 0 {
				log.Printf("wait err:%+v\n", syscall.Errno(ev.Data).Error())
				sk.setReadyIn()
				sk.setReadyOut()
				continue
			}

			if ev.Filter == syscall.EVFILT_READ {
				sk.setReadyIn()
			}

			if ev.Filter == syscall.EVFILT_WRITE {
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
}

func (p *kqueue) Add(fd uintptr, ops int) error {
	changes := [4]syscall.Kevent_t{}
	num := p.control(&changes, 0, fd, ops, true)
	_, err := syscall.Kevent(p.kfd, changes[:num], nil, nil)
	return err
}

func (p *kqueue) Delete(fd uintptr, ops int) error {
	changes := [4]syscall.Kevent_t{}
	num := p.control(&changes, 0, fd, ops, false)
	_, err := syscall.Kevent(p.kfd, changes[:num], nil, nil)
	return err
}

func (p *kqueue) Modify(fd uintptr, old, ops int) error {
	changes := [4]syscall.Kevent_t{}
	num := 0

	if old != 0 {
		// delete old
		num = p.control(&changes, 0, fd, old, false)
	}

	num = p.control(&changes, num, fd, ops, true)
	_, err := syscall.Kevent(p.kfd, changes[:num], nil, nil)
	return err
}

func (p *kqueue) control(changes *[4]syscall.Kevent_t, num int, fd uintptr, ops int, add bool) int {
	ident := uint64(fd)
	var flags uint16
	if add {
		flags = syscall.EV_ADD | syscall.EV_CLEAR
	} else {
		flags = syscall.EV_DELETE
	}

	if ops&op_IN != 0 {
		changes[num] = syscall.Kevent_t{Filter: syscall.EVFILT_READ, Ident: ident, Flags: flags}
		num++
	}

	if ops&op_OUT != 0 {
		changes[num] = syscall.Kevent_t{Filter: syscall.EVFILT_WRITE, Ident: ident, Flags: flags}
		num++
	}

	return num
}
