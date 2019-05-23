package main

import (
	"github.com/jeckbjy/nio"
	"log"
	"net"
)

func StartServer() {
	log.Printf("start server")

	selector, err := nio.New()
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", ":6789")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				break
			}

			_, err1 := selector.Add(conn, nio.OP_READ, nil)
			if err1 != nil {
				conn.Close()
				continue
			}
		}
	}()

	for {
		log.Printf("wait\n")
		keys, err := selector.Select()
		if err != nil {
			break
		}

		for _, key := range keys {
			switch {
			case key.Readable():
				bytes := make([]byte, 1024)
				n, err := key.Read(bytes)
				if err != nil {
					log.Printf("read: fd=%+v, err=%+v\n", key.Fd(), err)
					//continue
				} else {
					log.Printf("%s:%+v\n", bytes[:n], key.Fd())
				}
				_, err = key.Write([]byte("pong"))
				if err != nil {
					log.Printf("write: fd=%+v, err=%+v\n", key.Fd(), err)
				}
			}
		}
	}
}
