package main

import "fmt"

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func rotate_left(s []int, n int) {
	if len(s) == 0 || n <= 0 {
		return
	}

	n = n % len(s)
	reverse(s[:n])
	reverse(s[n:])
	reverse(s)
}

func main() {
	a := [...]int{1, 2, 3, 4, 5}
	reverse(a[:])
	fmt.Println(a)

	b := [...]int{1, 2, 3, 4, 5}
	rotate_left(b[:], 2)
	fmt.Println(b)
}
