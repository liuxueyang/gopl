package comma

import "testing"

func TestComma(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123", "123"},
		{"1234", "1,234"},
		{"12345", "12,345"},
		{"123456", "123,456"},
		{"1234567", "1,234,567"},
	}

	for _, test := range tests {
		result := comma(test.input)
		if result != test.expected {
			t.Errorf("comma(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
