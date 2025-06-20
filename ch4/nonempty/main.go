package main

import "fmt"

func main() {
	// Example usage
	strings := []string{"hello", "", "world", "", "golang"}
	result := nonempty(strings)
	fmt.Println(result)
}

func nonempty(lst []string) []string {
	i := 0
	for _, s := range lst {
		if len(s) > 0 {
			lst[i] = s
			i++
		}
	}
	return lst[:i]
}
