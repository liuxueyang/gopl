package main

import (
	"io"
	"log"
	"net"
	"os"
)

// This program is a simple TCP client, similar to a basic version of netcat. Here’s what it does:

// Connects to a TCP server at localhost:8080.
// Starts a goroutine to copy data from the server (conn) to the standard output (os.Stdout).
// When the server closes the connection, it logs a message and signals the main goroutine via a channel.
// In the main goroutine, it copies data from standard input (os.Stdin) to the server (conn) using the mustCopy function.
// Closes the connection after input is done, then waits for the background goroutine to finish.
// Summary:
// It lets you type input and send it to the server, while also displaying any data received from the server.
// It handles both directions concurrently and exits cleanly when either side closes the connection.

func main() {
	tcp_addr := net.TCPAddr{
		Port: 8080,
	}

	conn, err := net.DialTCP("tcp", nil, &tcp_addr)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		// Closing the read half causes the background goroutine's call to
		// io.Copy to return a "read from closed connection" error.
		// mustCopy(os.Stdout, conn)

		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			// we ignore the error returned by io.Copy
			log.Printf("io.Copy error: %v", err)
		}
		log.Println("Connection closed by server")
		conn.CloseRead()
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
