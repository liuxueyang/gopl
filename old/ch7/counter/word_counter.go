package counter

import (
	"bufio"
	"strconv"
)

type WordCounter int

func (w *WordCounter) Write(p []byte) (int, error) {
	var cnt int
	n := len(p)
	advance, token, err := bufio.ScanWords(p, true)
	for len(token) > 0 {
		cnt++
		p = p[advance:]
		advance, token, err = bufio.ScanWords(p, true)
	}
	*w += WordCounter(cnt)

	return n, err
}

func (w WordCounter) String() string {
	return strconv.FormatInt(int64(w), 10)
}
