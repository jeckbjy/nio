package main

import "flag"

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
