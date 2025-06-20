package printints

import (
	"bytes"
	"fmt"
)

func IntsToString(ints []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')

	n := len(ints)

	for i, v := range ints {
		fmt.Fprintf(&buf, "%d", v)

		if i != n-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteByte(']')

	return buf.String()
}
