package main

import (
	"log"
	"net"

	"github.com/jeckbjy/nio"
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
				log.Printf("%s:%+v\n", bytes[:n], key.Fd())
				key.Conn().Write([]byte("pong"))
			}
		}
	}
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net"

// 	"github.com/jeckbjy/nio"
// )

// func StartServer() {
// 	log.Printf("start server")
// 	l, err := net.Listen("tcp", ":6789")
// 	if err != nil {
// 		panic(err)
// 	}

// 	selector, err := nio.New()
// 	if err != nil {
// 		panic(err)
// 	}

// 	if _, err := selector.Add(l, nio.OP_ACCEPT, nil); err != nil {
// 		panic(err)
// 	}

// 	for {
// 		keys, err := selector.Select()
// 		if err != nil {
// 			panic(err)
// 			break
// 		}

// 		for _, key := range keys {
// 			switch {
// 			case key.Acceptable():
// 				log.Printf("accept\n")
// 				//for {
// 				fmt.Printf("aaa\n")
// 				// Accept依然会阻塞?
// 				conn, err := key.Listener().Accept()
// 				if err != nil {
// 					//panic(err)
// 					fmt.Printf("err:%+v\n", err)
// 					break
// 				}
// 				fmt.Printf("add conn\n")
// 				selector.Add(conn, nio.OP_READ, nil)
// 				//}
// 				fmt.Printf("accept finish\n")
// 			case key.Readable():
// 				bytes := make([]byte, 1024)
// 				n, err := key.Conn().Read(bytes)
// 				if err != nil {
// 					panic(err)
// 				}
// 				fmt.Printf("read:%s,send:pong\n", bytes[:n])
// 				key.Conn().Write([]byte("pong"))
// 			}
// 		}
// 	}
// }
