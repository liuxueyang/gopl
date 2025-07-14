package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleConn(conn)
	}
}

const timeoutDuration = 10 * time.Second

// using select statement, add a timeout to the echo server so that it
// disconnects any client that does not send a message within 10 seconds
func handleConn(c net.Conn) {
	defer c.Close()

	timeout := time.NewTimer(timeoutDuration)
	defer timeout.Stop()

	input := bufio.NewScanner(c)
	inputChan := make(chan string)

	go func() {
		defer close(inputChan)
		for input.Scan() {
			inputChan <- input.Text()
		}
	}()

	for {
		select {
		case <-timeout.C:
			fmt.Fprintln(c, "Timeout: No activity for 10 seconds")
			return
		case text, ok := <-inputChan:
			if !ok {
				return
			}

			// If we reach here, it means input was received before the timeout
			timeout.Stop()

			// Reset the timeout timer
			timeout.Reset(timeoutDuration)
			go echo(c, text, 1*time.Second)
		}
	}
}
