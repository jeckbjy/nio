package main

import (
	"log"
	"net"

	"github.com/jeckbjy/nio"
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
				log.Printf("%s:%+v\n", bytes[:n], key.Fd())
				key.Conn().Write([]byte("ping"))
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

// func StartClient() {
// 	log.Printf("start client")
// 	conn, err := net.Dial("tcp", "localhost:6789")
// 	if err != nil {
// 		panic(err)
// 	}

// 	selector, err := nio.New()
// 	if err != nil {
// 		panic(err)
// 	}

// 	selector.Add(conn, nio.OP_READ, nil)

// 	conn.Write([]byte("ping"))

// 	// buffer := make([]byte, 1024*1024*20)
// 	// for i := 0; i < len(buffer); i += 4 {
// 	// 	copy(buffer[i:i+4], []byte("ping"))
// 	// }
// 	// offset := -1
// 	// counts := 0

// 	for {
// 		fmt.Printf("wait\n")
// 		keys, err := selector.Select()
// 		if err != nil {
// 			panic(err)
// 			break
// 		}

// 		fmt.Printf("events:%+v\n", len(keys))

// 		for _, key := range keys {
// 			switch {
// 			case key.Readable():
// 				//fmt.Printf("begin read\n")
// 				bytes := make([]byte, 1024)
// 				n, err := key.Conn().Read(bytes)
// 				if err != nil {
// 					break
// 				}
// 				fmt.Printf("data:%s\n", bytes[:n])
// 				key.Conn().Write([]byte("ping"))
// 				//for {
// 				//	n, err := key.Conn().Read(bytes)
// 				//	if err != nil {
// 				//		break
// 				//	}
// 				//	if n < 1024 {
// 				//		fmt.Printf("read data:%+s\n", bytes[:n])
// 				//		break
// 				//	} else {
// 				//		// read all
// 				//		fmt.Printf("overflow:%s\n", bytes)
// 				//	}
// 				//}

// 				//fmt.Printf("after read\n")

// 				// if offset == -1 {
// 				// 	key.Conn().Write([]byte("ping"))
// 				// 	counts++
// 				// 	if counts == 100 {
// 				// 		fmt.Printf("begin write big data\n")
// 				// 		counts = 0
// 				// 		offset = 0
// 				// 		o, err := key.Conn().Write(buffer)
// 				// 		if err != nil {
// 				// 			fmt.Printf("err:%+v\n", err)
// 				// 			break
// 				// 		}
// 				// 		offset = o
// 				// 		if offset != len(buffer) {
// 				// 			selector.ModifyIf(key.Conn(), nio.OP_WRITE, true)
// 				// 		}
// 				// 	}
// 				// }
// 				// case key.Writable():
// 				// 	fmt.Printf("can write")
// 				// 	off, err := key.Conn().Write(buffer[offset:])
// 				// 	if err != nil {
// 				// 		fmt.Printf("err:%+v\n", err)
// 				// 		break
// 				// 	}
// 				// 	offset += off
// 				// 	if offset == len(buffer) {
// 				// 		selector.ModifyIf(key.Conn(), nio.OP_WRITE, false)
// 				// 	}
// 			}
// 		}
// 	}
// }

// //func ReadAll(reader io.Reader) ([]byte, error) {
// //	var buf bytes.Buffer
// //	var err error
// //	count := 0
// //	data := make([]byte, 1024)
// //	for {
// //		n, e := reader.Read(data)
// //		if e != nil {
// //			if e != io.EOF {
// //				err = e
// //			}
// //			break
// //		}
// //	}
// //}

// //func ReadAll(reader io.Reader) (b []byte, err error) {
// //
// //	//var buf bytes.Buffer
// //	//// If the buffer overflows, we will get bytes.ErrTooLarge.
// //	//// Return that as an error. Any other panic remains.
// //	//defer func() {
// //	//	e := recover()
// //	//	if e == nil {
// //	//		return
// //	//	}
// //	//	if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
// //	//		err = panicErr
// //	//	} else {
// //	//		panic(e)
// //	//	}
// //	//}()
// //	//
// //	//buf.Grow(512)
// //	//
// //	//
// //	//_, err = buf.ReadFrom(r)
// //	//return buf.Bytes(), err
// //}
