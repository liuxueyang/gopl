package basename2

import "strings"

func basename(s string) string {
	idx := strings.LastIndex(s, "/")
	if idx != -1 {
		s = s[idx + 1:]
	}

	if idx = strings.LastIndex(s, "."); idx != -1 {
		s = s[:idx]
	}

	return s
}
