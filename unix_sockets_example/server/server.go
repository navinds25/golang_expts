package main

import (
	"net"
	"log"
	"fmt"
)

func dataHandler(c net.Conn) {

	buf := make([]byte, 512)
	nr, err := c.Read(buf)

	if err != nil {
		return
	}

	data := string(buf[0:nr])
	fmt.Println(data)
}

func main() {
	l, err := net.Listen("unix", "/tmp/example.sock")
	if err != nil{
		log.Fatal(err)
		return
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		dataHandler(fd)
		fd.Close()
		l.Close()
	}
}
