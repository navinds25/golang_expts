package main

import (
	"net"
	"log"
)

func main() {
	arg := "ping"
	c, err := net.Dial("unix", "/tmp/example.sock")

	if err != nil {
		panic(err)
	}

	_, err = c.Write([]byte(arg))

	if err != nil {
		log.Println(err)
	}
}
