package counter

import (
	"fmt"
	"testing"
)

func TestByteCounter(t *testing.T) {
	var bc ByteCounter
	_, _ = fmt.Fprintf(&bc, "123")
	fmt.Printf("bc=%s\n", &bc)

	_, _ = fmt.Fprintf(&bc, "abcde")
	fmt.Printf("bc=%s\n", &bc)
}

func TestWordCounter(t *testing.T) {
	var wc WordCounter
	_, _ = fmt.Fprintf(&wc, "hello")
	fmt.Printf("wc=%s\n", &wc)

	_, _ = fmt.Fprintf(&wc, "")
	fmt.Printf("wc=%s\n", &wc)

	_, _ = fmt.Fprintf(&wc, "world w w    ")
	fmt.Printf("wc=%s\n", &wc)
}

func TestLineCounter(t *testing.T) {
	var lc LineCounter
	_, _ = fmt.Fprintf(&lc, "hello")
	fmt.Printf("lc=%s\n", &lc)

	_, _ = fmt.Fprintf(&lc, "")
	fmt.Printf("lc=%s\n", &lc)

	_, _ = fmt.Fprintf(&lc, `hello

	`)
	fmt.Printf("lc=%s\n", &lc)

	_, _ = fmt.Fprintf(&lc, `
hello
`)
	fmt.Printf("lc=%s\n", &lc)
}
