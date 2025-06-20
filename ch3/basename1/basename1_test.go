package basename1

import "testing"

func TestBaseName(t *testing.T) {
	// Test cases
	tests := []struct {
		input    string
		expected string
	}{
		{"a/b/c.go", "c"},
		{"a/b/c", "c"},
		{"a/b/c.", "c"},
		{"a/b/c..", "c."},
		{"a/b/c.go.txt", "c.go"},
	}

	for _, test := range tests {
		result := basename(test.input)
		if result != test.expected {
			panic("Test failed")
		}
	}
}
