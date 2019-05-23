package main

import (
	"github.com/jeckbjy/nio"
	"log"
	"net"
	"time"
)

func main() {
	runAccept()
}

func runAccept() {
	l, err := net.Listen("tcp", ":6789")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	fd, _ := nio.GetFd(l)
	nio.SetNonblock(fd, true)
	log.Printf("hope nonblocking")
	l.(*net.TCPListener).AcceptTCP()
	log.Printf("after accept")
}

func runRead() {
	l, err := net.Listen("tcp", ":6789")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		fd, _ := nio.GetFd(conn)
		nio.SetNonblock(fd, true)

		log.Printf("wait read\n")
		bytes := make([]byte, 1024)
		n, err := nio.Read(fd, bytes)
		//n, err := conn.Read(bytes)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("read:%s\n", bytes[:n])
	}()

	conn, err := net.Dial("tcp", "localhost:6789")
	if err != nil {
		log.Fatal(err)
	}

	fd, _ := nio.GetFd(conn)

	nio.SetNonblock(fd, true)

	log.Printf("sleep for write")
	time.Sleep(time.Second)
	log.Printf("do write")
	nio.Write(fd, []byte("ping"))
	//conn.Write([]byte("ping"))
	time.Sleep(time.Second)
}
