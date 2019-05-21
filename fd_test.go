package nio

import (
	"net"
	"os"
	"testing"
)

// filer describes an object that has ability to return os.File.
type filer interface {
	// File returns a copy of object's file descriptor.
	File() (*os.File, error)
}

func TestFD(t *testing.T) {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		t.Error(err)
		return
	}

	defer l.Close()

	f, err := l.(filer).File()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("listen fd:%+v", f.Fd())

	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		t.Error(err)
		return
	}

	f1, err := conn.(filer).File()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("channel fd:%+v", f1.Fd())
}
