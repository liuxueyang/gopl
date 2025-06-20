package main

import "fmt"

func reverse_array(arr *[32]int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func main() {
	var arr [32]int

	for i := range len(arr) {
		arr[i] = i + 1
	}

	reverse_array(&arr)
	fmt.Println(arr)
}
