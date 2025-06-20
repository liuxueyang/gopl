package comma2

import "bytes"
import "strings"

func comma2(s string) string {
	var buf bytes.Buffer

	var pre string

	if s[0] == '-' || s[0] == '+' {
		pre += string(s[0])
		s = s[1:]
	}

	var last string

	if idx := strings.Index(s, "."); idx != -1 {
		last = s[idx:]
		s = s[:idx]
	}

	j := 0
	for i := len(s) - 1; i >= 0; i-- {
		if j > 0 && j % 3 == 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(s[i])
		j++
	}

	mid := buf.Bytes()
	for i, j := 0, len(mid) - 1; i < j; i, j = i + 1, j - 1 {
		mid[i], mid[j] = mid[j], mid[i]
	}

	result := pre + string(mid) + last
	
	return result
}
