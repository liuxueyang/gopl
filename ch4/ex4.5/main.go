package main

func eliminate_duplicates(arr []string) []string {
	if len(arr) <= 1 {
		return arr
	}

	i := 1
	for j := 1; j < len(arr); j++ {
		if arr[j] == arr[i-1] {
			continue
		}
		arr[i] = arr[j]
		i++
	}

	return arr[:i]
}

func main() {
	strings := []string{"a", "b", "b", "c", "c", "c", "d"}
	result := eliminate_duplicates(strings)
	for _, v := range result {
		print(v, " ")
	}
	println()
}
