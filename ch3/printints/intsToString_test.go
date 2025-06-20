package printints

import "testing"

func TestIntsToString(t *testing.T) {
	tests := []struct {
		input    []int
		expected string
	}{
		{[]int{}, "[]"},
		{[]int{1}, "[1]"},
		{[]int{1, 2}, "[1, 2]"},
		{[]int{1, 2, 3}, "[1, 2, 3]"},
	}

	for _, test := range tests {
		result := IntsToString(test.input)
		if result != test.expected {
			t.Errorf("IntsToString(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}
