package comma1

import "testing"

func TestComma1(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1234567890", "1,234,567,890"},
		{"123456", "123,456"},
		{"123", "123"},
		{"1", "1"},
		{"", ""},
	}

	for _, test := range tests {
		result := comma1(test.input)
		if result != test.expected {
			t.Errorf("comma1(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}
