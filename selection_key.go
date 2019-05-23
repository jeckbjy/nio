package nio

import (
	"net"
)

/*
type Channel struct {
	fd uintptr
	interests int
	ready int
}

TCPListener
TCPConn
UnixListener
UnixConn
UDPConn
*/
// TODO: SelectionKey 改成Channel,实现net.Listener和net.Conn
// Read:读需要全部读完
// Write:自动维护状态,自动缓存未写完数据
// 上层只需监听读,写由库维护
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

func (sk *SelectionKey) Fd() uintptr {
	return sk.fd
}

//func (sk *SelectionKey) Acceptable() bool {
//	return sk.ready&OP_ACCEPT != 0
//}

func (sk *SelectionKey) Readable() bool {
	return sk.ready&OP_READ != 0
}

func (sk *SelectionKey) Writable() bool {
	return sk.ready&OP_WRITE != 0
}

func (sk *SelectionKey) isInterests(ops int) bool {
	return sk.interests&ops != 0
}

func (sk *SelectionKey) setReadyIn() {
	if sk.isInterests(OP_READ) {
		sk.ready |= OP_READ
	}
	//else if sk.isInterests(OP_ACCEPT) {
	//	sk.ready |= OP_ACCEPT
	//}
}

func (sk *SelectionKey) setReadyOut() {
	if sk.isInterests(OP_WRITE) {
		sk.ready |= OP_WRITE
	}
}

func (sk *SelectionKey) InterestOps() int {
	return sk.interests
}

func (sk *SelectionKey) ReadyOps() int {
	return sk.ready
}

func (sk *SelectionKey) Channel() interface{} {
	return sk.channel
}

func (sk *SelectionKey) Data() interface{} {
	return sk.data
}

func (sk *SelectionKey) Accept() (net.Conn, error) {
	return nil, ErrNotSupport
}

func (sk *SelectionKey) Read(b []byte) (int, error) {
	return Read(sk.fd, b)
}

func (sk *SelectionKey) Write(b []byte) (int, error) {
	return Write(sk.fd, b)
}
