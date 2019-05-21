package main

import (
	"fmt"
	"github.com/jeckbjy/nio"
	"log"
	"net"
)

func StartServer() {
	log.Printf("start server")
	l, err := net.Listen("tcp", ":6789")
	if err != nil {
		panic(err)
	}

	selector, err := nio.New()
	if err != nil {
		panic(err)
	}

	if _, err := selector.Add(l, nio.OP_ACCEPT, nil); err != nil {
		panic(err)
	}

	for {
		keys, err := selector.Select()
		if err != nil {
			break
		}

		for _, key := range keys {
			switch {
			case key.Acceptable():
				log.Printf("accept\n")
				conn, err := key.Listener().Accept()
				if err != nil {
					panic(err)
				}
				selector.Add(conn, nio.OP_READ, nil)
			case key.Readable():
				bytes := make([]byte, 1024)
				n, err := key.Conn().Read(bytes)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", bytes[:n])
				key.Conn().Write([]byte("pong"))
			}
		}
	}
}
