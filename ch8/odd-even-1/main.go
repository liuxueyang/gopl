package main

import "fmt"

func main() {
	ch := make(chan int, 1)

	for i := 1; i < 10; i++ {
		select {
		case ch <- i:
		case val := <-ch:
			fmt.Println(val)
			fmt.Println(i)
		}
	}
}
