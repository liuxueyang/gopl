package main

import "fmt"

func main() {
	odds := make(chan int)
	msg := make(chan bool)

	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Println(i)
			odds <- i
			<-msg
		}
		close(odds)
	}()

	for i := range odds {
		fmt.Println(i + 1)
		msg <- true
	}
}
