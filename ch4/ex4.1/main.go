package main

import (
	"crypto/sha256"
	"fmt"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func count_bits_sha256(arr1 *[32]uint8, arr2 *[32]uint8) uint {
	var ans uint = 0

	for i := range arr1 {
		ans += uint(pc[arr1[i]^arr2[i]])
	}
	return ans
}

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("a"))
	diff_bits := count_bits_sha256(&c1, &c2)

	fmt.Printf("diff_bits = %d\n", diff_bits)
	// for i := range pc {
	// 	fmt.Printf("%d %d\n", i, pc[i])
	// }
}
