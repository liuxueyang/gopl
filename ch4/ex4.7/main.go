package main

func reverse_utf8(s []byte) []byte {
	if len(s) == 0 {
		return s
	}

	lst := []rune(string(s))
	for i, j := 0, len(lst)-1; i < j; i, j = i+1, j-1 {
		lst[i], lst[j] = lst[j], lst[i]
	}
	return []byte(string(lst))
}

func main() {
	// Example usage
	s := []byte("Hello, 世界")
	reversed := reverse_utf8(s)
	println(string(reversed)) // Output: "界世 ,olleH"
}
