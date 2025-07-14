package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type Client struct {
	name string
	ch   chan<- string
}

var (
	entering = make(chan Client)
	leaving  = make(chan Client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[Client]bool)

	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}
		case cli := <-entering:
			clients[cli] = true

			for other := range clients {
				if other != cli {
					cli.ch <- other.name + " has joined the chat"
				}
			}
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

const timeoutDuration = 10 * time.Minute
const inputPrompt = "Enter your nickname: "

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	ch <- inputPrompt
	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) // 设置初始读取超时

	input := bufio.NewScanner(conn)
	if !input.Scan() { // Read the client's name
		conn.Close()
		return
	}

	who := input.Text()

	conn.SetReadDeadline(time.Time{}) // 清除读取超时

	ch <- "You are " + who

	messages <- who + " has arrived"

	cli := Client{name: who, ch: ch}
	entering <- cli

	inputChan := make(chan string)

	go func() {
		defer close(inputChan)
		for input.Scan() {
			inputChan <- input.Text()
		}
	}()

	timeout := time.NewTimer(timeoutDuration)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			messages <- fmt.Sprintf("%s has been disconnected due to inactivity.", who)
			leaving <- cli
			conn.Close()
			return
		case text, ok := <-inputChan:
			if !ok {
				leaving <- cli
				messages <- who + " has left"
				conn.Close()
				return
			}

			timeout.Stop()                 // Stop the timer if input is received
			timeout.Reset(timeoutDuration) // Reset the timer
			messages <- who + ": " + text
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		if msg == inputPrompt {
			fmt.Fprint(conn, msg) // Use Fprint to avoid adding a newline
		} else {
			fmt.Fprintln(conn, msg)
		}
	}
}
