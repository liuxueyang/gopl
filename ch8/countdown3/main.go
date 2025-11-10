package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	abort := make(chan struct{})

	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	// works from Go 1.23
	ticker := time.Tick(1 * time.Second)
	for i := 10; i > 0; i-- {
		select {
		case <-ticker:
			fmt.Println(i)
		case <-abort:
			fmt.Println("Countdown aborted!")
			return
		}
	}
	fmt.Println("Countdown finished!")
}
