package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	seen := make(map[string]bool)
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		line := input.Text()
		if _, ok := seen[line]; !ok {
			seen[line] = true
			fmt.Println(line)
		}
	}
}
