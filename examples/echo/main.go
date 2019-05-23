package main

import "flag"

// 例子还有问题,经常会出现bad file descriptor错误
func main() {
	var client bool
	flag.BoolVar(&client, "client", false, "run client")
	flag.Parse()

	if client {
		StartClient()
	} else {
		StartServer()
	}
}
