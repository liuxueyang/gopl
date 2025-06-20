package anagrams

import "testing"

func TestAnagrams(t *testing.T) {
	tests := []struct {
		s, t string
		want bool
	}{
		{"listen", "silent", true},
		{"triangle", "integral", true},
		{"apple", "pale", false},
		{"rat", "car", false},
	}

	for _, test := range tests {
		got := anagrams(test.s, test.t)
		if got != test.want {
			t.Errorf("anagrams(%q, %q) = %v; want %v", test.s, test.t, got, test.want)
		}
	}
}
