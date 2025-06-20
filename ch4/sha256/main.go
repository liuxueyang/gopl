package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))

	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)

	s := "2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881"
	fmt.Printf("%d\n", len(s)/2)
}
