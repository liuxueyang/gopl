package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var portNumber *int = flag.Int("port", 8080, "server port number")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *portNumber))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
