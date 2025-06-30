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

	fmt.Println("Press Enter to abort the countdown...")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for countdown := 10; countdown > 0; countdown-- {
		select {
		case <-abort:
			fmt.Println("Countdown aborted!")
			return
		case <-ticker.C:
			fmt.Println(countdown)
		}
	}
	fmt.Println("Countdown finished!")
}
