package basename2

import "testing"

func TestBasename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a/b/c.go", "c"},
		{"a/b/c", "c"},
		{"a/b/c.", "c"},
		{"a/b/c.go.txt", "c.go"},
	}

	for _, test := range tests {
		result := basename(test.input)
		if result != test.expected {
			t.Errorf("basename(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
