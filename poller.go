package nio

const maxEventNum = 1024

type poller interface {
	Open() error
	Close() error

	Wakeup() error
	Wait(s *Selector, cb SelectCB, msec int) error

	Add(fd uintptr, ops int) error
	Modify(fd uintptr, ops int) error
	Delete(fd uintptr, ops int) error
}
