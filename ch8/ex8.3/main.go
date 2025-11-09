package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	tcpAddr := net.TCPAddr{
		Port: 8080,
	}

	conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Printf("io.Copy error: %v", err)
		}
		log.Println("Connection closed by server")
		done <- struct{}{}
	}()

	mustCopy(conn, os.Stdin)
	conn.CloseWrite()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
