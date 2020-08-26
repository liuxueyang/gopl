package main

import (
	"fmt"
	"practice/if/counter"
)

func TestCounters() {
	var bc counter.ByteCounter
	fmt.Fprintf(&bc, "123")
	fmt.Printf("bc=%s\n", &bc)

	fmt.Fprintf(&bc, "abcde")
	fmt.Printf("bc=%s\n", &bc)

	var wc counter.WordCounter
	fmt.Fprintf(&wc, "hello")
	fmt.Printf("wc=%s\n", &wc)

	fmt.Fprintf(&wc, "")
	fmt.Printf("wc=%s\n", &wc)

	fmt.Fprintf(&wc, "world w w    ")
	fmt.Printf("wc=%s\n", &wc)

	var lc counter.LineCounter
	fmt.Fprintf(&lc, "hello")
	fmt.Printf("lc=%s\n", &lc)

	fmt.Fprintf(&lc, "")
	fmt.Printf("lc=%s\n", &lc)

	fmt.Fprintf(&lc, `hello

	`)
	fmt.Printf("lc=%s\n", &lc)

	fmt.Fprintf(&lc, `
hello
`)
	fmt.Printf("lc=%s\n", &lc)

}

func main() {

}
