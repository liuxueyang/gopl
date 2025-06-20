package main

import "fmt"

func main() {
	s := "This is a $test string with $variables."
	fmt.Println("Original string:", s)

	// Using the myf function to format variables
	result := expand(s, myf)
	fmt.Println("Expanded string:", result)
}

func myf(s string) string {
	return fmt.Sprintf("[%s]", s)
}

func expand(s string, f func(string) string) string {
	if len(s) == 0 {
		return s
	}

	var result string

	for i := 0; i < len(s); i++ {
		if s[i] == '$' {
			j := i + 1
			var variable string
			for j < len(s) && (s[j] >= 'a' && s[j] <= 'z' || s[j] >= 'A' && s[j] <= 'Z' || s[j] >= '0' && s[j] <= '9') {
				variable += string(s[j])
				j++
			}
			result += f(variable)
			i = j - 1
		} else {
			result += string(s[i])
		}
	}

	return result
}
