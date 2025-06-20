package comma2

import "testing"

func TestComma2(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1234567890", "1,234,567,890"},
		{"-1234567890", "-1,234,567,890"},
		{"+1234567890", "+1,234,567,890"},
		{"1234567.89", "1,234,567.89"},
		{"-1234567.89", "-1,234,567.89"},
		{"+1234567.89", "+1,234,567.89"},
	}

	for _, test := range tests {
		result := comma2(test.input)
		if result != test.expected {
			t.Errorf("comma2(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
