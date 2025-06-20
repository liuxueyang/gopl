package main

import "unicode"

func eliminate_space(lst []byte) []byte {
	if len(lst) == 0 {
		return lst
	}

	s := string(lst)
	su := []rune(s)

	i := 1
	for j := 1; j < len(su); j++ {
		if unicode.IsSpace(su[j]) && unicode.IsSpace(su[i-1]) {
			su[i-1] = ' '
			continue
		}
		su[i] = su[j]
		i++
	}
	su = su[:i]
	s = string(su)
	return []byte(s)
}

func main() {
	// Example usage
	input := []byte("This  is   a   test   string.")
	result := eliminate_space(input)
	println(string(result)) // Output: "This is a test string."

	test_eliminate_space()
	test_eliminate_space_with_utf8()
}

// Write more tests for the eliminate_space function. Include cases with UTF-8 characters, multiple spaces, and no spaces.
func test_eliminate_space() {
	tests := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte("This  is   a   test   string."), []byte("This is a test string.")},
		{[]byte("  Leading and trailing spaces  "), []byte(" Leading and trailing spaces ")},
		{[]byte("NoSpacesHere"), []byte("NoSpacesHere")},
		{[]byte("   Multiple    spaces   in between   "), []byte(" Multiple spaces in between ")},
		{[]byte("  \t\n  "), []byte(" ")}, // Test with whitespace characters
	}

	for _, test := range tests {
		result := eliminate_space(test.input)
		if string(result) != string(test.expected) {
			println("Test failed for input:", string(test.input))
			println("Expected:", string(test.expected), "Got:", string(result))
		} else {
			println("Test passed for input:", string(test.input))
		}
	}
}

// Write more tests for the eliminate_space function. Include cases with Chinese characters, multiple spaces, and no spaces.
func test_eliminate_space_with_utf8() {
	tests := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte("This  is   a   test   string."), []byte("This is a test string.")},
		{[]byte("  Leading and trailing spaces  "), []byte(" Leading and trailing spaces ")},
		{[]byte("NoSpacesHere"), []byte("NoSpacesHere")},
		{[]byte("   Multiple    spaces   in between   "), []byte(" Multiple spaces in between ")},
		{[]byte("  \t\n  "), []byte(" ")},   // Test with whitespace characters
		{[]byte("你好  世界"), []byte("你好 世界")}, // Chinese characters with space
	}

	for _, test := range tests {
		result := eliminate_space(test.input)
		if string(result) != string(test.expected) {
			println("Test failed for input:", string(test.input))
			println("Expected:", string(test.expected), "Got:", string(result))
		} else {
			println("Test passed for input:", string(test.input))
		}
	}
}
