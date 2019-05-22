package nio

import (
	"net"
	"testing"

	"github.com/jeckbjy/nio/internal"
)

func TestNonBlock(t *testing.T) {
	l, err := net.Listen("tcp", ":6789")
	if err != nil {
		t.Fatalf("%+v\n", err)
	}

	fd, err := getFd(l)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}

	// SetNonblock no effect
	err = internal.SetNonblock(fd, true)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}

	t.Logf("hope non blocking\n")
	_, err = l.Accept()
	if err != nil {
		t.Fatalf("err=%+v\n", err)
	}
	t.Logf("finish")
}
