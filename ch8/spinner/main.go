package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	res := fib(n)
	fmt.Printf("\rFibonacci(%d) = %d\n", n, res)
	// The spinner will continue to run until the program exits.
	// To see the spinner in action, you can run the program and wait for it to
	// complete the Fibonacci calculation.
	// The spinner will show a rotating character while the calculation is in progress.
	// You can stop the program by pressing Ctrl+C.
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
