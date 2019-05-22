// +build linux,noepoll !darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package nio

import (
	"github.com/jeckbjy/nio/goselect"
	"os"
	"time"
)

func newPoller() poller {
	return &pselect{}
}

type pselect struct {
	fdmax uintptr
	fds   []uintptr
	rset  goselect.FDSet
	wset  goselect.FDSet
	eset  goselect.FDSet
	pr    *os.File
	pw    *os.File
}

func (p *pselect) Open() error {
	p.fdmax = -1
	p.rset.Zero()
	p.wset.Zero()
	p.eset.Zero()
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	p.pr = r
	p.pw = w
	p.rset.Set(r.Fd())
	p.fdmax = r.Fd()
	p.fds = append(p.fds, r.Fd())

	return nil
}

func (p *pselect) Close() error {
	e1 := p.pr.Close()
	e2 := p.pw.Close()
	if e1 != nil {
		return e1
	}
	return e2
}

func (p *pselect) Wakeup() error {
	_, err := p.pw.Write([]byte("0"))
	return err
}

func (p *pselect) Wait(s *Selector, cb SelectCB, msec int) error {
	for {
		_, err := goselect.Select(int(p.fdmax+1), &p.rset, &p.wset, &p.eset, time.Millisecond*time.Duration(msec))
		if err != nil {
			if isTemporaryError(err) {
				continue
			}

			return err
		}

		if p.rset.IsSet(p.pr.Fd()) {
			// drain all
			bytes := [64]byte{}
			for {
				_, err := p.pr.Read(bytes[:64])
				if err != nil {
					if isTemporaryError(err) {
						continue
					}
					break
				}
			}
		}

		for fd, sk := range s.keys {
			sk.reset()

			if p.eset.IsSet(fd) {
				// error
				sk.setReadyIn()
				sk.setReadyOut()
				continue
			}

			if p.rset.IsSet(fd) {
				sk.setReadyIn()
			}

			if p.wset.IsSet(fd) {
				sk.setReadyOut()
			}
		}

		return nil
	}
}

func (p *pselect) Add(fd uintptr, ops int) error {
	if ops&op_IN != 0 {
		p.rset.Set(fd)
	}

	if ops&op_OUT != 0 {
		p.wset.Set(fd)
	}

	if fd >= p.fdmax {
		p.fdmax = fd
	}

	p.fds = append(p.fds, fd)

	return nil
}

func (p *pselect) Delete(fd uintptr, ops int) error {
	if ops&op_IN != 0 {
		p.rset.Clear(fd)
	}

	if ops&op_OUT != 0 {
		p.wset.Clear(fd)
	}

	// get max fd
	fdmax := uintptr(0)
	index := -1
	for i, f := range p.fds {
		if f == fd {
			index = i
		} else if f > fdmax {
			fdmax = f
		}
	}

	p.fdmax = fdmax
	if index != -1 {
		p.fds = append(p.fds[:index], p.fds[index+1:]...)
	}

	return nil
}

func (p *pselect) Modify(fd uintptr, old, ops int) error {
	oldi := old&op_IN != 0
	oldo := old&op_OUT != 0
	newi := ops&op_IN != 0
	newo := ops&op_OUT != 0

	if oldi != newi {
		if newi {
			p.rset.Set(fd)
		} else {
			p.rset.Clear(fd)
		}
	}

	if oldo != newo {
		if newo {
			p.wset.Set(fd)
		} else {
			p.wset.Clear(fd)
		}
	}

	return nil
}
