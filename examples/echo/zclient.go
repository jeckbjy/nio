package main

import (
	"log"
	"net"

	"github.com/jeckbjy/nio"
)

func StartClient() {
	log.Printf("start client")

	selector, err := nio.New()
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", "localhost:6789")
	if err != nil {
		panic(err)
	}

	selector.Add(conn, nio.OP_READ, nil)

	conn.Write([]byte("ping"))

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
				_, err = key.Write([]byte("ping"))
				if err != nil {
					log.Printf("write: fd=%+v, err=%+v\n", key.Fd(), err)
				}
			}
		}
	}
}
