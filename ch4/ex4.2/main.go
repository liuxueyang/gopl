package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

var sha_number = flag.Int("s", 256, "sha algorithm: 256, 384, 512")

func main() {
	flag.Parse()
	if *sha_number != 256 && *sha_number != 384 && *sha_number != 512 {
		fmt.Fprintf(os.Stderr, "Invalid SHA number: %d. Use 256, 384, or 512.\n", *sha_number)
		os.Exit(1)
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if len(line) == 0 {
			continue
		}

		if *sha_number == 256 {
			// Create a new SHA-256 hash
			output := sha256.Sum256([]byte(line))
			// Print the SHA-256 hash in hexadecimal format
			fmt.Printf("%x\n", output)
		} else if *sha_number == 384 {
			// Create a new SHA-384 hash
			output := sha512.Sum384([]byte(line))
			// Print the SHA-384 hash in hexadecimal format
			fmt.Printf("%x\n", output)
		} else if *sha_number == 512 {
			// Create a new SHA-512 hash
			output := sha512.Sum512([]byte(line))
			// Print the SHA-512 hash in hexadecimal format
			fmt.Printf("%x\n", output)
		}
	}
}
