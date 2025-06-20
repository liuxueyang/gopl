package anagrams

import "slices"

func anagrams(s string, t string) bool {
	sb := []byte(s)
	tb := []byte(t)

	if len(sb) != len(tb) {
		return false
	}

	slices.Sort(sb)
	slices.Sort(tb)

	for i := range sb {
		if sb[i] != tb[i] {
			return false
		}
	}

	return true
}
