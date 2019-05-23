package nio

import (
	"fmt"
	"net"
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

	// ErrNotSupport not support
	ErrNotSupport = fmt.Errorf("not support")
)

const (
	//OP_ACCEPT = 0x01
	//OP_CONNECT = 0x02 // not support
	OP_READ  = 0x04
	OP_WRITE = 0x08
	op_IN    = OP_READ
	op_OUT   = OP_WRITE
)

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

func (s *Selector) Wakeup() error {
	return s.poll.Wakeup()
}

func (s *Selector) Add(channel interface{}, ops int, data interface{}) (*SelectionKey, error) {
	fd, err := GetFd(channel)
	if err != nil {
		return nil, err
	}

	if s.keys[fd] != nil {
		return nil, ErrRegistered
	}

	if err := SetNonblock(fd, true); err != nil {
		return nil, err
	}

	if err := s.poll.Add(fd, ops); err != nil {
		return nil, err
	}

	sk := &SelectionKey{channel: channel, data: data, fd: fd, interests: ops}
	s.keys[fd] = sk
	return sk, nil
}

func (s *Selector) Delete(channel interface{}) error {
	sk, err := s.getSelectionKey(channel)
	if err != nil {
		return err
	}

	delete(s.keys, sk.fd)
	return s.poll.Delete(sk.fd, sk.interests)
}

func (s *Selector) Modify(channel interface{}, ops int) error {
	sk, err := s.getSelectionKey(channel)
	if err != nil {
		return err
	}

	if ops == sk.interests {
		return nil
	}

	old := sk.interests
	sk.interests = ops

	return s.poll.Modify(sk.fd, old, ops)
}

// ModifyXOR 切换某个状态开关,通常用于OP_WRITE状态控制
func (s *Selector) ModifyXOR(channel interface{}, ops int) error {
	if ops == 0 {
		return nil
	}

	sk, err := s.getSelectionKey(channel)
	if err != nil {
		return err
	}

	old := sk.interests
	sk.interests ^= ops

	return s.poll.Modify(sk.fd, old, sk.interests)
}

// ModifyIf 添加或删除某个状态
func (s *Selector) ModifyIf(channel interface{}, ops int, add bool) error {
	if ops == 0 {
		return nil
	}

	sk, err := s.getSelectionKey(channel)
	if err != nil {
		return err
	}

	old := sk.interests
	if add {
		sk.interests |= ops
	} else {
		sk.interests &^= ops
	}

	return s.poll.Modify(sk.fd, old, sk.interests)
}

func (s *Selector) getSelectionKey(channel interface{}) (*SelectionKey, error) {
	fd, err := GetFd(channel)
	if err != nil {
		return nil, err
	}

	sk := s.keys[fd]
	if sk == nil {
		return nil, ErrNotRegistered
	}

	return sk, nil
}

// TODO:自己实现net.Listener
func (s *Selector) Listen(network, address string, options ...interface{}) (net.Listener, error) {
	return net.Listen(network, address)
	//l, err := net.Listen(network, address)
	//if err != nil {
	//	return l, err
	//}
	//
	//var data interface{}
	//if len(options) > 0 {
	//	data = options[0]
	//}
	//
	//_, err1 := s.Add(l, OP_ACCEPT, data)
	//if err1 != nil {
	//	return nil, err1
	//}
	//
	//return l, err
}

func (s *Selector) Dial(network, address string, options ...interface{}) (net.Conn, error) {
	return net.Dial(network, address)
}
