package main

import (
	"fmt"
	"time"
)

func main() {
	ruby := make(chan string)

	go stage(ruby)

	ruby <- "Hello, Ruby!"
	fmt.Println(<-ruby)
	close(ruby) // ensure to close the channel to prevent goroutine leak
	time.Sleep(1 * time.Second)
	select {}
}

func stage(ruby chan string) {
	for message := range ruby {
		fmt.Println(message)
		ruby <- "Received: " + message
	}
	fmt.Println("Stage goroutine exiting")
}
