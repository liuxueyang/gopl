package counter

import (
	"bufio"
	"strconv"
)

type LineCounter int

func (l *LineCounter) Write(p []byte) (int, error) {
	n := len(p)

	advance, _, err := bufio.ScanLines(p, true)
	var cnt int
	for advance > 0 {
		cnt++
		p = p[advance:]
		if len(p) > 0 {
			advance, _, err = bufio.ScanLines(p, true)
		} else {
			break
		}
	}
	*l += LineCounter(cnt)

	return n, err
}

func (l *LineCounter) String() string {
	return strconv.FormatInt(int64(*l), 10)
}
