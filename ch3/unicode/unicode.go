package main

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
)

func Sloth(a string) (b string) {

	s := a
	el := '…'
	res := make([]rune, 0, 100)
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if !unicode.IsPunct(r) {
			res = append(res, r, el, el)
		} else {
			l := len(res)
			var high int
			high = l - 2
			if high < 0 {
				high = 0
			}
			res = res[:high]
			res = append(res, r)
		}
		i += size
	}
	b = string(res)
	fmt.Println(string(res))

	return
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	s := os.Args[1]
	Sloth(s)
}
