// +build darwin dragonfly freebsd netbsd openbsd

package nio

import (
	"syscall"
)

func newPoller() poller {
	return &kqueue{}
}

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
	return nil
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
			fd := ev.Ident
			sk := s.getSelectionKey(uintptr(fd))
			if sk == nil {
				continue
			}

			sk.reset()

			if ev.Filter == syscall.EVFILT_READ || ev.Flags&(syscall.EV_ERROR|syscall.EV_EOF) != 0 {
				if sk.isInterests(OP_ACCEPT) {
					sk.ready |= OP_ACCEPT
				} else if sk.isInterests(OP_READ) {
					sk.ready |= OP_READ
				}
			}

			if ev.Filter == syscall.EVFILT_WRITE || ev.Flags&(syscall.EV_ERROR|syscall.EV_EOF) != 0 {
				if sk.isInterests(OP_WRITE) {
					sk.ready |= OP_WRITE
				} else if sk.isInterests(OP_CONNECT) {
					sk.ready |= OP_WRITE
				}
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
	return p.control(fd, true, ops)
}

func (p *kqueue) Modify(fd uintptr, ops int) error {
	return p.control(fd, true, ops)
}

func (p *kqueue) Delete(fd uintptr, ops int) error {
	return p.control(fd, false, ops)
}

func (p *kqueue) control(fd uintptr, add bool, events int) error {
	changes := make([]syscall.Kevent_t, 0, 2)
	var flags uint16
	if add {
		// EV_CLEAR:Edge Triggered
		flags = syscall.EV_ADD | syscall.EV_CLEAR
	} else {
		flags = syscall.EV_DELETE
	}

	if events&(OP_ACCEPT|OP_READ) != 0 {
		changes = append(changes, syscall.Kevent_t{Ident: uint64(fd), Flags: flags, Filter: syscall.EVFILT_READ})
	}

	if events&(OP_WRITE|OP_CONNECT) != 0 {
		changes = append(changes, syscall.Kevent_t{Ident: uint64(fd), Flags: flags, Filter: syscall.EVFILT_WRITE})
	}

	_, err := syscall.Kevent(p.kfd, changes, nil, nil)
	return err
}

//func toKevent(events int, add bool) []syscall.Kevent_t {
//	//kevents := make([]syscall.Kevent_t, 0, 2)
//	//var flags uint16
//	//if add {
//	//	flags = syscall.EV_ADD
//	//	if (events & EventOneShot) != 0 {
//	//		flags |= syscall.EV_ONESHOT
//	//	}
//	//
//	//	if (events & EventEdgeTriggered) != 0 {
//	//		flags |= syscall.EV_CLEAR
//	//	}
//	//} else {
//	//	flags = syscall.EV_DELETE
//	//}
//	//
//	//if (events & EventRead) != 0 {
//	//	kevents = append(kevents, syscall.Kevent_t{Flags: flags, Filter: syscall.EVFILT_READ})
//	//}
//	//
//	//if (events & EventWrite) != 0 {
//	//	kevents = append(kevents, syscall.Kevent_t{Flags: flags, Filter: syscall.EVFILT_WRITE})
//	//}
//	//
//	//return kevents
//}
//
//func toEvent(ev *syscall.Kevent_t) int {
//	events := 0
//	flags := ev.Flags
//	filter := ev.Filter
//
//	if (flags & syscall.EV_ERROR) != 0 {
//		events |= EventErr
//	}
//
//	// Set EventHup for any EOF flag. Below will be more precise detection
//	// of what exatcly HUP occured.
//	if (flags & syscall.EV_EOF) != 0 {
//		events |= EventHup
//	}
//
//	if filter == syscall.EVFILT_READ {
//		events |= EventRead
//		if (flags & syscall.EV_EOF) != 0 {
//			events |= EventReadHup
//		}
//	}
//
//	if filter == syscall.EVFILT_WRITE {
//		events |= EventWrite
//		if (flags & syscall.EV_EOF) != 0 {
//			events |= EventWriteHup
//		}
//	}
//
//	return events
//}
