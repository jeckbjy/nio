package nio

import (
	"fmt"
	"net"
	"syscall"
)

var (
	// ErrRegistered is returned by Poller Start() method to indicate that
	// connection with the same underlying file descriptor was already
	// registered within the poller instance.
	ErrRegistered = fmt.Errorf("file descriptor is already registered in poller instance")

	// ErrNotRegistered is returned by Poller Stop() and Resume() methods to
	// indicate that connection with the same underlying file descriptor was
	// not registered before within the poller instance.
	ErrNotRegistered = fmt.Errorf("file descriptor was not registered before in poller instance")
)

const (
	OP_ACCEPT  = 0x01
	OP_CONNECT = 0x02 // not support
	OP_READ    = 0x04
	OP_WRITE   = 0x08
	op_CLOSED  = 0x10 // socket closed
)

type SelectionKey struct {
	channel   interface{} // net.Conn or net.Listener
	data      interface{} // attachment
	fd        uintptr     // file descriptor
	interests int         // registered ops
	ready     int         // ready ops
}

func (sk *SelectionKey) reset() {
	sk.ready = 0
}

func (sk *SelectionKey) Closed() bool {
	return sk.ready&op_CLOSED != 0
}

func (sk *SelectionKey) Acceptable() bool {
	return sk.ready&OP_ACCEPT != 0
}

func (sk *SelectionKey) Connectable() bool {
	return sk.ready&OP_CONNECT != 0
}

func (sk *SelectionKey) Readable() bool {
	return sk.ready&OP_READ != 0
}

func (sk *SelectionKey) Writable() bool {
	return sk.ready&OP_WRITE != 0
}

func (sk *SelectionKey) isInterests(ops int) bool {
	return sk.interests&ops != 0
}

func (sk *SelectionKey) InterestOps() int {
	return sk.interests
}

func (sk *SelectionKey) ReadyOps() int {
	return sk.ready
}

func (sk *SelectionKey) Listener() net.Listener {
	return sk.channel.(net.Listener)
}

func (sk *SelectionKey) Conn() net.Conn {
	return sk.channel.(net.Conn)
}

func (sk *SelectionKey) Channel() interface{} {
	return sk.channel
}

func (sk *SelectionKey) Data() interface{} {
	return sk.data
}

type SelectCB func(sk *SelectionKey)

type SelectOptions struct {
	Timeout  int      // 毫秒,-1表示永久
	Callback SelectCB // 设置回调,将会直接调用回调函数,而不会返回[]*SelectionKey数组
}

func New() (*Selector, error) {
	poll := newPoller()
	if poll == nil {
		return nil, fmt.Errorf("create poller fail")
	}

	if err := poll.Open(); err != nil {
		return nil, err
	}

	s := &Selector{keys: make(map[uintptr]*SelectionKey), poll: poll}
	return s, nil
}

// https://www.cnblogs.com/pingh/p/3224990.html
type Selector struct {
	poll      poller
	keys      map[uintptr]*SelectionKey
	readyKeys []*SelectionKey
}

func (s *Selector) Add(conn interface{}, ops int, data interface{}) (*SelectionKey, error) {
	fd, err := getFd(conn)
	if err != nil {
		return nil, err
	}

	if s.keys[fd] != nil {
		return nil, ErrRegistered
	}

	if err := syscall.SetNonblock(int(fd), true); err != nil {
		return nil, err
	}

	if err := s.poll.Add(fd, ops); err != nil {
		return nil, err
	}

	sk := &SelectionKey{channel: conn, data: data, fd: fd, interests: ops}
	s.keys[fd] = sk
	return sk, nil
}

func (s *Selector) Delete(conn interface{}) error {
	fd, err := getFd(conn)
	if err != nil {
		return err
	}

	sk := s.keys[fd]
	if sk == nil {
		return ErrNotRegistered
	}

	delete(s.keys, fd)
	return s.poll.Delete(fd, sk.interests)
}

func (s *Selector) Modify(conn interface{}, ops int) error {
	fd, err := getFd(conn)
	if err != nil {
		return err
	}

	sk := s.keys[fd]
	if sk == nil {
		return ErrNotRegistered
	}

	return s.poll.Modify(fd, ops)
}

func (s *Selector) Wakeup() error {
	return s.poll.Wakeup()
}

func (s *Selector) Select(ops ...SelectOptions) ([]*SelectionKey, error) {
	s.readyKeys = s.readyKeys[:0]
	var err error
	if len(ops) == 0 {
		err = s.poll.Wait(s, nil, -1)
	} else {
		err = s.poll.Wait(s, ops[0].Callback, ops[0].Timeout)
	}

	if err != nil {
		return nil, err
	}

	return s.readyKeys, nil
}

func (s *Selector) getSelectionKey(fd uintptr) *SelectionKey {
	return s.keys[fd]
}
