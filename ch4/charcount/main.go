package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	input := bufio.NewReader(os.Stdin)
	for {
		r, sz, err := input.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && sz == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[sz]++
	}

	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}

	fmt.Printf("\nlength\tcount\n")
	for i, n := range utflen {
		if n > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Fprintf(os.Stderr, "invalid UTF-8 characters: %d\n", invalid)
		os.Exit(1)
	}
}
