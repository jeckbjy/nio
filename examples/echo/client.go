package main

import (
	"fmt"
	"github.com/jeckbjy/nio"
	"log"
	"net"
)

func StartClient() {
	log.Printf("start client")
	conn, err := net.Dial("tcp", "localhost:6789")
	if err != nil {
		panic(err)
	}

	selector, err := nio.New()
	if err != nil {
		panic(err)
	}

	selector.Add(conn, nio.OP_READ, nil)

	conn.Write([]byte("ping"))

	for {
		keys, err := selector.Select()
		if err != nil {
			break
		}

		for _, key := range keys {
			switch {
			case key.Readable():
				bytes := make([]byte, 1024)
				n, err := key.Conn().Read(bytes)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", bytes[:n])
				key.Conn().Write([]byte("ping"))
			}
		}
	}
}
