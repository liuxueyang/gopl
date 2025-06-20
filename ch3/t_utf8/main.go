package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	var x1 = (1 << 7) - 1
	fmt.Println(x1)

	var x2 = 0b110_00000_10_000000
	fmt.Println(x2)
	var x3 = 0b110_11111_10_111111
	fmt.Println(x3)

	s := "Hello, 世界"
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%d: %c\n", i, r)
		i += size
	}
}
