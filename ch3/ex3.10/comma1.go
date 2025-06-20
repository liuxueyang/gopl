package comma1

import "bytes"

func comma1(s string) string {
	var buf bytes.Buffer
	var j = 0
	n := len(s)

	for i := n - 1; i >= 0; i-- {
		if j > 0 && j%3 == 0 {
			buf.WriteByte(',')
		}

		buf.WriteByte(s[i])
		j++
	}

	// Reverse the buffer
	var result = buf.Bytes()
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
